from ..interfaces.message import Message
from ..constants.message_codes import MessageCode
from ..interfaces.protocol import Protocol


class AskWinnersList(Message):
    def __init__(self, agency_id: int):
        self.agency_id = agency_id

    def get_message_code(self) -> MessageCode:
        return MessageCode.ASK_WINNERS_LIST

    def send_to_lottery_central(self, lottery_central):
        lottery_central.give_winners_list(self.agency_id)

    def send_to_agency(self, protocol: Protocol):
        raise NotImplementedError("AskWinnersList is not sent to agencies")
