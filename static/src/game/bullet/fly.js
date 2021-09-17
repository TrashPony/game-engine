import {gameStore} from "../store";
import {Scene} from "../create";
import {GetOffsetShadowByMapLvl, GetSpriteSizeByMapLvl} from "../map/height_offset";
import {ShortDirectionRotateTween} from "../utils/utils";
import {SmokeExplosion} from "../weapon/explosion";
import {GetGlobalPos} from "../map/gep_global_pos";
import {MoveSprite} from "../utils/move_sprite";
import {createFairy} from "./create";

function AddBulletMoveBufferData(data) {
  if (!gameStore.gameReady) return;

  let bullet = gameStore.bullets[data.id];

  if (!bullet) {
    bullet = {}
    gameStore.bullets[data.id] = bullet;
  }

  let pos = GetGlobalPos(data.x, data.y, data.m);
  data.x = pos.x;
  data.y = pos.y;

  if (!bullet.bufferMoveTick) {
    bullet.bufferMoveTick = [];
  }
  bullet.bufferMoveTick.push(data);

  // показываем прошлый тик что бы компенсировать сетевые лаги
  if (!bullet.updaterPos && bullet.bufferMoveTick.length >= 0) {
    bullet.updaterPos = true;
  }
}

function FlyBullet(bullet) {

  let data = bullet.bufferMoveTick.shift();
  if (!data) {
    return;
  }

  if (!bullet.sprite) {
    return
  }

  let scale = GetSpriteSizeByMapLvl(data.z, 1 / 5, 0.05);
  MoveSprite(bullet.sprite, data.x, data.y, data.ms, scale.x);

  let shadowPos = GetOffsetShadowByMapLvl(data.z, data.x, data.y, 0, 2.5, data.m);
  MoveSprite(bullet.shadow, shadowPos.x, shadowPos.y, data.ms, scale.x);

  ShortDirectionRotateTween(bullet.sprite, data.r, data.ms);
  ShortDirectionRotateTween(bullet.shadow, data.r, data.ms);

  if (bullet.fairy) {
    bullet.fairy.emitter.start()
    let speed = Phaser.Math.Distance.Between(data.x, data.y, bullet.sprite.x, bullet.sprite.y)

    MoveSprite(bullet.fairy, data.x, data.y, data.ms, null);
    bullet.fairy.emitter.setEmitterAngle({min: 180 + data.r, max: 180 + data.r});
    bullet.fairy.emitter.setSpeed({min: speed * 5, max: speed * 10})
  }
}

function FlyLaser(data) {
  if (!gameStore.gameReady) return;

  let infoBullet = gameStore.gameTypes.ammo[data.type_id];

  let laserName = infoBullet.name;
  let pos = GetGlobalPos(data.x, data.y, data.m);
  let targetPos = GetGlobalPos(data.to_x, data.to_y, data.m);

  let color = ["blue"];
  let scale = 0.1;
  let quantity = 128;

  if (laserName === "small_lens") {
    color = ['red'];
  }
  if (laserName === "build") {
    color = ['yellow'];
    scale = 0.03;
    quantity = 32;
  }
  if (laserName === "missile_defense") {
    color = ['red'];
    scale = 0.03;
    quantity = 32;
  }
  if (laserName === "heal") {
    color = ['green'];
    scale = 0.03;
    quantity = 32;
  }
  if (laserName === "transport_block") {
    color = ['red', 'blue'];
  }

  let line = new Phaser.Geom.Line(pos.x, pos.y, targetPos.x, targetPos.y);
  let particles = Scene.add.particles('flares');
  let emitter = particles.createEmitter({
    frame: color,
    x: 0, y: 0,
    scale: {start: scale, end: 0},
    alpha: {start: 1, end: 0},
    speed: {min: -5, max: 5},
    quantity: quantity,
    emitZone: {source: line},
    blendMode: 'SCREEN',
    lifespan: 300,
  });
  particles.setDepth(40);

  SmokeExplosion(targetPos.x, targetPos.y, 500, 15, 5, 10, null, null, 40);
  setTimeout(function () {
    emitter.stop();
  }, 100);

  setTimeout(function () {
    emitter.killAll();
    particles.destroy();
  }, 300);
}

export {FlyBullet, FlyLaser, AddBulletMoveBufferData}
