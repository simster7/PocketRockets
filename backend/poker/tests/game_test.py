from typing import List, Callable

from ..engine.card import Card
from ..engine.action import Action
from ..engine.game import Game
from ..engine.player import Player

import json


class GameMock(Game):

    @staticmethod
    def shuffle_deck(deck: List[Card]) -> List[Card]:
        """
        Don't shuffle the deck for testing purposes
        """
        return deck


class TestGame:
    """
    Test a full game
    """

    @staticmethod
    def expect_fail(func: Callable) -> bool:
        try:
            func()
            return False
        except Exception:
            return True

    def test_seating_and_standing(self, snapshot):
        """
        Seating mechanics test
        """
        game = GameMock(1, 2)
        snapshot.assert_match(str(game.__dict__))
        simon = Player("Simon")
        snapshot.assert_match(str(simon.__dict__))
        simon.stack = 100
        game.sit_player(simon, 0)
        snapshot.assert_match(str(simon.__dict__))
        snapshot.assert_match(str(game.__dict__))
        hersh = Player("Hersh")
        hersh.stack = 100
        game.sit_player(hersh, 5)
        game.stand_player(simon, 0)
        snapshot.assert_match(str(simon.__dict__))
        snapshot.assert_match(str(game.__dict__))

        assert self.expect_fail(lambda: game.stand_player(hersh, 2))   # Stand player from incorrect seat
        jarry = Player("Jarry")
        jarry.stack = 100
        assert self.expect_fail(lambda: game.sit_player(jarry, 5))     # Sit player in occupied seat
        assert self.expect_fail(lambda: game.stand_player(simon, 0))   # Stand player already standing

    def test_game_basic(self, snapshot):
        """
        Basic game test: six players, single round, no split pots, no all-ins, no option for big blind
        """

        game = GameMock(1, 2)
        grace = Player("Grace")
        grace.stack = 100
        game.sit_player(grace, 0)
        jason = Player("Jason")
        jason.stack = 100
        game.sit_player(jason, 1)
        simon = Player("Simon")
        simon.stack = 100
        game.sit_player(simon, 3)
        hersh = Player("Hersh")
        hersh.stack = 100
        game.sit_player(hersh, 4)
        chien = Player("Chien")
        chien.stack = 100
        game.sit_player(chien, 5)
        jarry = Player("Jarry")
        jarry.stack = 100
        game.sit_player(jarry, 6)
        snapshot.assert_match(str(game.__dict__))

        game.deal_hand()
        snapshot.assert_match(str(game.game_state.__dict__))
        game.take_action(chien, Action(Action.Actions.bet, 5))
        game.take_action(jarry, Action(Action.Actions.call))
        game.take_action(grace, Action(Action.Actions.fold))
        snapshot.assert_match(str(game.get_player_state(grace).__dict__))
        game.take_action(jason, Action(Action.Actions.call))
        game.take_action(simon, Action(Action.Actions.call))
        game.take_action(hersh, Action(Action.Actions.fold))
        snapshot.assert_match(str(game.game_state.__dict__))
        snapshot.assert_match(str(game.get_player_state(simon).__dict__))
        game.take_action(simon, Action(Action.Actions.check))
        game.take_action(chien, Action(Action.Actions.check))
        game.take_action(jarry, Action(Action.Actions.bet, 10))
        game.take_action(jason, Action(Action.Actions.fold))
        snapshot.assert_match(str(game.get_player_state(jason).__dict__))
        game.take_action(simon, Action(Action.Actions.call))
        game.take_action(chien, Action(Action.Actions.call))
        snapshot.assert_match(str(game.game_state.__dict__))
        game.take_action(simon, Action(Action.Actions.check))
        game.take_action(chien, Action(Action.Actions.check))
        game.take_action(jarry, Action(Action.Actions.check))
        snapshot.assert_match(str(game.game_state.__dict__))
        snapshot.assert_match(str(game.get_player_state(jarry).__dict__))
        game.take_action(simon, Action(Action.Actions.bet, 10))
        game.take_action(chien, Action(Action.Actions.fold))
        game.take_action(jarry, Action(Action.Actions.call))
        snapshot.assert_match(str(game.get_player_state(simon).__dict__))
        snapshot.assert_match(str(game.game_state.__dict__))
        assert self.expect_fail(lambda: game.take_action(simon, Action(Action.Actions.check)))  # Hand is over
        game.deal_hand()
        snapshot.assert_match(str(game.game_state.__dict__))

    def test_game_multiround(self, snapshot):
        """
        Basic game test: three players, multi round, no split pots, no all-ins, no option for big blind
        """

        game = GameMock(1, 2)
        jason = Player("Jason")
        jason.stack = 100
        game.sit_player(jason, 2)
        simon = Player("Simon")
        simon.stack = 100
        game.sit_player(simon, 5)
        chien = Player("Chien")
        chien.stack = 100
        game.sit_player(chien, 7)
        snapshot.assert_match(str(game.__dict__))

        game.deal_hand()
        game.take_action(jason, Action(Action.Actions.bet, 5))
        assert self.expect_fail(lambda: game.take_action(simon, Action(Action.Actions.check)))  # Can't check a call
        game.take_action(simon, Action(Action.Actions.call))
        game.take_action(chien, Action(Action.Actions.call))
        assert game.is_hand_active()
        assert self.expect_fail(lambda: game.take_action(chien, Action(Action.Actions.check)))  # Can't play out of turn
        game.take_action(simon, Action(Action.Actions.check))
        game.take_action(chien, Action(Action.Actions.check))
        game.take_action(jason, Action(Action.Actions.check))
        assert self.expect_fail(lambda: game.take_action(simon, Action(Action.Actions.bet, 1000)))  # Can't bet >stack
        game.take_action(simon, Action(Action.Actions.bet, 10))
        game.take_action(chien, Action(Action.Actions.call))
        game.take_action(jason, Action(Action.Actions.fold))
        game.take_action(simon, Action(Action.Actions.check))
        game.take_action(chien, Action(Action.Actions.check))
        assert not game.is_hand_active()

        grace = Player("Grace")
        grace.stack = 100
        game.sit_player(grace, 8)
        snapshot.assert_match(str(game.__dict__))

        game.deal_hand()
        snapshot.assert_match(str(game.game_state.__dict__))
        assert self.expect_fail(lambda: game.take_action(chien, Action(Action.Actions.call)))  # Can't play out of turn
        assert self.expect_fail(lambda: game.take_action(simon, Action(Action.Actions.call)))  # Can't play out of turn
        assert self.expect_fail(lambda: game.take_action(grace, Action(Action.Actions.call)))  # Can't play out of turn
        game.take_action(jason, Action(Action.Actions.call))
        game.take_action(simon, Action(Action.Actions.bet, 10))
        game.take_action(chien, Action(Action.Actions.fold))
        game.take_action(grace, Action(Action.Actions.call))
        game.take_action(jason, Action(Action.Actions.call))
        snapshot.assert_match(str(game.game_state.__dict__))
        assert self.expect_fail(lambda: game.take_action(jason, Action(Action.Actions.call)))  # Can't play out of turn
        assert self.expect_fail(lambda: game.take_action(simon, Action(Action.Actions.call)))  # Can't play out of turn
        game.take_action(grace, Action(Action.Actions.check))
        game.take_action(jason, Action(Action.Actions.check))
        game.take_action(simon, Action(Action.Actions.bet, 10))
        game.take_action(grace, Action(Action.Actions.fold))
        snapshot.assert_match(str(game.get_player_state(grace).__dict__))
        game.take_action(jason, Action(Action.Actions.call))
        game.take_action(jason, Action(Action.Actions.check))
        game.take_action(simon, Action(Action.Actions.check))
        game.take_action(jason, Action(Action.Actions.check))
        game.take_action(simon, Action(Action.Actions.check))
        snapshot.assert_match(str(game.get_player_state(simon).__dict__))
        assert not game.is_hand_active()
        snapshot.assert_match(str(game.game_state.__dict__))
        snapshot.assert_match(str(game.__dict__))
