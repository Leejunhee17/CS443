import json

from sqlalchemy import Column, ForeignKey, Integer, String

from . import Base


class Transaction(Base):
    __tablename__ = "tx"
    """ Header """
    hash = Column(String(255), primary_key=True)
    version = Column(String(16), nullable=False)
    sign = Column(String(16), nullable=False)
    """ Payload """
    sender = Column(String(1024), primary_key=False)
    receiver = Column(String(1024), primary_key=False)
    timestamp = Column(String(255), primary_key=False)
    amount = Column(Integer, nullable=False)
    block_s = Column("block", String(255), ForeignKey('block.hash'))

    @staticmethod
    def new_transaction(version: str, sender: str, receiver: str,
                        timestamp: str, amount: int):
        return Transaction(version, sender, receiver, timestamp, amount)

    def __init__(self, version: str, sender: str, receiver: str,
                 timestamp: str, amount: int):
        """
        Initialize transaction
        """
        self.version = version
        self.sender = sender
        self.receiver = receiver
        self.timestamp = timestamp
        self.amount = amount

    def __repr__(self):
        return "<Transaction: {}>".format(self.hash)

    def get_header(self) -> dict:
        """
        Return Header Dictionary
        """
        return {
            "type": "transaction",
            "hash": self.hash,
            "version": self.version,
            "sign": self.sign
        }

    def get_payload_dict(self) -> dict:
        """
        Return Payload Dict
        """
        payload = {
            "sender": self.sender,
            "receiver": self.receiver,
            "timestamp": self.timestamp,
            "amount": self.amount
        }

        return payload

    def get_payload(self) -> str:
        """
        Return Payload String
        """
        payload = self.get_payload_dict()

        return json.dumps(payload, sort_keys=True)

    def set_hash(self, hash: str):
        """
        Step 1: Set siginiture of a transaction
        """
        self.hash = hash

    def set_sign(self, sign: str):
        """
        Step 2: Set siginiture of a transaction
        """
        self.sign = sign
