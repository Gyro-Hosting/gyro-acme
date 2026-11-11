# certbot-dns-gyro

Gyro DNS authenticator plugin for [Certbot](https://certbot.eff.org): DNS-01 challenges, wildcard certificates.

## Install
```sh
pip install certbot-dns-gyro
```
The package is on PyPI: https://pypi.org/project/certbot-dns-gyro/. Install it into the same Python environment as Certbot.

**Snap Certbot cannot load this plugin.** Install Certbot with pip in a virtual environment instead (or use acme.sh):
```sh
python3 -m venv /opt/certbot
/opt/certbot/bin/pip install certbot certbot-dns-gyro
ln -s /opt/certbot/bin/certbot /usr/local/bin/certbot
```

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
