import * as React from 'react';

import { Player, Card } from '../utils/PlayerState';
import PokerPlayer from './Player';
import { Message } from './PokerClient';

interface IProp {
	seatNumber: number;
	buttonPosition: number;
	cards?: Card[];
	player: Player;
	showSitButton: boolean;
	isCurrentTurn: boolean;
	sendMessage: (message: Message) => void;
}

interface IState {}

class PlayerContainer extends React.Component<IProp, IState> {
	public constructor(props: IProp) {
		super(props);
	}

	private turnStyle = {
		borderStyle: 'solid',
		borderWidth: '1px'
	};

	public render() {
		return (
			<div style={this.props.isCurrentTurn ? this.turnStyle : {}}>
				{this.props.buttonPosition == this.props.seatNumber ? '(D)' : ''}
				<PokerPlayer player={this.props.player} cards={this.props.cards} />
				{!this.props.player && this.props.showSitButton ? (
					<button
						onClick={() =>
							this.props.sendMessage({
								type: 'manage',
								value: {
									item: 'sit',
									value: '' + this.props.seatNumber
								}
							})}
					>
						Sit here
					</button>
				) : (
					''
				)}
			</div>
		);
	}
}

export default PlayerContainer;
