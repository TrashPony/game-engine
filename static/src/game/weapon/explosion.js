import {gameStore} from "../store";
import {Scene} from "../create";
import {GetGlobalPos} from "../map/gep_global_pos";
import {PlayPositionSound} from "../sound/play_sound";

function ExplosionBullet(data) {
  if (!gameStore.gameReady) return;

  let infoBullet = gameStore.gameTypes.ammo[data.type_id];
  let pos = GetGlobalPos(data.x, data.y, data.m);

  if (infoBullet.name === "big_missile_bullet" || infoBullet.name === "aim_big_missile_bullet") {
    PlayPositionSound(['explosion_1x2', 'explosion_5x2'], null, pos.x, pos.y);
    ExplosionRing(pos.x, pos.y, 250, infoBullet.area_covers, data.z);
    FireExplosion(pos.x, pos.y, 500, 4, 40, 25, 0, 360, data.z);
    SmokeExplosion(pos.x, pos.y, 500, 4, 40, 25, 0, 360, data.z);
    return
  }

  if (infoBullet.name === "small_missile_bullet" || infoBullet.name === "aim_small_missile_bullet") {
    PlayPositionSound(['explosion_1x2', 'explosion_5x2'], null, pos.x, pos.y);
    ExplosionRing(pos.x, pos.y, 250, infoBullet.area_covers, data.z);
    FireExplosion(pos.x, pos.y, 500, 4, 25, 25, 0, 360, data.z);
    SmokeExplosion(pos.x, pos.y, 500, 4, 25, 25, 0, 360, data.z);
    return
  }

  if (infoBullet.name === "piu-piu") {
    PlayPositionSound(['explosion_3x2'], null, pos.x, pos.y);
    FireExplosion(pos.x, pos.y, 250, 15, 20, 10, 0, 360, data.z);
    SmokeExplosion(pos.x, pos.y, 250, 15, 20, 10, 0, 360, data.z);
    return;
  }

  if (infoBullet.name === "piu-piu_2") {
    PlayPositionSound(['explosion_4x3'], null, pos.x, pos.y);
    FireExplosion(pos.x, pos.y, 75, 2, 10, 1, 0, 360, data.z);
    SmokeExplosion(pos.x, pos.y, 75, 1, 10, 1, 0, 360, data.z);
    return;
  }

  if (infoBullet.name === "ballistics_artillery_bullet") {
    PlayPositionSound(['explosion_4x2'], null, pos.x, pos.y);
    ExplosionRing(pos.x, pos.y, 150, infoBullet.area_covers, data.z);
    FireExplosion(pos.x, pos.y, 500, 10, 30, 55, 0, 360, data.z);
    SmokeExplosion(pos.x, pos.y, 1000, 10, 30, 35, 0, 360, data.z);
    return
  }

  SmokeExplosion(pos.x, pos.y, 1000, 20, 15, 15, 0, 360, data.z);
}

function FireExplosion(x, y, time, count, scale, speed, minAngle, maxAngle, z, reverseScale) {
  Explosion(x, y, time, count, scale, ["yellow"], speed, minAngle, maxAngle, z, 'ADD', reverseScale)
}

function RedExplosion(x, y, time, count, scale, speed, minAngle, maxAngle, z, reverseScale) {
  Explosion(x, y, time, count, scale, ["red"], speed, minAngle, maxAngle, z, 'ADD', reverseScale)
}

function BlueExplosion(x, y, time, count, scale, speed, minAngle, maxAngle, z, reverseScale) {
  Explosion(x, y, time, count, scale, ["blue"], speed, minAngle, maxAngle, z, 'ADD', reverseScale)
}

function BlueExplosionSCREEN(x, y, time, count, scale, speed, minAngle, maxAngle, z, reverseScale) {
  Explosion(x, y, time, count, scale, ["blue"], speed, minAngle, maxAngle, z, 'SCREEN', reverseScale)
}

function SmokeExplosion(x, y, time, count, scale, speed, minAngle, maxAngle, z, reverseScale) {
  Explosion(x, y, time, count, scale, ["smoke_puff"], speed, minAngle, maxAngle, z, 'SCREEN', reverseScale)
}

function GreenExplosion(x, y, time, count, scale, speed, minAngle, maxAngle, z, reverseScale) {
  Explosion(x, y, time, count, scale, ["green"], speed, minAngle, maxAngle, z, 'ADD', reverseScale)
}

function ExplosionRing(x, y, time, radius, z) {

  let border = Scene.make.sprite({
    x: x,
    y: y,
    key: 'explosion_border',
    add: true
  });
  border.setOrigin(0.5);
  border.setDisplaySize(0, 0);
  border.setDepth(z);

  Scene.tweens.add({
    targets: border,
    props: {
      displayHeight: {value: radius * 3, duration: time, ease: 'Linear'},
      displayWidth: {value: radius * 3, duration: time, ease: 'Linear'},
      alpha: {value: 0, duration: time, ease: 'Linear', delayTime: time},
    },
    repeat: 0,
    onComplete: function () {
      border.destroy();
    }
  });

  setInterval(function () {
    border.destroy()
  }, time)
}

// setInterval(function () {
//   if (Scene) {
//console.log(Scene.children)
//console.log(Scene.children.list)
//console.log(Scene.children.systems.updateList)
//console.log(Scene.children.systems.tweens)
//console.log(Scene.children.systems.sound.sounds)
//   }
// }, 1000)

function Explosion(x, y, time, count, scale, type, speed, minAngle, maxAngle, z, blendMode, reverseScale) {

  if (reverseScale) {
    scale = {start: scale / 200, end: scale / 25}
  } else {
    scale = {start: scale / 50, end: scale / 100}
  }

  let explosion = Scene.add.particles('flares');
  let emitter = explosion.createEmitter({
    frame: type,
    x: x,
    y: y,
    gravityY: 0,
    speed: {min: 0, max: speed},
    scale: scale,
    angle: {min: minAngle, max: maxAngle},
    alpha: {start: 1, end: 0, ease: 'Quad.easeIn'},
    lifespan: {min: time / 2, max: time},
  });

  if (z < 0) z = 1;
  explosion.setDepth(z);

  if (blendMode) emitter.setBlendMode(blendMode);

  emitter.setQuantity(count);
  emitter.explode();

  setTimeout(function () {
    explosion.destroy();
  }, time)
}

export {
  ExplosionBullet,
  FireExplosion,
  RedExplosion,
  BlueExplosion,
  GreenExplosion,
  SmokeExplosion,
  ExplosionRing,
  Explosion,
  BlueExplosionSCREEN
}
