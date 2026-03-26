import socket
import logging
import threading

from .socket.safe_socket import SafeSocket
from .lottery_central import LotteryCentral
from .monitors.bets_storage_monitor import BetsStorageMonitor
from .monitors.lottery_event_monitor import LotteryEventMonitor


class Server:
    def __init__(self, port, listen_backlog, agencies_amount: int):
        # Initialize server socket
        self._server_safe_socket = SafeSocket()
        self._server_safe_socket.bind(("", port))
        self._server_safe_socket.listen(listen_backlog)
        self._server_safe_socket.settimeout(0.2)
        self._closed = False
        self._stop_requested = False
        self._bets_storage_monitor = BetsStorageMonitor()
        self._lottery_event_monitor = LotteryEventMonitor(agencies_amount)
        self._threads = []

    def __enter__(self):
        return self

    def __exit__(self, exc_type, exc, tb):
        self.close()
        return False

    def close(self):
        if self._closed:
            return
        self._closed = True
        self._lottery_event_monitor.abort()
        self._server_safe_socket.close()
        logging.info("action: close | result: success")

    def request_shutdown(self):
        self._stop_requested = True

    def run(self):
        """
        Server loop

        Server that accepts new connections and spawns a thread per client.
        After all clients finish, the server exits.
        """

        while not self._stop_requested:
            client_sock = self.__accept_new_connection()
            if client_sock is None:
                continue
            t = threading.Thread(
                target=self.__handle_client_connection,
                args=(client_sock,),
                daemon=True,
            )
            self._threads.append(t)
            t.start()

        for t in self._threads:
            t.join()

    def __handle_client_connection(self, client_socket):
        try:
            lottery_central = LotteryCentral(
                client_socket,
                self._bets_storage_monitor,
                self._lottery_event_monitor,
            )
            lottery_central.start()
        except Exception as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        finally:
            client_socket.close()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        logging.info("action: accept_connections | result: in_progress")
        try:
            client_safe_socket, addr = self._server_safe_socket.accept()
            logging.info(
                f"action: accept_connections | result: success | ip: {addr[0]}"
            )
            return client_safe_socket
        except socket.timeout:
            return None
        except OSError as e:
            if not self._stop_requested and not self._closed:
                logging.error(f"action: accept_connections | result: fail | error: {e}")
            return None
