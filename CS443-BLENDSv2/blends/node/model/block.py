import json

from sqlalchemy import Column, ForeignKey, Integer, String
from sqlalchemy.orm import relationship

from . import Base


class Block(Base):
    __tablename__ = 'block'
    """ Header """
    hash = Column(String(255), primary_key=True)
    version = Column(String(16), nullable=False)
    """ Payload """
    parent = Column(String(255), ForeignKey('block.hash'))
    timestamp = Column(String(255), nullable=False)
    miner = Column(String(1024), nullable=False)
    difficulty = Column(Integer, nullable=False)
    nonce = Column(Integer, nullable=False)

    txs = relationship("Transaction", backref="block")

    @staticmethod
    def new_block(version: str, parent: str, timestamp: str, miner: str,
                  difficulty: int):
        """
        Generate empty block
        """
        return Block(version, parent, timestamp, miner, difficulty, 0)

    def __init__(self, version: str, parent: str, timestamp: str, miner: str,
                 difficulty: str, nonce: int):
        """
        Initialize empty block
        """
        self.version = version
        self.parent = parent
        self.timestamp = timestamp
        self.miner = miner
        self.difficulty = difficulty
        self.nonce = nonce

    def __repr__(self):
        return "<Block: {}>".format(self.hash)

    def get_header(self) -> dict:
        """
        Return Header Dictionary
        """
        return {"type": "block", "hash": self.hash, "version": self.version}

    def get_payload_dict(self) -> dict:
        """
        Return Payload Dict
        """
        payload = {
            "parent":
            self.parent,
            "timestamp":
            self.timestamp,
            "miner":
            self.miner,
            "difficulty":
            self.difficulty,
            "nonce":
            self.nonce,
            "transactions": [{
                "header": tx.get_header(),
                "payload": tx.get_payload_dict()
            } for tx in self.txs]
        }
        return payload

    def get_payload(self) -> str:
        """
        Return Payload String
        """
        payload = self.get_payload_dict()
        return json.dumps(payload, sort_keys=True)

    def set_nonce(self, nonce: int):
        """
        Step 1: change nonce
        """
        self.nonce = nonce

    def set_hash(self, hash):
        """
        Step 2: set hash
        """
        self.hash = hash
