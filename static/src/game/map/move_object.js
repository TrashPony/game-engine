import {gameStore} from "../store";
import {AddMoveBufferData} from "../utils/add_move_buffer_data";
import {GetGlobalPos} from "./gep_global_pos";
import {MoveSprite} from "../utils/move_sprite";
import {ShortDirectionRotateTween} from "../utils/rotate_sprite";

function AddObjectMoveBufferData(data) {
  let obj = gameStore.objects[data.id];
  if (obj) {
    AddMoveBufferData(data, obj)
  }
}

function ObjectTo(object) {
  if (!object || !object.objectSprite) return;

  let data = object.bufferMoveTick.shift();
  if (!data) {
    return;
  }

  let pos = GetGlobalPos(data.x, data.y, gameStore.map.id);
  data.x = pos.x;
  data.y = pos.y;

  MoveSprite(object.objectSprite, data.x, data.y, data.ms);
  ShortDirectionRotateTween(object.objectSprite, data.r, data.ms);

  if (object.objectSprite.weapon) {
    MoveSprite(object.objectSprite.weapon, data.x, data.y, data.ms);
  }

  if (object.objectSprite.weaponShadow) {
    MoveSprite(object.objectSprite.weaponShadow, data.x + object.objectSprite.xShadow, data.y + object.objectSprite.yShadow, data.ms);
  }

  if (object.objectSprite.shadow) {
    MoveSprite(object.objectSprite.shadow, data.x + object.x_shadow_offset, data.y + object.y_shadow_offset, data.ms);
    ShortDirectionRotateTween(object.objectSprite.shadow, data.r, data.ms);
  }
}

export {AddObjectMoveBufferData, ObjectTo}
