import * as React from 'react';

import {Player} from '../utils/PlayerState'

interface IProp {
    player: Player;
}

interface IState {
    occupied: boolean;

}

class PokerPlayer extends React.Component<IProp, IState> {

    public constructor(props: IProp) {
        super(props)
        if (props.player === undefined) {
            this.state = {
                occupied: false
            }
        } else {
            this.state = {
                occupied: true
            }
        }
    }

    public render() {
        if (!this.state.occupied) {
            return (<div>
                Empty
            </div>)
        }
        return (
            <div>
                {this.props.player.name}: {this.props.player.stack}
            </div>
        )
    }
}

export default PokerPlayer;
