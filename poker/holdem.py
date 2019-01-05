ANTE = 0.1

SUITS = {0: "Spade",
         1: "Heart",
         2: "Club",
         3: "Diamond"}

VALUES = {**{val: str(val) for val in range(2, 11)},
          **{11: "J", 12: "Q", 13: "K", 14: "A"}}

class Card:
    def __init__(self, v_id, s_id):
        self.v_id = v_id
        self.value = VALUES[v_id] 
        self.s_id = s_id
        self.suit = SUITS[s_id] 

class Player:
    def __init__(self, stack):
        self.stack = stack
        #All-in check

    def receive(self, card_1, card_2):
        self.hand = [card_1, card_2] 
    def bet(self, amount):
        return
    def call(self, amount):
        return
    def fold(self):
        return

class Game:
    def __init__(self, players):
        self.players = players
        self.deck = []
        for v_id in range(2, 15):
            for s_id in range(4):
                self.deck.append(Card(v_id, s_id))
        self.deck.shuffle()
        #It's random so dealing this way is fine
        self.deal()
    def deal(self):
        for player in self.players:
            player.receive(self.deck.pop(), self.deck.pop())

# Here 'hand' should just be a list of cards

def straight_flush_check(hand):
    #TODO
    return False, None

def quad_check(hand):
    values = [card.v_id for card in hand]
    frequencies = [v_id for v_id in set(values) if values.count(v_id) == 4]
    hit = len(frequencies) == 1 
    if not hit:
        return False, None
    return hit, (frequencies[0])

def boat_check(hand):
    values = [card.v_id for card in hand]
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
    suits = [card.s_id for card in hand]
    frequencies = [s_id for s_id in set(suits) if suits.count(s_id) >= 5] 
    hit = len(frequencies) == 1 
    if not hit:
        return hit, None
    suit = frequencies[0]
    order = sorted(hand, lambda x: -x)
    return hit, (card.v_id for card in order if card.s_id == suit)[:5]

def straight_check(hand):
    #TODO
    return False, None

def trip_check(hand):
    values = [card.v_id for card in hand]
    frequencies = [v_id for v_id in set(values) if values.count(v_id) >= 3]
    hit = len(frequencies) > 0
    trip = max(frequencies)
    remaining = sorted(list(set(values) - {trip}))
    kicker1 = remaining.pop() 
    kicker2 = remaining.pop()
    return hit, (trip, kicker1, kicker2)

def two_pair_check(hand):
    values = [card.v_id for card in hand]
    frequencies = sorted([v_id for v_id in set(values) if values.count(v_id) >= 2], lambda x: -x)
    hit = len(frequencies) >= 2
    if not hit:
        return hit, None
    pair1 = frequencies[0]
    pair2 = frequencies[1]
    remaining = sorted(list(set(values) - {pair1, pair2}))
    kicker = remaining.pop()
    return hit, (pair1, pair2, kicker)

def pair_check(hand):
    values = [card.v_id for card in hand]
    frequencies = [v_id for v_id in set(values) if values.count(v_id) >= 2]
    hit = len(frequencies) >= 1 
    if not hit:
        return hit, None
    pair = max(frequencies)
    remaining = sorted(list(set(values) - {pair}))
    kickers = [remaining.pop(), remaining.pop(), remaining.pop()]
    return hit, (pair1, pair2, *kickers)

def high_card_check(hand):
    values = tuple(sorted(list(set([card.v_id for card in hand])), lambda x: -x))
    return True, values[:5]

HAND_CHECKS = [straight_flush_check, quad_check, boat_check, flush_check, straight_check,
               trip_check, two_pair_check, pair_check, high_card_check]

def calculate_hand(hand):
    for index, check in enumerate(HAND_CHECKS):
        hit, state = check(hand)
        if hit:
            score = len(HAND_CHECKS) - index
            return (score, *state)

def compare(hand1, hand2):
    score1, score2 = calucate_hand(hand1), calucate_hand(hand2)
    if score1 > score2:
        return 1
    elif score1 < score2:
        return -1
    else:
        return 0

