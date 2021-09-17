import {gameStore} from "../store";
import {removeAllObj} from "../radar/object";
import {deleteBullet} from "../bullet/remove";

function RemoveOldMap() {

  for (let i in gameStore.mapsState) {
    removeMap(gameStore.mapsState[i])
  }

  for (let i in gameStore.clouds) {
    gameStore.clouds[i].destroy();
    delete gameStore.clouds[i];
  }

  for (let i in gameStore.bullets) {
    deleteBullet(i)
  }

  for (let i in gameStore.spawns) {
    if (gameStore.spawns.hasOwnProperty(i)) {
      gameStore.spawns[i].destroy();
      delete gameStore.spawns[i];
    }
  }
}

function removeMap(mapState) {
  if (mapState.bmdTerrain && mapState.bmdTerrain.bmd) {
    mapState.bmdTerrain.bmd.clear();
    mapState.bmdTerrain.bmd.destroy();
  }

  removeAllObj();
  for (let obj of mapState.staticObjects) {
    if (obj.objectSprite) {
      if (obj.objectSprite.shadow) obj.objectSprite.shadow.destroy();
      obj.objectSprite.destroy();

      if (obj.objectSprite.emitter) {
        obj.objectSprite.emitter.emitter.stop();
        obj.objectSprite.emitter.destroy();
      }
    }
  }
  mapState.staticObjects = [];
}

export {RemoveOldMap}
