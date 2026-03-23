from abc import ABC, abstractmethod

from ..constants.result_codes import ResultCode
from ..interfaces.message import Message
from ..bet_batch import BetBatch


class Protocol(ABC):
    """Abstract base class for protocol implementations.

    Defines the contract for communication protocols between server and clients.
    """

    @abstractmethod
    def receive_message(self) -> Message:
        """Receive and parse a message from the client."""
        pass

    @abstractmethod
    def receive_bet_batch(self) -> BetBatch:
        """Receive a complete bet batch from the client."""
        pass

    @abstractmethod
    def send_message(self, message: Message) -> None:
        """Send a message to the client."""
        pass

    @abstractmethod
    def send_confirm_bet_batch(self, agency_id: int, bet_amount: int, result_code: ResultCode) -> None:
        """Send a confirm bet batch message to the client."""
        pass
