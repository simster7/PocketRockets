from enum import Enum
from typing import Optional


class Action:
    class Actions(Enum):
        check = 0
        call = 1
        fold = 2
        bet = 3

    action: Actions
    value: Optional[int]

    def __init__(self, action: Actions, value: Optional[int] = None):
        self.action = action
        if self.action == Action.Actions.bet and value is None:
            raise Exception("A bet value is required when betting")
        self.value = value

    def __str__(self):
        return "{'action': '" + self.action.name + "' , 'value': " + str(self.value) + "}"

    def __repr__(self):
        return "{'action': '" + self.action.name + "' , 'value': " + str(self.value) + "}"
