import threading
import queue

import requests

from ... import PORT
from ..crypto import get_hash
from ..model import Block
from ..parser import Parser
from ..util import now


class Miner:
    def __init__(self,
                 block: Block,
                 mining: threading.Event,
                 main_queue: queue.Queue = None):
        self.parser = Parser()
        self.block = block
        self.mining = mining
        self.main_queue = main_queue

    def mine_block(self):
        idx = 0
        while True:
            idx += 1
            if idx % 100000 == 0 and not self.mining.is_set():
                return False
            """
            Step 1: Your mining code
            """
            diff = self.block.difficulty
            self.block.nonce = idx
            hash_ = get_hash(self.block.get_payload())
            self.block.set_hash(hash_)
            if idx % 1000000 == 0:
                self.block.timestamp = now()
                print(hash_)
            if(int(hash_, 16) < pow(2, 512) / pow(2, (20 + diff))):
                return True

    def mine(self):
        if self.mine_block():
            self.new_block_mint()

    def new_block_mint(self):
        if self.main_queue:
            self.main_queue.put(self.block)
