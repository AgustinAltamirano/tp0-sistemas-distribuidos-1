from abc import ABC, abstractmethod
from ..constants.message_codes import MessageCode


class Message(ABC):
    @abstractmethod
    def get_message_code(self) -> MessageCode:
        pass

    @abstractmethod
    def send_to_lottery_central(self, lottery_central):
        pass

    @abstractmethod
    def send_to_agency(self, protocol):
        pass
