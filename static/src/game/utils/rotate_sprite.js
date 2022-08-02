import {Scene} from "../create";
import store from "../../store/store";

function ShortDirectionRotateTween(sprite, needAngle, time) {

  if (!sprite || Math.round(needAngle) === Math.round(sprite.angle)) return false;

  if (sprite.rotateTween && sprite.rotateTween.data && sprite.rotateTween.data[0].key === "angle" && Math.round(sprite.rotateTween.data[0].end) === needAngle) {
    return;
  }

  if (!store.getters.getSettings.MoveAndRotateTween) {
    sprite.setAngle(needAngle);
    return;
  }

  sprite.rotateTween = Scene.tweens.add({
    targets: sprite,
    angle: '+=' + Phaser.Math.Angle.ShortestBetween(sprite.angle, Phaser.Math.Angle.WrapDegrees(needAngle)),
    duration: time,
    repeat: 0,
  });

  return true
}

function ShortDirectionRotateTweenGroup(sprites, needAngle, time, currentAngle) {

  if (Math.round(needAngle) === Math.round(currentAngle)) return false;

  if (sprites[0].rotateTween && sprites[0].rotateTween.data[0].key === "angle" && Math.round(sprites[0].rotateTween.data[0].end) === needAngle) {
    return;
  }

  if (!store.getters.getSettings.MoveAndRotateTween) {
    for (let sprite of sprites) {
      if (sprite) sprite.setAngle(needAngle);
    }

    return;
  }

  sprites[0].rotateTween = Scene.tweens.add({
    targets: sprites,
    angle: '+=' + Phaser.Math.Angle.ShortestBetween(currentAngle, Phaser.Math.Angle.WrapDegrees(needAngle)),
    duration: time,
    repeat: 0,
  });

  return true
}

function rotatePoint(x, y, angle) {
  let alpha = angle * Math.PI / 180
  let newX = (x)*Math.cos(alpha) - (y)*Math.sin(alpha)
  let newY = (x)*Math.sin(alpha) + (y)*Math.cos(alpha)

  return {x: newX, y: newY}
}

export {
  ShortDirectionRotateTween,
  ShortDirectionRotateTweenGroup,
  rotatePoint
}
