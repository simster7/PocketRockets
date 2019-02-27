from ..engine.game import Game
from ..engine.state import Action
from ..engine.player import Player


class ConsoleGame(Game):

    # TESTING ONLY
    def get_state_testing(self):
        return self.game_state

    # TESTING ONLY
    def take_action_testing(self, action: Action) -> None:
        self.game_state = self.game_state.take_action(action)


if __name__ == '__main__':

    def prompt(game_state):
        player = game_state.get_acting_player()
        bet_round = game_state.get_round()
        lead_action = game_state.get_lead_action()
        lead_player = game_state.get_leading_player()
        print()
        print("Current players: {}".format(game_state.get_players()))
        print("Player hands: {}".format([[str(card) for card in game_state.get_player_cards(i)] if game_state.get_player_cards(i) else "Fold" for i in range(len(game_state.players))]))
        print()
        print("Community cards: {}".format([str(card) for card in game_state.get_community_cards()]))
        print("Pot: {}".format(game_state.get_pot()))
        print()
        print("Current round: {}".format(bet_round))
        print("Lead action: {}: {}".format(lead_player.name, lead_action))
        print("Acting as player: {}".format(player.name))
        print("With hand: {}".format([str(card) for card in game_state.get_player_cards(game_state.get_acting_index())]))
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


    game = ConsoleGame(1, 2)
    simon = Player("Simon")
    simon.set_stack(100)
    game.sit_player(simon, 0)
    hersh = Player("Hersh")
    hersh.set_stack(100)
    game.sit_player(hersh, 1)
    chien = Player("Chien")
    chien.set_stack(100)
    game.sit_player(chien, 2)
    jarry = Player("Jarry")
    jarry.set_stack(100)
    game.sit_player(jarry, 3)
    grace = Player("Grace")
    grace.set_stack(100)
    game.sit_player(grace, 6)
    jason = Player("Jason")
    jason.set_stack(100)
    game.sit_player(jason, 7)
    while True:
        game.deal_hand()
        while not game.game_state.is_hand_over():
            game.take_action_testing(prompt(game.get_state_testing()))
