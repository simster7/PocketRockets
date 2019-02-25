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
        elif player.get_seat_number() != seat_number:
            raise Exception("Player seat number out of sync")

        self.seats[seat_number] = player

    def stand_player(self, seat_number: int) -> None:
        self.seats[seat_number] = None

    def deal_hand(self) -> None:
        players = [player for player in self.seats if player is not None]
        deck = [Card(i) for i in range(52)]
        shuffle(deck)  # TODO: replace this with more e n t r o p y

        self.game_state = GameState()
        self.game_state.start_hand(players, self.big_blind, self.small_blind, deck)

    def take_action(self, seat_number: int, action: Action) -> None:
        if self.game_state.is_hand_over():
            raise Exception("Trying to take action when hand is over")
        if seat_number != self.get_acting_seat():
            raise Exception("Player ({}) is playing out of turn".format(self.seats[seat_number]))
        self.game_state = self.game_state.take_action(action)

    def is_hand_active(self) -> bool:
        if not self.game_state:
            return False
        return not self.game_state.is_hand_over()

    def get_player_state(self, seat_number: int) -> PlayerState:
        if self.game_state:
            state_seat_number = seat_number - sum([1 if not seat else 0 for seat in self.seats[:seat_number]])
            return PlayerState(self.game_state.get_round(), self.game_state.get_lead_action(),
                               self.game_state.get_leading_player(), self.game_state.get_players(),
                               self.game_state.get_player_cards(state_seat_number),
                               self.game_state.get_community_cards(), self.game_state.get_acting_player(),
                               self.game_state.get_end_game())

    # TESTING ONLY
    def get_state_testing(self):
        return self.game_state

    # TESTING ONLY
    def take_action_testing(self, action: Action) -> None:
        self.game_state = self.game_state.take_action(action)

    def get_acting_seat(self) -> int:
        state_index = self.game_state.get_acting_index()
        count = -1
        index = -1
        while count != state_index:
            index += 1
            while not self.seats[index]:
                index += 1
            count += 1
        return index


if __name__ == '__main__':

    def prompt(game_state):
        player = game_state.get_acting_player()
        bet_round = game_state.get_round()
        lead_action = game_state.get_lead_action()
        lead_player = game_state.get_leading_player()
        print()
        print("Current players: {}".format(game_state.get_players()))
        print("Player hands: {}".format([game_state.get_player_cards(i) for i in range(len(game_state.players))]))
        print()
        print("Community cards: {}".format(game_state.get_community_cards()))
        print("Pot: {}".format(game_state.get_pot()))
        print()
        print("Current round: {}".format(bet_round))
        print("Lead action: {}: {}".format(lead_player.name, lead_action))
        print("Acting as player: {}".format(player.name))
        print("With hand: {}".format(game_state.get_player_cards(game_state.get_acting_index())))
        print("""
            F - Fold
            C - Check
            L - Call
            [Number] - {} 
        """.format("Call {} and raise [Number]".format(
            lead_action.value) if lead_action.action == Action.Actions.bet else "Bet [Number]"))
        while True:
            action = input("Enter action: ")
            if action == "F":
                return Action(Action.Actions.fold)
            elif action == "C":
                return Action(Action.Actions.check)
            elif action == "L":
                return Action(Action.Actions.call)
            elif action.replace('.', '', 1).isdigit():
                return Action(Action.Actions.bet, float(action))


    game = Game(1, 2)
    simon = Player("Simon")
    simon.set_stack(100)
    simon.set_seat_number(0)
    game.sit_player(simon, 0)
    hersh = Player("Hersh")
    hersh.set_stack(100)
    hersh.set_seat_number(1)
    game.sit_player(hersh, 1)
    chien = Player("Chien")
    chien.set_stack(100)
    chien.set_seat_number(2)
    game.sit_player(chien, 2)
    jarry = Player("Jarry")
    jarry.set_stack(100)
    jarry.set_seat_number(3)
    game.sit_player(jarry, 3)
    grace = Player("Grace")
    grace.set_stack(100)
    grace.set_seat_number(6)
    game.sit_player(grace, 6)
    jason = Player("Jason")
    jason.set_stack(100)
    jason.set_seat_number(7)
    game.sit_player(jason, 7)
    while True:
        game.deal_hand()
        while not game.game_state.is_hand_over():
            game.take_action_testing(prompt(game.get_state_testing()))
