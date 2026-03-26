from enum import Enum


class MessageCode(Enum):
    REGISTER_BET_BATCH = 1
    CONFIRM_BET_BATCH = 2
    FINALIZE_BETS = 3
    ASK_WINNERS_LIST = 4
    GIVE_WINNERS_LIST = 5
