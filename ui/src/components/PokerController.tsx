import * as React from 'react';

import { Message } from './PokerClient';

interface IProp {
	sendMessage: (message: Message) => void;
}

interface IState {
	betAmount: number;
	playerName: string;
	addStackAmount: number;
	seatNumber?: number;
}

class PokerController extends React.Component<IProp, IState> {
	public constructor(props: IProp) {
		super(props);
		this.state = {
			betAmount: 0,
			playerName: '',
			addStackAmount: 0,
			seatNumber: undefined
		};
	}

	public render() {
		return (
			<table>
				<tbody>
					<tr>
						<td>
							<button
								onClick={() => this.props.sendMessage({ type: 'action', value: { action: 'fold' } })}
							>
								Fold
							</button>
						</td>
						<td>
							<button
								onClick={() => this.props.sendMessage({ type: 'action', value: { action: 'check' } })}
							>
								Check
							</button>
						</td>
						<td>
							<input
								value={this.state.playerName}
								onChange={(evt) => this.updatePlayerName(evt.target.value)}
								type="text"
							/>
							<button
								onClick={() =>
									this.props.sendMessage({
										type: 'manage',
										value: {
											item: 'player_name',
											value: this.state.playerName
										}
									})}
							>
								Player name
							</button>
						</td>
					</tr>
					<tr>
						<td>
							<button
								onClick={() => this.props.sendMessage({ type: 'action', value: { action: 'call' } })}
							>
								Call
							</button>
						</td>
						<td>&nbsp;</td>
						<td>
							<input
								value={this.state.seatNumber}
								onChange={(evt) => this.updateSeatNumber(+evt.target.value)}
								type="text"
								disabled={true}
							/>
							<button
								onClick={() =>
									this.props.sendMessage({
										type: 'manage',
										value: {
											item: 'sit',
											value: '' + this.state.seatNumber
										}
									})}
								disabled={true}
							>
								Sit
							</button>
						</td>
					</tr>
					<tr>
						<td>
							<input
								value={this.state.betAmount}
								onChange={(evt) => this.updateBetAmount(+evt.target.value)}
								type="text"
							/>
							<button
								onClick={() =>
									this.props.sendMessage({
										type: 'action',
										value: { action: 'bet', value: this.state.betAmount }
									})}
							>
								Bet
							</button>
						</td>
						<td>
							<button
								onClick={() =>
									this.props.sendMessage({
										type: 'manage',
										value: {
											item: 'deal',
											value: ''
										}
									})}
							>
								Deal
							</button>
						</td>
						<td>
							<input
								value={this.state.addStackAmount}
								onChange={(evt) => this.updateAddStackAmount(+evt.target.value)}
								type="text"
							/>
							<button
								onClick={() =>
									this.props.sendMessage({
										type: 'manage',
										value: {
											item: 'add_stack',
											value: '' + this.state.addStackAmount
										}
									})}
							>
								Add to stack
							</button>
						</td>
					</tr>
				</tbody>
			</table>
		);
	}

	private updateBetAmount(betAmount: number) {
		this.setState({ betAmount });
	}

	private updatePlayerName(playerName: string) {
		this.setState({ playerName });
	}

	private updateAddStackAmount(addStackAmount: number) {
		this.setState({ addStackAmount });
	}

	private updateSeatNumber(seatNumber: number) {
		this.setState({ seatNumber });
	}
}

export default PokerController;
