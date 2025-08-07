# wges: WireGuard Easy Setup

Bu proje, WireGuard VPN sunucusu ve istemcilerini kolayca oluşturmak için terminal tabanlı bir Go uygulamasıdır.

## Özellikler
- `wges create server wg0`: /etc/wireguard dizininde wg0.conf dosyası oluşturur, sunucu için private/public key üretir, rastgele bir private IP bloğundan IP seçer.
- `wges create client clientname`: /etc/wireguard/clients dizininde client için yapılandırma ve anahtarlar oluşturur.

## Kurulum
1. Go kurulu olmalı (`go version` ile kontrol edebilirsiniz).
2. Proje dizininde terminal açın.
3. Gerekli izinler için root olarak çalıştırın veya sudo kullanın.

## Kullanım
```sh
sudo ./wges create server wg0
sudo ./wges create client clientname
```

## Adımlar
1. Go ile terminal tabanlı uygulama iskeleti oluşturulacak.
2. WireGuard anahtarları Go ile üretilecek.
3. Sunucu ve istemci yapılandırma dosyaları oluşturulacak.
4. Komut satırı argümanları ile yönetim sağlanacak.
5. /etc/wireguard ve /etc/wireguard/clients dizinleri kontrol edilecek/oluşturulacak.
6. Private IP bloğundan rastgele IP seçilecek.
7. README ve örnek kullanım eklenecek.

## Gereksinimler
- Go
- WireGuard kurulu olmalı (`wg`, `wg-quick` komutları erişilebilir olmalı)

## Notlar
- Anahtar üretimi için Go'nun crypto kütüphanesi kullanılacak.
- Dosya oluşturma ve yazma işlemleri için root yetkisi gerekebilir.

---

Devamında: Ana uygulama iskeleti ve komut işleyici eklenecek.
