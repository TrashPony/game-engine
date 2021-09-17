import {gameStore} from "../store";
import {Scene} from "../create";
import {PositionAttachSprite, ShortDirectionRotateTween} from "../utils/utils";

function SetBodyAngle(unit, angle, time, ignoreOldTween, shadowDist) {
  SetShadowAngle(unit, time, shadowDist);
  if (angle > 180) {
    angle -= 360
  }

  ShortDirectionRotateTween(unit.sprite, angle, time);
}

function RotateUnitGun(id, slot, rotate, ms) {
  if (!gameStore.gameReady) return;

  let unit = gameStore.units[id];
  if (!unit || !unit.sprite || !unit.weapons[slot]) return;

  ShortDirectionRotateTween(unit.weapons[slot].weapon, rotate - unit.sprite.angle, ms);
  ShortDirectionRotateTween(unit.weapons[slot].shadow, rotate - unit.sprite.angle, ms);

  let connectWeapons = PositionAttachSprite(45 - unit.sprite.angle, Scene.shadowXOffset * 2);
  shadowTime(unit.weapons[slot].shadow, connectWeapons.x + unit.weapons[slot].xAttach, connectWeapons.y + unit.weapons[slot].yAttach, ms);
}

function SetShadowAngle(unit, time, shadowDist) {
  let shadowAngle = 45 - unit.sprite.angle;
  let connectPoints = PositionAttachSprite(shadowAngle, shadowDist);
  shadowTime(unit.sprite.bodyShadow, connectPoints.x, connectPoints.y);
}

function shadowTime(sprite, newX, newY, rotateTime = 10) {
  Scene.tweens.add({
    targets: sprite,
    props: {
      x: {value: newX, duration: rotateTime, ease: 'Linear'},
      y: {value: newY, duration: rotateTime, ease: 'Linear'}
    },
    repeat: 0,
  });
}

export {SetBodyAngle, RotateUnitGun, SetShadowAngle}
