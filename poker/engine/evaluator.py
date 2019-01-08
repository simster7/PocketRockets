if __name__ == '__main__':
    from card import Card
else:
    from .card import Card

def check_straight_flush(hand):
    """
    Returns True if the hand conatins a straight flush.
    Returns composition of four of a kind for tie-breaking, it is structured as ([rank id])
    """
    suits = [card.suit_id for card in hand]
    frequencies = [s_id for s_id in set(suits) if suits.count(s_id) >= 5] 
    if not len(frequencies) == 1:
        return False, None
    suit = frequencies[0]
    cards = [card for card in hand if card.suit_id == suit]
    return check_straight(cards)

def check_four_of_a_kind(hand):
    """
    Returns True if the hand conatins a four of a kind, hand could be better than a four of a kind and check_four_of_a_kind would still return true.
    Returns composition of four of a kind for tie-breaking, it is structured as ([rank id])
    """
    values = [card.rank_id for card in hand]
    frequencies = [v_id for v_id in set(values) if values.count(v_id) == 4]
    if not len(frequencies) == 1:
        return False, None
    return True, (frequencies[0],)

def check_full_house(hand):
    """
    Returns True if the hand conatins a full house, hand could be better than a full house and check_full_house would still return true.
    Returns composition of full house for tie-breaking, it is structured as ([trip rank id], [pair rank id])
    """
    trip_check = check_three_of_a_kind(hand)
    pair_check = check_pair(hand)
    if trip_check[0] and pair_check[0]:
        return True, (trip_check[1][0], pair_check[1][0])
    if trip_check[0]:
        new_hand = [card for card in hand if card.rank_id != trip_check[1][0]]
        second_trip_check = check_three_of_a_kind(new_hand)
        if second_trip_check[0]:
            return True, (trip_check[1][0], second_trip_check[1][0])
    return False, None

def check_flush(hand):
    """
    Returns True if the hand conatins a flush, hand could be better than a flush and check_flush would still return true.
    Returns high card for tie-breaking, it is structured as ([high card rank id],)
    """
    suits = [card.suit_id for card in hand]
    frequencies = [s_id for s_id in set(suits) if suits.count(s_id) >= 5] 
    if not len(frequencies) == 1:
        return False, None
    suit = frequencies[0]
    ordered_flush = sorted([card.rank_id for card in hand if card.suit_id == suit], reverse=True)[:5]
    return True, tuple(ordered_flush)

def check_straight(hand):
    """
    Returns True if the hand conatins a straight, hand could be better than a straight and check_straight would still return true.
    Returns high card for tie-breaking, it is structured as ([high card rank id],)
    """
    values = sorted(list(set([card.rank_id for card in hand])), reverse=True)
    if len(values) <  5:
        return False, None
    if 12 in values:
        values.append(-1)
    for i in range(len(values) - 4):
        constraints = all([(values[j] - values[j+1] == 1) for j in range(i, i+4)])
        if constraints:
            return True, (values[i],)

    return False, None

def check_three_of_a_kind(hand):
    """
    Returns True if the hand conatins a three of a kind, hand could be better than a three of a kind and check_three_of_a_kind would still return true.
    Returns remaining cards for tie-breaking, it is structured as ([three of a kind rank id], *[sorted kicker rank ids])
    If hand is actually better than a three of a kind, this function returns undefined behavior.
    """
    values = [card.rank_id for card in hand]
    frequencies = sorted([v_id for v_id in set(values) if values.count(v_id) == 3], reverse=True)
    if not len(frequencies) >= 1:
        return False, None
    trip = frequencies[0]
    remaining = sorted([value for value in values if value not in frequencies], reverse=True)[:2]
    return True, (trip, *remaining)

def check_two_pair(hand):
    """
    Returns True if the hand conatins two pairs, hand could be better than two pair and check_two_pair would still return true.
    Returns remaining cards for tie-breaking, it is structured as ([high pair rank id], [low pair rank id], *[sorted kicker rank ids])
    If hand is actually better than two pair, this function returns undefined behavior.
    """
    values = [card.rank_id for card in hand]
    frequencies = sorted([v_id for v_id in set(values) if values.count(v_id) == 2], reverse=True)
    if not len(frequencies) >= 2:
        return False, None
    pair1 = frequencies[0]
    pair2 = frequencies[1]
    remaining = sorted([value for value in values if value not in [pair1, pair2]], reverse=True)[:1]
    return True, (pair1, pair2, remaining[0])

def check_pair(hand):
    """
    Returns True if the hand conatins at only one pair, hand could be better than one pair and check_pair would still return true.
    Returns remaining cards for tie-breaking, it is structured as ([pair rank id], *[sorted kicker rank ids])
    """
    values = [card.rank_id for card in hand]
    frequencies = [rank_id for rank_id in set(values) if values.count(rank_id) == 2]
    if not len(frequencies) == 1:
        return False, None
    pair = frequencies[0]
    remaining = sorted([value for value in values if value != pair], reverse=True)[:3]
    return True, (pair, *remaining)

def check_high_card(hand):
    """
    Always returns True, because hand is always at least high card good. Returns ordered cards for tie-breaking
    """
    values = sorted([card.rank_id for card in hand], reverse=True)
    return True, values[:5]


HAND_CHECKS = [ check_straight_flush,   # 9
                check_four_of_a_kind,   # 8
                check_full_house,       # 7
                check_flush,            # 6
                check_straight,         # 5
                check_three_of_a_kind,  # 4
                check_two_pair,         # 3
                check_pair,             # 2
                check_high_card         # 1
            ]

def calculate_hand(hand):
    POSSIBLE_HANDS = len(HAND_CHECKS)
    for index, check in enumerate(HAND_CHECKS):
        hit, state = check(hand)
        if hit:
            score = POSSIBLE_HANDS - index
            return (score, *state)
