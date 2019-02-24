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

    def get_room(self, room_name):
        if room_name in self.rooms:
            return self.rooms[room_name]
        else:
            new_room = Game(1, 2)
            self.rooms[room_name] = new_room
            return new_room
