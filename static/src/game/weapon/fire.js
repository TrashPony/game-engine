import {BlueExplosion, FireExplosion, SmokeExplosion} from "./explosion";
import {GetGlobalPos} from "../map/gep_global_pos";
import {PlayPositionSound} from "../sound/play_sound";
import {gameStore} from "../store";
import {Scene} from "../create";

function FireWeapon(data) {
  if (!gameStore.gameReady) return;

  let weaponInfo = gameStore.gameTypes.weapons[data.type_id]
  let pos = GetGlobalPos(data.x, data.y, gameStore.map.id);
  if (!Scene.cameras.main.worldView.contains(pos.x, pos.y)) {
    return;
  }

  if ((weaponInfo.name === "explores_weapon_1") && data.ap !== 0) {
    BlueExplosion(pos.x, pos.y, 200, 1, 15, 100, data.r - 10, data.r + 10, data.z, null, {
      start: data.ap / 100,
      end: 0,
      ease: 'Quad.easeIn'
    });
  }

  if (weaponInfo.name === "explores_weapon_1" && data.ap === 0) {
    BlueExplosion(pos.x, pos.y, 250, 3, 40, 350, data.r - 5, data.r + 5, data.z);
  }

  if (weaponInfo.name === "replic_weapon_2" || weaponInfo.name === "laser_turret_weapon") {
    PlayPositionSound(['firearms_1', 'firearms_2'], null, pos.x, pos.y, false, 0.5, 'fire_weapon');
    FireExplosion(pos.x, pos.y, 150, 2, 30, 15, data.r - 25, data.r + 25, data.z);
  }

  if (weaponInfo.name === 'reverses_weapon_1') {
    SmokeExplosion(pos.x, pos.y, 500, 7, 10, 100, data.r - 5, data.r + 5, data.z);
    SmokeExplosion(pos.x, pos.y, 500, 3, 10, 75, data.r - 185, data.r - 185, data.z);
  }
}

export {FireWeapon}
