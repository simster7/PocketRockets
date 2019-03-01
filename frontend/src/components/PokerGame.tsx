import * as React from 'react';

interface IProp {
    roomName: string;

}

interface IState {

}

class PokerGame extends React.Component<IProp, IState> {
  private connection: WebSocket;

  public componentDidMount(){
    this.connection = new WebSocket('ws://localhost:8000/ws/poker/room/' + this.props.roomName + '/');
    this.connection.onmessage = function(e) {
        let data = JSON.parse(e.data);
        let message = data['message'];

        // document.querySelector('#chat-log').value += (message + '\n');
        // document.querySelector('#chat-log').scrollTop = document.querySelector('#chat-log').scrollHeight;
    };;
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
