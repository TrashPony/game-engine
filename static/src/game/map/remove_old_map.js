import {gameStore} from "../store";
import {deleteBullet} from "../bullet/remove";
import {removeUnit} from "../unit/remove";
import {Scene} from "../create";
import {removeAllObj, removeDynamicObject} from "../watch/remove";

function RemoveOldMap() {

  for (let i in gameStore.mapsState) {
    removeMap(gameStore.mapsState[i])
  }

  for (let i in gameStore.bullets) {
    deleteBullet(i)
  }

  for (let i in gameStore.units) {
    removeUnit(gameStore.units[i])
  }

  if (Scene && Scene.children) {
    for (let obj of Scene.children.systems.displayList.list) {

      if (!obj) continue;

      if (obj.type === 'ParticleEmitterManager' && obj.name === "bullet") {
        obj.destroy();
      }

      if (obj.frame && obj.frame.texture) {
        if (obj.frame.texture.key.includes("craters") || obj.frame.texture.key.includes("bullets")) {
          obj.destroy();
        }
      }
    }
  }
}

function removeMap(mapState) {
  if (mapState.bmdTerrain && mapState.bmdTerrain.bmd) {
    mapState.bmdTerrain.bmd.clear();
    mapState.bmdTerrain.bmd.destroy();

    mapState.bmdTerrain.bmdObject.clear();
    mapState.bmdTerrain.bmdObject.destroy();
  }

  removeAllObj();
  for (let obj of mapState.staticObjects) {
    if (obj.objectSprite) {
      removeDynamicObject(obj)
    }
  }
  mapState.staticObjects = [];
}

export {RemoveOldMap}
