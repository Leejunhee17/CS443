import threading

from .miner import Miner
from ..model import Block


class MinerThread(threading.Thread, Miner):
    def __init__(self, block: Block, mining: threading.Event, queue):
        threading.Thread.__init__(self)
        self.miner = Miner(block, mining, queue)

    def run(self):
        self.miner.mine()
