from .bet_batch import BetBatch
from .protocol.lottery_central_protocol import LotteryCentralProtocol
from .socket.safe_socket import SafeSocket
from . import utils
from .messages.confirm_bet_batch import ConfirmBetBatch
from .constants.result_codes import ResultCode
import logging


class LotteryCentral:
    def __init__(self, safe_socket: SafeSocket):
        self._protocol = LotteryCentralProtocol(safe_socket)

    def start(self):
        while True:
            try:
                message = self._protocol.receive_message()
                message.send_to_lottery_central(self)
            except ConnectionError:
                break

    def register_bet_batch(self, bet_batch: BetBatch):
        result_code = ResultCode.SUCCESS
        try:
            utils.store_bets(bet_batch.bets)
            logging.info(
                f"action: apuesta_recibida | result: success | cantidad: {len(bet_batch.bets)}"
            )
        except Exception:
            result_code = ResultCode.ERROR
            logging.info(
                f"action: apuesta_recibida | result: fail | cantidad: {len(bet_batch.bets)}"
            )

        confirmation_message = ConfirmBetBatch(
            agency_id=bet_batch.agency_id,
            bet_amount=len(bet_batch.bets),
            result_code=result_code,
        )
        self._protocol.send_message(confirmation_message)
