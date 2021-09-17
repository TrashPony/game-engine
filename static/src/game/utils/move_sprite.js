import {Scene} from "../create";

function MoveSprite(sprite, x, y, ms, size) {
  if (size) {
    sprite.moveTween = Scene.tweens.add({
      targets: sprite,
      x: x,
      y: y,
      scale: size,
      ease: 'Linear',
      duration: ms,
      // onComplete: function () {
      //   tween.remove();
      //   tween = null;
      // }
    });
  } else {
    sprite.moveTween = Scene.tweens.add({
      targets: sprite,
      x: x,
      y: y,
      ease: 'Linear',
      duration: ms,
      // onComplete: function () {
      //   tween.remove();
      //   tween = null;
      // }
    });
  }
}

export {MoveSprite}
