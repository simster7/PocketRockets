from math import floor
from enum import Enum
from .evaluator import rank_hands


class Action:
    actions = Enum("Actions", "check raise call fold bet")

    def __init__(self, action, value=None):
        self.action = action
        if self.action == Action.actions.bet and value == None:
            raise Exception("A bet value is requiered when betting")
        self.value = value

    def __str__(self):
        return (
            self.action.name
            if self.action != Action.actions.bet
            else self.action.name + " " + str(self.value)
        )


class GameState:
    def __init__(self, old=None):
        if old:
            self.copy(old)

    def start_hand(
        self, players, big_blind, small_blind, deck, starting_pot=0
    ):
        self.big_blind = big_blind
        self.players = players
        self.bet_vector = [small_blind, big_blind] + [0] * (
            len(self.players) - 2
        )
        self.players[0].make_bet(small_blind)
        self.players[1].make_bet(big_blind)
        self.fold_vector = [False] * len(self.players)
        self.all_in = [0] * len(self.players)
        self.pot = starting_pot
        self.remaining_pot = 0
        self.acting_player = 2 % len(self.players)
        self.leading_player = 1
        self.deck = deck
        self.round_over = False
        self.round = 0

    def copy(self, old_game_state):
        self.big_blind = old_game_state.big_blind
        self.players = old_game_state.players
        self.bet_vector = old_game_state.bet_vector
        self.fold_vector = old_game_state.fold_vector
        self.all_in = old_game_state.all_in
        self.pot = old_game_state.pot
        self.remaining_pot = old_game_state.remaining_pot
        self.acting_player = old_game_state.acting_player
        self.leading_player = old_game_state.leading_player
        self.deck = old_game_state.deck
        self.round = old_game_state.round
        self.round_over = False

    def get_remaining_pot(self):
        return self.remaining_pot

    def is_hand_over(self):
        return (sum(self.fold_vector) == len(self.players) - 1) or (
            self.is_round_over() and self.round == 4
        )

    def get_round(self):
        return self.round

    def get_players(self):
        return self.players

    def get_player_cards(self, player_index):
        if self.fold_vector[player_index]:
            return None
        return [
            self.deck[player_index],
            self.deck[player_index + len(self.players)],
        ]

    def get_community_cards(self):
        cards = []
        curr_round = self.get_round()
        num_players = len(self.players)
        if curr_round >= 1:
            cards.extend(self.deck[2 * num_players + 1 : 2 * num_players + 4])
        if curr_round >= 2:
            cards.extend([self.deck[2 * num_players + 5]])
        if curr_round >= 3:
            cards.extend([self.deck[2 * num_players + 7]])
        return cards

    def get_acting_index(self):
        return self.acting_player

    def get_acting_player(self):
        return self.players[self.acting_player]

    def get_leading_player(self):
        return self.players[self.leading_player]

    def get_pot(self):
        return self.pot

    def get_lead_action(self):
        if self.bet_vector[self.leading_player] == 0:
            return Action(Action.actions.check)
        else:
            return Action(
                Action.actions.bet, self.bet_vector[self.leading_player]
            )

    def end_game(self):
        print()
        print("==== END OF HAND ====")
        showdown = [
            i for i in range(len(self.players)) if not self.fold_vector[i]
        ]
        if sum(self.fold_vector) == len(self.players) - 1:
            assert (
                len(showdown) == 1
            ), "Bug: more than one player left on a fold win condition"
            winner = self.players[showdown[0]]
            winner.recieve_pot(self.pot)
            print(winner.name, "won due to folds")
        else:
            showdown_hands = [
                (
                    player_index,
                    self.get_player_cards(player_index)
                    + self.get_community_cards(),
                )
                for player_index in showdown
            ]
            ranked_hands = rank_hands(showdown_hands)
            top_rank = ranked_hands[0].rank
            top_hands = [
                hand for hand in ranked_hands if hand.rank == top_rank
            ]
            winners = [self.players[hand.player_index] for hand in top_hands]
            losers = [hand for hand in ranked_hands if hand.rank != top_rank]
            split = self.pot // len(winners)
            self.remaining_pot = self.pot % len(winners)
            winner_hand = ranked_hands[0].hand_name
            for winner in winners:
                winner.recieve_pot(split)
                print(winner.name, "won with a", winner_hand)
            print(
                "Other showdown hands:",
                ", ".join(
                    [
                        self.players[r_hand.player_index].name
                        + " had a "
                        + r_hand.hand_name
                        for r_hand in losers
                    ]
                ),
            )
        return self

    def move_acting_player(self):
        self.acting_player = (self.acting_player + 1) % len(self.players)
        while (
            self.fold_vector[self.acting_player] and not self.is_round_over()
        ):
            self.acting_player = (self.acting_player + 1) % len(self.players)
        if self.is_round_over():
            self.acting_player = 0
            self.leading_player = 0
            while self.fold_vector[self.acting_player]:
                self.acting_player = (self.acting_player + 1) % len(
                    self.players
                )
                self.leading_player = self.acting_player
            self.pot += sum(self.bet_vector)
            self.bet_vector = [0] * len(self.players)
            self.round += 1
        if self.is_hand_over():
            self.end_game()

    def is_round_over(self):
        full_circle = self.acting_player == self.leading_player
        if full_circle:
            if (
                self.round == 0
                and self.bet_vector[self.leading_player] == self.big_blind
            ):
                return False
        return full_circle or self.round_over

    def take_action(self, action, action_param=None):
        full_circle = self.acting_player == self.leading_player
        limped_in = (
            full_circle
            and self.round == 0
            and self.bet_vector[self.leading_player] == self.big_blind
        )
        if action.action == Action.actions.fold:
            self.fold_vector[self.acting_player] = True
        if action.action == Action.actions.check:
            if limped_in:
                self.round_over = True
            if (
                self.get_lead_action().action == Action.actions.bet
                and not full_circle
            ):
                raise Exception(
                    "Illegal game state: player can't check when there is a bet"
                )
        if action.action == Action.actions.call:
            if limped_in:
                raise Exception(
                    "Illegal game stat: player can't call their own bet"
                )
            to_call = (
                self.bet_vector[self.leading_player]
                - self.bet_vector[self.acting_player]
            )
            if self.players[self.acting_player].stack < to_call:
                raise Exception(
                    "Illegal game state: player doesn't have enough chips to call"
                )
            self.players[self.acting_player].make_bet(to_call)
            self.bet_vector[self.acting_player] += to_call
        if action.action == Action.actions.bet:
            lead_action = self.get_lead_action()
            to_call = 0
            if lead_action.action == Action.actions.bet:
                to_call = max(
                    lead_action.value - self.bet_vector[self.acting_player], 0
                )
            current_player = self.players[self.acting_player]
            if current_player.stack < action.value:
                raise Exception(
                    "Illegal game state: player doesn't have enough chips to make that bet"
                )
            self.players[self.acting_player].make_bet(action.value + to_call)
            self.bet_vector[self.acting_player] += action.value + to_call
            self.leading_player = self.acting_player
        self.move_acting_player()
        return GameState(self)
