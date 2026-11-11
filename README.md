# Gyro DNS ACME plugins

Get **Let's Encrypt** (or any ACME CA) certificates — **wildcards included** — automatically with the **DNS-01**
challenge on **Gyro DNS** (ns1/ns2.gyrodns.net).

*Türkçe açıklama aşağıda.*

| Tool | Plugin | Where |
|---|---|---|
| acme.sh | `dns_gyro` | [`acme.sh/dns_gyro.sh`](acme.sh/) |
| Certbot | `dns-gyro` | [`certbot-dns-gyro/`](certbot-dns-gyro/) |
| Caddy | `dns.providers.gyro` | [Gyro-Hosting/caddy-dns-gyro](https://github.com/Gyro-Hosting/caddy-dns-gyro) |
| Traefik, lego CLI | `exec` + `gyro-acme-exec` | [Gyro-Hosting/libdns-gyro](https://github.com/Gyro-Hosting/libdns-gyro) |
| Go (libdns) | `github.com/Gyro-Hosting/libdns-gyro` | [Gyro-Hosting/libdns-gyro](https://github.com/Gyro-Hosting/libdns-gyro) |
| Kubernetes | cert-manager webhook | [Gyro-Hosting/cert-manager-webhook-gyro](https://github.com/Gyro-Hosting/cert-manager-webhook-gyro) |
| Anything else | REST API | below |

## 1. Create an ACME key
Gyro DNS service page → **ACME / Let's Encrypt** → *Create key* — or with the API / CLI:
```sh
gyro dns acme-key create <service-id> --name web-01
```
An ACME key:
- can **only** add and remove `_acme-challenge` TXT records (TTL 60) — it cannot read your zones, change any other
  record, or reach invoices, services or account data;
- can be limited to selected zones and to IP addresses;
- is shown once; delete it in the panel if it leaks (one key per server/cluster is a good habit).

Your domain must be a zone in your Gyro DNS service (name servers ns1/ns2.gyrodns.net). Hosting (cPanel) zones are
managed from cPanel and are not supported.

## 2. Use it
**acme.sh**
```sh
curl -fsSL -o ~/.acme.sh/dnsapi/dns_gyro.sh https://raw.githubusercontent.com/Gyro-Hosting/gyro-acme/main/acme.sh/dns_gyro.sh
export GYRO_Token="wck_..."
acme.sh --issue --dns dns_gyro -d example.com -d '*.example.com'
```
**Certbot**
```sh
pip install certbot-dns-gyro   # PyPI; Snap Certbot cannot load it, use pip/venv Certbot
printf 'dns_gyro_token = wck_...\n' > /etc/letsencrypt/gyro.ini && chmod 600 /etc/letsencrypt/gyro.ini
certbot certonly --authenticator dns-gyro --dns-gyro-credentials /etc/letsencrypt/gyro.ini -d example.com -d '*.example.com'
```
**Caddy**, **Traefik**, **Kubernetes**: see the linked repositories.

## REST API (any client)
Base `https://www.gyrohosting.com/api/v1/client/`, header `Authorization: Bearer <ACME key>`.
```
POST   cloud/dns/acme-challenge   {"name":"_acme-challenge.example.com","value":"<digest>"}   add a value (others stay)
GET    cloud/dns/acme-challenge?name=...&value=...                                          {"ready": true} when ns1 and ns2 publish it
DELETE cloud/dns/acme-challenge?name=...&value=...                                          remove the value
```
Full reference: https://www.gyrohosting.com/developers/

## Türkçe
Gyro DNS ile **Let's Encrypt joker sertifikayı** otomatik alın. Gyro DNS hizmet sayfasındaki **ACME / Let's Encrypt**
bölümünden bir ACME anahtarı oluşturun; bu anahtar yalnız `_acme-challenge` TXT kaydı ekleyip silebilir (zone okuyamaz,
başka kayda, hizmete, faturaya erişemez; isterseniz zone ve IP ile sınırlanır). Sonra acme.sh, Certbot, Caddy, Traefik ya
da Kubernetes (cert-manager) eklentisini kullanın. Adım adım rehber: Gyro Hosting Bilgi Merkezi → Gyro DNS →
"Gyro DNS ile Let's Encrypt joker sertifika".

## License
MIT
