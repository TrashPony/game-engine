import {GrabCamera} from './interface/grabCamera'
import {gameStore} from "./store";
import store from '../store/store'
import {CreateAllMaps, MapSize} from "./map/createMap";
import {MoveTo} from "./unit/move";
import {UpdatePosBars} from "./interface/status_layer";
import {FlyBullet} from "./bullet/fly";

let connect = null;

function gameInit(scene) {

  if (location.href.includes('lobby')) {
    gameStore.reload = true;
    gameStore.unitReady = false;
    return
  }

  if (!gameStore.gameDataInit.sendRequest) {
    gameStore.gameDataInit.sendRequest = true;
    store.dispatch("sendSocketData", JSON.stringify({
      event: "InitBattle",
      service: "battle",
    }))
  }

  if (gameStore.reload || !gameStore.gameDataInit.data) return;
  if (!gameStore.gameReady) {
    connect = store.getters.getWSConnectState;

    store.commit({
      type: 'setVisibleLoader',
      visible: true,
      text: 'Строим карту...'
    });

    if (!gameStore.mapDrawing) {
      gameStore.mapDrawing = true;
      setTimeout(function () {
        CreateAllMaps(scene);
        gameStore.gameReady = true;
        gameStore.mapDrawing = false;
      })
    }

    connect.ws.send(JSON.stringify({
      event: "StartLoad",
      service: "battle",
    }));

    store.commit({
      type: 'setVisibleLoader',
      visible: false,
    });
  }
}

function update() {

  gameInit(this)
  if (!gameStore.gameReady) return;

  GrabCamera(this); // функцуия для перетаскивания карты мышкой /* Магия */

  if (!connect || !connect.connect) return;

  for (let id of gameStore.selectUnits) {
    let unit = gameStore.units[id];
    if (unit) {
      unit.selectSprite.setTint(0x00ff00);
      unit.selectSprite.visible = true
    }
  }

  for (let i in gameStore.units) {
    let unit = gameStore.units[i];

    if (!unit.sprite.moveTween || !unit.sprite.moveTween.isPlaying() || unit.bufferMoveTick.length > 0) {
      if (unit.updaterPos) {
        MoveTo(unit, 64);
      }
    }

    if (unit.speed > 0) {
      UpdatePosBars(unit.sprite, unit.body.max_hp, unit.hp, 10, null, this, 'unit', unit.id, 'hp', 50);
    }
  }

  for (let i in gameStore.bullets) {
    let bullet = gameStore.bullets[i];
    if (bullet.bufferMoveTick.length > 0) {
      if (bullet.updaterPos) FlyBullet(bullet)
    }
  }

  // todo только для тестов
  if (gameStore.inputType === "wasd") {
    connect.ws.send(JSON.stringify({
      event: "i",
      service: "battle",
      select_units: gameStore.selectUnits,
      w: this.wasd.up.isDown,
      a: this.wasd.left.isDown,
      s: this.wasd.down.isDown,
      d: this.wasd.right.isDown,
      x: Math.round(this.game.input.mousePointer.worldX - MapSize),
      y: Math.round(this.game.input.mousePointer.worldY - MapSize),
      fire: this.game.input.activePointer.leftButtonDown() && !gameStore.HoldAttackMouse,
    }));
  }
}

export {update}
