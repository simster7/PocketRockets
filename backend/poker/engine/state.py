from __future__ import annotations

from dataclasses import dataclass
from typing import List, Optional, Tuple

from backend.api.data_classes import EndGameState, PlayerState
from .action import Action
from .card import Card
from .evaluator import rank_hands
from .player import Player


@dataclass
class ActionResponse:
    player_index: int
    folded: bool
    bet: int
    end_game: Optional[EndGameState]


class GameState:
    players: List[Player]
    button_position: int
    fold_vector: List[bool]
    bet_vector: List[int]
    pot: int
    deck: List[Card]
    round: int
    acting_player: int
    leading_player: int
    is_hand_active: bool

    def __init__(self, old=None):
        if old:
            self.__copy(old)
        else:
            self.round = -1
            self.leading_player = -1
            self.acting_player = -1
            self.players = []
            self.fold_vector = []
            self.is_hand_active = False
            self.pot = 0

    def start_hand(self, players: List[Player], button_position: int, big_blind: int, small_blind: int,
                   deck: List[Card]) -> Tuple[ActionResponse, ...]:

        # TODO: Add support for straddles
        self.players = players
        self.button_position = button_position
        self.fold_vector = [not player or player.sitting_out for player in self.players]
        self.bet_vector = [0] * 9
        self.pot = 0
        self.deck = deck
        self.round = 0

        small_blind_index = self.__get_n_active_player_from(button_position, 1)
        big_blind_index = self.__get_n_active_player_from(button_position, 2)
        utg_blind_index = self.__get_n_active_player_from(button_position, 3)
        self.bet_vector[big_blind_index] = big_blind
        self.bet_vector[small_blind_index] = small_blind
        self.players[small_blind_index].last_action = Action(Action.Actions.bet, small_blind)
        self.players[big_blind_index].last_action = Action(Action.Actions.bet, big_blind)
        
        self.acting_player = utg_blind_index
        self.leading_player = big_blind_index

        self.is_hand_active = True

        return ActionResponse(small_blind_index, False, small_blind, None), ActionResponse(big_blind_index, False, big_blind, None)

    def take_action(self, action: Action) -> Tuple[GameState, ActionResponse]:
        if action.action == Action.Actions.fold:
            self.fold_vector[self.acting_player] = True
            end_game_state = self.__move_acting_player()
            return GameState(self), ActionResponse(self.acting_player, True, 0, end_game_state)
        if action.action == Action.Actions.check:
            if self.__get_lead_action().action == Action.Actions.bet:
                raise Exception("Illegal game state: player can\'t check when there is a bet")
            end_game_state = self.__move_acting_player()
            return GameState(self), ActionResponse(self.acting_player, False, 0, end_game_state)
        if action.action == Action.Actions.call:
            to_call = self.bet_vector[self.leading_player] - self.bet_vector[self.acting_player]
            if self.players[self.acting_player].stack < to_call:
                raise Exception("Illegal game state: player doesn\'t have enough chips to call")
            self.bet_vector[self.acting_player] += to_call
            end_game_state = self.__move_acting_player()
            return GameState(self), ActionResponse(self.acting_player, False, to_call, end_game_state)
        if action.action == Action.Actions.bet:
            lead_action = self.__get_lead_action()
            to_call = 0
            if lead_action.action == Action.Actions.bet:
                to_call = max(lead_action.value - self.bet_vector[self.acting_player], 0)
            if self.players[self.acting_player].stack < action.value:
                raise Exception("Illegal game state: player doesn\'t have enough chips to make that bet")
            self.bet_vector[self.acting_player] += action.value + to_call
            self.leading_player = self.acting_player
            end_game_state = self.__move_acting_player()
            return GameState(self), ActionResponse(self.acting_player, False, action.value + to_call, end_game_state)

    def get_player_state(self, player: Player) -> PlayerState:
        # TODO: Refactor this, this is horrible
        return PlayerState(self.round, self.__get_leading_player(),
                           self.__get_acting_player(), self.players,
                           self.__get_player_cards(
                               player.seat_number) if player and player.seat_number else [],
                           self.__get_community_cards(),
                           self.__get_end_game_state() if self.is_hand_active else None,
                           player.seat_number if player and player.seat_number else None,
                           self.button_position, self.pot)

    def get_acting_index(self) -> int:
        return self.acting_player

    def __copy(self, old_game_state: GameState):
        self.players = old_game_state.players
        self.button_position = old_game_state.button_position
        self.fold_vector = old_game_state.fold_vector
        self.bet_vector = old_game_state.bet_vector
        self.pot = old_game_state.pot
        self.deck = old_game_state.deck
        self.round = old_game_state.round
        self.acting_player = old_game_state.acting_player
        self.leading_player = old_game_state.leading_player
        self.is_hand_active = old_game_state.is_hand_active

    def __move_acting_player(self) -> Optional[EndGameState]:
        self.acting_player = (self.acting_player + 1) % len(self.players)
        while self.fold_vector[self.acting_player] and not self.__is_round_over():
            self.acting_player = (self.acting_player + 1) % len(self.players)

        if self.__is_round_over():
            self.acting_player = self.__get_n_active_player_from(self.button_position, 1)
            self.leading_player = self.__get_n_active_player_from(self.button_position, 1)
            self.pot += sum(self.bet_vector)
            self.bet_vector = [0] * 9
            self.round += 1

        if self.__is_hand_over():
            self.is_hand_active = False
            return self.__get_end_game_state()
        return None

    @property
    def __num_players_in_hand(self) -> int:
        return sum([1 if player else 0 for player in self.players])

    def __get_n_active_player_from(self, base: int, n: int) -> int:
        index = base
        count = 0
        while count != n:
            index = (index + 1) % 9
            while not self.players[index] or self.fold_vector[index]:
                index = (index + 1) % 9
            count += 1
        return index

    def __is_hand_over(self) -> bool:
        return (sum([0 if fold else 1 for fold in self.fold_vector]) == 1) or (self.__is_round_over() and self.round == 4)

    def __is_round_over(self) -> bool:
        # TODO hard-code option for big blind when fold-around
        return self.acting_player == self.leading_player

    def __get_player_cards(self, player_index: int) -> Optional[List[Card]]:
        if not self.is_hand_active or self.fold_vector[player_index]:
            return []
        return [self.deck[player_index], self.deck[player_index + self.__num_players_in_hand]]

    def __get_community_cards(self) -> List[Card]:
        if not self.is_hand_active:
            return []
        cards = []
        curr_round = self.round
        num_players = self.__num_players_in_hand
        if curr_round >= 1:
            cards.extend(self.deck[2 * num_players + 1: 2 * num_players + 4])
        if curr_round >= 2:
            cards.extend([self.deck[2 * num_players + 5]])
        if curr_round >= 3:
            cards.extend([self.deck[2 * num_players + 7]])
        return cards

    def __get_acting_player(self) -> Optional[Player]:
        if not self.is_hand_active or self.acting_player < 0 or self.acting_player > 9:
            return None
        return self.players[self.acting_player]

    def __get_leading_player(self) -> Optional[Player]:
        if not self.is_hand_active or self.leading_player < 0 or self.leading_player > 9:
            return None
        return self.players[self.leading_player]

    def __get_pot(self) -> int:
        return self.pot

    def __get_lead_action(self) -> Optional[Action]:
        if not self.is_hand_active:
            return None
        if self.bet_vector[self.leading_player] == 0:
            return Action(Action.Actions.check)
        else:
            return Action(Action.Actions.bet, self.bet_vector[self.leading_player])

    def __get_end_game_state(self) -> Optional[EndGameState]:
        if not self.__is_hand_over():
            return None
        showdown: List[int] = [i for i in range(len(self.players)) if not self.fold_vector[i]]
        if sum(self.fold_vector) == len(self.players) - 1:
            assert len(showdown) == 1, "Bug: more than one player left on a fold win condition"
            winner = self.players[showdown[0]]
            return EndGameState([(winner, self.pot)], "folds", [])
        else:
            showdown_hands = [(player_index, self.__get_player_cards(player_index) + self.__get_community_cards()) for
                              player_index in showdown]
            ranked_hands = rank_hands(showdown_hands)
            winner = self.players[ranked_hands[0].player_index]
            return EndGameState([(winner, self.pot)], "showdown", ranked_hands)

    def __repr__(self):
        return str(self.__dict__)
