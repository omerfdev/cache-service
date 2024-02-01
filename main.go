package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var collection *mongo.Collection // MongoDB koleksiyonunu global olarak tanımlayın
func connectToMongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("cachedb").Collection("customer")
}

func init() {
	connectToMongoDB()
}

// Veriyi önbellekte saklamak için bir map kullanacağız
var cache = make(map[string]interface{})

// Anahtar ve değer çiftini saklamak için bir veri yapısı
type Data struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}
// Anahtar değerine karşılık gelen veriyi önce önbellekte arar, ardından önbellekte bulunamazsa veritabanından alır
func getValue(key string) (interface{}, error) {
	// Önce önbellekte anahtarı kontrol edin
	if val, ok := cache[key]; ok {
		fmt.Println("cache getirdi.")
		return val, nil
	}

	// Önbellekte bulunamazsa veritabanından al
	var result Data
	filter := bson.M{"key": key}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		fmt.Println("db getirdi.")
		return nil, err
	}

	// Değeri önbelleğe ekleyin
	// addToCache(result.Key, result.Value, 5*time.Second) // Önbellekte 5 dakika sakla
	fmt.Println("cache eklendi.")
	return result.Value, nil
}

// Anahtar ve değer çiftini önbelleğe ekler
func addToCache(key string, value interface{}, duration time.Duration) {
	cache[key] = value
	// Belirli bir süre sonra önbellekten kaldır
	go func() {
		time.Sleep(duration)
		delete(cache, key)
	}()
}

// Anahtar ve değeri döndüren endpoint
func getValueHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	key := params["key"]
	value,_ := getValue(key)
	if value == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// 2 adet veri geri döndü interface ve error  err yazarsam kullanmam lazım kullanmak istemiyorsam _  yazmam yeterlidir.
	json.NewEncoder(w).Encode(value)
}

// Anahtar ve değeri önbelleğe ekleyen endpoint
func setValueHandler(w http.ResponseWriter, r *http.Request) {
	var data Data
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Önbelleğe ekle
	addToCache(data.Key, data.Value, 5*time.Second) // Önbellekte 5 dakika sakla
	// MongoDB'ye kaydet
	_, err = collection.InsertOne(context.Background(), data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func main() {
	r := mux.NewRouter()

	// Endpoint tanımlamaları
	r.HandleFunc("/get/{key}", getValueHandler).Methods("GET")
	r.HandleFunc("/set", setValueHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}
