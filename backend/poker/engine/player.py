from typing import Optional

from .action import Action


class Player:
    name: str
    stack: int
    seat_number: Optional[int]
    folded: bool
    last_action: Optional[Action]
    sitting_out: bool

    def __init__(self, _name: str) -> None:
        self.name = _name
        self.stack = 0
        self.seat_number = None
        self.folded = False
        self.last_action = None
        self.sitting_out = False

    def make_bet(self, bet_size: int) -> int:
        assert self.stack >= bet_size, "Player may not bet amount larger than their stack"
        self.stack -= bet_size
        return bet_size

    def receive_pot(self, pot_size: int) -> int:
        self.stack += pot_size
        return pot_size

    def __repr__(self):
        return str(self.__dict__)
