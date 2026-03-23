import socket


class SafeSocket:
    def __init__(self):
        self._socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    @classmethod
    def _from_socket(cls, sock: socket.socket):
        safe_socket = cls.__new__(cls)
        safe_socket._socket = sock
        return safe_socket

    def bind(self, address):
        self._socket.bind(address)

    def listen(self, backlog: int):
        self._socket.listen(backlog)

    def settimeout(self, timeout: float):
        self._socket.settimeout(timeout)

    def accept(self):
        client_sock, addr = self._socket.accept()
        return SafeSocket._from_socket(client_sock), addr

    def read(self, n: int) -> bytes:
        if n < 0:
            raise ValueError("n must be non-negative")
        if n == 0:
            return b""

        chunks = []
        bytes_read = 0

        while bytes_read < n:
            bytes_left = n - bytes_read
            chunk = self._socket.recv(bytes_left)
            if chunk == b"":
                raise ConnectionError(
                    "socket connection closed before reading expected bytes"
                )
            chunks.append(chunk)
            bytes_read += len(chunk)

        return b"".join(chunks)

    def write(self, data: bytes) -> int:
        total_sent = 0

        while total_sent < len(data):
            sent = self._socket.send(data[total_sent:])
            if sent == 0:
                raise ConnectionError("socket connection broken during write")
            total_sent += sent

        return total_sent

    def close(self):
        self._socket.close()
