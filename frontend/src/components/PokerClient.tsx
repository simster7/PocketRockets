import * as React from 'react';

import { parsePlayerStateString, PlayerState, Action } from '../utils/PlayerState';
import PokerGame from './PokerGame';
import PokerController from './PokerController';

export interface Message {
	type: string;
	value: ManageMessage | Action | JSON;
}

export interface ManageMessage {
	item: string;
	value: string;
}

interface IProp {
	roomName: string;
}

interface IState {
	playerState?: PlayerState;
}

class PokerClient extends React.Component<IProp, IState> {
	private connection: WebSocket;

	public constructor(props: IProp) {
		super(props);
		this.state = {};

		this.sendMessage = this.sendMessage.bind(this);
	}

	public componentDidMount() {
		this.connection = new WebSocket('ws://localhost:8000/ws/poker/room/' + this.props.roomName + '/');
		this.connection.onmessage = function(this: PokerClient, e: any) {
			let rawData: string = e.data;
			//let data = "{'message': {'bet_round': 0, 'lead_player': {'name': 'Jarry', 'stack': 98, 'seat_number': 3, 'folded': False, 'last_action': {'action': 'bet' , 'value': 2}, 'sitting_out': False}, 'acting_player': {'name': 'Simon', 'stack': 99, 'seat_number': 6, 'folded': False, 'last_action': {'action': 'bet' , 'value': 1}, 'sitting_out': False}, 'current_players': [None, None, None, {'name': 'Jarry', 'stack': 98, 'seat_number': 3, 'folded': False, 'last_action': {'action': 'bet' , 'value': 2}, 'sitting_out': False}, None, None, {'name': 'Simon', 'stack': 99, 'seat_number': 6, 'folded': False, 'last_action': {'action': 'bet' , 'value': 1}, 'sitting_out': False}, None, None], 'player_cards': [{'card_id': 14}, {'card_id': 22}], 'community_cards': [], 'end_game': None}}";
			// data = data.replace(/'/g, '"').replace(/True/g, 'true').replace(/False/g, 'false').replace(/None/g, 'null');
			let data: Message = JSON.parse(rawData);
			console.log(data);
			if (data.type === 'game_update') {
				// TODO PLAYER STATE
				let playerState: PlayerState = parsePlayerStateString(data.value as JSON);
				console.log(playerState);
				// let playerState: PlayerState = data.value as PlayerState;
				this.setState({ playerState });
			}
		}.bind(this);
	}

	private sendMessage = (message: Message) => {
		if (this.connection) {
			this.connection.send(JSON.stringify(message));
		}
	};

	public render() {
		if (this.state.playerState === undefined) {
			return <div>Loading...</div>;
		} else {
			return (
				<div>
					<PokerGame playerState={this.state.playerState} sendMessage={this.sendMessage} />
					<br />
					<PokerController sendMessage={this.sendMessage} />
				</div>
			);
		}
	}
}

export default PokerClient;
