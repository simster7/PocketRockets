from math import floor
from collections import deque
from enum import Enum
from evaluator import rank_hands


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

    def start_hand(
        self, players, big_blind, small_blind, deck
    ):
        self.deck = deck
        self.folded = []
        self.all_in = {}
        self.players = [p for p in players if p.stack > 0]
        self.side_pots = deque()
        self.acted = deque()
        self.to_act = deque(list(range(len(self.players))))
        self.to_act.rotate(-2)
        self.players[0].make_bet(small_blind)
        self.players[1].make_bet(big_blind)
        self.last_to_act = 1
        self.bet_vector = [small_blind, big_blind] + [0] * (len(self.players) - 2)
        self.current_pot = small_blind + big_blind 
        self.side_pots.append(0)
        self.round = 0
        self.min_bet = big_blind
        self.current_bet = big_blind
        print(self.to_act)

    def get_round(self):
        return self.round

    def get_players(self):
        return self.players

    def get_pot(self):
        return self.side_pots

    def get_current_pot(self):
        return self.current_pot

    def get_acted_queue(self):
        return [self.players[i] for i in self.acted]

    def get_queue(self):
        return [self.players[i] for i in self.to_act]

    def get_player_cards(self, player_index):
        if player_index in self.folded:
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
        if self.to_act:
            return self.to_act[0]
        elif self.acted:
            return min(self.acted)

    def get_acting_player(self):
        return self.players[self.get_acting_index()]

    def end_game(self):
        showdown = [
            i for i in range(len(self.players)) if i not in self.folded
        ]
        showdown_hands = [
            (
                player_index,
                self.get_player_cards(player_index)
                + self.get_community_cards(),
            )
            for player_index in showdown
        ]
        print(f"Showdown: {showdown_hands}")
        ranked_hands = deque(rank_hands(showdown_hands))
        payoffs = {player: 0 for player in showdown}
        side_bet_index = 0
        winners = []
        while self.side_pots:
            if not winners:
                top_hands = [ranked_hands.popleft()]
                top_rank = top_hands[0].rank
                while ranked_hands[0].rank == top_rank:
                    top_hands.append(ranked_hands.popleft())
                winners = sorted([hand.player_index for hand in top_hands])
            print(f"Winners: {[self.players[i] for i in winners]}")
            split_pot = [self.all_in[w] for w in winners if w in self.all_in]
            print(f"Split Pot: {split_pot}")
            pot = 0
            if split_pot:
                min_side_pots = min(split_pot) + 1 - side_bet_index
                side_bet_index += min_side_pots
            else:
                min_side_pots = len(self.side_pots)
            total = 0
            for _ in range(min_side_pots):
                total += self.side_pots.popleft()
            N = len(winners) 
            split = total // N 
            extra = total % N 
            for player in winners:
                payoffs[player] += split
            if extra:
                receipient = self.last_to_act
                while receipient not in winners:
                    receipient = (receipient + 1) % len(self.players)
                start = winners.index(receipient)
                for i in range(extra):
                    payoffs[winners[start+i]] += 1
            winners = [w for w in winners if w not in self.all_in or self.all_in[w] >= min_side_pots]

        for i, value in payoffs.items():
            self.players[i].recieve_pot(value)
        for player in self.players:
            player.reset()
            print("Result:", player)
        return self

    def resolve_pot(self):
        all_in_players = []
        remaining_players = []

        # Determine who is all-in based on the current rounds betting
        for player_index, bet in enumerate(self.bet_vector):
            if bet:
                if self.players[player_index].stack == 0:
                    all_in_players.append((bet, player_index))
                else:
                    remaining_players.append((bet, player_index))

        if not all_in_players:
            self.side_pots[-1] += self.current_pot
            return

        print(f"All in: {[self.players[i] for _, i in all_in_players]}")
        print(f"Still betting: {[self.players[i] for _, i in remaining_players]}")

        # Sort the all in players by stack size
        all_in_players.sort(reverse=True) 

        # Work out how individual side pots should be distributed given stack sizes
        prev_bet = 0
        while all_in_players:
            total_players = len(all_in_players) + len(remaining_players)
            bet, player_index = all_in_players.pop()
            # With each new all-in, we compute the marginal increase in pot share
            pot_size = (bet - prev_bet) * total_players
            if pot_size > 0:
                if prev_bet == 0:
                    self.side_pots[-1] += pot_size
                else:
                    self.side_pots.append(pot_size) 
                prev_bet = bet
            self.all_in[player_index] = len(self.side_pots) - 1

        print(f"Side Pots: {self.side_pots}")

        if remaining_players:
            assert(all([bet == remaining_players[0][0] for bet, _ in remaining_players]))
            pot_size = (remaining_players[0][0] - prev_bet) * len(remaining_players)
            self.side_pots.append(pot_size)


    def take_action(self, action):
        # This is the index of the current player
        player_index = self.to_act.popleft()

        if action.action == Action.actions.fold:
            self.folded.append(player_index)

        elif action.action == Action.actions.check:
            # The only case where you can check when the current bet is > 0
            # is when everyone limps to the big blind
            preflop_edge_case = (
                self.round == 0
                and self.current_bet == self.min_bet 
                and not self.to_act 
            )
            
            if self.current_bet == 0 or preflop_edge_case:
                self.acted.append(player_index)
            elif self.current_bet > 0:
                raise Exception(
                    "Illegal game state: player can't check when there is a bet"
                )

        elif action.action == Action.actions.call:
            not_all_in = self.current_bet < self.players[player_index].stack
            if not_all_in:
                to_call = self.current_bet - self.bet_vector[player_index]
                self.acted.append(player_index)
            else:
                to_call = self.players[player_index].stack
            self.players[player_index].make_bet(to_call)
            self.current_pot += to_call
            self.bet_vector[player_index] += to_call 

        elif action.action == Action.actions.bet:
            not_all_in = action.value < self.players[player_index].stack
            if not action.value:
                raise Exception("Bet has no value")
            if not_all_in:
                if self.current_bet > 0 and action.value < 2 * self.current_bet: 
                    raise Exception("Raise is too small")
                elif self.current_bet == 0 and action.value < self.min_bet:
                    raise Exception("Bet is too small")
                self.current_bet = action.value
            else:
                if self.current_bet > 0 and action.value < 2 * self.current_bet: 
                    raise Exception("Should call instead of bet")
                self.current_bet = self.players[player_index].stack
            self.players[player_index].make_bet(self.current_bet)
            self.current_pot += self.current_bet
            self.bet_vector[player_index] = self.current_bet 
            while self.acted:
                self.to_act.append(self.acted.popleft())
            if not_all_in:
                self.acted.append(player_index)
            self.last_to_act = player_index

        if not self.to_act:
            self.resolve_pot()
            self.round += 1
            if len(self.acted) <= 1 or self.round == 4:
                return self, True 
            else:
                self.to_act = deque(list(sorted(self.acted)))
                self.acted = deque()
                self.current_bet = 0
                self.current_pot = 0
                self.bet_vector = [0] * len(self.to_act)
            return self, False

        return self, False 
