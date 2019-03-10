import * as React from 'react';

import { Player } from '../utils/PlayerState';

interface IProp {
	player: Player;
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
			</div>
		);
	}
}

export default PokerPlayer;
