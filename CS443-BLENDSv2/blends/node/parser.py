import json
from typing import Optional

from .model import Block, Transaction


class Parser:
    def parse_transaction(self, tx_json: str) -> Optional[Transaction]:
        """
        Step 1
        """
        raise NotImplementedError

        return None

    def parse_block(self, block_json: str) -> Optional[Block]:
        """
        Step 2
        """
        raise NotImplementedError

        return None

    def dump_transaction(self, transaction: Transaction) -> Optional[str]:
        """
        Step 3
        """
        raise NotImplementedError
        return None

    def dump_block(self, block: Block) -> Optional[str]:
        """
        Step 4
        """
        raise NotImplementedError
        return None
