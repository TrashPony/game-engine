import {gameStore} from "../../store";
import {GetGlobalPos} from "../gep_global_pos";
import {Scene} from "../../create";
import {ShortDirectionRotateTweenGroup} from "../../utils/rotate_sprite";

function CreateSpriteTurret(turretState, turretBase, scene) {

  turretBase.setDepth(1);
  let pos = GetGlobalPos(turretState.x, turretState.y, gameStore.map.id);

  for (let weaponState of turretState.weapons) {

    let weapon = scene.make.image({
      x: pos.x + weaponState.real_x_attach,
      y: pos.y + weaponState.real_y_attach,
      scale: turretState.scale / 100,
      key: 'sprites',
      add: true,
      frame: weaponState.weapon_texture,
    });
    weapon.setOrigin(weaponState.x_anchor, weaponState.y_anchor);
    weapon.setDepth(turretState.height);

    let xShadow = scene.shadowXOffset * 0.5;
    let yShadow = scene.shadowYOffset * 0.5;

    let weaponShadow = scene.make.image({
      x: pos.x + weaponState.real_x_attach + xShadow,
      y: pos.y + weaponState.real_y_attach + yShadow,
      scale: turretState.scale / 100,
      key: 'sprites',
      add: true,
      frame: weaponState.weapon_texture,
    });
    weaponShadow.setDepth(turretState.height - 1);
    weaponShadow.setOrigin(weaponState.x_anchor, weaponState.y_anchor);
    weaponShadow.setAlpha(0.4);
    weaponShadow.setTint(0x000000);

    weapon.angle = weaponState.gun_rotate;
    weaponShadow.angle = weaponState.gun_rotate;

    turretBase.weapon = weapon;
    turretBase.weaponShadow = weaponShadow;
    turretBase.xShadow = xShadow;
    turretBase.yShadow = yShadow;
  }

  return turretBase
}

function RotateTurretGun(path) {
  if (gameStore.objects[path.id]) {
    let turret = gameStore.objects[path.id].objectSprite;

    let angle = path.r - turret.angle
    if (angle > 180) {
      angle -= 360
    }

    if (Math.round(turret.weapon.angle) === Math.round(angle)) {
      return;
    }

    if (!Scene.cameras.main.worldView.contains(turret.x, turret.y)) {
      turret.weapon.setAngle(path.r)
      turret.weaponShadow.setAngle(path.r)
    } else {
      turret.rotateTween = ShortDirectionRotateTweenGroup(
        [turret.weapon, turret.weaponShadow],
        path.r,
        path.ms,
        turret.weapon.angle,
      );
    }
  }
}

export {CreateSpriteTurret, RotateTurretGun}
