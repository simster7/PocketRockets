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
        elif action.isdigit():
            return game_state.take_action(Action(Action.actions.bet, int(action)))

class Game:

    def __init__(self, small_blind_amount, big_blind_amount, prompt=prompt):
        self.seats = [None] * 9
        self.button_position = 0
        self.small_blind = small_blind_amount
        self.big_blind = big_blind_amount
        self.prompt = prompt
        self.remaining_pot = 0

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
        game_state.start_hand(players, self.big_blind, self.small_blind, deck,
                              starting_pot=self.remaining_pot)
        while not game_state.is_hand_over():
            game_state = self.prompt(game_state)
        self.remaining_pot = game.get_remaining_pot()

    def deal_fake_hand(self):
        players = [player for player in self.seats if player != None]
        deck = [Card(i) for i in range(52)]
        shuffle(deck) # TODO: replace this with more e n t r o p y
        print(deck)
        AS = [i for i, card in enumerate(deck) if card.suit == 'S' and card.rank == 'A'][0]
        AH = [i for i, card in enumerate(deck) if card.suit == 'H' and card.rank == 'A'][0] 
        AC = [i for i, card in enumerate(deck) if card.suit == 'C' and card.rank== 'A'][0] 
        AD = [i for i, card in enumerate(deck) if card.suit == 'D' and card.rank == 'A'][0]
        indices = [0,1,3,4]
        for i, j in zip(indices, [AS, AH, AC, AD]):
            deck[i], deck[j] = deck[j], deck[i]
        print(deck) 
        game_state = GameState()
        game_state.start_hand(players, self.big_blind, self.small_blind, deck,
                              starting_pot=self.remaining_pot)
        while not game_state.is_hand_over():
            game_state = self.prompt(game_state)
        self.remaining_pot = game_state.get_remaining_pot()


if __name__ == '__main__':
    game = Game(1, 2)
    game.sit_player(Player("Simon", 100), 0)
    game.sit_player(Player("Hersh", 100), 1)
    game.sit_player(Player("Chien", 100), 2)
    #game.sit_player(Player("Jarry", 100), 3)
    #game.sit_player(Player("Grace", 100), 6)
    #game.sit_player(Player("Jason", 100), 7)
    while True:
        game.deal_fake_hand()
