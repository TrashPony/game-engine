import {gameStore} from "../store";
import {GetGlobalPos} from "../map/gep_global_pos";
import {MoveSprite} from "../utils/move_sprite";
import {SetBodyAngle} from "./rotate";
import {Scene} from "../create";
import {AddMoveBufferData} from "../utils/add_move_buffer_data";

function AddUnitMoveBufferData(data) {
  let unit = gameStore.units[data.id];
  if (unit) {
    AddMoveBufferData(data, unit)
  }
}

function MoveTo(unit, ms, view) {
  if (!gameStore.units.hasOwnProperty(unit.id)) return;

  let path = unit.bufferMoveTick.shift();
  if (!path) {
    return;
  }

  unit.x = path.x;
  unit.y = path.y;

  let pos = GetGlobalPos(path.x, path.y, gameStore.map.id);
  MoveSprite(unit.sprite, pos.x, pos.y, ms, null, false, false);
  let needZ = 0;

  let shadowDist = Scene.shadowXOffset * 1.5
  if (unit.body) needZ = unit.body.height + path.z + 1
  SetBodyAngle(unit, path.r, ms, true, shadowDist);

  if (needZ !== unit.sprite.depth) {
    unit.sprite.setDepth(needZ);
  }

  unit.speedDirection = path.d;
  unit.speed = path.s;
  unit.angularVelocity = path.av;
  unit.animateSpeed = path.a;
  unit.rotate = path.r;
}

export {MoveTo, AddUnitMoveBufferData}
