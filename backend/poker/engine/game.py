from dataclasses import dataclass
from random import shuffle
from typing import List, Optional

from .card import Card
from .player import Player
from .state import GameState, EndGameState
from .action import Action


@dataclass
class PlayerState:
    bet_round: int
    lead_player: Player
    current_players: List[Player]
    player_cards: List[Card]
    community_cards: List[Card]
    acting_player: Player
    end_game: EndGameState


class Game:
    seats: List[Optional[Player]]
    button_position: int
    small_blind: int
    big_blind: int
    game_state: GameState

    def __init__(self, small_blind_amount: int, big_blind_amount: int) -> None:
        self.seats = [None] * 9
        self.button_position = 0
        self.small_blind = small_blind_amount
        self.big_blind = big_blind_amount
        self.game_state = None

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
        if self.game_state.is_hand_over():
            raise Exception("Trying to take action when hand is over")
        if player.seat_number != self.game_state.get_acting_index():
            raise Exception("Player ({}) is playing out of turn".format(self.seats[player.seat_number]))
        self.game_state = self.game_state.take_action(action)
        if self.game_state.is_hand_over():
            self.game_state.process_end_game()

    def get_player_state(self, player: Player) -> PlayerState:
        if self.game_state:
            return PlayerState(self.game_state.round, self.game_state.get_leading_player(),
                               self.game_state.players, self.game_state.get_player_cards(player.seat_number),
                               self.game_state.get_community_cards(), self.game_state.get_acting_player(),
                               self.game_state.get_end_game_state())

    def deal_hand(self) -> None:

        self.button_position = (self.button_position + 1) % 9
        while not self.seats[self.button_position]:
            self.button_position = (self.button_position + 1) % 9

        deck = [Card(i) for i in range(52)]
        deck = self.shuffle_deck(deck)

        self.game_state = GameState()
        self.game_state.start_hand(self.seats, self.button_position,  self.big_blind, self.small_blind, deck)

    def is_hand_active(self) -> bool:
        if not self.game_state:
            return False
        return not self.game_state.is_hand_over()

    @staticmethod
    def shuffle_deck(deck: List[Card]) -> List[Card]:
        deck_copy = list(deck)
        shuffle(deck_copy)      # TODO: replace this with more e n t r o p y
        return deck_copy
