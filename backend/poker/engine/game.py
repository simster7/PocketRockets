from dataclasses import dataclass
from random import shuffle
from typing import List, Optional

from .card import Card
from .player import Player
from .state import GameState, Action, EndGameState


@dataclass
class PlayerState:
    bet_round: int
    lead_action: Action
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
        elif player.get_seat_number() is not None:
            raise Exception("Player is already sitting")

        player.set_seat_number(seat_number)
        self.seats[seat_number] = player

    def stand_player(self, player: Player, seat_number: int) -> None:
        if not 0 <= seat_number < 9:
            raise Exception("Invalid seat number")
        elif not self.seats[seat_number]:
            raise Exception("Seat is already empty")
        elif self.seats[seat_number].get_seat_number() != seat_number:
            raise Exception("Incorrect player/seat number combination")

        player.set_seat_number(None)
        self.seats[seat_number] = None

    def take_action(self, player: Player, action: Action) -> None:
        if self.game_state.is_hand_over():
            raise Exception("Trying to take action when hand is over")
        if player.get_seat_number() != self.get_acting_seat():
            raise Exception("Player ({}) is playing out of turn".format(self.seats[player.get_seat_number()]))
        self.game_state = self.game_state.take_action(action)
        if self.game_state.is_hand_over():
            self.game_state.process_end_game()

    def get_player_state(self, player: Player) -> PlayerState:
        if self.game_state:
            state_seat_number = player.get_seat_number() - sum([1 if not seat else 0 for seat in
                                                                self.seats[:player.get_seat_number()]])
            return PlayerState(self.game_state.get_round(), self.game_state.get_lead_action(),
                               self.game_state.get_leading_player(), self.game_state.get_players(),
                               self.game_state.get_player_cards(state_seat_number),
                               self.game_state.get_community_cards(), self.game_state.get_acting_player(),
                               self.game_state.get_end_game_state())

    def deal_hand(self) -> None:
        players = [player for player in self.seats if player is not None]
        self.button_position = (self.button_position + 1) % len(players)
        players = players[self.button_position:] + players[:self.button_position]
        deck = [Card(i) for i in range(52)]
        deck = self.shuffle_deck(deck)

        self.game_state = GameState()
        self.game_state.start_hand(players, self.big_blind, self.small_blind, deck)

    def is_hand_active(self) -> bool:
        if not self.game_state:
            return False
        return not self.game_state.is_hand_over()

    def get_acting_seat(self) -> int:
        state_index = (self.game_state.get_acting_index() + self.button_position)\
                      % sum([1 if seat else 0 for seat in self.seats])
        count = -1
        index = -1
        while count != state_index:
            index += 1
            while not self.seats[index]:
                index += 1
            count += 1
        return index

    @staticmethod
    def shuffle_deck(deck: List[Card]) -> List[Card]:
        deck_copy = list(deck)
        shuffle(deck_copy)      # TODO: replace this with more e n t r o p y
        return deck_copy

