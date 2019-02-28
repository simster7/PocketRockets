from enum import Enum


class Action:
    class Actions(Enum):
        check = 0
        call = 1
        fold = 2
        bet = 3

    def __init__(self, action, value=None):
        self.action = action
        if self.action == Action.Actions.bet and value is None:
            raise Exception("A bet value is required when betting")
        self.value = value

    def __str__(self):
        return self.action.name if self.action != Action.Actions.bet else self.action.name + " " + str(self.value)

    def __repr__(self):
        return str(self.__dict__)
