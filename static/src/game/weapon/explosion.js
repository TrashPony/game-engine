import {gameStore} from "../store";
import {Scene} from "../create";
import {GetGlobalPos} from "../map/gep_global_pos";
import {PlayPositionSound} from "../sound/play_sound";

function ExplosionBullet(data) {
  if (!gameStore.gameReady) return;

  let infoBullet = gameStore.gameTypes.ammo[data.type_id];
  let pos = GetGlobalPos(data.x, data.y, gameStore.map.id);
  if (!Scene.cameras.main.worldView.contains(pos.x, pos.y) && !infoBullet.force_animate) {
    return;
  }

  if (infoBullet.name === "piu-piu_2") {
    PlayPositionSound(['explosion_4x3'], null, pos.x, pos.y, false,undefined, 'explosions');
    FireExplosion(pos.x, pos.y, 75, 2, 10, 1, 0, 360, data.z);
    SmokeExplosion(pos.x, pos.y, 75, 1, 10, 1, 0, 360, data.z);
    return;
  }

  SmokeExplosion(pos.x, pos.y, 1000, 20, 15, 15, 0, 360, data.z);
}

function FireExplosion(x, y, time, count, scale, speed, minAngle, maxAngle, z, reverseScale, alpha) {
  Explosion(x, y, time, count, scale, ["yellow"], speed, minAngle, maxAngle, z, 'ADD', reverseScale, alpha)
}

function RedExplosion(x, y, time, count, scale, speed, minAngle, maxAngle, z, reverseScale, alpha) {
  Explosion(x, y, time, count, scale, ["red"], speed, minAngle, maxAngle, z, 'ADD', reverseScale, alpha)
}

function BlueExplosion(x, y, time, count, scale, speed, minAngle, maxAngle, z, reverseScale, alpha) {
  Explosion(x, y, time, count, scale, ["blue"], speed, minAngle, maxAngle, z, 'ADD', reverseScale, alpha)
}

function BlueExplosionSCREEN(x, y, time, count, scale, speed, minAngle, maxAngle, z, reverseScale, alpha) {
  Explosion(x, y, time, count, scale, ["blue"], speed, minAngle, maxAngle, z, 'SCREEN', reverseScale, alpha)
}

function SmokeExplosion(x, y, time, count, scale, speed, minAngle, maxAngle, z, reverseScale, alpha) {
  Explosion(x, y, time, count, scale, ["smoke_puff"], speed, minAngle, maxAngle, z, 'SCREEN', reverseScale, alpha)
}

function GreenExplosion(x, y, time, count, scale, speed, minAngle, maxAngle, z, reverseScale, alpha) {
  Explosion(x, y, time, count, scale, ["green"], speed, minAngle, maxAngle, z, 'ADD', reverseScale, alpha)
}

let effectsParticles;
let explosions = [];

function getEffectsParticles() {
  if (!effectsParticles) {
    effectsParticles = Scene.add.particles('flares');
    effectsParticles.setDepth(400);
  }

  return effectsParticles
}

function Explosion(x, y, time, count, scale, type, speed, minAngle, maxAngle, z, blendMode, reverseScale, alpha) {

  if (reverseScale) {
    scale = {start: scale / 200, end: scale / 25}
  } else {
    scale = {start: scale / 50, end: scale / 100}
  }

  if (!alpha) {
    alpha = {start: 1, end: 0}
  }

  let angle = getRandomArbitrary(minAngle, maxAngle)

  let emitter
  if (explosions.length === 0) {
    emitter = getEffectsParticles().createEmitter({
      frame: type,
      gravityY: 0,
      speed: {min: 0, max: speed},
      scale: scale,
      angle: angle,
      alpha: alpha,
      lifespan: {min: time / 2, max: time},
    });
  } else {
    emitter = explosions.shift()
    emitter.setFrame(type);
    emitter.setSpeed({min: 0, max: speed})
    emitter.setScale(scale);
    emitter.setAngle(angle);
    emitter.setAlpha(alpha);
    emitter.setLifespan({min: time / 2, max: time});
  }

  if (blendMode) {
    emitter.setBlendMode(blendMode);
  } else {
    emitter.setBlendMode(undefined);
  }

  emitter.explode(count, x, y);
  Scene.time.addEvent({
    delay: time,                // ms
    callback: function () {
      emitter.stop()
      explosions.push(emitter);
    },
    loop: false
  });
}

function getRandomArbitrary(min, max) {
  return Math.random() * (max - min) + min;
}

export {
  ExplosionBullet,
  FireExplosion,
  RedExplosion,
  BlueExplosion,
  GreenExplosion,
  SmokeExplosion,
  Explosion,
  BlueExplosionSCREEN,
  getEffectsParticles,
  explosions,
}
