import {Scene} from "../create";
import {gameStore} from "../store";

function CreateBullet(data, infoBullet) {

  if (!gameStore.radarWork) return;

  let bullet = {bufferMoveTick: []};

  bullet.sprite = Scene.make.image({
    x: data.x, y: data.y, key: 'sprites', add: true, frame: infoBullet.name,
  });
  bullet.sprite.setOrigin(0.5);
  bullet.sprite.setScale(0.2);
  bullet.sprite.setAngle(data.r);

  gameStore.bullets[data.id] = bullet;
  bullet.sprite.setDepth(399);

  return bullet
}

function GetCacheBulletSprite(data, infoBullet) {
  return CreateBullet(data, infoBullet)
}

export {CreateBullet, GetCacheBulletSprite}
