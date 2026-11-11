# certbot-dns-gyro

Gyro DNS authenticator plugin for [Certbot](https://certbot.eff.org): DNS-01 challenges, wildcard certificates.

## Install
```sh
pip install "certbot-dns-gyro @ git+https://github.com/Gyro-Hosting/gyro-acme@v0.1.0#subdirectory=certbot-dns-gyro"
```
(Install into the same Python environment as Certbot. Snap Certbot users: use `pip` Certbot or acme.sh.)

## Credentials
```ini
# /etc/letsencrypt/gyro.ini  (chmod 600)
dns_gyro_token = wck_...
# dns_gyro_api = https://www.gyrohosting.com/api/v1/client/   (optional)
```
## Use
```sh
certbot certonly --authenticator dns-gyro --dns-gyro-credentials /etc/letsencrypt/gyro.ini \
  -d example.com -d '*.example.com'
```
The plugin waits until both Gyro DNS name servers publish the record; the default
`--dns-gyro-propagation-seconds` is 10.

Türkçe: Gyro DNS ACME anahtarını `dns_gyro_token` olarak bir ini dosyasına yazın (izin 600) ve
`--authenticator dns-gyro` ile kullanın. Joker sertifika için hem `ornek.com` hem `*.ornek.com` isteyin.
