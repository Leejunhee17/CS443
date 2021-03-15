from typing import Dict

from ... import DIFFICULTY, GENESIS_HASH, REWARD
from ..model import Block, Transaction
from .dbmanager import DBManager
from .verifier import Verifier


class Blockchain:
    def __init__(self, db_connection_string: str,
                 difficulty: int = DIFFICULTY):

        self.verifier = Verifier()
        self.dbmanager = DBManager(db_connection_string)

    def append(self, block: Block) -> bool:
        if self.verifier.verify_block(block) and self.validate_block(block):
            self.dbmanager.insert_block(block)
            return True
        else:
            return False

    def validate_transaction(self, tx: Transaction, block: Block) -> bool:
        header = block.get_header()
        payloadDict = block.get_payload_dict()
        txHeader = tx.get_header()
        txPayloadDict = tx.get_payload_dict()
        hash_ = header["hash"]
        txhash = txHeader["hash"]
        sender = txPayloadDict["sender"]
        amount = txPayloadDict["amount"]
        b = block
        while(True):
            payloadDict = b.get_payload_dict()
            txs = payloadDict["transactions"]
            for itx in txs:
                if(int(txhash, 16) == int(itx["header"]["hash"], 16)):
                    print("This transaction is already in block chain")
                    return False
            b = self.dbmanager.search_block(b.parent)
            if(b == None):
                break
        balance = self.get_balance()
        if sender not in balance:
            print("Sender is not in balance")
            return False
        if balance[sender] < amount:
            print("Sender's credit is less than amount")
            return False
        return True

    def validate_block(self, block: Block) -> bool:
        """
        check difficulty
        """
        header = block.get_header()
        payloadDict = block.get_payload_dict()
        txs = payloadDict["transactions"]
        if(payloadDict["parent"] != GENESIS_HASH):
            if(self.dbmanager.search_block(payloadDict["parent"]) == None):
                print("No Parent Block")
                return False
        pbalance = self.dbmanager.get_block_balance(payloadDict["parent"])
        miner = payloadDict["miner"]
        for tx in txs:
            txPayload = tx["payload"]
            sender = txPayload["sender"]
            receiver = txPayload["receiver"]
            amount = txPayload["amount"]
            if(sender not in pbalance):
                print("No Sender")
                return False
            if(pbalance[sender] < amount or pbalance[sender] < 0):
                print("Sender No Money")
                return False
            pbalance[sender] -= amount
            if(receiver not in pbalance):
                pbalance[receiver] = amount
            else:
                pbalance[receiver] += amount
        return True

    def get_balance(self) -> Dict[str, int]:
        block = self.dbmanager.get_current()
        balance = self.dbmanager.get_block_balance(block.hash)
        return balance

    def get_current(self):
        return self.dbmanager.get_current()
