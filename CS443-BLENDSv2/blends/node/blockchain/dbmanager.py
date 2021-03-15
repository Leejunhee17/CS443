import collections
from typing import Dict, List, Optional

from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker

from ... import GENESIS_HASH, REWARD
from ..model import Block, Transaction


class DBManager:
    def __init__(self, connection_string: str):
        self.engine = create_engine(connection_string)
        self.sessionmaker = sessionmaker(bind=self.engine)
        self.session = None

    def get_session(self):
        if not self.session:
            self.session = self.sessionmaker()
        return self.session

    def search_block(self, block_hash: str) -> Optional[Block]:
        """
        Step 1
        """
        session = self.get_session()
        for block in session.query(Block).all():
            if block.get_header()["hash"] == block_hash:
                return block
        return None

    def get_height(self, block_hash: str) -> Optional[int]:
        """
        Step 2
        """
        block = self.search_block(block_hash)
        if(block == None):
            return None
        height = 0
        while(True):
            block = self.search_block(block.parent)
            if(block == None):
                break
            height += 1
        return height

    def get_current(self) -> Block:
        """
        Step 3
        """
        session = self.get_session()
        max = 0
        cur = session.query(Block).all()[0]
        for block in session.query(Block).all():
            if(self.get_height(block.get_header()["hash"]) >= max):
                cur = block
        return cur

    def get_longest(self) -> List[Block]:
        """
        Step 4
        """
        session = self.get_session()
        cur = self.get_current()
        longest = []
        while(True):
            longest.append(cur)
            cur = self.search_block(cur.parent)
            if(cur == None):
                break
        longest.reverse()
        return longest

    def search_transaction(self, tx_hash: str) -> Optional[Transaction]:
        """
        Step 5
        """
        session = self.get_session()
        for tx in session.query(Transaction).all():
            if tx.get_header()["hash"] == tx_hash:
                return tx
        return None

    def get_block_balance(self, block_hash: str) -> Dict[str, int]:
        """
        Step 6
        """
        block = self.search_block(block_hash)
        balance = dict()
        if(block == None):
            return balance
        while(True):
            payload = block.get_payload_dict()
            miner = payload["miner"]
            if miner not in balance:
                balance[miner] = REWARD
            else:
                balance[miner] += REWARD
            txs = payload["transactions"]
            for tx in txs:
                txpayload = tx["payload"]
                sender = txpayload["sender"]
                receiver = txpayload["receiver"]
                amount = txpayload["amount"]
                if sender not in balance:
                     balance[sender] = -amount
                else:
                    balance[sender] -= amount
                if receiver not in balance:
                    balance[receiver] = amount
                else:
                    balance[receiver] += amount
            block = self.search_block(block.parent)
            if(block == None):
                break
        return balance

    def insert_block(self, block: Block) -> bool:
        session = self.get_session()
        session.add(block)
        session.commit()

        return True
