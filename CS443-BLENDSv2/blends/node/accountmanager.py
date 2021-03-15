from .crypto import create_secret_key, get_pk, load_secret_key, sign
import os


class AccountManager:
    def __init__(self, key_path: str):
        self.key_path = key_path
        if not os.path.exists(key_path):
            create_secret_key(key_path)
        self._secret_key = load_secret_key(self.key_path)
        self.public_key = get_pk(self._secret_key)

    def sign(self, msg: str) -> str:
        return sign(self._secret_key, msg)

    def get_public_key(self) -> str:
        return self.public_key
