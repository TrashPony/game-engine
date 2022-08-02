import {Scene} from "../create";
import store from "../../store/store";

function MoveSprite(sprite, x, y, ms, size, forceAnimate, smooth) {

  if (!store.getters.getSettings.MoveAndRotateTween) {
    sprite.setPosition(x, y);
    if (size) sprite.setScale(size)
    return;
  }

  if (sprite.x === x && sprite.y === y && (!size || sprite.scale === size)) {
    return
  }

  if (sprite.moveTween &&
    sprite.moveTween.data[0].key === "x" && sprite.moveTween.data[0].end === x &&
    sprite.moveTween.data[1].key === "y" && sprite.moveTween.data[1].end === y) {
    return;
  }

  let newTween;
  if (forceAnimate || Scene.cameras.main.worldView.contains(x, y)) {
    newTween = Scene.tweens.add({
      targets: sprite,
      duration: ms,
      x: x,
      y: y,
    });
  } else {
    sprite.setPosition(x, y);
  }

  if (size && sprite.scale !== size) sprite.setScale(size)
  if (newTween) {
    sprite.moveTween = newTween
  }
}

function movePlugin(sprite, x, y, ms) {
  if (!sprite.moveToPlugin) {
    sprite.moveToPlugin = Scene.plugins.get('rexmovetoplugin').add(sprite, {})
  }

  let speed = Phaser.Math.Distance.Between(sprite.x, sprite.y, x, y) * (1000 / ms)
  sprite.moveToPlugin.setSpeed(speed);
  sprite.moveToPlugin.moveTo(x, y);
}

export {MoveSprite, movePlugin}
