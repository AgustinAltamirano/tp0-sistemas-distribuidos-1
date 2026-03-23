from ..interfaces.message import Message
from ..constants.message_codes import MessageCode
from ..interfaces.protocol import Protocol
from ..constants.result_codes import ResultCode


class ConfirmBetBatch(Message):
    def __init__(
        self,
        agency_id: int,
        bet_amount: int,
        result_code: ResultCode,
    ):
        self.agency_id = agency_id
        self.bet_amount = bet_amount
        self.result_code = result_code

    def get_message_code(self) -> MessageCode:
        return MessageCode.CONFIRM_BET_BATCH

    def send_to_lottery_central(self, lottery_central):
        raise NotImplementedError("ConfirmBetBatch is not sent to lottery central")

    def send_to_agency(self, protocol: Protocol):
        protocol.send_confirm_bet_batch(
            self.agency_id, self.bet_amount, self.result_code
        )
