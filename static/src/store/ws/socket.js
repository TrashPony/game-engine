import router from "../../router/router";
import {gameStore} from "../../game/store";
import {RemoveOldMap} from "../../game/map/remove_old_map";
import {CreateGame} from "../../game/create";
import {BinaryReader} from "./binary_reader";
import {messagesQueue} from "../../game/update";
import {CreateDynamicObjects} from "../../game/watch/create";

export default function createSocketPlugin(WS) {
  return store => {
    WS.binaryType = "arraybuffer";
    WS.onopen = function () {
      console.log("Connection chat opened..." + this.readyState);
      store.commit({
        type: 'setWSConnectState', connect: true, error: false, reconnect: false, pending: false, ws: WS,
      });
    };

    WS.onclose = function (i) {

      store.commit({
        type: 'setWSConnectState', connect: false, error: true, pending: false, ws: WS,
      });

      if (!WS.noRedirect) {
        if (location.href.includes('lobby') || location.href.includes('global') || location.href.includes('gate')) {
          router.push('/login');
        }

        console.log("socket close" + i);
      } else {
        console.log("socket close soft");
      }
    };

    WS.onmessage = function (msg) {
      if (msg.data instanceof ArrayBuffer) {

        if (gameStore.exitTab) {
          return;
        }

        if (gameStore.gameReady) {
          messagesQueue.push(new Uint8Array((msg.data)));
        } else {
          BinaryReader(new Uint8Array((msg.data), store))
        }
      } else {
        WSReader(JSON.parse(msg.data), store, WS)
      }
    };
  };
}

function WSReader(response, store, ws) {
  //console.log(response)

  //if (response.hasOwnProperty('event')) logMsg(response.event, response)
  //if (response.hasOwnProperty('e')) logMsg(response.e, response)
  if (response.hasOwnProperty("error")) {
    addError(response.event + ": " + response.error, store);
    return;
  }

  if (response.event === "gd") {
    store.commit({
      type: 'setHandBook', description_items: response.data.description_items,
    });

    if (response.data.user_interface === "") return;
    store.commit({
      type: 'setInterfaceState',
      user_interface: JSON.parse(response.data.user_interface),
      allow_window_save: response.data.allow_windows,
    });

    if (response.data.settings) {
      store.commit({
        type: 'setGameSettings', settings: response.data.settings,
      });

      gameStore.gameSettings = response.data.settings
    }

    store.commit({
      type: 'setInitGame', init: true,
    });
  }

  if (response.event === "GetPlayers") {
    store.commit({
      type: 'setGameUserSate',
      player: response.data.player,
      gameUser: response.data.game_user,
      userPlayers: response.data.players,
    });
    gameStore.player = response.data.player;

    store.commit({
      type: 'setServerTime', data: response.data.st,
    });
  }

  if (response.event === "RefreshDynamicObj") {
    CreateDynamicObjects(response.data);
  }

  if (response.event === "setGameTypes") {
    gameStore.gameTypes = response.data
    gameStore.unitSpriteSize = response.data.unit_sprite_size
    gameStore.mapBinItems = response.data.map_bin_items
    gameStore.mapBinItems = Object.assign({}, ...Object.entries(gameStore.mapBinItems).map(([a, b]) => ({[b]: a})))
  }

  if (response.event === "CreateLobbySession") {
    store.commit({
      type: 'setLobbyState', state: response.data.lobby_session,
    });
    store.commit({
      type: 'setShortInfoMaps', maps: response.data.maps,
    });
  }

  if (response.event === "ToBattle") {
    CreateGame();
    if (!location.href.includes('global')) router.push('/global');
  }

  if (response.event === "ToLobby") {
    if (!location.href.includes('lobby')) router.push('/lobby');
  }

  if (response.event === "InitBattle") {

    if (response.hasOwnProperty("error")) return;
    RemoveOldMap();

    gameStore.gameReady = false;
    gameStore.unitReady = false;

    response.data = JSON.parse(response.data)

    gameStore.maps = response.data.maps;
    gameStore.player = response.data.player;
    gameStore.playerNames = response.data.player_names;

    store.commit({
      type: 'setPlayerRole', role: response.data.spectrum ? 'Spectrum' : '',
    });

    gameStore.BattleData = response.data.battle
    gameStore.gameDataInit.data = true;
  }

  if (response.event === "gsgs") {
    console.log(response.data)
    if (response.data.ig) {
      CreateGame();
      router.push('/global');
    }
  }
}

function addError(error, store) {
  store.commit({
    type: 'addNotification', id: error, removeSec: 5, html: `<span class="importantly">${error}</span>`,
  });
}
