from ..engine.util import hand_generator
from ..engine.evaluator import *
from ..engine.card import RANK_MAP


class TestCheckHighCard:
    """
    check_high_card
    """

    def test_high_card_correct(self):
        """
        High card evaluated correctly
        """
        pair_hand = hand_generator('8S 7H TH KD 4C')
        result = check_high_card(pair_hand)
        assert result == (True, [11, 8, 6, 5, 2])

    def test_high_card_correct_seven(self):
        """
        High card evaluated correctly even for seven cards
        """
        pair_hand = hand_generator('8S 7H TH KD 4C 3C 2C')
        result = check_high_card(pair_hand)
        assert result == (True, [11, 8, 6, 5, 2])

class TestCheckPair:
    """
    check_pair
    """

    def test_one_pair_correct(self):
        """
        One pair correctly evaluated
        """
        pair_hand = hand_generator('8S 8H TH KD 4C')
        result = check_pair(pair_hand)
        assert result == (True, (6, 11, 8, 2))

    def test_one_pair_correct_with_seven(self):
        """
        One pair correctly evaluated with seven cards
        """
        pair_hand = hand_generator('8S 8H KH JD 4C AH 9C')
        result = check_pair(pair_hand)
        assert result == (True, (6, 12, 11, 9))

    def test_no_pair_not_evaluated(self):
        """
        No pair not evaluated, even with seven cards
        """
        pair_hand = hand_generator('8S 7H KH TD QC 2S 4H')
        result = check_pair(pair_hand)
        assert result == (False, None)

class TestCheckTwoPair:
    """
    check_two_pair
    """

    def test_two_pair_correct(self):
        """
        Two pair correctly evaluated
        """
        pair_hand = hand_generator('8S 8H TH TD 4C')
        result = check_two_pair(pair_hand)
        print(result)
        assert result == (True, (8, 6, 2))

    def test_two_pair_correct_with_seven(self):
        """
        Two pair correctly evaluated with seven cards
        """
        pair_hand = hand_generator('8S 8H TH TD 4C 6C 9S')
        result = check_two_pair(pair_hand)
        print(result)
        assert result == (True, (8, 6, 7))

    def test_two_pair_correct_with_seven_and_three_fair(self):
        """
        Two pair correctly evaluated with seven cards and precence of three pair
        """
        pair_hand = hand_generator('8S 8H TH TD 9C 9S 6C')
        result = check_two_pair(pair_hand)
        print(result)
        assert result == (True, (8, 7, 6))

    def test_no_pair_not_evaluated(self):
        """
        No two pair not evaluated, even with seven cards
        """
        pair_hand = hand_generator('8S 7H KH TD QC 2S 4H')
        result = check_two_pair(pair_hand)
        print(result)
        assert result == (False, None)

    def test_no_pair_not_evaluated_with_pair(self):
        """
        No two pair not evaluated, even with seven cards and a single pair
        """
        pair_hand = hand_generator('8S 8H KH TD QC 2S 4H')
        result = check_two_pair(pair_hand)
        print(result)
        assert result == (False, None)

class TestCheckThreeOfAKind:
    """
    check_three_of_a_kind
    """

    def test_three_of_a_kind_correct(self):
        """
        Three of a kind correctly evaluated
        """
        pair_hand = hand_generator('8S 8H 8C TD 4C')
        result = check_three_of_a_kind(pair_hand)
        print(result)
        assert result == (True, (6, 8, 2))

    def test_three_of_a_kind_correct_with_seven(self):
        """
        Three of a kind correctly evaluated, even with seven cards
        """
        pair_hand = hand_generator('8S 8H 8C TD 4C JC 2S')
        result = check_three_of_a_kind(pair_hand)
        print(result)
        assert result == (True, (6, 9, 8))

    def test_no_three_of_a_kind_evaluated_with_seven(self):
        """
        No three of a kind evaluated, even with seven cards
        """
        pair_hand = hand_generator('8S 8H 7C TD 4C JC 2S')
        result = check_three_of_a_kind(pair_hand)
        assert result == (False, None)

    def test_three_of_a_kind_correct_with_seven_and_pair(self):
        """
        Three of a kind correctly evaluated, even with seven cards and extra pair (full house)
        !! This is undefined behavior so if test fails in the future it is safe to modify/delete !!
        """
        pair_hand = hand_generator('8S 8H 8C TD TC JC 2S')
        result = check_three_of_a_kind(pair_hand)
        print(result)
        assert result == (True, (6, 9, 8))

    def test_three_of_a_kind_correct_with_seven_and_three(self):
        """
        Three of a kind correctly evaluated, even with seven cards and three of a kind (full house).
        !! This is undefined behavior so if test fails in the future it is safe to modify/delete !!
        """
        pair_hand = hand_generator('8S 8H 8C TD TC TS 2S')
        result = check_three_of_a_kind(pair_hand)
        print(result)
        assert result == (True, (8, 0))

