import {BlueExplosion, FireExplosion, SmokeExplosion} from "../weapon/explosion";
import {GetGlobalPos} from "./gep_global_pos";
import {gameStore} from "../store";
import {PlayPositionSound} from "../sound/play_sound";
import {Scene} from "../create";

function ObjectDead(data) {
  if (!gameStore.gameReady) return;

  let pos = GetGlobalPos(data.x, data.y, gameStore.map.id);

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
  SmokeExplosion(x, y, 500, 3, obj.scale * 3, 50, 0, 360, obj.height + 1)
  PlayPositionSound(['crunch'], null, x, y, false, 0.35);
}

function DestroyBuildObject(obj, startX, startY) {
  for (let i = 0; i < 2; i++) {

    let x = startX + Math.random() * 25;
    let y = startY + Math.random() * 25;

    Scene.time.addEvent({
      delay: 250 * i,
      callback: function () {
        if (i === 0) BlueExplosion(x, y, 1000, 8, obj.scale, 50, 0, 360, obj.height);
        if (i === 1) FireExplosion(x, y, 1000, 8, obj.scale, 50, 0, 360, obj.height);
        SmokeExplosion(x, y, 1500, 4, obj.scale * 2, 50, 0, 360, obj.height - 2);
      },
    });
  }
}

function DestroyExplosionUnit(x, y) {
  BlueExplosion(x, y, 1000, 8, 25, 50, 0, 360, 25);
  FireExplosion(x, y, 1000, 8, 40, 50, 0, 360, 24);
  SmokeExplosion(x, y, 1500, 4, 45, 50, 0, 360, 23);
  PlayPositionSound(['explosion_1', 'explosion_2'], null, x, y);
}

export {ObjectDead}
