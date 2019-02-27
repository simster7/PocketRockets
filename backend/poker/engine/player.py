from typing import Optional


class Player:
    name: str
    stack: int
    seat_number: Optional[int]

    def __init__(self, _name: str) -> None:
        self.name = _name
        self.stack = 0
        self.seat_number = None

    def make_bet(self, bet_size: int) -> int:
        assert self.stack >= bet_size, "Player may not bet amount larger than their stack"
        self.stack -= bet_size
        return bet_size

    def receive_pot(self, pot_size: int) -> int:
        self.stack += pot_size
        return pot_size

    def set_stack(self, stack: int) -> None:
        self.stack = stack

    def get_stack(self) -> int:
        return self.stack

    def set_name(self, name: str) -> None:
        self.name = name

    def get_name(self) -> str:
        return self.name

    def set_seat_number(self, seat_number: Optional[int]) -> None:
        self.seat_number = seat_number

    def get_seat_number(self) -> int:
        return self.seat_number

    def __repr__(self):
        return str(self.__dict__)
