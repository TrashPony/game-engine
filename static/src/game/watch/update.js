import {gameStore} from "../store";
import {CreateMapBar, createObjectBars} from "../interface/status_layer";
import {Scene} from "../create";
import {intFromBytes} from "../../utils";
import {MoveSprite} from "../utils/move_sprite";

function UpdateObject(type, id, updateMsg) {

  if (type === "unit") {
    UpdateUnit(id, updateMsg);
  }

  if (type === "object" || type === "dynamic_objects") {
    UpdateDynamicObject(id, updateMsg);
    createObjectBars(id);
  }
}

function UpdateDynamicObject(id, update) {

  let obj = gameStore.objects[id];
  if (obj && obj.objectSprite) {
    let oldScale = obj.scale;

    obj.work = intFromBytes(update.slice(0, 1)) === 1;
    obj.scale = intFromBytes(update.slice(1, 2));
    obj.hp = intFromBytes(update.slice(4, 8));
    obj.max_hp = intFromBytes(update.slice(8, 12));
    obj.current_energy = intFromBytes(update.slice(12, 16));
    obj.max_energy = intFromBytes(update.slice(16, 20));
    obj.owner_id = intFromBytes(update.slice(20, 24));

    let xs = intFromBytes(update.slice(2, 3));
    let ys = intFromBytes(update.slice(3, 4));
    let gt = intFromBytes(update.slice(24, 28));

    obj.height = intFromBytes(update.slice(28, 29));
    obj.team_id = intFromBytes(update.slice(29, 30));
    obj.view_range = intFromBytes(update.slice(30, 34));

    let lengthGeoData = intFromBytes(update.slice(34, 35));
    if (lengthGeoData > 0) {

      let stopByte = 35;
      let geoData = [];

      for (; stopByte < update.length;) {

        geoData.push({
          x: intFromBytes(update.slice(stopByte, stopByte + 4)),
          y: intFromBytes(update.slice(stopByte + 4, stopByte + 8)),
          radius: intFromBytes(update.slice(stopByte + 8, stopByte + 12)),
        })

        stopByte += 12
      }

      obj.geo_data = geoData;
    }

    if (!obj.build && oldScale !== obj.scale) {

      obj.objectSprite.setDepth(obj.height);

      if (obj.objectSprite.shadow) {
        MoveSprite(obj.objectSprite.shadow, obj.objectSprite.x + xs, obj.objectSprite.y + ys)
        obj.objectSprite.shadow.setDepth(obj.height - 1);
      }
    }
  }
}

function UpdateUnit(id, update) {
  let unit = gameStore.units[id];
  if (unit) {

    unit.hp = intFromBytes(update.slice(0, 4));
    let new_team_id = intFromBytes(update.slice(4, 8));
    unit.body.range_view = intFromBytes(update.slice(8, 12));
    unit.body.range_radar = intFromBytes(update.slice(12, 16));
    unit.view_mode = intFromBytes(update.slice(16, 17));

    CreateMapBar(unit.sprite, unit.body.max_hp, unit.hp, 10, null, Scene, 'unit', unit.id, 'hp', 50);
  }
}

export {UpdateObject}
