from random import shuffle
from typing import List, Dict

from .card import Card
from .player import Player
from .state import GameState, Action


class Game:

    def __init__(self, small_blind_amount: int, big_blind_amount: int) -> None:
        self.seats: List[Player] = [None] * 9
        self.button_position: int = 0
        self.small_blind: int = small_blind_amount
        self.big_blind: int = big_blind_amount
        self.game_state: GameState = None

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
        if seat_number != self.get_acting_seat():
            raise Exception("Player ({}) is playing out of turn".format(self.seats[seat_number]))
        self.game_state = self.game_state.take_action(action)

    def get_player_state(self, seat_number: int) -> Dict:
        print("Requesting player state for", seat_number)
        if self.game_state:
            state_seat_number = seat_number - sum([1 if not seat else 0 for seat in self.seats[:seat_number]])
            return {"bet_round": self.game_state.get_round(),
                    "lead_action": self.game_state.get_lead_action(),
                    "lead_player": self.game_state.get_leading_player(),
                    "current_players": self.game_state.get_players(),
                    "player_cards": self.game_state.get_player_cards(state_seat_number),
                    "community_cards": self.game_state.get_community_cards(),
                    "acting_player": self.game_state.get_acting_player()
                    }

    def get_acting_seat(self) -> int:
        state_index = self.game_state.get_acting_index()
        count = 0
        index = 0
        while count != state_index:
            while not self.seats[index]:
                index += 1
            count += 1
        return count


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
            lead_action.value) if lead_action.action == Action.actions.bet else "Bet [Number]"))
        while True:
            action = input("Enter action: ")
            if action == "F":
                return Action(Action.actions.fold)
            elif action == "C":
                return Action(Action.actions.check)
            elif action == "L":
                return Action(Action.actions.call)
            elif action.replace('.', '', 1).isdigit():
                return Action(Action.actions.bet, float(action))


    game = Game(1, 2)
    game.sit_player(Player("Simon", 100), 0)
    game.sit_player(Player("Hersh", 100), 1)
    game.sit_player(Player("Chien", 100), 2)
    game.sit_player(Player("Jarry", 100), 3)
    game.sit_player(Player("Grace", 100), 6)
    game.sit_player(Player("Jason", 100), 7)
    while True:
        game.deal_hand()
        while not game.game_state.is_hand_over():
            game.take_action(prompt(game.get_state()))
