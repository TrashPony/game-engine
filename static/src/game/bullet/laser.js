import {gameStore} from "../store";
import {GetGlobalPos} from "../map/gep_global_pos";
import {Scene} from "../create";
import {getEffectsParticles, SmokeExplosion} from "../weapon/explosion";

function FlyLaser(data) {
  if (!gameStore.gameReady) return;

  let pos = GetGlobalPos(data.x, data.y, gameStore.map.id);
  let targetPos = GetGlobalPos(data.to_x, data.to_y, gameStore.map.id);
  if (!Scene.cameras.main.worldView.contains(pos.x, pos.y) && !Scene.cameras.main.worldView.contains(targetPos.x, targetPos.y)) {
    return;
  }

  let color = ["blue"];
  let scale = 0.13;
  let quantity = 128;
  let lifeSpan = 300

  let line = new Phaser.Geom.Line(pos.x, pos.y, targetPos.x, targetPos.y);
  let emitter = getEffectsParticles().createEmitter({
    frame: color,
    x: 0,
    y: 0,
    scale: {start: scale * 2, end: 0},
    alpha: {start: 1, end: 0},
    speed: {min: -5, max: 5},
    quantity: quantity / 4,
    emitZone: {source: line},
    blendMode: 'SCREEN',
    lifespan: lifeSpan,
  });

  SmokeExplosion(targetPos.x, targetPos.y, 500, 15, 5, 10, null, null, 40);
  emitter.explode((quantity * 4) / 3);
  emitter.stop()

  Scene.time.addEvent({
    delay: lifeSpan,
    callback: function () {
      getEffectsParticles().emitters.remove(emitter)
    },
    loop: false
  });
}

export {FlyLaser}
