import {urls} from "../const";
import createSocketPlugin from "./ws/socket";
import store from './store'
import Vue from 'vue'

const mutations = {
  reconnectWS() {
    const WS = new WebSocket(urls.socketURL);
    let plugin = createSocketPlugin(WS);
    plugin(store);
  },
  setWSConnectState(state, payload) {
    state.wsConnectState.connect = payload.connect;
    state.wsConnectState.error = payload.error;
    state.wsConnectState.ws = payload.ws;
  },
  setVisibleLoader(state, payload) {
    Vue.set(state.EndLoad, 'visible', payload.visible);
    Vue.set(state.EndLoad, 'text', payload.text);
  },
  setGameUserSate(state, payload) {
    state.GameUser = payload.gameUser;
    state.UserPlayers = payload.userPlayers;
    state.CurrentPlayer = payload.player;
  },
  setLobbyState(state, payload) {
    state.LobbyState = payload.state
  },
  setShortInfoMaps(state, payload) {
    state.ShortInfoMaps = payload.maps
  },
  setPlayerRole(state, payload) {
    Vue.set(state.Game, 'role', payload.role);
  },
  /** WINDOWS MANAGER **/
  setWindowState(state, payload) {

    if (!state.Interface.state) {
      state.Interface.state = {}
    }

    if (!state.Interface.state[state.Interface.resolution]) {
      state.Interface.state[state.Interface.resolution] = {};
    }

    state.Interface.state[state.Interface.resolution][payload.id] = payload.state;

    if (payload.state.open) {
      this.commit({
        type: 'setWindowZIndex',
        id: payload.id,
      });
    } else {
      if (!state.Interface.allowIDs.hasOwnProperty(payload.id)) {
        delete state.Interface.state[state.Interface.resolution][payload.id]
      }
    }
  },
  setWindowZIndex(state, payload) {
    let windowIndex = state.Interface.openQueue.indexOf(payload.id);
    if (windowIndex > -1) {
      state.Interface.openQueue.splice(windowIndex, 1);
    }

    /** это немного не повьюшному но вот как то так, это позволяет выводить открываемые/выделяемые окна на передний план **/
    if (state.Interface.state && state.Interface.state[state.Interface.resolution]) {
      for (let id in state.Interface.state[state.Interface.resolution]) {
        if (state.Interface.state[state.Interface.resolution].hasOwnProperty(id)) {
          let modal = document.getElementById(id);
          if (modal) {
            let z = state.Interface.openQueue.indexOf(id);
            if (z < 1) {
              z = 1;
            }
            modal.style.zIndex = z;
          }
        }
      }
    }

    state.Interface.openQueue.push(payload.id)
    if (state.Interface.openQueue.length > 15) {
      state.Interface.openQueue.shift();
    }

    let modal = document.getElementById(payload.id);
    if (modal) modal.style.zIndex = state.Interface.openQueue.indexOf(payload.id);
  },
  toggleWindow(state, payload) {

    if (state.NeedOpenComponents.hasOwnProperty(payload.id)) {
      payload.open = !state.NeedOpenComponents[payload.id].open

      // если окно на заднем плане то делаем его на передний
      if (state.Interface.openQueue[state.Interface.openQueue.length - 1] !== payload.id) {
        payload.open = true;
      }
    } else {
      payload.open = true;
    }

    if (payload.forceOpen) {
      payload.open = true;
    }

    if (payload.forceClose) {
      payload.open = false;
    }

    let window = {
      id: payload.id,
      component: payload.component,
      open: payload.open,
      meta: payload.meta,
    };

    if (payload.open) {
      this.commit({
        type: 'setWindowZIndex',
        id: payload.id,
      });
    } else {
      state.Interface.openQueue.splice(state.Interface.openQueue.indexOf(payload.id), 1);
    }

    Vue.set(state.NeedOpenComponents, payload.id, window);

    if (!payload.open && !state.Interface.allowIDs.hasOwnProperty(payload.id)) {
      Vue.delete(state.NeedOpenComponents, payload.id)
    }
  },
  removeAllWindowsByComponentName(state, payload) {
    for (let i in state.NeedOpenComponents) {
      if (state.NeedOpenComponents.hasOwnProperty(i) && state.NeedOpenComponents[i].component === payload.component) {
        this.commit({
          type: 'toggleWindow',
          id: state.NeedOpenComponents[i].id,
          component: state.NeedOpenComponents[i].component,
          forceClose: true,
        });
      }
    }
  },
  setInterfaceState(state, payload) {

    Vue.set(state.Interface, 'state', payload.user_interface);
    Vue.set(state.Interface, 'allowIDs', payload.allow_window_save);

    if (payload.user_interface && payload.user_interface[state.Interface.resolution]) {
      for (let id in payload.user_interface[state.Interface.resolution]) {

        if (!payload.user_interface[state.Interface.resolution].hasOwnProperty(id)) continue;

        Vue.set(state.NeedOpenComponents, id, {
          id: id,
          open: payload.user_interface[state.Interface.resolution][id].open
        });
      }
    }
  },
};

export default mutations;
