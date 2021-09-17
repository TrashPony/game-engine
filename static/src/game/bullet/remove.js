import {gameStore} from "../store";

function deleteBullet(id) {

  let bullet = gameStore.bullets[id]
  if (!bullet) return ;

  if (bullet.debugText) {
    bullet.debugText.setVisible(false);
  }

  if (bullet.shadow) {
    bullet.shadow.setVisible(false);
  }

  if (bullet.fairy) {
    bullet.fairy.destroy();
  }

  if (bullet.emitter) { // метеориты
    bullet.emitter.killAll()
    bullet.emitter.stop();

    if (bullet.shadow) {
      bullet.shadow.emitter.killAll()
      bullet.shadow.emitter.stop();
      bullet.shadow.destroy();
    }

    bullet.sprite.destroy();
    delete gameStore.bullets[id];
    return;
  }

  delete gameStore.bullets[id];

  if (bullet.sprite) {
    bullet.sprite.setVisible(false);

    if (!gameStore.cacheSpriteBullets.hasOwnProperty(bullet.sprite.frame.name)) {
      gameStore.cacheSpriteBullets[bullet.sprite.frame.name] = [];
    }

    gameStore.cacheSpriteBullets[bullet.sprite.frame.name].push(bullet);
  }
}

export {deleteBullet}
