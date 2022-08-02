import {Grab_camera} from './interface/grab_camera'
import {gameStore} from "./store";
import store from '../store/store'
import {CreateAllMaps, MapSize} from "./map/create_map";
import {MoveTo} from "./unit/move";
import {UpdatePosBars} from "./interface/status_layer";
import {FlyBullet} from "./bullet/fly";
import {ServerGunDesignator} from "./interface/server_gun_designator";
import {Scene} from "./create";
import {BinaryReader} from "../store/ws/binary_reader";
import {ObjectTo} from "./map/move_object";
import {rotatePoint} from "./utils/rotate_sprite";
import {initMiniMap, setPositionMapRectangle} from "./interface/mini_map";

let connect = null;
let userUnit = false;

let messagesQueue = [];

function parseMsg() {
  while (messagesQueue.length > 0) {
    BinaryReader(messagesQueue.shift(), store)
  }
}

function gameInit(scene) {

  if (location.href.includes('lobby')) {
    gameStore.reload = true;
    gameStore.unitReady = false;
    return
  }

  if (!gameStore.gameDataInit.sendRequest) {
    gameStore.gameDataInit.sendRequest = true;
    gameStore.radarWork = false;
    store.dispatch("sendSocketData", JSON.stringify({
      event: "InitBattle", service: "battle",
    }))
  }

  if (gameStore.reload || !gameStore.gameDataInit.data) return;
  if (!gameStore.gameReady) {
    connect = store.getters.getWSConnectState;

    store.commit({
      type: 'setVisibleLoader', visible: true, text: 'Строим карту...'
    });

    if (!gameStore.mapDrawing) {
      gameStore.mapDrawing = true;
      setTimeout(function () {
        CreateAllMaps(scene);
        gameStore.gameReady = true;
        gameStore.mapDrawing = false;

        connect.ws.send(JSON.stringify({
          event: "StartLoad", service: "battle",
        }));
      })
    }

    store.commit({
      type: 'setVisibleLoader', visible: false,
    });
  }
}

function update() {
  parseMsg()

  if (gameStore.exitTab) {
    return;
  }

  if (!gameStore.gameReady) {
    gameInit(this)
    return;
  }

  if (!gameStore.gameSettings.follow_camera || !userUnit) Grab_camera(this);
  if (!connect || !connect.connect) return;

  processUnit(this)
  processBullets(this)
  processObjects(this)
  sendInputsToBack(this)
  initMiniMap()
  setPositionMapRectangle();
}

function processUnit(scene) {
  let findUserUnit = false
  for (let i in gameStore.units) {
    let unit = gameStore.units[i];
    let view = Scene.cameras.main.worldView.contains(unit.sprite.x, unit.sprite.y)

    if (!unit.sprite.moveTween || !unit.sprite.moveTween.isPlaying() || unit.bufferMoveTick.length > 0) {
      if (unit.updaterPos) {
        MoveTo(unit, 64, view);
      }
    }

    processUnitWeapon(unit)
    if (gameStore.player.id === unit.owner_id) {
      findUserUnit = unit
    }

    UpdatePosBars(unit.sprite, unit.body.max_hp, unit.hp, 10, null, scene, 'unit', unit.id, 'hp', 50);
  }

  userUnit = findUserUnit
  if (!findUserUnit) {
    gameStore.unitReady = false
    ServerGunDesignator({hide: true})
  }
}

function processUnitWeapon(unit) {
  for (let i in unit.weapons) {

    let newPos = rotatePoint(unit.weapons[i].weapon.xAttach, unit.weapons[i].weapon.yAttach, unit.sprite.angle)

    unit.weapons[i].weapon.setPosition(unit.sprite.x + (newPos.x * unit.sprite.scale), unit.sprite.y + (newPos.y * unit.sprite.scale));
    if (unit.weapons[i].weapon.scale !== unit.sprite.scale) unit.weapons[i].weapon.setScale(unit.sprite.scale)
    if (unit.weapons[i].weapon.depth !== unit.sprite.depth + 1) unit.weapons[i].weapon.setDepth(unit.sprite.depth + 1)
  }
}

function processBullets() {
  for (let i in gameStore.bullets) {
    let bullet = gameStore.bullets[i];
    if (bullet.bufferMoveTick.length > 0 && bullet.updaterPos) {
      FlyBullet(bullet)
    }
  }
}

function processObjects(scene) {
  for (let i in gameStore.objects) {
    let obj = gameStore.objects[i];
    if (!obj) continue;

    if (!obj.moveTween || !obj.moveTween.isPlaying() || obj.bufferMoveTick.length > 0) {
      if (obj.updaterPos) ObjectTo(obj)
    }

    if (obj.max_energy > 0) {
      UpdatePosBars(obj.objectSprite, obj.max_energy / 100, obj.current_energy / 100, 7,
        0x00ffd6, Scene, 'object', obj.id, 'energy', 5);
    }
  }
}

function sendInputsToBack(scene) {
  if (gameStore.inputType === "wasd") {

    let leftKeyDown = false;
    let rightKeyDown = false;
    let upKeyDown = false;
    let downKeyDown = false;
    let fire = false
    let pos = {};

    let getXY = function (pointer, yOffset = 0) {

      yOffset = yOffset * (1 / Scene.cameras.main.zoom)

      let x = Math.round((Scene.cameras.main.worldView.left + (pointer.x / Scene.cameras.main.zoom)) - MapSize);
      let y = Math.round((Scene.cameras.main.worldView.top + ((pointer.y + yOffset) / Scene.cameras.main.zoom)) - MapSize);

      return {x: x, y: y, worldX: pointer.x, worldY: pointer.y + yOffset}
    }

    let getFire = function (pointer) {
      return (pointer.leftButtonDown() && (!gameStore.HoldAttackMouse && !gameStore.MouseMoveInterface))
    }

    pos = getXY(Scene.game.input.activePointer);
    fire = getFire(Scene.game.input.activePointer);

    gameStore.oldTargetPos = pos
    connect.ws.send(JSON.stringify({
      event: "i",
      service: "battle",
      select_units: gameStore.selectUnits,
      w: scene.wasd.up.isDown || upKeyDown,
      a: scene.wasd.left.isDown || leftKeyDown,
      s: scene.wasd.down.isDown || downKeyDown,
      d: scene.wasd.right.isDown || rightKeyDown,
      x: pos.x,
      y: pos.y,
      fire: fire,
    }));
  }
}

function checkViewMode(objOwnerID) {
  return !userUnit || (userUnit.view_mode === 0 || userUnit.owner_id === objOwnerID);
}

export {update, messagesQueue, userUnit}
