import json
from ast import literal_eval
from dataclasses import dataclass

from asgiref.sync import async_to_sync
from channels.generic.websocket import WebsocketConsumer

from .engine.action import Action
from .engine.game import Game
from .engine.player import Player
from .manager import get_manager, Manager


@dataclass
class ManageMessage:
    item: str
    value: str


class FrontEndConsumer(WebsocketConsumer):
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
        self.broadcast_game_update()

    def disconnect(self, close_code):
        # Leave room group
        async_to_sync(self.channel_layer.group_discard)(
            self.room_group_name,
            self.channel_name
        )

    def receive(self, text_data):
        raw_message = json.loads(text_data)
        print(raw_message)
        # raw_message = json.loads(text_data_json['message'])
        #
        # message: Message = Message(raw_message['type'], raw_message['value'])
        #
        # if message.type == 'manage':
        #     message.value = ManageMessage(message.value['item'], message.value['value'])
        # elif message.type == 'action':
        #     message.value = Action(message.value['action'], message.value['value'] if message.value['value'] else None)

        # if not self.player and not message == 'manage' and not message['value']['item'] == 'player_name':
        #     self.send_error('Must choose player name first')

        if raw_message['type'] == 'manage':
            manage_message = ManageMessage(raw_message['value']['item'], raw_message['value']['value'])
            if manage_message.item == 'player_name':
                self.player = self.manager.get_player(manage_message.value)
                self.send_info("You have logged in as: " + self.player.name)
                return
            elif manage_message.item == 'sit':
                desired_seat_number = int(manage_message.value)
                if self.game.seats[desired_seat_number]:
                    self.send_error("This seat is already taken")
                    return
                self.game.sit_player(self.player, desired_seat_number)
                self.send_info("You sat down at seat number " + str(desired_seat_number))
                return
            elif manage_message.item == 'add_stack':
                desired_chips = int(manage_message.value)
                self.player.stack = self.player.stack + desired_chips
                self.send_info("You have bought {} chips, your stack size is now {} "
                               .format(desired_chips, self.player.stack))
            elif manage_message.item == 'deal':
                if self.game.is_hand_active():
                    self.send_error("There is already an active hand being played")
                    return
                self.game.deal_hand()
                self.broadcast_game_update()
                return

        elif raw_message['type'] == 'action':
            action = Action(Action.Actions[raw_message['value']['action']],
                            int(raw_message['value']['value']) if 'value' in raw_message['value'] else None)
            print(action)
            if not self.game.is_hand_active():
                self.send_error("There is no active hand, deal a new one")
                return
            if self.player and self.player.seat_number is not None:
                if self.player.seat_number == self.game.get_acting_seat():
                    if action.action == Action.Actions.fold:
                        self.game.take_action(self.player, Action(Action.Actions.fold))
                    elif action.action == Action.Actions.check:
                        self.game.take_action(self.player, Action(Action.Actions.check))
                    elif action.action == Action.Actions.call:
                        self.game.take_action(self.player, Action(Action.Actions.call))
                    elif action.action == Action.Actions.bet:
                        self.game.take_action(self.player, Action(Action.Actions.bet, int(action.value)))
                    else:
                        self.send_error("Please enter a valid action")
                        return
                    self.broadcast_game_update()

    def broadcast_game_update(self):
        async_to_sync(self.channel_layer.group_send)(
            self.room_group_name,
            {
                'type': 'update',
            }
        )

    def update(self, event=None):
        if self.game.game_state:  # TODO Abstraction violation
            data = literal_eval(str(self.game.get_player_state(self.player).__dict__))
            self.send(text_data=json.dumps({
                'type': 'game_update',
                'value': data
            }))

    def send_error(self, error):
        self.send(text_data=json.dumps({
            'type': 'error',
            'value': error
        }))

    def send_info(self, info):
        self.send(text_data=json.dumps({
            'type': 'info',
            'value': info
        }))
