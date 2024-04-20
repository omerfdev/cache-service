# Önbellek Servisi

Bu basit Go programı, önbellekteki verilerin depolanması ve çağrılması için bir HTTP API sağlar. Veriler önce önbellekte aranır, bulunamazsa MongoDB veritabanından alınır ve önbelleğe eklenir.

## Nasıl Çalışır

1. **Kurulum**: Projeyi klonlayın ve Go'nun ve MongoDB'nin yüklü olduğundan emin olun.
2. **Veritabanı Ayarları**: `connectToMongoDB` fonksiyonunu kullanarak MongoDB'ye bağlanın ve `cachedb` veritabanı ile `customer` koleksiyonunu tanımlayın.
3. **Servisi Başlatın**: `go run main.go` komutuyla servisi başlatın.
4. **Veri Ekleme ve Alma**: Tarayıcıdan veya bir API test aracından `/set` endpoint'ine bir POST isteği yaparak veri ekleyin veya `/get/{key}` endpoint'ine bir GET isteği yaparak veri alın.

## API Dökümantasyonu

### `/set` Endpoint'i

- **Method**: POST
- **Gövde**: JSON formatında `{ "key": "ANAHTAR", "value": "DEĞER" }` şeklinde veri gönderilir.
- **Cevap**: Eklenen veriyi JSON formatında döndürür.

Örnek istek:
POST http://localhost:8080/set

/get/{key} Endpoint'i
Method: GET
Parametreler:
key: Alınacak verinin anahtarı.
Cevap:
Belirtilen anahtara sahip veriyi JSON formatında döndürür.
Veri bulunamazsa HTTP status kodu 404 döndürür.
