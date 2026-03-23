from ..interfaces.message import Message
from ..constants.message_codes import MessageCode
from ..interfaces.protocol import Protocol
from .register_bet_batch import RegisterBetBatch


class MessageFactory:
    def __init__(self):
        pass

    def get_message(self, message_code: MessageCode, protocol: Protocol) -> Message:
        if message_code == MessageCode.REGISTER_BET_BATCH:
            return self._get_register_bet_batch_message(protocol)

        raise ValueError(f"Unknown message code: {message_code}")

    def _get_register_bet_batch_message(self, protocol: Protocol) -> RegisterBetBatch:
        bet_batch = protocol.receive_bet_batch()
        return RegisterBetBatch(bet_batch)
