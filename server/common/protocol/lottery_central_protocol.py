from ..constants.message_codes import MessageCode
from ..interfaces.message import Message
from ..interfaces.protocol import Protocol
from ..constants.result_codes import ResultCode
from ..messages.message_factory import MessageFactory
from ..socket.safe_socket import SafeSocket
from ..utils import Bet
from ..bet_batch import BetBatch


# messageCode(2 bytes) + agencyId(2 bytes) + betAmount(4 bytes) + betBatchSize(4 bytes) = 12 bytes
REGISTER_BET_BATCH_BASE_SIZE = 12

# document(4 bytes) + birthdate(10 bytes) + number(4 bytes)
# + firstNameLength(2 bytes) + lastNameLength(2 bytes) = 22 bytes
BET_BASE_SIZE = 22

# messageCode(2 bytes) + agencyId(2 bytes) + winnersAmount(4 bytes) = 8 bytes
WINNERS_LIST_BASE_SIZE = 8

NETWORK_ENDIANNESS = "big"


class LotteryCentralProtocol(Protocol):
    def __init__(self, safe_socket: SafeSocket):
        self._socket = safe_socket
        self._message_factory = MessageFactory()

    def receive_message(self) -> Message:
        message_code = self._receive_message_code()
        return self._message_factory.get_message(message_code, self)

    def receive_bet_batch(self) -> BetBatch:
        header = self._socket.read(REGISTER_BET_BATCH_BASE_SIZE - 2)
        agency_id, bet_amount, bet_batch_size = self._deserialize_bet_batch_header(
            header
        )
        bets_buffer = self._socket.read(bet_batch_size)
        return self._deserialize_bet_batch_from_buffer(
            bets_buffer, bet_amount, agency_id
        )

    def send_message(self, message: Message):
        message.send_to_agency(self)

    def send_confirm_bet_batch(
        self, agency_id: int, bet_amount: int, result_code: ResultCode
    ):
        buf = bytearray(9)
        buf[0:2] = MessageCode.CONFIRM_BET_BATCH.value.to_bytes(
            2, byteorder=NETWORK_ENDIANNESS
        )
        buf[2:4] = agency_id.to_bytes(2, byteorder=NETWORK_ENDIANNESS)
        buf[4:8] = bet_amount.to_bytes(4, byteorder=NETWORK_ENDIANNESS)
        buf[8] = result_code.value
        self._socket.write(bytes(buf))

    def receive_finalize_bets(self) -> int:
        return self._receive_int(2)

    def receive_ask_winners_list(self) -> int:
        return self._receive_int(2)

    def send_give_winners_list(self, agency_id: int, winners: list) -> None:
        winners_amount = len(winners)
        buffer = bytearray(WINNERS_LIST_BASE_SIZE + winners_amount * 4)
        buffer[0:2] = MessageCode.GIVE_WINNERS_LIST.value.to_bytes(
            2, byteorder=NETWORK_ENDIANNESS
        )
        buffer[2:4] = agency_id.to_bytes(2, byteorder=NETWORK_ENDIANNESS)
        buffer[4:8] = winners_amount.to_bytes(4, byteorder=NETWORK_ENDIANNESS)
        for i, doc in enumerate(winners):
            offset = WINNERS_LIST_BASE_SIZE + i * 4
            buffer[offset : offset + 4] = doc.to_bytes(4, byteorder=NETWORK_ENDIANNESS)
        self._socket.write(bytes(buffer))

    def _receive_message_code(self) -> MessageCode:
        message_code_value = self._receive_int(2)
        return MessageCode(message_code_value)

    def _deserialize_bet_batch_header(self, buf: bytes) -> tuple[int, int, int]:
        agency_id = int.from_bytes(buf[0:2], byteorder=NETWORK_ENDIANNESS)
        bet_amount = int.from_bytes(buf[2:6], byteorder=NETWORK_ENDIANNESS)
        bet_batch_size = int.from_bytes(buf[6:10], byteorder=NETWORK_ENDIANNESS)
        return agency_id, bet_amount, bet_batch_size

    def _deserialize_bet_batch_from_buffer(
        self, buffer: bytes, bet_amount: int, agency_id: int
    ) -> BetBatch:
        bets = []
        offset = 0
        for _ in range(bet_amount):
            bet, offset = self._deserialize_bet_from_buffer(buffer, offset, agency_id)
            bets.append(bet)
        return BetBatch(agency_id, bets)

    def _deserialize_bet_from_buffer(
        self, buffer: bytes, offset: int, agency_id: int
    ) -> tuple[Bet, int]:
        document = int.from_bytes(
            buffer[offset : offset + 4], byteorder=NETWORK_ENDIANNESS
        )
        offset += 4
        birthdate = buffer[offset : offset + 10].decode("utf-8")
        offset += 10
        number = int.from_bytes(
            buffer[offset : offset + 4], byteorder=NETWORK_ENDIANNESS
        )
        offset += 4
        first_name_size = int.from_bytes(
            buffer[offset : offset + 2], byteorder=NETWORK_ENDIANNESS
        )
        offset += 2
        last_name_size = int.from_bytes(
            buffer[offset : offset + 2], byteorder=NETWORK_ENDIANNESS
        )
        offset += 2
        first_name = buffer[offset : offset + first_name_size].decode("utf-8")
        offset += first_name_size
        last_name = buffer[offset : offset + last_name_size].decode("utf-8")
        offset += last_name_size
        bet = Bet(agency_id, first_name, last_name, document, birthdate, number)
        return bet, offset

    def _receive_int(self, n_bytes: int) -> int:
        if n_bytes <= 0:
            raise ValueError("n_bytes must be positive")
        network_bytes = self._socket.read(n_bytes)
        return int.from_bytes(network_bytes, byteorder=NETWORK_ENDIANNESS, signed=False)
