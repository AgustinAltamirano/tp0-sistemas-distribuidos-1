from ..interfaces.message import Message
from ..constants.message_codes import MessageCode
from ..interfaces.protocol import Protocol


class FinalizeBets(Message):
    def __init__(self, agency_id: int):
        self.agency_id = agency_id

    def get_message_code(self) -> MessageCode:
        return MessageCode.FINALIZE_BETS

    def send_to_lottery_central(self, lottery_central):
        lottery_central.finalize_bets(self.agency_id)

    def send_to_agency(self, protocol: Protocol):
        raise NotImplementedError("FinalizeBets is not sent to agencies")
