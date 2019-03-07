import json

from asgiref.sync import async_to_sync
from channels.generic.websocket import WebsocketConsumer

from .engine.game import Game
from .engine.player import Player
from .manager import get_manager, Manager


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

    def broadcast_game_update(self):
        async_to_sync(self.channel_layer.group_send)(
            self.room_group_name,
            {
                'type': 'update',
            }
        )

    def update(self, event=None):
        print("sent")
        self.send(text_data=json.dumps({
            'message': self.game.get_player_state(self.player).__dict__
        }))
