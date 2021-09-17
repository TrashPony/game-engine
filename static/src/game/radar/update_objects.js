import {gameStore} from "../store";
import {CreateMapBar, createObjectBars} from "../interface/status_layer";
import {MoveSprite} from "../utils/move_sprite";

function UpdateObject(mark, updateMsg) {

  // if (mark.to === "unit") {
  //   UpdateUnit(mark.id, updateMsg);
  // }

  if (mark.to === "dynamic_objects") {
    UpdateDynamicObject(mark, updateMsg);
    createObjectBars(mark.id);
  }
}

function UpdateDynamicObject(mark, object) {

  if (gameStore.objects[mark.id]) {

    let obj = gameStore.objects[mark.id];

    if (object.gd) {
      for (let i in object.gd) {
        object.gd[i] = JSON.parse(object.gd[i])
      }
      obj.geo_data = object.gd;
    }

    if (object.s) obj.scale = object.s;
    if (object.hp) obj.hp = object.hp;
    if (object.mhp) obj.max_hp = object.mhp;
    if (object.c) obj.complete = object.c;
    if (object.rr) obj.range_radar = object.rr;
    if (object.ce) obj.current_energy = object.ce;
    if (object.me) obj.max_energy = object.me;
    if (object.ow) obj.owner_id = object.ow;
    if (object.w !== undefined) obj.work = Boolean(object.w)

    if (!obj.build) {
      if (!obj.build) {
        if (object.h) obj.objectSprite.setDepth(object.h);
      }

      if (obj.objectSprite.shadow) {

        if (!obj.build) {
          if (object.h) obj.objectSprite.shadow.setDepth(object.h - 2);
        }

        if (object.s && object.gt) {
          MoveSprite(obj.objectSprite, obj.objectSprite.x, obj.objectSprite.y, object.gt, object.s / 100);
        }

        if (object.xs && object.ys && object.gt && object.s) {
          MoveSprite(obj.objectSprite.shadow, obj.objectSprite.x + object.xs, obj.objectSprite.y + object.ys, object.gt, object.s / 100);
        }
      }
    }
  }
}

export {UpdateObject}
