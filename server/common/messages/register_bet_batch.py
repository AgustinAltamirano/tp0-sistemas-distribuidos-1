from ..interfaces.message import Message
from ..constants.message_codes import MessageCode
from ..interfaces.protocol import Protocol


class RegisterBetBatch(Message):
    def __init__(self, bet_batch):
        self.bet_batch = bet_batch

    def get_message_code(self) -> MessageCode:
        return MessageCode.REGISTER_BET_BATCH

    def send_to_lottery_central(self, lottery_central):
        lottery_central.register_bet_batch(self.bet_batch)

    def send_to_agency(self, protocol: Protocol):
        raise NotImplementedError("RegisterBetBatch is not sent to agencies")
