import router from "../../router/router";
import {gameStore} from "../../game/store";
import {RemoveOldMap} from "../../game/map/remove_old_map";
import {CreateGame} from "../../game/create";
import {RadarWork} from "../../game/radar/work";
import {CreateDynamicObjects} from "../../game/radar/object";
import {AddUnitMoveBufferData} from "../../game/unit/move";
import {RotateUnitGun} from "../../game/unit/rotate";
import {AddBulletMoveBufferData, FlyLaser} from "../../game/bullet/fly";
import {FireWeapon} from "../../game/weapon/fire";
import {ExplosionBullet} from "../../game/weapon/explosion";

export default function createSocketPlugin(WS) {
  return store => {

    WS.onopen = function () {
      console.log("Connection chat opened..." + this.readyState);

      store.commit({
        type: 'setWSConnectState',
        connect: true,
        error: false,
        reconnect: false,
        ws: WS,
      });
    };

    WS.onclose = function (i) {
      store.commit({
        type: 'setWSConnectState',
        connect: false,
        error: true,
        ws: WS,
      });

      if (location.href.includes('lobby') || location.href.includes('global') || location.href.includes('gate')) {
        router.push('/login');
      }

      console.log("socket закрыт" + i);
    };

    WS.onmessage = function (msg) {
      ChatReader(JSON.parse(msg.data), store, WS)
    };
  };
}


function ChatReader(response, store, ws) {
  //console.log(response)
  if (response.hasOwnProperty("error")) alert(response.event + ": " + response.error);

  if (response.event === "GetPlayers") {
    store.commit({
      type: 'setGameUserSate',
      player: response.data.player,
      gameUser: response.data.game_user,
      userPlayers: response.data.players,
    });
  }

  if (response.event === "setGameTypes") {
    gameStore.gameTypes = response.data
    gameStore.unitSpriteSize = response.data.unit_sprite_size
  }

  if (response.event === "CreateLobbySession") {
    store.commit({
      type: 'setLobbyState',
      state: response.data.lobby_session,
    });
    store.commit({
      type: 'setShortInfoMaps',
      maps: response.data.maps,
    });
  }

  if (response.event === "ToBattle") {
    CreateGame();
    router.push('/global');
  }

  if (response.event === "InitBattle") {

    gameStore.gameReady = false;
    gameStore.unitReady = false;

    RemoveOldMap();

    gameStore.maps = response.data.maps;
    gameStore.player = response.data.player;

    store.commit({
      type: 'setPlayerRole',
      role: response.data.spectrum ? 'Spectrum' : '',
    });

    gameStore.gameDataInit.data = true;
  }

  if (response.event === "RefreshDynamicObj") {
    CreateDynamicObjects(response.data);
  }

  if (response.event === "rw") {
    RadarWork(response.data.ev)
  }

  if (response.event === "um") {
    for (let path of response.data) {
      AddUnitMoveBufferData(path)
    }
  }

  if (response.event === "uwr") {
    for (let rotatePath of response.data) {
      RotateUnitGun(rotatePath.id, rotatePath.slot_number, rotatePath.rotate, rotatePath.ms);
    }
  }

  if (response.event === "fb") {
    for (let flyBulletPath of response.data) {
      AddBulletMoveBufferData(flyBulletPath)
    }
  }

  if (response.event === "fl") {
    for (let flyLaserPath of response.data) {
      FlyLaser(flyLaserPath)
    }
  }

  if (response.event === "ufw") {
    for (let fireWeapon of response.data) {
      FireWeapon(fireWeapon)
    }
  }

  if (response.event === "eb") {
    for (let explosion of response.data) {
      ExplosionBullet(explosion)
    }
  }
}
