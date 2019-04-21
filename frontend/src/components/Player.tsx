import * as React from 'react';

import { Player, Card } from '../utils/PlayerState';
import { cardIdToString } from 'src/utils/CardConverter';

interface IProp {
	player: Player;
	cards?: Card[];
}

interface IState {}

class PokerPlayer extends React.Component<IProp, IState> {
	public constructor(props: IProp) {
		super(props);
		console.log(props);
	}

	public render() {
		if (!this.props.player) {
			return <div>Empty</div>;
		}
		return (
			<div>
				{this.props.player.name}: {this.props.player.stack}
				{this.props.cards && this.props.cards.length == 2 ? (
					cardIdToString(this.props.cards[0].cardId) + ' ' + cardIdToString(this.props.cards[1].cardId)
				) : (
					''
				)}
				<br />
				{this.props.player.lastAction ? (
					this.props.player.lastAction.action +
					(this.props.player.lastAction.value ? ' ' + this.props.player.lastAction.value : '')
				) : (
					''
				)}
			</div>
		);
	}
}

export default PokerPlayer;
