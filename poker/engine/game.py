from player import Player
from card import Card
from evaluator import rank_hands
from state_new import GameState, Action

from random import shuffle


def prompt(game_state):
    player = game_state.get_acting_player()
    bet_round = game_state.get_round()
    print()
    print(f"Current players: {game_state.get_players()}")
    player_hands = [
        game_state.get_player_cards(i)
        for i in range(len(game_state.players))
    ]
    print(f"Player hands: {player_hands}")
    print()
    print(f"Community cards: {game_state.get_community_cards()}")
    print()
    print(f"Current round: {bet_round}")
    print(f"Acting as player: {player.name}")
    print(
        "With hand: {}".format(
            game_state.get_player_cards(game_state.get_acting_index())
        )
    )
    side_pots = [x for x in game_state.get_pot()]
    if len(side_pots) > 1:
        print(f"Fixed Pots: {side_pots[:-1]}")
    print(
        """
        F - Fold
        C - Check
        L - Call
        [Number] - {} 
    """.format("Bet [Number]")
    )
    print(f"Active Pot: {side_pots[-1]}")
    print(f"Current Pot: {game_state.get_current_pot()}")
    print(f"Acted Players: {game_state.get_acted_queue()}")
    print(f"Players to act: {game_state.get_queue()}")
    while True:
        action = input("Enter action: ")
        if action == "F":
            return game_state.take_action(Action(Action.actions.fold))
        elif action == "C":
            return game_state.take_action(Action(Action.actions.check))
        elif action == "L":
            return game_state.take_action(Action(Action.actions.call))
        elif action.isdigit():
            return game_state.take_action(
                Action(Action.actions.bet, int(action))
            )


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
            raise Exception(
                "Cannot sit player ({}) on an occupied seat".format(player)
            )

    def deal_hand(self):
        players = [player for player in self.seats if player != None]
        deck = [Card(i) for i in range(52)]
        shuffle(deck)  # TODO: replace this with more e n t r o p y
        game_state = GameState()
        game_state.start_hand(
            players,
            self.big_blind,
            self.small_blind,
            deck
        )
        hand_over = False
        while not hand_over: 
            game_state, hand_over = self.prompt(game_state)

    def deal_fake_hand(self):
        players = [player for player in self.seats if player != None]
        deck = [Card(i) for i in range(52)]
        shuffle(deck)  # TODO: replace this with more e n t r o p y
        print(deck)
        AS = [
            i
            for i, card in enumerate(deck)
            if card.suit == "S" and card.rank == "A"
        ][0]
        AH = [
            i
            for i, card in enumerate(deck)
            if card.suit == "H" and card.rank == "A"
        ][0]
        AC = [
            i
            for i, card in enumerate(deck)
            if card.suit == "C" and card.rank == "A"
        ][0]
        AD = [
            i
            for i, card in enumerate(deck)
            if card.suit == "D" and card.rank == "A"
        ][0]
        indices = [0, 1, 3, 4]
        for i, j in zip(indices, [AS, AH, AC, AD]):
            deck[i], deck[j] = deck[j], deck[i]
        print(deck)
        game_state = GameState()
        game_state.start_hand(
            players,
            self.big_blind,
            self.small_blind,
            deck
        )
        hand_over = False
        while not hand_over: 
            game_state, hand_over = self.prompt(game_state)
        game_state.end_game()

if __name__ == "__main__":
    game = Game(1, 2)
    game.sit_player(Player("Simon", 100), 0)
    game.sit_player(Player("Hersh", 90), 1)
    game.sit_player(Player("Chien", 80), 2)
    # game.sit_player(Player("Jarry", 100), 3)
    # game.sit_player(Player("Grace", 100), 6)
    # game.sit_player(Player("Jason", 100), 7)
    while True:
        game.deal_fake_hand()
