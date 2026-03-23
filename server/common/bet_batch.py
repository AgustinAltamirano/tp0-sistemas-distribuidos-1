from .utils import Bet


class BetBatch:
    def __init__(self, agency_id: int, bets: list[Bet]):
        self.agency_id = agency_id
        self.bets = bets
