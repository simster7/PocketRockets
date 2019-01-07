if __name__ == '__main__':
    from card import Card
else:
    from .card import Card

def straight_flush_check(hand):
    suits = [card.s_id for card in hand]
    frequencies = [s_id for s_id in set(suits) if suits.count(s_id) >= 5] 
    flush_hit = len(frequencies) == 1 
    if not flush_hit:
        return flush_hit, None
    suit = frequencies[0]
    cards = [card for card in hand if card.s_id == suit]
    return straight_check(cards)

def quad_check(hand):
    values = [card.rank for card in hand]
    frequencies = [v_id for v_id in set(values) if values.count(v_id) == 4]
    hit = len(frequencies) == 1 
    if not hit:
        return False, None
    return hit, (frequencies[0])

def boat_check(hand):
    values = [card.rank for card in hand]
    trips = [v_id for v_id in set(values) if values.count(v_id) >= 3]
    if not trips:
        return False, None
    top_trip = max(trips)
    pairs = [v_id for v_id in set(values) if values.count(v_id) >= 2 and v_id != top_trip]
    if not pairs:
        return False, None
    top_pair = max(pairs)
    return True, (top_trip, top_pair)

def flush_check(hand):
    suits = [card.suit for card in hand]
    frequencies = [s_id for s_id in set(suits) if suits.count(s_id) >= 5] 
    hit = len(frequencies) == 1 
    if not hit:
        return hit, None
    suit = frequencies[0]
    order = sorted(hand, lambda x: -x)
    return hit, (card.rank for card in order if card.suit == suit)[:5]

def check_straight(hand):
    """
    Returns True if the hand conatins a straight, hand could be better than a straight and check_straight would still return true.
    Returns high card for tie-breaking, it is structured as ([high card rank id],)
    If hand is actually better than a three of a kind, this function returns undefined behavior.
    """
    values = sorted(list(set([card.rank_id for card in hand])), key=lambda x: -x)
    if len(values) <  5:
        return False, None
    if 12 in values:
        values.append(-1)
    for i in range(len(values) - 4):
        constraints = all([(values[j] - values[j+1] == 1) for j in range(i, i+4)])
        print(constraints, i, values)
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
    frequencies = sorted([v_id for v_id in set(values) if values.count(v_id) == 3], key=lambda x: -x)
    if not len(frequencies) >= 1:
        return False, None
    trip = frequencies[0]
    remaining = sorted([value for value in values if value not in frequencies], key=lambda x: -x)[:2]
    return True, (trip, *remaining)

def check_two_pair(hand):
    """
    Returns True if the hand conatins two pairs, hand could be better than two pair and check_two_pair would still return true.
    Returns remaining cards for tie-breaking, it is structured as ([high pair rank id], [low pair rank id], *[sorted kicker rank ids])
    If hand is actually better than two pair, this function returns undefined behavior.
    """
    values = [card.rank_id for card in hand]
    frequencies = sorted([v_id for v_id in set(values) if values.count(v_id) == 2], key=lambda x: -x)
    if not len(frequencies) >= 2:
        return False, None
    pair1 = frequencies[0]
    pair2 = frequencies[1]
    remaining = sorted([value for value in values if value not in [pair1, pair2]], key=lambda x: -x)[:1]
    return True, (pair1, pair2, remaining[0])

def check_pair(hand):
    """
    Returns True if the hand conatins at only one pair, hand could be better than one pair and check_pair would still return true.
    Returns remaining cards for tie-breaking, it is structured as ([pair rank id], *[sorted kicker rank ids])
    """
    values = [card.rank_id for card in hand]
    frequencies = [v_id for v_id in set(values) if values.count(v_id) == 2]
    if not len(frequencies) == 1:
        return False, None
    pair = frequencies[0]
    remaining = sorted([value for value in values if value != pair], key=lambda x: -x)[:3]
    return True, (pair, *remaining)

def check_high_card(hand):
    """
    Always returns True, because hand is always at least high card good. Returns ordered cards for tie-breaking
    """
    values = sorted([card.rank_id for card in hand], key=lambda x: -x)
    return True, values[:5]


# HAND_CHECKS = [straight_flush_check, quad_check, boat_check, flush_check, straight_check,
#                trip_check, two_pair_check, check_pair, high_card_check]

# def calculate_hand(hand):
#     for index, check in enumerate(HAND_CHECKS):
#         hit, state = check(hand)
#         if hit:
#             score = len(HAND_CHECKS) - index
#             return (score, *state)
