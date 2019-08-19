import * as React from 'react';
import {Action, GetPlayerStateRequest, PlayerState} from "../../api/v1/apis_pb";
import {PokerServiceClient, ServiceError} from "../../api/v1/apis_pb_service";


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

    private client: PokerServiceClient;

    public constructor(props: IProp) {
        super(props);
        this.state = {active: false};
        this.client = new PokerServiceClient("http://localhost:1234", null);
        const savePlayerState = this.savePlayerState.bind(this);
        const request = new GetPlayerStateRequest();
        request.setPlayerid(2);
        request.setGameid(0);
        this.client.getPlayerState(request, (err: ServiceError, message: PlayerState) => {
            console.log(err);
            savePlayerState(message);
        });
    }

    public savePlayerState(playerState: PlayerState): void {
        console.log("playerState got called");
        console.log(playerState);
        this.setState({
            playerState
        })
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
