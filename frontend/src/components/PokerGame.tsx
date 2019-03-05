import * as React from 'react';

import {parsePlayerStateString, PlayerState} from '../PlayerState'



interface IProp {
    roomName: string;

}

interface IState {
    playerState?: PlayerState;

}

class PokerGame extends React.Component<IProp, IState> {
  private connection: WebSocket;

  public constructor(props: IProp) {
      super(props)
      this.state = {
          
      }
  }

  public componentDidMount(){
    this.connection = new WebSocket('ws://localhost:8000/ws/poker/room/' + this.props.roomName + '/');
    this.connection.onmessage = function(this: PokerGame, e: any) {
        let test = "{'message': {'bet_round': 0, 'lead_player': {'name': 'Jarry', 'stack': 98, 'seat_number': 3, 'folded': False, 'last_action': {'action': 'bet' , 'value': 2}, 'sitting_out': False}, 'acting_player': {'name': 'Simon', 'stack': 99, 'seat_number': 6, 'folded': False, 'last_action': {'action': 'bet' , 'value': 1}, 'sitting_out': False}, 'current_players': [None, None, None, {'name': 'Jarry', 'stack': 98, 'seat_number': 3, 'folded': False, 'last_action': {'action': 'bet' , 'value': 2}, 'sitting_out': False}, None, None, {'name': 'Simon', 'stack': 99, 'seat_number': 6, 'folded': False, 'last_action': {'action': 'bet' , 'value': 1}, 'sitting_out': False}, None, None], 'player_cards': [{'card_id': 14}, {'card_id': 22}], 'community_cards': [], 'end_game': None}}";
        test = test.replace(/'/g, '\"').replace(/True/g, "true").replace(/False/g, "false").replace(/None/g, "null");
        console.log(test)
        //let data = JSON.parse(e.data);
        let data = JSON.parse(test);
        console.log(data)
        let message = data['message'];
        console.log(message)
        let playerState: PlayerState = parsePlayerStateString(message);
        this.setState({playerState})



        // document.querySelector('#chat-log').value += (message + '\n');
        // document.querySelector('#chat-log').scrollTop = document.querySelector('#chat-log').scrollHeight;
    }.bind(this);
    console.log("done")
    console.log(this.state.playerState);
  }
  
  public render() {
    return (
        <div>
            <table>
                <tbody>
                    <tr>
                        <td>&nbsp;</td>
                        <td>Player 1</td>
                        <td>&nbsp;</td>
                        <td>Player 2</td>
                        <td>&nbsp;</td>
                        <td>Player 3</td>
                        <td>&nbsp;</td>
                    </tr>
                    <tr>
                        <td>Player 9</td>
                        <td>&nbsp;</td>
                        <td>&nbsp;</td>
                        <td>&nbsp;</td>
                        <td>&nbsp;</td>
                        <td>&nbsp;</td>
                        <td>Player 4</td>
                    </tr>
                    <tr>
                        <td>Player 8</td>
                        <td>&nbsp;</td>
                        <td>&nbsp;</td>
                        <td>&nbsp;</td>
                        <td>&nbsp;</td>
                        <td>&nbsp;</td>
                        <td>Player 5</td>
                    </tr>
                    <tr>
                        <td>&nbsp;</td>
                        <td>&nbsp;</td>
                        <td>Player 7</td>
                        <td>&nbsp;</td>
                        <td>Player 6</td>
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
