// config.js
const config = {
    production: {
        websocketUrl: "wss://tervicketactoe.onrender.com/ws",
    },
    development: {
        websocketUrl: "ws://localhost:5000/ws",
    },
};

const environment = 'production'; 

export default config[environment];
