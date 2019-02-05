from player import Player
from card import Card
from enum import Enum

from random import shuffle

ACTIONS = Enum('Actions', 'none check fold bet')

class Action:

    def __init__(self, action, value = None):
        self.action = action 

        if self.action == ACTIONS.bet and value == None:
            raise Exception("A bet value is requiered when betting")

        self.value = value

    def requires_action(self, lead_action):
        if self.action == ACTIONS.none:
            return True

        if self.action == ACTIONS.fold:
            return False

        if lead_action.action == ACTIONS.check:
            if self.action == ACTIONS.check:
                return False

        if lead_action.action == ACTIONS.bet:
            if self.action != ACTIONS.bet:
                return True
            if lead_action.value != self.value:
                return True
            else:
                return False

    def __eq__(self, other):
        if self.action != other.action:
            return False
        else:
            if self.action == ACTIONS.bet:
                return self.value == other.value
            return True
    
    def __str__(self):
        return self.action.name if self.action != ACTIONS.bet else self.action.name + " " + str(self.value)


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
        # TODO: shift players for button
        num_players = len(players)
        deck = [Card(i) for i in range(52)]
        shuffle(deck) # TODO: replace this with more e n t r o p y

        # Deal cards
        for i, player in enumerate(players):
            player_hole_cards = [deck[i], deck[i + num_players]]
            player.recieve_cards(player_hole_cards)

        for bet_round in ["PRE", "FLOP", "TURN", "RIVER"]:
            if bet_round == "PRE":
                round_bets = [Action(ACTIONS.bet, self.small_blind), Action(ACTIONS.bet, self.big_blind)] + ([Action(ACTIONS.none)] * (num_players - 2))
                acting_player = 2 % num_players
                lead_player = 1
            else:
                round_bets = [Action(ACTIONS.none) if old_action != Action(ACTIONS.fold) else old_action for old_action in round_bets]
                acting_player = 0
                lead_player = 0

            while any(map(lambda x: x.requires_action(round_bets[lead_player]), round_bets)):
                if not round_bets[acting_player].requires_action(round_bets[lead_player]):
                    acting_player = (acting_player + 1) % num_players

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
