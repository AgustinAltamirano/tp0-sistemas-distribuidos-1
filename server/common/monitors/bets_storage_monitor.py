import threading
from typing import Generator

from ..utils import Bet, store_bets, load_bets


class BetsStorageMonitor:
    def __init__(self):
        self._lock = threading.Lock()

    def store_bets(self, bets):
        with self._lock:
            store_bets(bets)

    def load_bets(self) -> Generator[Bet, None, None]:
        return load_bets()
