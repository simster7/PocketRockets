SUIT_MAP = {
    0: 'S',
    1: 'H',
    2: 'C',
    3: 'D',
}

RANK_MAP = {
    0: '2',
    1: '3',
    2: '4',
    3: '5',
    4: '6',
    5: '7',
    6: '8',
    7: '9',
    8: 'T',
    9: 'J',
    10: 'Q',
    11: 'K',
    12: 'A',
}


class Card:

    @staticmethod
    def to_card_id(rank: int, suit: int) -> int:
        if rank not in RANK_MAP.values() or suit not in SUIT_MAP.values():
            raise Exception('Card is malformed')

        _id = 0
        for rank_id, card_rank in RANK_MAP.items():
            if rank == card_rank:
                _id += rank_id
        for suit_id, card_suit in SUIT_MAP.items():
            if suit == card_suit:
                _id += (suit_id * 13)
        return _id

    def __init__(self, _id: int) -> None:
        if _id < 0 or _id >= 52 or type(_id) != int:
            raise Exception('Card id must be an integer in [0, 51]')
        self.card_id = _id

    @property
    def suit_id(self) -> int:
        return self.card_id // 13

    @property
    def suit(self) -> str:
        return SUIT_MAP[self.suit_id]

    @property
    def rank_id(self) -> int:
        return self.card_id % 13

    @property
    def rank(self) -> str:
        return RANK_MAP[self.rank_id]

    def __str__(self):
        return '{}{}'.format(self.rank, self.suit)

    def __repr__(self):
        return 'Card: ' + str(self)
