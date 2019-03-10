import * as React from 'react';

import { PlayerState } from '../utils/PlayerState';
import PokerPlayer from './Player';

interface IProp {
	playerState: PlayerState;
}

interface IState {}

class PokerGame extends React.Component<IProp, IState> {
	public constructor(props: IProp) {
		super(props);
		this.state = {};
	}

	public render() {
		console.log(this.props.playerState);
		return (
			<div>
				<table>
					<tbody>
						<tr>
							<td>&nbsp;</td>
							<td>
								{this.props.playerState.buttonPosition == 0 ? '(D)' : ''}{' '}
								<PokerPlayer player={this.props.playerState.currentPlayers[0]} />
							</td>
							<td>&nbsp;</td>
							<td>
								{this.props.playerState.buttonPosition == 1 ? '(D)' : ''}{' '}
								<PokerPlayer player={this.props.playerState.currentPlayers[1]} />
							</td>
							<td>&nbsp;</td>
							<td>
								{this.props.playerState.buttonPosition == 2 ? '(D)' : ''}{' '}
								<PokerPlayer player={this.props.playerState.currentPlayers[2]} />
							</td>
							<td>&nbsp;</td>
						</tr>
						<tr>
							<td>
								{this.props.playerState.buttonPosition == 8 ? '(D)' : ''}{' '}
								<PokerPlayer player={this.props.playerState.currentPlayers[8]} />
							</td>
							<td>&nbsp;</td>
							<td>&nbsp;</td>
							<td>&nbsp;</td>
							<td>&nbsp;</td>
							<td>&nbsp;</td>
							<td>
								{this.props.playerState.buttonPosition == 3 ? '(D)' : ''}{' '}
								<PokerPlayer player={this.props.playerState.currentPlayers[3]} />
							</td>
						</tr>
						<tr>
							<td>
								{this.props.playerState.buttonPosition == 7 ? '(D)' : ''}{' '}
								<PokerPlayer player={this.props.playerState.currentPlayers[7]} />
							</td>
							<td>&nbsp;</td>
							<td>&nbsp;</td>
							<td>&nbsp;</td>
							<td>&nbsp;</td>
							<td>&nbsp;</td>
							<td>
								{this.props.playerState.buttonPosition == 4 ? '(D)' : ''}{' '}
								<PokerPlayer player={this.props.playerState.currentPlayers[4]} />
							</td>
						</tr>
						<tr>
							<td>&nbsp;</td>
							<td>&nbsp;</td>
							<td>
								{this.props.playerState.buttonPosition == 6 ? '(D)' : ''}{' '}
								<PokerPlayer player={this.props.playerState.currentPlayers[6]} />
							</td>
							<td>&nbsp;</td>
							<td>
								{this.props.playerState.buttonPosition == 5 ? '(D)' : ''}{' '}
								<PokerPlayer player={this.props.playerState.currentPlayers[5]} />
							</td>
							<td>&nbsp;</td>
							<td>&nbsp;</td>
						</tr>
					</tbody>
				</table>
			</div>
		);
	}
}

export default PokerGame;
