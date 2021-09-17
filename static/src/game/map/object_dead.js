import {BlueExplosion, ExplosionRing, FireExplosion, SmokeExplosion} from "../weapon/explosion";
import {GetGlobalPos} from "./gep_global_pos";
import {gameStore} from "../store";
import {PlayPositionSound} from "../sound/play_sound";

function ObjectDead(data) {
  let pos = GetGlobalPos(data.x, data.y, data.m);

  if (data.t === "object") {
    let obj = gameStore.objects[data.id];

    if (!obj) obj = gameStore.removeObjects[data.id]

    if (obj) {
      if (obj.build) {
        PlayPositionSound(['explosion_1', 'explosion_2', 'explosion_3'], null, pos.x, pos.y);
        DestroyBuildObject(obj, pos.x, pos.y)
      } else {
        DestroyDynamicObject(obj, pos.x, pos.y)
      }
    }
  }

  if (data.t === "unit") {
    DestroyExplosionUnit(pos.x, pos.y)
  }
}

function DestroyDynamicObject(obj, x, y) {
  SmokeExplosion(x, y, 500, 3, obj.scale * 3, 50, 0, 360, obj.height)
  PlayPositionSound(['explosion_3x2'], null, x, y);
}

function DestroyBuildObject(obj, startX, startY) {
  for (let i = 0; i < 2; i++) {

    let x = startX + Math.random() * 25;
    let y = startY + Math.random() * 25;

    setTimeout(function () {
      ExplosionRing(x, y, 250, 50, obj.height);
      if (i === 0) BlueExplosion(x, y, 1000, 8, obj.scale, 50, 0, 360, obj.height);
      if (i === 1) FireExplosion(x, y, 1000, 8, obj.scale, 50, 0, 360, obj.height);
      SmokeExplosion(x, y, 1500, 4, obj.scale * 2, 50, 0, 360, obj.height - 2);
    }, 250 * i)
  }
}

function DestroyExplosionUnit(x, y) {
  ExplosionRing(x, y, 250, 50, 25);
  BlueExplosion(x, y, 1000, 8, 25, 50, 0, 360, 25);
  FireExplosion(x, y, 1000, 8, 40, 50, 0, 360, 24);
  SmokeExplosion(x, y, 1500, 4, 45, 50, 0, 360, 23);
  PlayPositionSound(['explosion_1', 'explosion_2'], null, x, y);
}

export {ObjectDead}
