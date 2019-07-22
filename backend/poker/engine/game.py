from random import shuffle
from typing import List, Optional

from backend.api.data_classes import PlayerState
from .action import Action
from .card import Card
from .player import Player
from .state import GameState


class Game:
    seats: List[Optional[Player]]
    button_position: int
    small_blind: int
    big_blind: int
    game_state: GameState
    active_hand: bool

    def __init__(self, small_blind_amount: int, big_blind_amount: int) -> None:
        self.seats = [None] * 9
        self.button_position = 0
        self.small_blind = small_blind_amount
        self.big_blind = big_blind_amount
        self.game_state = GameState()
        self.active_hand = False

    def sit_player(self, player: Player, seat_number: int) -> None:
        if not 0 <= seat_number < 9:
            raise Exception("Invalid seat number")
        elif self.seats[seat_number]:
            raise Exception("Cannot sit player ({}) on an occupied seat".format(player))
        elif player.seat_number is not None:
            raise Exception("Player is already sitting")

        player.seat_number = seat_number
        self.seats[seat_number] = player

    def stand_player(self, player: Player, seat_number: int) -> None:
        if not 0 <= seat_number < 9:
            raise Exception("Invalid seat number")
        elif not self.seats[seat_number]:
            raise Exception("Seat is already empty")
        elif self.seats[seat_number].seat_number != seat_number:
            raise Exception("Incorrect player/seat number combination")

        player.seat_number = None
        self.seats[seat_number] = None

    def take_action(self, player: Player, action: Action) -> None:
        if not self.is_hand_active():
            raise Exception("Trying to take action when hand is over")
        if player.seat_number != self.game_state.get_acting_index():
            raise Exception("Player ({}) is playing out of turn".format(self.seats[player.seat_number]))

        proposed_last_action = action
        proposed_state, action_response = self.game_state.take_action(action)

        player.folded = action_response.folded
        player.make_bet(action_response.bet)
        player.last_action = proposed_last_action

        if action_response.end_game:
            self.active_hand = False
            for winner, win_amount in action_response.end_game.winners:
                winner.receive_pot(win_amount)

    def get_player_state(self, player: Player) -> PlayerState:
        # TODO: Refactor this, this is horrible
        return self.game_state.get_player_state(player)

    def deal_hand(self) -> None:
        self.button_position = (self.button_position + 1) % 9
        while not self.seats[self.button_position]:
            self.button_position = (self.button_position + 1) % 9
        deck = [Card(i) for i in range(52)]
        deck = self.shuffle_deck(deck)

        action_responses = self.game_state.start_hand(self.seats, self.button_position, self.big_blind,
                                                      self.small_blind, deck)

        for action in action_responses:
            self.seats[action.player_index].make_bet(action.bet)

        self.active_hand = True

    def is_hand_active(self) -> bool:
        return self.active_hand

    def get_acting_seat(self) -> int:
        return self.game_state.acting_player

    @staticmethod
    def shuffle_deck(deck: List[Card]) -> List[Card]:
        deck_copy = list(deck)
        shuffle(deck_copy)  # TODO: replace this with more e n t r o p y
        return deck_copy
