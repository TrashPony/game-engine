import {gameStore} from "../store";
import {Scene} from "../create";
import {ShortDirectionRotateTween, ShortDirectionRotateTweenGroup} from "../utils/rotate_sprite";

function SetBodyAngle(unit, angle, time, ignoreOldTween, shadowDist, force) {
  if (angle > 180) {
    angle -= 360
  }

  if (Math.round(unit.sprite.angle) === Math.round(angle) && !force) {
    return;
  }

  if (!Scene.cameras.main.worldView.contains(unit.sprite.x, unit.sprite.y)) {
    unit.sprite.setAngle(angle)
    return;
  }

  ShortDirectionRotateTween(unit.sprite, angle, time);
}

function RotateUnitGun(id, ms, rotate, slot) {
  if (!gameStore.gameReady) return;

  let unit = gameStore.units[id];
  if (!unit || !unit.sprite || !unit.weapons || !unit.weapons[slot]) return;

  if (Math.round(unit.weapons[slot].weapon.angle) === Math.round(rotate)) {
    return;
  }

  if (!Scene.cameras.main.worldView.contains(unit.sprite.x, unit.sprite.y)) {
    unit.weapons[slot].weapon.setAngle(rotate)
  } else {
    ShortDirectionRotateTweenGroup([unit.weapons[slot].weapon], rotate, ms, unit.weapons[slot].weapon.angle)
  }
}

export {SetBodyAngle, RotateUnitGun}
