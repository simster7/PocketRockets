import json

from asgiref.sync import async_to_sync
from channels.generic.websocket import WebsocketConsumer

from backend.poker.engine.game import Game, PlayerState
from .engine.player import Player
from .engine.action import Action
from .manager import get_manager, Manager


class TextPokerConsumer(WebsocketConsumer):
    manager: Manager
    game: Game
    player: Player

    def connect(self):
        self.room_name = self.scope['url_route']['kwargs']['room_name']
        self.room_group_name = 'room_%s' % self.room_name

        # Join room group
        async_to_sync(self.channel_layer.group_add)(
            self.room_group_name,
            self.channel_name
        )

        self.manager = get_manager()
        self.game = self.manager.get_room(self.room_name)
        self.player = None

        self.accept()

        self.send_message_to_user("Please login with: /login [name]")

    def disconnect(self, close_code):
        # Leave room group
        async_to_sync(self.channel_layer.group_discard)(
            self.room_group_name,
            self.channel_name
        )

    def receive(self, text_data):
        text_data_json = json.loads(text_data)
        message = text_data_json['message']

        if not self.player and not message[:6] == "/login":
            self.send_message_to_user("You must first login with: /login [name]")
            return

        if message:
            if message[0] == "/":
                args = message[1:].split(" ")
                if args[0] == "login":
                    self.player = self.manager.get_player(" ".join(args[1:]))
                    self.send_message_to_user("You have logged in as: " + self.player.get_name())
                    self.update()
                    return
                if args[0] == "sit":
                    desired_seat_number = int(args[1])
                    if self.game.seats[desired_seat_number]:
                        self.send_message_to_user("This seat is already taken")
                        return
                    self.game.sit_player(self.player, desired_seat_number)
                    self.send_message_to_user("You sat down at seat number " + str(desired_seat_number))
                    return
                if args[0] == "deal":
                    if self.game.is_hand_active():
                        self.send_message_to_user("There is already an active hand being played")
                        return
                    self.game.deal_hand()
                    self.broadcast_game_update()
                    return
                if args[0] == "buy":
                    desired_chips = int(args[1])
                    self.player.set_stack(self.player.get_stack() + desired_chips)
                    self.send_message_to_user("You have bought {} chips, your stack size is now {} "
                                              .format(desired_chips, self.player.get_stack()))
                    return
            else:
                if not self.game.is_hand_active():
                    self.send_message_to_user("There is no active hand, deal a new one with: /deal")
                if self.player and self.player.get_seat_number() is not None:
                    if self.player.get_seat_number() == self.game.get_acting_seat():
                        if message == "F":
                            self.game.take_action(self.player, Action(Action.Actions.fold))
                        elif message == "C":
                            self.game.take_action(self.player, Action(Action.Actions.check))
                        elif message == "L":
                            self.game.take_action(self.player, Action(Action.Actions.call))
                        elif message.isdigit():
                            self.game.take_action(self.player, Action(Action.Actions.bet, int(message)))
                        else:
                            self.send_message_to_user("Please enter a valid action")
                            return
                        self.broadcast_game_update()

    def broadcast_game_update(self):
        async_to_sync(self.channel_layer.group_send)(
            self.room_group_name,
            {
                'type': 'update',
            }
        )

    def send_message_to_user(self, message):
        self.send(text_data=json.dumps({
            'message': message
        }))

    def update(self, event=None):
        if self.player and self.player.get_seat_number() is not None:
            game_string = self.get_personal_game_string(self.game.get_player_state(self.player))
            self.send(text_data=json.dumps({
                'message': game_string
            }))

    # This is a huge abstraction barrier violation, only for testing purposes.
    @staticmethod
    def get_full_game_state_string(game_state):
        if not game_state:
            return ""
        out = ""
        player = game_state.get_acting_player()
        bet_round = game_state.get_round()
        lead_action = game_state.get_lead_action()
        lead_player = game_state.get_leading_player()
        out += "\n"
        out += "Current players: {}".format(game_state.get_players()) + "\n"
        out += "Player hands: {}".format(
            [game_state.get_player_cards(i) for i in range(len(game_state.players))]) + "\n"
        out += "\n"
        out += "Community cards: {}".format(game_state.get_community_cards()) + "\n"
        out += "\n"
        out += "\n"
        out += "Current round: {}".format(bet_round) + "\n"
        out += "Lead action: {}: {}".format(lead_player.name, lead_action) + "\n"
        out += "Acting as player: {}".format(player.name) + "\n"
        out += "With hand: {}".format(game_state.get_player_cards(game_state.get_acting_index())) + "\n"
        out += """
            F - Fold
            C - Check
            L - Call
            [Number] - {} 
        """.format("Call {} and raise [Number]".format(
            lead_action.value) if lead_action.action == Action.Actions.bet else "Bet [Number]") + "\n"
        return out

    @staticmethod
    def get_personal_game_string(player_state: PlayerState) -> str:
        if not player_state:
            return ""
        out = ""
        if player_state.end_game:
            win = player_state.end_game
            out += "=== END OF HAND ===" + "\n"
            if win.condition == "showdown":
                out += win.winner.name + " won with a " + win.hands[0].hand_name + "\n"
            elif win.condition == "folds":
                out += win.winner.name + " won due to folds" + "\n"
            out += "\n\n\nUse /deal to deal a new hand"
        bet_round = player_state.bet_round
        lead_action = player_state.lead_action
        lead_player = player_state.lead_player
        acting_player = player_state.acting_player
        out += "\n"
        out += "Current players: {}".format(player_state.current_players) + "\n"
        out += "Lead action: {}: {}".format(lead_player.name, lead_action) + "\n"
        out += "Current round: {}".format(bet_round) + "\n"
        out += "Current turn: {}".format(acting_player.name) + "\n"
        out += "\n"
        out += "Your hand: {}".format([str(card) for card in player_state.player_cards]) + "\n"
        out += "Community cards: {}".format([str(card) for card in player_state.community_cards]) + "\n"
        out += "\n"
        out += """
            F - Fold
            C - Check
            L - Call
            [Number] - {} 
        """.format("Call {} and raise [Number]".format(
            lead_action.value) if lead_action.action == Action.Actions.bet else "Bet [Number]") + "\n"
        return out
