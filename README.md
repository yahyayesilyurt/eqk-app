# Earthquake Monitoring App

Bu uygulama, belirli bir büyüklüğün üzerinde olan depremleri görselleştirmek için geliştirilmiş bir web uygulamasıdır. Backend kısmı Golang, frontend kısmı React ile geliştirilmiştir. Veritabanı olarak Mongodb kullanılmaktadır. Dockerize edilmiştir, docker compose ile çalıştırılabilir.

## İçindekiler

- [Earthquake Monitoring App](#earthquake-monitoring-app)
  - [İçindekiler](#i̇çindekiler)
  - [Nasıl Çalışır](#nasıl-çalışır)
  - [Kurulum](#kurulum)
  - [Kullanım](#kullanım)
  - [Dikkat Edilmesi Gerekenler](#dikkat-edilmesi-gerekenler)
  - [Geliştiriciden Mesaj](#geliştiriciden-mesaj)

## Nasıl Çalışır

Backend, MongoDB'ye bağlanır, USGS'den deprem verilerini alır ve belirli bir büyüklüğün üzerinde olan depremleri filtreler, yani büyüklüğü 4 ve üzeri olan depremleri anormal deprem sayar ve veritabanına kaydeder. Frontend, bu filtrelenmiş verileri alır ve harita üzerinde pin'lerle görselleştirir. 28 saniye boyunca haritada pinli bir şekilde depremler gösterilir. Pinlerin üstüne tıklanınca depremler hakkında büyüklük, enlem ve boylam bilgileri görünür. Her 30 saniyede bir kez veriler yeniden çekilir ve güncel verilerle sayfa yeniden renderlanır. Backendde iki adet script bulunur. Bunlardan biri dakika başı bir kez random deprem verisi üretir. Eğer üretilen veri 4 ve üzeri büyüklükteyse veritabanına kaydedilir ve haritada gösterilir. Diğer script ise kullanıcıdan deprem verisi ister. Kullanıcının girdiği veriler beklenen deprem özelliklerine uyuyorsa veritabanına kaydedilir ve haritada pinlenir.

## Kurulum

1. Projenin klonunu alın:

```
git clone <proje_url>
```

2. Docker ve Docker Compose yüklü olmalıdır.

## Kullanım

1. Terminal veya komut istemcisini açın.
2. Projenin kök dizinine gidin.
3. Aşağıdaki komutu çalıştırarak Docker Compose ile uygulamayı başlatın:

```
docker-compose up
```

4. Uygulama başladıktan sonra, tarayıcınızda `http://localhost:3000` adresine giderek uygulamayı görebilirsiniz.

## Dikkat Edilmesi Gerekenler

- Projeyi başlatmadan önce internet bağlantınızın olduğundan emin olun, çünkü backend USGS'den canlı deprem verilerini alacaktır.
- Docker Compose kullanarak projeyi başlattığınızda, MongoDB de dahil olmak üzere tüm bileşenler birlikte çalışacaktır.
- Uygulamayı durdurmak için terminalde `Ctrl + C` kombinasyonunu kullanabilirsiniz.
- Docker Compose kullanarak projeyi başlattığınızda, MongoDB verileri otomatik olarak saklayacak ve veriler korunacaktır.

## Geliştiriciden Mesaj

Uygulamada beklenen özelliklerden biri olan Apache Flink kullanamamamın sebebi projenin sınav dönemime denk gelmesiydi. Daha önce Apache Flink kullanmamıştım ve araştırmama rağmen vaktim çok kısıtlı olduğu için projede uygulayacak kadar Apache Flink'e hakim olamadım. Buna alternatif olarak deprem verilerinin hepsini kaydetmek yerine büyüklüğü 4 ve üzeri olan depremleri anormal deprem sayarak bu depremleri veritabanına kaydettirdim. Daha sonra da bu depremleri arayüzde harita üzerinde pinledim.
