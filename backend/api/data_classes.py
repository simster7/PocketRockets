from dataclasses import dataclass
from typing import Optional, List, Tuple

from backend.poker.engine.card import Card
from backend.poker.engine.evaluator import RankedHand
from backend.poker.engine.player import Player


@dataclass
class EndGameState:
    winners: List[Tuple[Player, int]]
    condition: str
    hands: List[RankedHand]


@dataclass
class PlayerState:
    bet_round: int
    lead_player: Optional[Player]
    acting_player: Optional[Player]
    current_players: List[Optional[Player]]
    player_cards: Optional[List[Card]]
    community_cards: Optional[List[Card]]
    end_game: Optional[EndGameState]
    player_seat: Optional[int]
    button_position: int
    pot: int


