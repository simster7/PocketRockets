from ..engine.util import hand_generator
from ..engine.evaluator import *


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

