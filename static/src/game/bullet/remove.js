import {gameStore} from "../store";

function deleteBullet(id) {
  let bullet = gameStore.bullets[id]
  if (!bullet || !bullet.sprite) {
    return;
  }

  bullet.sprite.destroy();
  delete gameStore.bullets[id];
}

export {deleteBullet}
