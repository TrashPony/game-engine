import {FireExplosion} from "./explosion";
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

  if (weaponInfo.name === "replic_weapon_2" || weaponInfo.name === "laser_turret_weapon") {
    PlayPositionSound(['firearms_1', 'firearms_2'], null, pos.x, pos.y, false, 0.5, 'fire_weapon');
    FireExplosion(pos.x, pos.y, 150, 2, 30, 15, data.r - 25, data.r + 25, data.z);
  }
}

export {FireWeapon}
