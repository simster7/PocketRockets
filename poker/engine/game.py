from .player import Player
from .card import Card
from .evaluator import rank_hands
from .state import GameState, Action

from random import shuffle

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
    """.format("Call {} and raise [Number]".format(lead_action.value) if lead_action.action == Action.actions.bet else "Bet [Number]"))
    while True:
        action = input("Enter action: ")
        if action == "F":
            return game_state.take_action(Action(Action.actions.fold))
        elif action == "C":
            return game_state.take_action(Action(Action.actions.check))
        elif action == "L":
            return game_state.take_action(Action(Action.actions.call))
        elif action.replace('.','',1).isdigit():
            return game_state.take_action(Action(Action.actions.bet, float(action)))

class Game:

    def __init__(self, small_blind_amount, big_blind_amount, prompt = prompt):
        self.seats = [None] * 9
        self.button_position = 0
        self.small_blind = small_blind_amount
        self.big_blind = big_blind_amount
        self.prompt = prompt

    def sit_player(self, player, seat_number):
        assert seat_number >= 0 and seat_number < 9, "Invalid seat number"
        if self.seats[seat_number] == None:
            self.seats[seat_number] = player
        else:
            raise Exception("Cannot sit player ({}) on an occupied seat".format(player))

    def deal_hand(self):
        players = [player for player in self.seats if player != None]
        deck = [Card(i) for i in range(52)]
        shuffle(deck) # TODO: replace this with more e n t r o p y

        game_state = GameState()
        game_state.start_hand(players, self.big_blind, self.small_blind, deck)

        while not game_state.is_hand_over():
            game_state = self.prompt(game_state)





if __name__ == '__main__':
    game = Game(1, 2)
    game.sit_player(Player("Simon", 100), 0)
    game.sit_player(Player("Hersh", 100), 1)
    game.sit_player(Player("Chien", 100), 2)
    game.sit_player(Player("Jarry", 100), 3)
    game.sit_player(Player("Grace", 100), 6)
    game.sit_player(Player("Jason", 100), 7)
    while True:
        game.deal_hand()
