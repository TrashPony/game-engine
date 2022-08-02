import {CreateObject} from "../map/objects";
import {gameStore} from "../store";
import {Scene} from "../create";
import {createObjectBars} from "../interface/status_layer";
import {ParseObject} from "../map/create_map";
import {GetCacheBulletSprite} from "../bullet/create";
import {GetGlobalPos} from "../map/gep_global_pos";
import {CreateNewUnit} from "../unit/unit";
import {base64ToArrayBuffer, intFromBytes} from "../../utils";
import {removeDynamicObjects} from "./remove";

function CreateRadarObject(type, object) {
  if (!gameStore.radarWork) return;

  if (type === "unit") {
    CreateNewUnit(object)
  }

  if (type === "object" || type === "dynamic_objects") {
    createDynamicObject(ParseObject(object));
  }

  if (type === "bullet") {
    object = {
      type_id: intFromBytes(object.slice(0, 1)),
      id: intFromBytes(object.slice(1, 5)),
      x: intFromBytes(object.slice(5, 9)),
      y: intFromBytes(object.slice(9, 13)),
      z: intFromBytes(object.slice(13, 17)),
      r: intFromBytes(object.slice(17, 21)),
    }

    let bullet = gameStore.bullets[object.id]
    if (bullet && bullet.sprite) {
      return;
    }

    let pos = GetGlobalPos(object.x, object.y, gameStore.map.id);
    object.x = pos.x;
    object.y = pos.y;

    let infoBullet = gameStore.gameTypes.ammo[object.type_id];

    if (infoBullet && infoBullet.name === "drone_destroy" || !infoBullet) {
      return
    }

    GetCacheBulletSprite(object, infoBullet);
  }
}

function createDynamicObject(object) {
  if (gameStore.objects[object.id]) return;

  if (object.animate) {
    CreateAnimate(object, object.x, object.y, Scene);
  } else {
    CreateObject(object, object.x, object.y, true, Scene);
  }

  gameStore.objects[object.id] = object;
  createObjectBars(object.id)
}

function CreateDynamicObjects(dynamicObjects) {
  removeDynamicObjects();
  for (let id in dynamicObjects) {
    let bytes = base64ToArrayBuffer(dynamicObjects[id])
    createDynamicObject(ParseObject(bytes));
  }
}

export {
  CreateRadarObject,
  CreateDynamicObjects,
  createDynamicObject,
}
