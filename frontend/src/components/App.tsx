import * as React from 'react';
import PokerClient from "./PokerClient";

class App extends React.Component {
    render() {
        return (<div>
                <PokerClient roomName="test" />
            </div>
        );
    }
}

export default App;