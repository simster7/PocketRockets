export interface PlayerState {
    betRound: number;
    leadPlayer?: Player;
    actingPlayer?: Player;
    currentPlayers: Player[];
    playerCards: Card[];
    communityCards: Card[];
    endGame?: EndGameState;
}

export interface Player {
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

export interface Action {
    action: string;
    value?: number;
}

export function parsePlayerStateString(data: JSON): PlayerState {
    console.log(data['current_players'])
    return {
        betRound: data['bet_round'],
        leadPlayer: parsePlayerString(data['lead_player']),
        actingPlayer: parsePlayerString(data['acting_player']),
        currentPlayers: data['current_players'].map(parsePlayerString),
        playerCards: data['player_cards'].map(parseCardString),
        communityCards: data['community_cards'].map(parseCardString),
        endGame: parseEndGameStateString(data['end_game'])
    }
}

function parsePlayerString(data: JSON): Player | undefined {
    if (!data) {
        return undefined;
    }
    return {
        name: data['name'],
        stack: data['stack'],
        seatNumber: data['seat_number'],
        folded: data['folded'],
        lastAction: parseActionString(data['last_action']),
        sittingOut: data['sitting_out']
    }
}

function parseActionString(data: JSON): Action {
    return {
        action: data['action'],
        value: data['value']
    }
}

function parseCardString(data: JSON): Card {
    return {
        cardId: data['card_id'],
    }
}

function parseEndGameStateString(data: JSON): EndGameState | undefined {
    if (!data) {
        return undefined;
    }
    return {
        winners: data['winners'].map(parsePlayerString),
        condition: data['condition']
    }
}
//{'bet_round': 0, 'lead_player': {'name': 'asdfasdf', 'stack': 98, 'seat_number': 3, 'folded': False, 'last_action': {'action': <Actions.bet: 3>, 'value': 2}, 'sitting_out': False}, 'acting_player': {'name': 'asfd', 'stack': 99, 'seat_number': 7, 'folded': False, 'last_action': {'action': <Actions.bet: 3>, 'value': 1}, 'sitting_out': False}, 'current_players': [None, None, None, {'name': 'asdfasdf', 'stack': 98, 'seat_number': 3, 'folded': False, 'last_action': {'action': <Actions.bet: 3>, 'value': 2}, 'sitting_out': False}, None, None, None, {'name': 'asfd', 'stack': 99, 'seat_number': 7, 'folded': False, 'last_action': {'action': <Actions.bet: 3>, 'value': 1}, 'sitting_out': False}, None], 'player_cards': [{'card_id': 27}, {'card_id': 26}], 'community_cards': [], 'end_game': None}
