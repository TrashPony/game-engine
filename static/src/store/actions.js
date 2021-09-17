const actions = {
  sendSocketData({commit, getters}, payload) {

    let reconnect = false;
    let waitConnect = setInterval(function () {
      if (getters.getWSConnectState.connect) {
        getters.getWSConnectState.ws.send(payload);
        clearInterval(waitConnect)
      } else {
        // todo отвал по таймауту или кол-во попыток
        if (reconnect) {
          commit({type: 'reconnectWS'});
        } else {
          reconnect = true
        }
      }
    }, 100)
  }
};


export default actions;
