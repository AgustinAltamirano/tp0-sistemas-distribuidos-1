from ..interfaces.message import Message
from ..constants.message_codes import MessageCode
from ..interfaces.protocol import Protocol
from .register_bet_batch import RegisterBetBatch
from .finalize_bets import FinalizeBets
from .ask_winners_list import AskWinnersList


class MessageFactory:
    def __init__(self):
        pass

    def get_message(self, message_code: MessageCode, protocol: Protocol) -> Message:
        match message_code:
            case MessageCode.REGISTER_BET_BATCH:
                return self._get_register_bet_batch_message(protocol)
            case MessageCode.FINALIZE_BETS:
                return self._get_finalize_bets_message(protocol)
            case MessageCode.ASK_WINNERS_LIST:
                return self._get_ask_winners_list_message(protocol)
            case _:
                raise ValueError(f"Unknown message code: {message_code}")

    def _get_register_bet_batch_message(self, protocol: Protocol) -> RegisterBetBatch:
        bet_batch = protocol.receive_bet_batch()
        return RegisterBetBatch(bet_batch)

    def _get_finalize_bets_message(self, protocol: Protocol) -> FinalizeBets:
        agency_id = protocol.receive_finalize_bets()
        return FinalizeBets(agency_id)

    def _get_ask_winners_list_message(self, protocol: Protocol) -> AskWinnersList:
        agency_id = protocol.receive_ask_winners_list()
        return AskWinnersList(agency_id)
