from __future__ import annotations

from typing import Dict

from .engine.player import Player
from .engine.game import Game


def get_manager() -> Manager:
    if Manager.instance:
        return Manager.instance
    manager = Manager()
    Manager.instance = manager
    return manager


class Manager:
    instance: Manager = None
    rooms: Dict[str, Game]
    players: Dict[str, Player]

    def __init__(self):
        self.rooms = {}
        self.players = {}

    def get_room(self, room_name: str) -> Game:
        if room_name in self.rooms:
            return self.rooms[room_name]
        else:
            new_room = Game(1, 2)
            self.rooms[room_name] = new_room
            return new_room

    def get_player(self, player_name: str) -> Player:
        if player_name in self.players:
            return self.players[player_name]
        else:
            new_player = Player(player_name)
            self.players[player_name] = new_player
            return new_player
