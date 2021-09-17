import {BlueExplosion, FireExplosion, SmokeExplosion} from "./explosion";
import {GetGlobalPos} from "../map/gep_global_pos";
import {PlayPositionSound} from "../sound/play_sound";
import {gameStore} from "../store";

function FireWeapon(data) {
  if (!gameStore.gameReady) return;

  let weaponInfo = gameStore.gameTypes.weapons[data.type_id]
  let pos = GetGlobalPos(data.x, data.y, data.m);

  if (weaponInfo.name === "explores_weapon_1" || weaponInfo.name === "big_laser") {
    PlayPositionSound(['big_laser_1_fire', 'big_laser_2_fire', 'big_laser_3_fire'], null, pos.x, pos.y);
    BlueExplosion(pos.x, pos.y, 250, 15, 20, 250, data.r - 15, data.r + 15, data.z);
  }

  if (weaponInfo.name === "explores_weapon_2" || weaponInfo.name === "explores_weapon_3" || weaponInfo.name === "explores_weapon_4" || weaponInfo.name === "small_laser") {
    if (weaponInfo.name === "explores_weapon_2" || weaponInfo.name === "small_laser") {
      PlayPositionSound(['laser_1_fire_1', 'laser_1_fire_2', 'laser_1_fire_3'], null, pos.x, pos.y);
    } else {
      PlayPositionSound(['laser_2_fire_1', 'laser_2_fire_2', 'laser_2_fire_3', 'laser_2_fire_4'], null, pos.x, pos.y);
    }
    FireExplosion(pos.x, pos.y, 700, 15, 10, 10, data.r - 5, data.r + 5, data.z);
  }

  if (weaponInfo.name === "replic_weapon_1" || weaponInfo.name === "artillery") {
    PlayPositionSound(['firearms_6', 'firearms_4'], null, pos.x, pos.y);
    SmokeExplosion(pos.x, pos.y, 1000, 15, 20, 30, data.r - 15, data.r + 15, data.z);
    FireExplosion(pos.x, pos.y, 200, 10, 20, 300, data.r - 5, data.r + 5, data.z);
  }

  if (weaponInfo.name === "replic_weapon_2") {
    PlayPositionSound(['firearms_1', 'firearms_2'], null, pos.x, pos.y);
    FireExplosion(pos.x, pos.y, 150, 5, 20, 15, data.r - 25, data.r + 25, data.z);
  }

  if (weaponInfo.name === "replic_weapon_3" || weaponInfo.name === "replic_weapon_4" || weaponInfo.name === "tank_gun") {
    if (weaponInfo.name === "replic_weapon_4") {
      PlayPositionSound(['firearms_5'], null, pos.x, pos.y);
    } else {
      PlayPositionSound(['firearms_3'], null, pos.x, pos.y);
    }

    SmokeExplosion(pos.x, pos.y, 300, 5, 20, 30, data.r - 15, data.r + 15, data.z);
    FireExplosion(pos.x, pos.y, 200, 5, 20, 300, data.r - 5, data.r + 5, data.z);
  }

  if (weaponInfo.name === 'reverses_weapon_1' || weaponInfo.name === 'reverses_weapon_2' || weaponInfo.name === 'reverses_weapon_3' ||
    weaponInfo.name === 'reverses_weapon_4' || weaponInfo.name === 'big_missile' || weaponInfo.name === 'small_missile') {
    PlayPositionSound(['missile_1', 'missile_2', 'missile_3'], null, pos.x, pos.y);
  }
}

export {FireWeapon}