class TestCheckStraight:
    """
    check_straight
    """

    def test_check_straight(self):
        """
        Straight check works correctly
        """
        pair_hand = hand_generator('2S 3H 4C 5D 6C')
        result = check_straight(pair_hand)
        assert result == (True, (4,))

    def test_check_straight_with_seven(self):
        """
        Straight check works correctly, even with seven cards
        """
        pair_hand = hand_generator('2S 3H 4C 5D 6C KC JS')
        result = check_straight(pair_hand)
        assert result == (True, (4,))

    def test_check_straight_with_seven_long(self):
        """
        Straight check works correctly, even with seven cards and straight is more than five cards
        """
        pair_hand = hand_generator('2S 3H 4C 5D 6C 7C JS')
        result = check_straight(pair_hand)
        assert result == (True, (5,))

    def test_check_straight_with_seven_ace(self):
        """
        Straight check works correctly, even with seven cards and is ace low
        """
        pair_hand = hand_generator('2S 3H 4C 5D AC 9C JS')
        result = check_straight(pair_hand)
        assert result == (True, (3,))

    def test_check_straight_with_seven_ace_high(self):
        """
        Straight check works correctly, even with seven cards and is ace high
        """
        pair_hand = hand_generator('TS QH KC 5D AC 9C JS')
        result = check_straight(pair_hand)
        assert result == (True, (12,))

    def test_check_straight_with_seven_no_straight(self):
        """
        No straight when there is not straight
        """
        pair_hand = hand_generator('8S QH KC 6D AC 9C JS')
        result = check_straight(pair_hand)
        assert result == (False, None)

    def test_check_straight_with_seven_no_straight_wrap(self):
        """
        No straight when there is straight wrap-around
        """
        for start in [9, 10, 11]:
            pair_hand = hand_generator(' '.join([RANK_MAP[card % 13] + 'S' for card in range(start, start + 5)]))
            result = check_straight(pair_hand)
            assert result == (False, None)

class TestCheckFlush:
    """
    check_flush
    """

    def test_check_flush(self):
        """
        Flush check works correctly
        """
        pair_hand = hand_generator('JS 3S TS 5S 6S')
        result = check_flush(pair_hand)
        assert result == (True, (9,))

    def test_check_flush_with_seven(self):
        """
        Flush check works correctly with seven cards, gets high card correctly
        """
        pair_hand = hand_generator('JS 3S TS 5S 6S AS 8C')
        result = check_flush(pair_hand)
        assert result == (True, (12,))

    def test_check_flush_with_no_flush(self):
        """
        Flush not detected when there is no flush
        """
        pair_hand = hand_generator('JS 3S TS 5S 6C AC 8C')
        result = check_flush(pair_hand)
        assert result == (False, None)

class TestCheckFullHouse:
    """
    check_full_house
    """

    def test_check_full_house(self):
        """
        Full house check works correctly
        """
        pair_hand = hand_generator('JS JC JD 6C 6S')
        result = check_full_house(pair_hand)
        assert result == (True, (9, 4))

    def test_check_full_house_seven(self):
        """
        Full house check works correctly with seven cards
        """
        pair_hand = hand_generator('AS AC AD 6C 6S 2S KC')
        result = check_full_house(pair_hand)
        assert result == (True, (12, 4))

    def test_check_full_house_seven_two_trip(self):
        """
        Full house check works correctly with seven cards and two three of a kinds
        """
        pair_hand = hand_generator('JS JC JD 6C 6S 6D KC')
        result = check_full_house(pair_hand)
        assert result == (True, (9, 4))

    def test_check_full_house_seven_two_pair(self):
        """
        No full house when there is no full house
        """
        pair_hand = hand_generator('JS JC KD 6C 6S 5D KC')
        result = check_full_house(pair_hand)
        assert result == (False, None)

class TestCheckFourOfAKind:
    """
    check_four_of_a_kind
    """

    def test_check_four_of_a_kind(self):
        """
        Four of a kind check works correctly
        """
        pair_hand = hand_generator('JS JC JD JH AS')
        result = check_four_of_a_kind(pair_hand)
        assert result == (True, (9))

    def test_check_four_of_a_kind_with_seven(self):
        """
        Four of a kind check works correctly with seven cards
        """
        pair_hand = hand_generator('2S 2C 2D 2H AS KC 7C')
        result = check_four_of_a_kind(pair_hand)
        assert result == (True, (0))

    def test_check_four_of_a_kind_with_seven_and_trips(self):
        """
        Four of a kind check works correctly with seven cards and extra trips
        """
        pair_hand = hand_generator('5S 5C 5D 5H AS AC AH')
        result = check_four_of_a_kind(pair_hand)
        assert result == (True, (3))

    def test_check_four_of_a_kind_with_seven_no_four(self):
        """
        No four of a kind when there is none
        """
        pair_hand = hand_generator('JS JC JD 6C 6S 6D KC')
        result = check_four_of_a_kind(pair_hand)
        assert result == (False, None)

class TestStraightFlush:
    """
    check_straight_flush
    """

    def test_check_straight_flush(self):
        """
        Straight flush check works correctly
        """
        pair_hand = hand_generator('4S 5S 6S 7S 8S')
        result = check_straight_flush(pair_hand)
        assert result == (True, (6,))

    def test_check_straight_flush_seven(self):
        """
        Straight flush check works correctly even with seven cards
        """
        pair_hand = hand_generator('7S 8S 9S TS JS AH TC')
        result = check_straight_flush(pair_hand)
        assert result == (True, (9,))

    def test_check_straight_flush_seven_no_flush(self):
        """
        No straight flush if there is no flush
        """
        pair_hand = hand_generator('7S 8S 9H TS JS AH TC')
        result = check_straight_flush(pair_hand)
        assert result == (False, None)

    def test_check_straight_flush_seven_no_straight(self):
        """
        No straight flush if there is no straight
        """
        pair_hand = hand_generator('6S 8S 9S TS JS AH TC')
        result = check_straight_flush(pair_hand)
        assert result == (False, None)
