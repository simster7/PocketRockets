import * as React from 'react';
import {Action, PlayerState} from "../../api/v1/apis_pb";


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
    active?: boolean;
    playerState?: PlayerState;
}

class PokerClient extends React.Component<IProp, IState> {

    public constructor(props: IProp) {
        super(props);
        this.state = {active: false};
    }

    public render() {
        if (this.state.active === false) {
            return <div>Loading...</div>;
        } else {
            return (
                <div>
                    Poker Game here
                </div>
            );
        }
    }
}

export default PokerClient;
