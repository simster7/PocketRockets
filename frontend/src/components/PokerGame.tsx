import * as React from 'react';

import { PlayerState } from '../utils/PlayerState';
import PlayerContainer from './PlayerContainer';
import { Message } from './PokerClient';
import TableInfo from './TableInfo';

interface IProp {
	playerState: PlayerState;
	sendMessage: (message: Message) => void;
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
								<PlayerContainer
									seatNumber={0}
									buttonPosition={this.props.playerState.buttonPosition}
									player={this.props.playerState.currentPlayers[0]}
									cards={
										this.props.playerState.playerSeat == 0 ? (
											this.props.playerState.playerCards
										) : (
											undefined
										)
									}
									sendMessage={this.props.sendMessage}
									isCurrentTurn={
										this.props.playerState.actingPlayer ? (
											this.props.playerState.actingPlayer.seatNumber === 0
										) : (
											false
										)
									}
									showSitButton={!this.props.playerState.playerSeat}
								/>
							</td>
							<td>&nbsp;</td>
							<td>
								<PlayerContainer
									seatNumber={1}
									buttonPosition={this.props.playerState.buttonPosition}
									player={this.props.playerState.currentPlayers[1]}
									cards={
										this.props.playerState.playerSeat == 1 ? (
											this.props.playerState.playerCards
										) : (
											undefined
										)
									}
									sendMessage={this.props.sendMessage}
									isCurrentTurn={
										this.props.playerState.actingPlayer ? (
											this.props.playerState.actingPlayer.seatNumber === 1
										) : (
											false
										)
									}
									showSitButton={!this.props.playerState.playerSeat}
								/>
							</td>
							<td>&nbsp;</td>
							<td>
								<PlayerContainer
									seatNumber={2}
									buttonPosition={this.props.playerState.buttonPosition}
									player={this.props.playerState.currentPlayers[2]}
									cards={
										this.props.playerState.playerSeat == 2 ? (
											this.props.playerState.playerCards
										) : (
											undefined
										)
									}
									sendMessage={this.props.sendMessage}
									isCurrentTurn={
										this.props.playerState.actingPlayer ? (
											this.props.playerState.actingPlayer.seatNumber === 2
										) : (
											false
										)
									}
									showSitButton={!this.props.playerState.playerSeat}
								/>
							</td>
							<td>&nbsp;</td>
						</tr>
						<tr>
							<td>
								<PlayerContainer
									seatNumber={8}
									buttonPosition={this.props.playerState.buttonPosition}
									player={this.props.playerState.currentPlayers[8]}
									cards={
										this.props.playerState.playerSeat == 8 ? (
											this.props.playerState.playerCards
										) : (
											undefined
										)
									}
									sendMessage={this.props.sendMessage}
									isCurrentTurn={
										this.props.playerState.actingPlayer ? (
											this.props.playerState.actingPlayer.seatNumber === 8
										) : (
											false
										)
									}
									showSitButton={!this.props.playerState.playerSeat}
								/>
							</td>
							<td>&nbsp;</td>
							<td>&nbsp;</td>
							<td>
								<TableInfo
									communityCards={this.props.playerState.communityCards}
									pot={this.props.playerState.pot}
								/>
							</td>
							<td>&nbsp;</td>
							<td>&nbsp;</td>
							<td>
								<PlayerContainer
									seatNumber={3}
									buttonPosition={this.props.playerState.buttonPosition}
									player={this.props.playerState.currentPlayers[3]}
									cards={
										this.props.playerState.playerSeat == 3 ? (
											this.props.playerState.playerCards
										) : (
											undefined
										)
									}
									sendMessage={this.props.sendMessage}
									isCurrentTurn={
										this.props.playerState.actingPlayer ? (
											this.props.playerState.actingPlayer.seatNumber === 3
										) : (
											false
										)
									}
									showSitButton={!this.props.playerState.playerSeat}
								/>
							</td>
						</tr>
						<tr>
							<td>
								<PlayerContainer
									seatNumber={7}
									buttonPosition={this.props.playerState.buttonPosition}
									player={this.props.playerState.currentPlayers[7]}
									cards={
										this.props.playerState.playerSeat == 7 ? (
											this.props.playerState.playerCards
										) : (
											undefined
										)
									}
									sendMessage={this.props.sendMessage}
									isCurrentTurn={
										this.props.playerState.actingPlayer ? (
											this.props.playerState.actingPlayer.seatNumber === 7
										) : (
											false
										)
									}
									showSitButton={!this.props.playerState.playerSeat}
								/>
							</td>
							<td>&nbsp;</td>
							<td>&nbsp;</td>
							<td>&nbsp;</td>
							<td>&nbsp;</td>
							<td>&nbsp;</td>
							<td>
								<PlayerContainer
									seatNumber={4}
									buttonPosition={this.props.playerState.buttonPosition}
									player={this.props.playerState.currentPlayers[4]}
									cards={
										this.props.playerState.playerSeat == 4 ? (
											this.props.playerState.playerCards
										) : (
											undefined
										)
									}
									sendMessage={this.props.sendMessage}
									isCurrentTurn={
										this.props.playerState.actingPlayer ? (
											this.props.playerState.actingPlayer.seatNumber === 4
										) : (
											false
										)
									}
									showSitButton={!this.props.playerState.playerSeat}
								/>
							</td>
						</tr>
						<tr>
							<td>&nbsp;</td>
							<td>&nbsp;</td>
							<td>
								<PlayerContainer
									seatNumber={6}
									buttonPosition={this.props.playerState.buttonPosition}
									player={this.props.playerState.currentPlayers[6]}
									cards={
										this.props.playerState.playerSeat == 6 ? (
											this.props.playerState.playerCards
										) : (
											undefined
										)
									}
									sendMessage={this.props.sendMessage}
									isCurrentTurn={
										this.props.playerState.actingPlayer ? (
											this.props.playerState.actingPlayer.seatNumber === 6
										) : (
											false
										)
									}
									showSitButton={!this.props.playerState.playerSeat}
								/>
							</td>
							<td>&nbsp;</td>
							<td>
								<PlayerContainer
									seatNumber={5}
									buttonPosition={this.props.playerState.buttonPosition}
									player={this.props.playerState.currentPlayers[5]}
									cards={
										this.props.playerState.playerSeat == 5 ? (
											this.props.playerState.playerCards
										) : (
											undefined
										)
									}
									sendMessage={this.props.sendMessage}
									isCurrentTurn={
										this.props.playerState.actingPlayer ? (
											this.props.playerState.actingPlayer.seatNumber === 5
										) : (
											false
										)
									}
									showSitButton={!this.props.playerState.playerSeat}
								/>
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
