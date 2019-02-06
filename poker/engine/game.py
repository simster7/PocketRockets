from player import Player
from card import Card
from enum import Enum

from random import shuffle

actions = Enum('Actions', 'check raise call fold bet')

class Action:
    def __init__(self, action, value = None):
        self.action = action 
        if self.action == actions.bet and value == None:
            raise Exception("A bet value is requiered when betting")
        self.value = value
    def __str__(self):
        return self.action.name if self.action != actions.bet else self.action.name + " " + str(self.value)

class GameState:

    def __init__(self, old = None):
        if old:
            self.copy(old)

    def start_hand(self, players, big_blind, small_blind, deck):
        self.players = players
        self.bet_vector = [small_blind, big_blind] + [0] * (len(self.players) - 2)
        self.fold_vector = [False] * len(self.players)
        self.pot = 0
        self.acting_player = 2 % len(self.players)
        self.leading_player = 1
        self.deck = []
        self.round = 0

    def copy(self, old_game_state):
        self.players = old_game_state.players 
        self.bet_vector = old_game_state.bet_vector 
        self.fold_vector = old_game_state.fold_vector 
        self.pot = old_game_state.pot 
        self.acting_player = old_game_state.acting_player 
        self.leading_player = old_game_state.leading_player 
        self.deck = old_game_state.deck 
        self.round = old_game_state.round 

    def is_hand_over(self):
        return (sum(self.fold_vector) == len(self.players) - 1) or (self.is_round_over() and self.round == 4)

    def is_round_over(self):
        return self.acting_player == self.leading_player

    def get_round(self):
        return self.round

    def get_player_cards(self, player_index):
        return [deck[i], deck[i + len(players)]]

    def get_acting_player(self):
        return self.players[self.acting_player]

    def get_leading_player(self):
        return self.players[self.leading_player]

    def get_lead_action(self):
        if self.bet_vector[self.leading_player] == 0:
            return Action(actions.check)
        else:
            return Action(actions.bet, self.bet_vector[self.leading_player])


    def take_action(self, action, action_param = None):
        if action.action == actions.check:
            if self.get_lead_action().action == actions.bet:
                raise Exception('Illegal game state')
            self.acting_player = (self.acting_player + 1) % len(self.players)
            return GameState(self)



def prompt(game_state):
    player = game_state.get_acting_player()
    bet_round = game_state.
    lead_action = game_state.get_lead_action()
    lead_player = game_state.get_leading_player()
    print("Current round: {}".format(bet_round))
    print("Lead action: {}: {}".format(lead_player.name, lead_action))
    print("Acting as player {}".format(player.name))
    print("""
        F - Fold
        C - Check
        [Number] - Bet [Number] (currently, use this to call, bet, and raise)
    """)
    while True:
        action = input("Enter action: ")
        if action == "F":
            return Action(ACTIONS.fold), False
        elif action == "C":
            return Action(ACTIONS.check), True
        elif action.replace('.','',1).isdigit():
            return Action(ACTIONS.bet, float(action)), True

class Game:

    def __init__(self, small_blind_amount, big_blind_amount):
        self.seats = [None] * 9
        self.button_position = 0
        self.small_blind = small_blind_amount
        self.big_blind = big_blind_amount

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

        for bet_round in ["PRE", "FLOP", "TURN", "RIVER"]:
            while not game_state.is_round_over():
                prompt(game_state)





game = Game(1, 2)
game.sit_player(Player("simon", 20), 0)
game.sit_player(Player("hersh", 20), 1)
game.sit_player(Player("chien", 20), 2)
game.sit_player(Player("jarry", 20), 3)
game.deal_hand()
print(actions.check == actions.check)
print(game.seats)
