import { ENGINE_METHOD_NONE } from "constants";

interface PlayerState {
    betRound: number;
    leadPlayer: Player;
    actingPlayer: Player;
    currentPlayers: Player[];
    playerCards: Card[];
    communityCards: Card[];
    endGame: EndGameState;
}

interface Player {
    name: string;
    stack: number;
    seatNumber?: number;
    folded: boolean;
    lastAction?: Action;
    sittingOut: boolean;

}

interface Card {
    cardId: number;
}

interface EndGameState {
    winners: Player[];
    condition: string;
}

enum Actions {
    check,
    call,
    fold,
    bet
}

interface Action {
    action: string;
    value?: number;
}

function parsePlayerStateString(json: JSON): PlayerState {
    return {
        betRound: data['bet_round'],
        leadPlayer: parsePlayerString(data['lead_player']),
        actingPlayer: parsePlayerString(data['acting_player']),
        currentPlayers: data['current_players'].map(player: JSON => parsePlayerString(player)),

    }
}

function parsePlayerString(json: JSON): Player {
    return {
        name: json['name'],
        stack: json['stack'],
        seatNumber: json['seat_number'],
        folded: json['folded'],
        lastAction: parseActionString(json['last_action']),
        sittingOut: json['sitting_out']
    }
}

function parseActionString(json: JSON): Action {
    return {
        action: json['action'],
        value: json['value']
    }
}

//{'bet_round': 0, 'lead_player': {'name': 'asdfasdf', 'stack': 98, 'seat_number': 3, 'folded': False, 'last_action': {'action': <Actions.bet: 3>, 'value': 2}, 'sitting_out': False}, 'acting_player': {'name': 'asfd', 'stack': 99, 'seat_number': 7, 'folded': False, 'last_action': {'action': <Actions.bet: 3>, 'value': 1}, 'sitting_out': False}, 'current_players': [None, None, None, {'name': 'asdfasdf', 'stack': 98, 'seat_number': 3, 'folded': False, 'last_action': {'action': <Actions.bet: 3>, 'value': 2}, 'sitting_out': False}, None, None, None, {'name': 'asfd', 'stack': 99, 'seat_number': 7, 'folded': False, 'last_action': {'action': <Actions.bet: 3>, 'value': 1}, 'sitting_out': False}, None], 'player_cards': [{'card_id': 27}, {'card_id': 26}], 'community_cards': [], 'end_game': None}
