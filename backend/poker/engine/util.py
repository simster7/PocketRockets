from typing import List

from .card import Card


def hand_generator(hand_string: str) -> List[Card]:
    card_strings = hand_string.split(' ')
    return [Card(Card.to_card_id(rank, suit)) for rank, suit in card_strings]