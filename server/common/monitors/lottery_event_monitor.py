import logging
import threading


class LotteryEventMonitor:
    def __init__(self, n_agencies: int):
        self._pending = n_agencies
        self._lock = threading.Lock()
        self._event = threading.Event()

    def wait(self):
        with self._lock:
            self._pending -= 1
            if self._pending == 0:
                logging.info("action: sorteo | result: success")
                self._event.set()
        self._event.wait()

    def abort(self):
        self._event.set()
