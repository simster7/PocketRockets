from poker.engine.player import Player
from .engine.game import Game

def get_manager():
    if Manager.instance:
        return Manager.instance
    manager = Manager()
    Manager.instance = manager
    return manager

class Manager:
    instance = None

    def __init__(self):
        self.rooms = {}
        self.players = {}

    def get_room(self, room_name):
        if room_name in self.rooms:
            return self.rooms[room_name]
        else:
            new_room = Game(1, 2)
            self.rooms[room_name] = new_room
            return new_room

    def get_player(self, player_name):
        if player_name in self.players:
            return self.players[player_name]
        else:
            new_player = Player(player_name)
            self.players[player_name] = new_player
            return new_player
