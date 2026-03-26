import logging

from .bet_batch import BetBatch
from .protocol.lottery_central_protocol import LotteryCentralProtocol
from .socket.safe_socket import SafeSocket
from . import utils
from .messages.confirm_bet_batch import ConfirmBetBatch
from .messages.give_winners_list import GiveWinnersList
from .constants.result_codes import ResultCode
from .monitors.bets_storage_monitor import BetsStorageMonitor
from .monitors.lottery_event_monitor import LotteryEventMonitor


class LotteryCentral:
    def __init__(
        self,
        safe_socket: SafeSocket,
        bets_storage_monitor: BetsStorageMonitor,
        lottery_event_monitor: LotteryEventMonitor,
    ):
        self._protocol = LotteryCentralProtocol(safe_socket)
        self._bets_storage_monitor = bets_storage_monitor
        self._lottery_event_monitor = lottery_event_monitor

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
            self._bets_storage_monitor.store_bets(bet_batch.bets)
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

    def finalize_bets(self, agency_id: int):
        logging.info(
            f"action: finalize_bets | result: success | agency_id: {agency_id}"
        )
        self._lottery_event_monitor.wait()

    def give_winners_list(self, agency_id: int):
        winners = [
            int(bet.document)
            for bet in self._bets_storage_monitor.load_bets()
            if bet.agency == agency_id and utils.has_won(bet)
        ]
        response = GiveWinnersList(agency_id=agency_id, winners=winners)
        self._protocol.send_message(response)
        logging.info(
            f"action: give_winners_list | result: success | agency_id: {agency_id} | cant_ganadores: {len(winners)}"
        )
