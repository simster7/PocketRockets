import * as React from 'react';

import { Card } from '../utils/PlayerState';
import { cardIdToString } from 'src/utils/CardConverter';

interface IProp {
	communityCards?: Card[];
	pot: number;
}

interface IState {}

class TableInfo extends React.Component<IProp, IState> {
	public constructor(props: IProp) {
		super(props);
	}

	public render() {
		return (
			<div>
				{this.props.communityCards ? (
					this.props.communityCards.map((card: Card) => cardIdToString(card.cardId)).join(' ')
				) : (
					''
				)}
				<br />
				Pot: {this.props.pot}
			</div>
		);
	}
}

export default TableInfo;
