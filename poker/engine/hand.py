class Hand:
    def __init__(self, cards, owner=None):
        """
        cards - A list of Card objects
        owner - A Player object
        """
        self.cards = cards
        self.owner = owner

    def cards(self):
        return self.cards

    def cards_with_community(self, community):
        return self.cards + community

    def owner(self):
        return self.owner

    def __str__(self):
        return "Hand({})".format(str(self.cards))
