import React from 'react';
import ReactDOM from 'react-dom';
import App from './components/App/App';
import * as serviceWorker from './utils/service-worker/serviceWorker';

function createApp() {
    return new Promise(
        resolve =>
            ReactDOM.render(
                <App/>,
                document.getElementById('root'),
                resolve
            )
    )
}

async function main() {
    await createApp();
    serviceWorker.register();
    return `DONE RENDERING`;
}

main()
    .then(console.log)
    .catch(console.error)
    .finally(
        console.log.bind(
            null,
            'EXITING'
        )
    );