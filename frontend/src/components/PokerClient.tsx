import * as React from 'react';

import {parsePlayerStateString, PlayerState} from '../utils/PlayerState'
import PokerGame from './PokerGame';



interface IProp {
    roomName: string;
}

interface IState {
    playerState?: PlayerState;
}

class PokerClient extends React.Component<IProp, IState> {
  private connection: WebSocket;

  public constructor(props: IProp) {
      super(props)
      this.state = {
          
      }
  }

  public componentDidMount(){
    this.connection = new WebSocket('ws://localhost:8000/ws/poker/room/' + this.props.roomName + '/');
    this.connection.onmessage = function(this: PokerClient, e: any) {
        console.log(e.data);
        let data = e.data;
        //let data = "{'message': {'bet_round': 0, 'lead_player': {'name': 'Jarry', 'stack': 98, 'seat_number': 3, 'folded': False, 'last_action': {'action': 'bet' , 'value': 2}, 'sitting_out': False}, 'acting_player': {'name': 'Simon', 'stack': 99, 'seat_number': 6, 'folded': False, 'last_action': {'action': 'bet' , 'value': 1}, 'sitting_out': False}, 'current_players': [None, None, None, {'name': 'Jarry', 'stack': 98, 'seat_number': 3, 'folded': False, 'last_action': {'action': 'bet' , 'value': 2}, 'sitting_out': False}, None, None, {'name': 'Simon', 'stack': 99, 'seat_number': 6, 'folded': False, 'last_action': {'action': 'bet' , 'value': 1}, 'sitting_out': False}, None, None], 'player_cards': [{'card_id': 14}, {'card_id': 22}], 'community_cards': [], 'end_game': None}}";
        data = data.replace(/'/g, '\"').replace(/True/g, "true").replace(/False/g, "false").replace(/None/g, "null");
        data = JSON.parse(data);
        console.log(typeof data)
        let message = data['message'];
        console.log(data)
        console.log(message)
        let playerState: PlayerState = parsePlayerStateString(message);
        console.log("done parse2")
        console.log(playerState);
        this.setState({playerState})
    }.bind(this);
  }

  public render() {
      if (this.state.playerState === undefined) {
          return <div>Loading...</div>;
      } else {
          return <PokerGame playerState={this.state.playerState} />;
      }
  }
}

export default PokerClient;
