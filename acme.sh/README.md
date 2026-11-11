# acme.sh: dns_gyro

```sh
curl -fsSL -o ~/.acme.sh/dnsapi/dns_gyro.sh https://raw.githubusercontent.com/Gyro-Hosting/gyro-acme/main/acme.sh/dns_gyro.sh
export GYRO_Token="wck_..."          # Gyro DNS ACME key
acme.sh --issue --dns dns_gyro -d example.com -d '*.example.com'
```
`GYRO_Token` (and `GYRO_Api` when set) are saved to `account.conf` for renewals. After adding a record the hook waits
until both Gyro DNS name servers publish it (normally under a minute), so `--dnssleep` is not needed.

Türkçe: Gyro DNS hizmet sayfasından ACME anahtarı oluşturun, `GYRO_Token` olarak verin ve `--dns dns_gyro` ile kullanın.
Yenilemeler otomatik çalışır.

Upstream: the same file is prepared for an acme.sh pull request (see `../upstream/acme.sh/`).
