from player import Player
from card import Card
from enum import Enum

from random import shuffle

class GameState:
    def __init__(self, players, big_blind, small_blind, deck):
        self.players = players
        self.bet_vector = [small_blind, big_blind] + [0] * (len(players) - 2)
        self.pot = 0
        self.acting_player = 2 % len(plyaers)
        self.leading_player = 1
        self.deck = []
        self.round = 0

    def is_round_over(self):
        return self.acting_player == self.leading_player

    def get_player_cards(self, player_index):
        return [deck[i], deck[i + len(players)]]

    def next(self, action):



def prompt(player, bet_round, lead_action, lead_player):
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

        game_state = GameState(players, big_blind_amount, small_blind_amount, deck)

        for bet_round in ["PRE", "FLOP", "TURN", "RIVER"]:
            while not game_state.is_round_over():
                if not round_bets[acting_player].requires_action(round_bets[lead_player]):
                    acting_player = (acting_player + 1) % num_players
                    continue

                round_bets[acting_player], takes_lead = prompt(players[acting_player],
                                                               bet_round,
                                                               round_bets[lead_player],
                                                               players[lead_player])
                if takes_lead:
                    lead_player = acting_player

                acting_player = (acting_player + 1) % num_players





game = Game(1, 2)
game.sit_player(Player("simon", 20), 0)
game.sit_player(Player("hersh", 20), 1)
game.sit_player(Player("chien", 20), 2)
game.sit_player(Player("jarry", 20), 3)
game.deal_hand()
print(game.seats)
