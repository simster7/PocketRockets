import * as React from 'react';

import {PlayerState} from '../utils/PlayerState'
import PokerPlayer from './Player';



interface IProp {
    playerState: PlayerState;
}

interface IState {
}

class PokerGame extends React.Component<IProp, IState> {

  public constructor(props: IProp) {
      super(props)
      this.state = {
          
      }
  }
  
  public render() {
    console.log(this.props.playerState)
    return (
        <div>
            <table>
                <tbody>
                    <tr>
                        <td>&nbsp;</td>
                        <td><PokerPlayer player={this.props.playerState.currentPlayers[0]} /></td>
                        <td>&nbsp;</td>
                        <td><PokerPlayer player={this.props.playerState.currentPlayers[1]} /></td>
                        <td>&nbsp;</td>
                        <td><PokerPlayer player={this.props.playerState.currentPlayers[2]} /></td>
                        <td>&nbsp;</td>
                    </tr>
                    <tr>
                        <td><PokerPlayer player={this.props.playerState.currentPlayers[8]} /></td>
                        <td>&nbsp;</td>
                        <td>&nbsp;</td>
                        <td>&nbsp;</td>
                        <td>&nbsp;</td>
                        <td>&nbsp;</td>
                        <td><PokerPlayer player={this.props.playerState.currentPlayers[3]} /></td>
                    </tr>
                    <tr>
                        <td><PokerPlayer player={this.props.playerState.currentPlayers[7]} /></td>
                        <td>&nbsp;</td>
                        <td>&nbsp;</td>
                        <td>&nbsp;</td>
                        <td>&nbsp;</td>
                        <td>&nbsp;</td>
                        <td><PokerPlayer player={this.props.playerState.currentPlayers[4]} /></td>
                    </tr>
                    <tr>
                        <td>&nbsp;</td>
                        <td>&nbsp;</td>
                        <td><PokerPlayer player={this.props.playerState.currentPlayers[6]} /></td>
                        <td>&nbsp;</td>
                        <td><PokerPlayer player={this.props.playerState.currentPlayers[5]} /></td>
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
