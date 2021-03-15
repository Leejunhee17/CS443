import json
import threading
import queue

from .. import DIFFICULTY, BLENDS_VERSION
from .accountmanager import AccountManager
from .blockchain import Blockchain
from .miner import MinerThread
from .model import Block, Transaction
from .networkmanager import NetworkManager
from .parser import Parser
from .util import now


class Node:
    def __init__(self,
                 broadcast_url: str,
                 key_path: str,
                 db_connection_string: str,
                 difficulty: int = DIFFICULTY):

        self.main_queue = queue.Queue()
        self.network_manager = NetworkManager(broadcast_url, self.main_queue)
        self.account_manager = AccountManager(key_path)
        self.parser = Parser()
        self.blockchain = Blockchain(db_connection_string, difficulty)
        self.mining = None
        self.difficulty = difficulty

        self.height = 0
        self.block = None

    def new_requests_received(self, requests: list):
        # must check the request is block or transaction
        # request is list of requests (from start!)
        # Your code starts here

        # you should have to call self.install_new_block(), self.stop_miner(), self.start_miner(), see below
        print(requests)

    def new_block_mint(self, block: Block):
        self.stop_miner()
        data = {}
        # Your code starts here
        # you should have to call self.install_new_block(), self.stop_miner(), self.start_miner(), see below

        print(Block)

        self.network_manager.publish(data)

    def new_transaction_issued(self, tx: Transaction):
        data = {}
        print(tx)
        # Your code starts here
        # you should have to call self.install_new_block(), self.stop_miner(), self.start_miner(), see below

        self.network_manager.publish(data)

    def install_new_block(self):
        parent_block = self.blockchain.get_current()
        self.height = self.blockchain.dbmanager.get_height(
            parent_block.hash) + 1
        self.block = Block.new_block(BLENDS_VERSION, parent_block.hash, now(),
                                     self.account_manager.get_public_key(),
                                     self.difficulty)

    def stop_miner(self):
        self.mining.clear()
        self.mining = None
        self.miner = None

    def start_miner(self):
        mining = threading.Event()
        self.mining = mining
        self.mining.set()
        miner = MinerThread(self.block, mining, self.main_queue)
        miner.start()

    def run(self):
        self.install_new_block()
        self.start_miner()
        self.network_manager.start()

        while True:
            result = self.main_queue.get()
            if type(result) == Block:
                self.new_block_mint(result)
            elif type(result) == Transaction:
                self.new_transaction_issued(result)
            elif type(result) == list:
                self.new_requests_received(result)
