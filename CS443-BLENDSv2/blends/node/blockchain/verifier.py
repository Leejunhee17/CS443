from ..crypto import get_hash, verify, RSA
from ..model import Block, Transaction
import json


class Verifier:
    def verify_transaction(self, transaction: Transaction) -> bool:
        """
        Step 1 : verify transaction
        """
        header = transaction.get_header()
        payloadDict = transaction.get_payload_dict()
        payload = transaction.get_payload()
        hash_ = header["hash"]
        ver = header["version"]
        sig = header["sign"]
        sender = payloadDict["sender"]
        receiver = payloadDict["receiver"]
        timestamp = payloadDict["timestamp"]
        amount = payloadDict["amount"]
        
        if(type(hash_) != type("")):
            print("Type Hash_ Error")
            print("Hash_:", type(hash_))
            return False
        if(type(ver) != type("")):
            print("Type Version Error")
            print("Version:", type(ver))
            return False
        if(type(sig) != type("")):
            print("Type Sign Error")
            print("Sign:", type(sig))
            return False
        if(type(sender) != type("")):
            print("Type Sender Error")
            print("Sender:", type(sender))
            return False
        if(type(receiver) != type("")):
            print("Type Receiver Error")
            print("Receiver:", type(receiver))
            return False
        if(type(timestamp) != type("")):
            print("Type Timestamp Error")
            print("Timestamp:", type(timestamp))
            return False
        if(type(amount) != type(1)):
            print("Type amount Error")
            print("Amount:", type(amount))
            return False
        
        if(int(get_hash(payload), 16) != int(hash_, 16)):
            print("Hash Error")
            print("get_hash(payload):", get_hash(payload))
            print("hash_:", hash_)
            return False
        
        sender_n = int(sender, 16)
        publicExponent = 0x10001
        rsaComponent = (sender_n, publicExponent)
        pk = RSA.construct(rsaComponent, True)
        if(not verify(pk, hash_, sig)):
            print("Sign Error")
            return False
        return True

    def verify_block(self, block: Block) -> bool:
        """
        Step 2 : verify block
        """
        header = block.get_header()
        payloadDict = block.get_payload_dict()
        payload = block.get_payload()
        hash_ = header["hash"]
        ver = header["version"]
        parent = payloadDict["parent"]
        timestamp = payloadDict["timestamp"]
        miner = payloadDict["miner"]
        diff = payloadDict["difficulty"]
        nonce = payloadDict["nonce"]
        txs = payloadDict["transactions"]
        
        if(type(hash_) != type("")):
            print("Type Hash_ Error")
            print("Hash_:", type(hash_))
            return False
        if(type(ver) != type("")):
            print("Type Version Error")
            print("Version:", type(ver))
            return False
        if(type(parent) != type("")):
            print("Type Parent Error")
            print("Parent:", type(parent))
            return False
        if(type(timestamp) != type("")):
            print("Type Timestamp Error")
            print("Timestamp:", type(timestamp))
            return False
        if(type(miner) != type("")):
            print("Type Miner Error")
            print("Miner:", type(miner))
            return False
        if(type(diff) != type(1)):
            print("Type Difficulty Error")
            print("Difficulty:", type(diff))
            return False
        if(type(nonce) != type(1)):
            print("Type Nonce Error")    
            print("Nonce:", type(nonce))
            return False
        if(type(txs) != type([])):
            print("Type List Error")
            print("Transactions:", type(txs))
            return False
        
        if(int(get_hash(payload), 16) != int(hash_, 16)):
            print("Hash Error")
            print("get_hash(payload):", get_hash(payload))
            print("hash_:", hash_)
            return False
        
        if(int(hash_, 16) >= pow(2, 512) / pow(2, (20 + diff))):
            print("Difficulty Error")
            return False
        
        for tx in txs:
            txHeader = tx["header"]
            txPayloadDict = tx["payload"]
            txPayload = json.dumps(txPayloadDict, sort_keys=True)
            hash_ = txHeader["hash"]
            ver = txHeader["version"]
            sig = txHeader["sign"]
            sender = txPayloadDict["sender"]
            receiver = txPayloadDict["receiver"]
            timestamp = txPayloadDict["timestamp"]
            amount = txPayloadDict["amount"]

            if(type(hash_) != type(ver) != type(sig) != type(sender) != type(receiver) != type(timestamp) != type("")):
                print("Type Str Error")
                print("hash_:", type(hash_))
                print("Version:", type(ver))
                print("Sign:", type(sig))
                print("Sender:", type(sender))
                print("receiver:", type(receiver))
                print("Timestamp:", type(timestamp))
                return False
            if(type(amount) != type(1)):
                print("Type Int Error")
                print("Amount:", type(amount))
                return False

            if(int(get_hash(txPayload), 16) != int(hash_, 16)):
                print("Hash Error")
                print("get_hash(payload):", get_hash(payload))
                print("hash_:", hash_)
                return False
            

        return True
