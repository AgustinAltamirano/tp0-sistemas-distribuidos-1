from abc import ABC, abstractmethod

from ..constants.result_codes import ResultCode
from ..interfaces.message import Message
from ..bet_batch import BetBatch


class Protocol(ABC):
    @abstractmethod
    def receive_message(self) -> Message:
        pass

    @abstractmethod
    def receive_bet_batch(self) -> BetBatch:
        pass

    @abstractmethod
    def send_message(self, message: Message) -> None:
        pass

    @abstractmethod
    def send_confirm_bet_batch(
        self, agency_id: int, bet_amount: int, result_code: ResultCode
    ) -> None:
        pass
