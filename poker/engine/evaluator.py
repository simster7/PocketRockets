from dataclasses import dataclass
from typing import List, Tuple, Optional, Callable, Dict

from .card import Card


@dataclass
class CheckResult:
    match: bool
    tiebreakers: Optional[Tuple[int, ...]]


class RankedHand:
    hand: List[Card]
    rank: Tuple[int, ...]
    player_index: int
    hand_name: str

    def __init__(self, hand: List[Card], rank: Tuple[int, ...], player_index: int):
        self.hand = hand
        self.rank = rank
        self.player_index = player_index
        self.hand_name = HAND_NAMES[self.rank[0]]

    def __str__(self):
        return str(self.player_index) + ": " + self.hand_name + " (" + str(self.hand) + ")"

    def __repr__(self):
        return str(self)


def check_straight_flush(hand: List[Card]) -> CheckResult:
    """
    Returns True if the hand contains a straight flush.
    Returns composition of four of a kind for tie-breaking, it is structured as ([rank id])
    """
    suits = [card.suit_id for card in hand]
    frequencies = [s_id for s_id in set(suits) if suits.count(s_id) >= 5]
    if not len(frequencies) == 1:
        return CheckResult(False, None)
    suit = frequencies[0]
    cards = [card for card in hand if card.suit_id == suit]
    return check_straight(cards)


def check_four_of_a_kind(hand: List[Card]) -> CheckResult:
    """
    Returns True if the hand contains a four of a kind, hand could be better than a four of a kind and
    check_four_of_a_kind would still return true.
    Returns composition of four of a kind for tie-breaking, it is structured as ([rank id])
    """
    values = [card.rank_id for card in hand]
    frequencies = [v_id for v_id in set(values) if values.count(v_id) == 4]
    if not len(frequencies) == 1:
        return CheckResult(False, None)
    return CheckResult(True, (frequencies[0],))


def check_full_house(hand: List[Card]) -> CheckResult:
    """
    Returns True if the hand contains a full house, hand could be better than a full house and check_full_house would
    still return true.
    Returns composition of full house for tie-breaking, it is structured as ([trip rank id], [pair rank id])
    """
    trip_check = check_three_of_a_kind(hand)
    pair_check = check_pair(hand)
    if trip_check.match and pair_check.match:
        return CheckResult(True, (trip_check.tiebreakers[0], pair_check.tiebreakers[0]))
    if trip_check.match:
        new_hand = [card for card in hand if card.rank_id != trip_check.tiebreakers[0]]
        second_trip_check = check_three_of_a_kind(new_hand)
        if second_trip_check.match:
            return CheckResult(True, (trip_check.tiebreakers[0], second_trip_check.tiebreakers[0]))
    return CheckResult(False, None)


def check_flush(hand: List[Card]) -> CheckResult:
    """
    Returns True if the hand contains a flush, hand could be better than a flush and check_flush would still return
    true.
    Returns high card for tie-breaking, it is structured as ([high card rank id],)
    """
    suits = [card.suit_id for card in hand]
    frequencies = [s_id for s_id in set(suits) if suits.count(s_id) >= 5]
    if not len(frequencies) == 1:
        return CheckResult(False, None)
    suit = frequencies[0]
    ordered_flush = sorted([card.rank_id for card in hand if card.suit_id == suit], reverse=True)[:5]
    return CheckResult(True, tuple(ordered_flush))


def check_straight(hand: List[Card]) -> CheckResult:
    """
    Returns True if the hand contains a straight, hand could be better than a straight and check_straight would still
    return true.
    Returns high card for tie-breaking, it is structured as ([high card rank id],)
    """
    values = sorted(list(set([card.rank_id for card in hand])), reverse=True)
    if len(values) < 5:
        return CheckResult(False, None)
    if 12 in values:
        values.append(-1)
    for i in range(len(values) - 4):
        constraints = all([(values[j] - values[j + 1] == 1) for j in range(i, i + 4)])
        if constraints:
            return CheckResult(True, (values[i],))

    return CheckResult(False, None)


