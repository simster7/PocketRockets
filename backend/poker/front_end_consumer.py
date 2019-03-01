import json

from asgiref.sync import async_to_sync
from channels.generic.websocket import WebsocketConsumer

from .engine.game import Game, PlayerState
from .engine.player import Player
from .engine.action import Action
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

    def disconnect(self, close_code):
        # Leave room group
        async_to_sync(self.channel_layer.group_discard)(
            self.room_group_name,
            self.channel_name
        )

    def update_game_state(self, event=None):
        if self.player:
            self.send(text_data=json.dumps(self.game.get_player_state(self.player)))
