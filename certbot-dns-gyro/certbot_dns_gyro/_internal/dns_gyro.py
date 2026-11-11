"""DNS Authenticator for Gyro DNS (ACME key: may only write _acme-challenge TXT records)."""
import logging
import time
from typing import Any, Callable, Optional

import requests
from certbot import errors
from certbot.plugins import dns_common
from certbot.plugins.dns_common import CredentialsConfiguration

logger = logging.getLogger(__name__)

DEFAULT_API = "https://www.gyrohosting.com/api/v1/client/"
USER_AGENT = "certbot-dns-gyro/0.1.0"


class Authenticator(dns_common.DNSAuthenticator):
    """Gyro DNS Authenticator: fulfils dns-01 challenges with the Gyro DNS API."""

    description = "Obtain certificates using a DNS TXT record (if you are using Gyro DNS)."

    def __init__(self, *args: Any, **kwargs: Any) -> None:
        super().__init__(*args, **kwargs)
        self.credentials: Optional[CredentialsConfiguration] = None

    @classmethod
    def add_parser_arguments(cls, add: Callable[..., None], default_propagation_seconds: int = 10) -> None:
        super().add_parser_arguments(add, default_propagation_seconds)
        add("credentials", help="Gyro DNS credentials INI file (dns_gyro_token = <ACME key>).")

    def more_info(self) -> str:
        return "This plugin configures a DNS TXT record to respond to a dns-01 challenge using the Gyro DNS API."

    def _setup_credentials(self) -> None:
        self.credentials = self._configure_credentials(
            "credentials", "Gyro DNS credentials INI file",
            {"token": "ACME key created in the Gyro DNS service page (ACME / Let's Encrypt)."})

    def _client(self) -> "_GyroClient":
        assert self.credentials is not None
        return _GyroClient(self.credentials.conf("token"), self.credentials.conf("api") or DEFAULT_API)

    def _perform(self, domain: str, validation_name: str, validation: str) -> None:
        self._client().add(validation_name, validation)

    def _cleanup(self, domain: str, validation_name: str, validation: str) -> None:
        try:
            self._client().remove(validation_name, validation)
        except errors.PluginError as e:
            logger.warning("Could not remove the TXT record %s: %s", validation_name, e)


class _GyroClient:
    def __init__(self, token: str, api: str) -> None:
        self.api = api if api.endswith("/") else api + "/"
        self.s = requests.Session()
        self.s.headers.update({"Authorization": "Bearer " + token, "Accept": "application/json", "User-Agent": USER_AGENT})

    def _call(self, method: str, name: str, value: str) -> dict:
        url = self.api + "cloud/dns/acme-challenge"
        try:
            if method == "POST":
                r = self.s.post(url, json={"name": name, "value": value}, timeout=60)
            else:
                r = self.s.request(method, url, params={"name": name, "value": value}, timeout=60)
        except requests.RequestException as e:
            raise errors.PluginError("Gyro DNS API unreachable: %s" % e)
        try:
            j = r.json()
        except ValueError:
            j = {}
        if r.status_code >= 300:
            err = (j.get("error") or {}) if isinstance(j, dict) else {}
            raise errors.PluginError("Gyro DNS API error %s: %s (%s)" % (r.status_code, err.get("message", r.text[:200]), err.get("code", "")))
        return j.get("data") or {}

    def add(self, name: str, value: str) -> None:
        self._call("POST", name, value)
        # wait until both authoritative name servers publish the value
        deadline = time.time() + 180
        while time.time() < deadline:
            if self._call("GET", name, value).get("ready"):
                logger.info("Gyro DNS: %s is live on both name servers", name)
                return
            time.sleep(5)
        logger.warning("Gyro DNS: name servers have not confirmed %s yet; continuing", name)

    def remove(self, name: str, value: str) -> None:
        self._call("DELETE", name, value)
