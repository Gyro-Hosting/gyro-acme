# Upstream contributions (prepared, not submitted)

## acme.sh
- File: `acme.sh/dns_gyro.sh` (follows the dnsapi conventions: `dns_gyro_info`, `_readaccountconf_mutable`, `_post`).
- PR target: `acmesh-official/acme.sh` branch `dev`; wiki entry for `dns_gyro` (Options: `GYRO_Token`).
- Checks to run before submitting: `shellcheck`, `shfmt -d`, the acme.sh DNS test workflow (needs a test domain + token secret).

## lego (Traefik's ACME library)
- Files: `lego/providers/dns/gyro/{gyro.go,gyro.toml,internal/client.go}` (env `GYRO_API_TOKEN`, `GYRO_BASE_URL`).
- Built inside lego v4.35.2 with `go generate ./...`; `lego --dns gyro` issued a Let's Encrypt staging wildcard certificate.
- Before a PR: unit tests (`gyro_test.go`, `internal/client_test.go` with fixtures) as lego requires, `make checks`.
- Once merged, Traefik users can use `provider: gyro` directly; until then `provider: exec` + `gyro-acme-exec`.
