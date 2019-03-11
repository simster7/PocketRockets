import * as React from 'react';

import { Player, Card } from '../utils/PlayerState';
import PokerPlayer from './Player';
import { cardIdToString } from 'src/utils/CardConverter';
import { Message } from './PokerClient';

interface IProp {
	seatNumber: number;
	buttonPosition: number;
	cards?: Card[];
	player: Player;
	sendMessage: (message: Message) => void;
}

interface IState {}

class PlayerContainer extends React.Component<IProp, IState> {
	public constructor(props: IProp) {
		super(props);
	}

	public render() {
		return (
			<div>
				{this.props.buttonPosition == this.props.seatNumber ? '(D)' : ''}
				<PokerPlayer player={this.props.player} />
				{this.props.cards ? (
					cardIdToString(this.props.cards[0].cardId) + ' ' + cardIdToString(this.props.cards[1].cardId)
				) : (
					''
				)}
				{!this.props.player ? (
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