def check_three_of_a_kind(hand: List[Card]) -> CheckResult:
    """
    Returns True if the hand contains a three of a kind, hand could be better than a three of a kind and
    check_three_of_a_kind would still return true.
    Returns remaining cards for tie-breaking, it is structured as ([three of a kind rank id], *[sorted kicker rank ids])
    If hand is actually better than a three of a kind, this function returns undefined behavior.
    """
    values = [card.rank_id for card in hand]
    frequencies = sorted([v_id for v_id in set(values) if values.count(v_id) == 3], reverse=True)
    if not len(frequencies) >= 1:
        return CheckResult(False, None)
    trip = frequencies[0]
    remaining = sorted([value for value in values if value not in frequencies], reverse=True)[:2]
    return CheckResult(True, (trip, *remaining))


def check_two_pair(hand: List[Card]) -> CheckResult:
    """
    Returns True if the hand contains two pairs, hand could be better than two pair and check_two_pair would still
    eturn true.
    Returns remaining cards for tie-breaking, it is structured as ([high pair rank id], [low pair rank id],
    *[sorted kicker rank ids])
    If hand is actually better than two pair, this function returns undefined behavior.
    """
    values = [card.rank_id for card in hand]
    frequencies = sorted([v_id for v_id in set(values) if values.count(v_id) == 2], reverse=True)
    if not len(frequencies) >= 2:
        return CheckResult(False, None)
    pair1 = frequencies[0]
    pair2 = frequencies[1]
    remaining = sorted([value for value in values if value not in [pair1, pair2]], reverse=True)[:1]
    return CheckResult(True, (pair1, pair2, remaining[0]))


def check_pair(hand: List[Card]) -> CheckResult:
    """
    Returns True if the hand contains at only one pair, hand could be better than one pair and check_pair would still
    return true.
    Returns remaining cards for tie-breaking, it is structured as ([pair rank id], *[sorted kicker rank ids])
    """
    values = [card.rank_id for card in hand]
    frequencies = [rank_id for rank_id in set(values) if values.count(rank_id) == 2]
    if not len(frequencies) == 1:
        return CheckResult(False, None)
    pair = frequencies[0]
    remaining = sorted([value for value in values if value != pair], reverse=True)[:3]
    return CheckResult(True, (pair, *remaining))


def check_high_card(hand: List[Card]) -> CheckResult:
    """
    Always returns True, because hand is always at least high card good. Returns ordered cards for tie-breaking
    """
    values = sorted([card.rank_id for card in hand], reverse=True)
    return CheckResult(True, tuple(values[:5]))


HAND_CHECKS: List[Callable] = [check_straight_flush,  # 9
                               check_four_of_a_kind,  # 8
                               check_full_house,  # 7
                               check_flush,  # 6
                               check_straight,  # 5
                               check_three_of_a_kind,  # 4
                               check_two_pair,  # 3
                               check_pair,  # 2
                               check_high_card  # 1
                               ]

HAND_NAMES: Dict[int, str] = {
    9: "Straight Flush",
    8: "Four Of A Kind",
    7: "Full House",
    6: "Flush",
    5: "Straight",
    4: "Three Of A Kind",
    3: "Two Pair",
    2: "Pair",
    1: "High Card"
}


def calculate_hand(hand: List[Card]) -> Tuple[int, ...]:
    possible_hands = len(HAND_CHECKS)
    for index, check in enumerate(HAND_CHECKS):
        check_result: CheckResult = check(hand)
        if check_result.match:
            score = possible_hands - index
            return (score, *check_result.tiebreakers)


def get_hand_rank_with_name(player_index: int, hand: List[Card]) -> RankedHand:
    return RankedHand(hand, calculate_hand(hand), player_index)


def rank_hands(list_of_hands: List[Tuple[int, List[Card]]]) -> List[RankedHand]:
    return list(sorted(map(lambda pair: get_hand_rank_with_name(*pair), list_of_hands), key=lambda hand: hand.rank,
                       reverse=True))
