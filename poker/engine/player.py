from hand import Hand


class Player:
    def __init__(self, _name, initial_stack):
        self.name = _name
        self.stack = initial_stack
        self.current_bet = 0
        self.hand = None

    def make_bet(self, bet_size):
        assert (
            self.stack >= bet_size
        ), "Player may not bet amount larger than their stack"
        self.stack -= bet_size
        self.current_bet += bet_size
        return bet_size

    def recieve_pot(self, pot_size):
        self.stack += pot_size
        return pot_size

    def recieve_cards(self, cards):
        self.hand = Hand(cards, self)

    def reset(self):
        self.current_bet = 0 

    def hand(self):
        return self.hand

    def stack(self):
        return self.stack

    def name(self):
        return self.name

    def __str__(self):
        return f"Player(Name: {self.name}, Stack: {self.stack}, Bet: {self.current_bet})" 

    def __repr__(self):
        return str(self)
