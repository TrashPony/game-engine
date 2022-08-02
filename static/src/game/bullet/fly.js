import {gameStore} from "../store";
import {GetGlobalPos} from "../map/gep_global_pos";
import {MoveSprite} from "../utils/move_sprite";
import {AddMoveBufferData} from "../utils/add_move_buffer_data";

function AddBulletMoveBufferData(data) {
  if (!gameStore.gameReady || data.m === 0) return;

  let bullet = gameStore.bullets[data.id];
  if (!bullet) {
    bullet = {}
    gameStore.bullets[data.id] = bullet;
  }

  let pos = GetGlobalPos(data.x, data.y, gameStore.map.id);
  data.x = pos.x;
  data.y = pos.y;

  AddMoveBufferData(data, bullet)
}

function FlyBullet(bullet) {

  let data = bullet.bufferMoveTick.shift();
  if (!data) {
    return;
  }

  if (!bullet.sprite) {
    return
  }

  MoveSprite(bullet.sprite, data.x, data.y, data.ms, null, false, false);
}

export {FlyBullet, AddBulletMoveBufferData}
