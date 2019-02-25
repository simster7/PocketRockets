class Player:

    def __init__(self, _name: str) -> None:
        self.name: str = _name
        self.stack: int = 0
        self.seat_number: int = None

    def make_bet(self, bet_size: int) -> int:
        assert self.stack >= bet_size, "Player may not bet amount larger than their stack"
        self.stack -= bet_size
        return bet_size

    def receive_pot(self, pot_size: int) -> int:
        self.stack += pot_size
        return pot_size

    def set_stack(self, stack: int):
        self.stack = stack

    def get_stack(self) -> int:
        return self.stack

    def set_name(self, name: str):
        self.name = name

    def get_name(self) -> str:
        return self.name

    def set_seat_number(self, seat_number: int):
        self.seat_number = seat_number

    def get_seat_number(self) -> int:
        return self.seat_number

    def __str__(self):
        return 'Player({}, {})'.format(self.name, self.stack)

    def __repr__(self):
        return str(self)
