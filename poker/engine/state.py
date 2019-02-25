from __future__ import annotations

from dataclasses import dataclass
from enum import Enum
from typing import List, Optional

from .card import Card
from .player import Player
from .evaluator import rank_hands, RankedHand


@dataclass
class EndGameState:
    winner: Player
    condition: str
    hands: List[RankedHand]


class Action:
    class Actions(Enum):
        check = 0
        call = 1
        fold = 2
        bet = 3

    def __init__(self, action, value=None):
        self.action = action
        if self.action == Action.Actions.bet and value is None:
            raise Exception("A bet value is required when betting")
        self.value = value

    def __str__(self):
        return self.action.name if self.action != Action.Actions.bet else self.action.name + " " + str(self.value)


class GameState:

    players: List[Player]
    bet_vector: List[int]
    fold_vector: List[bool]
    pot: int
    acting_player: int
    leading_player: int
    deck: List[Card]
    round: int

    def __init__(self, old=None):
        if old:
            self.copy(old)

    def start_hand(self, players: List[Player], big_blind: int, small_blind: int, deck: List[Card]):
        self.players = players
        self.bet_vector = [small_blind, big_blind] + [0] * (len(self.players) - 2)
        self.players[0].make_bet(small_blind)
        self.players[1].make_bet(big_blind)
        self.fold_vector = [False] * len(self.players)
        self.pot = 0
        self.acting_player = 2 % len(self.players)
        self.leading_player = 1
        self.deck = deck
        self.round = 0

    def copy(self, old_game_state: GameState):
        self.players = old_game_state.players
        self.bet_vector = old_game_state.bet_vector
        self.fold_vector = old_game_state.fold_vector
        self.pot = old_game_state.pot
        self.acting_player = old_game_state.acting_player
        self.leading_player = old_game_state.leading_player
        self.deck = old_game_state.deck
        self.round = old_game_state.round

    def is_hand_over(self) -> bool:
        return (sum(self.fold_vector) == len(self.players) - 1) or (self.is_round_over() and self.round == 4)

    def is_round_over(self) -> bool:
        # TODO hard-code option for big blind when fold-around
        return self.acting_player == self.leading_player

    def get_round(self) -> int:
        return self.round

    def get_players(self) -> List[Player]:
        return self.players

    def get_player_cards(self, player_index: int) -> Optional[List[Card]]:
        if self.fold_vector[player_index]:
            return None
        return [self.deck[player_index], self.deck[player_index + len(self.players)]]

    def get_community_cards(self) -> List[Card]:
        cards = []
        curr_round = self.get_round()
        num_players = len(self.players)
        if curr_round >= 1:
            cards.extend(self.deck[2 * num_players + 1: 2 * num_players + 4])
        if curr_round >= 2:
            cards.extend([self.deck[2 * num_players + 5]])
        if curr_round >= 3:
            cards.extend([self.deck[2 * num_players + 7]])
        return cards

    def get_acting_index(self) -> int:
        return self.acting_player

    def get_acting_player(self) -> Player:
        return self.players[self.acting_player]

    def get_leading_player(self) -> Player:
        return self.players[self.leading_player]

    def get_pot(self) -> int:
        return self.pot

    def get_lead_action(self) -> Action:
        if self.bet_vector[self.leading_player] == 0:
            return Action(Action.Actions.check)
        else:
            return Action(Action.Actions.bet, self.bet_vector[self.leading_player])

    def get_end_game(self) -> Optional[EndGameState]:
        if not self.is_hand_over():
            return None
        showdown: List[int] = [i for i in range(len(self.players)) if not self.fold_vector[i]]
        if sum(self.fold_vector) == len(self.players) - 1:
            assert len(showdown) == 1, "Bug: more than one player left on a fold win condition"
            winner = self.players[showdown[0]]
            winner.receive_pot(self.pot)
            return EndGameState(winner, "folds", [])
        else:
            showdown_hands = [(player_index, self.get_player_cards(player_index) + self.get_community_cards()) for
                              player_index in showdown]
            ranked_hands = rank_hands(showdown_hands)
            winner = self.players[ranked_hands[0].player_index]
            winner.receive_pot(self.pot)
            return EndGameState(winner, "showdown", ranked_hands)

    def move_acting_player(self) -> None:
        self.acting_player = (self.acting_player + 1) % len(self.players)
        while self.fold_vector[self.acting_player] and not self.is_round_over():
            self.acting_player = (self.acting_player + 1) % len(self.players)
        if self.is_round_over():
            self.acting_player = 0
            self.leading_player = 0
            while self.fold_vector[self.acting_player]:
                self.acting_player = (self.acting_player + 1) % len(self.players)
                self.leading_player = self.acting_player
            self.pot += sum(self.bet_vector)
            self.bet_vector = [0] * len(self.players)
            self.round += 1

    def take_action(self, action: Action) -> GameState:
        if action.action == Action.Actions.fold:
            self.fold_vector[self.acting_player] = True
            self.move_acting_player()
            return GameState(self)
        if action.action == Action.Actions.check:
            if self.get_lead_action().action == Action.Actions.bet:
                raise Exception("Illegal game state: player can\'t check when there is a bet")
            self.move_acting_player()
            return GameState(self)
        if action.action == Action.Actions.call:
            to_call = self.bet_vector[self.leading_player] - self.bet_vector[self.acting_player]
            if self.players[self.acting_player].stack < to_call:
                raise Exception("Illegal game state: player doesn\'t have enough chips to call")
            self.players[self.acting_player].make_bet(to_call)
            self.bet_vector[self.acting_player] += to_call
            self.move_acting_player()
            return GameState(self)
        if action.action == Action.Actions.bet:
            lead_action = self.get_lead_action()
            to_call = 0
            if lead_action.action == Action.Actions.bet:
                to_call = max(lead_action.value - self.bet_vector[self.acting_player], 0)
            if self.players[self.acting_player].stack < action.value:
                raise Exception("Illegal game state: player doesn\'t have enough chips to make that bet")
            self.players[self.acting_player].make_bet(action.value + to_call)
            self.bet_vector[self.acting_player] += action.value + to_call
            self.leading_player = self.acting_player
            self.move_acting_player()
            return GameState(self)
