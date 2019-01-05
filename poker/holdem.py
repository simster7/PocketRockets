# This class will be deleted as its contents are moved to other files 

# class Player:
#     def __init__(self, stack):
#         self.stack = stack
#         #All-in check

#     def receive(self, card_1, card_2):
#         self.hand = [card_1, card_2] 
#     def bet(self, amount):
#         return
#     def call(self, amount):
#         return
#     def fold(self):
#         return

# class Game:
#     def __init__(self, players):
#         self.players = players
#         self.deck = []
#         for v_id in range(2, 15):
#             for s_id in range(4):
#                 self.deck.append(Card(v_id, s_id))
#         self.deck.shuffle()
#         #It's random so dealing this way is fine
#         self.deal()
#     def deal(self):
#         for player in self.players:
#             player.receive(self.deck.pop(), self.deck.pop())

# def compare(hand1, hand2):
#     score1, score2 = calucate_hand(hand1), calucate_hand(hand2)
#     if score1 > score2:
#         return 1
#     elif score1 < score2:
#         return -1
#     else:
#         return 0

