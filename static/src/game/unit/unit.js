import {gameStore} from "../store";
import {Scene} from "../create";
import {RotateUnitGun} from "./rotate"
import {CreateMapBar} from "../interface/status_layer";
import {GetGlobalPos} from "../map/gep_global_pos";
import {mouseBodyOver} from "./select_units";
import {intFromBytes} from "../../utils";
import {FocusMS} from "../interface/focus_unit";

let UnitHeight = 25;

function CreateNewUnit(newUnit) {
  if (!gameStore.gameReady) return;

  newUnit = parseData(newUnit)

  let unit = gameStore.units[newUnit.id];
  if (!unit || !unit.sprite) {

    let pos = GetGlobalPos(newUnit.position_data.x, newUnit.position_data.y, gameStore.map.id);

    let bodyState = gameStore.gameTypes.bodies[newUnit.body_state.id]
    newUnit.body = JSON.parse(JSON.stringify(bodyState));
    newUnit.body.max_hp = newUnit.body_state.max_hp;
    unit = createUnit(newUnit, pos.x, pos.y, newUnit.body);

    gameStore.units[newUnit.id] = unit;
    createWeapons(unit)
    CreateMapBar(unit.sprite, newUnit.body.max_hp, unit.hp, 10, null, Scene, 'unit', unit.id, 'hp', 50);

    if (gameStore.gameSettings.follow_camera && gameStore.player.id === newUnit.owner_id) {
      FocusMS()
    }
  }
}

function parseData(unitData) {
  let unit = {
    id: intFromBytes(unitData.slice(0, 4)),
    owner_id: intFromBytes(unitData.slice(4, 8)),
    hp: intFromBytes(unitData.slice(8, 12)),
    body_state: {
      id: intFromBytes(unitData.slice(12, 13)),
      texture: "",
      max_hp: intFromBytes(unitData.slice(26, 30)),
    },
    team_id: intFromBytes(unitData.slice(13, 14)),
    position_data: {
      x: intFromBytes(unitData.slice(14, 18)),
      y: intFromBytes(unitData.slice(18, 22)),
      rotate: intFromBytes(unitData.slice(22, 26)),
    },
    weapon_slots_data: [],
  };

  let textureLength = intFromBytes(unitData.slice(30, 31))
  unit.body_state.texture = String.fromCharCode.apply(String, unitData.subarray(31, 31 + textureLength))
  let stopByte = 31 + textureLength;

  for (; stopByte < unitData.length;) {
    let slotData = {
      number: intFromBytes(unitData.slice(stopByte, stopByte + 1)),
      real_x_attach: intFromBytes(unitData.slice(stopByte + 1, stopByte + 5)),
      real_y_attach: intFromBytes(unitData.slice(stopByte + 5, stopByte + 9)),
      x_anchor: intFromBytes(unitData.slice(stopByte + 9, stopByte + 10)) / 100,
      y_anchor: intFromBytes(unitData.slice(stopByte + 10, stopByte + 11)) / 100,
      gun_rotate: intFromBytes(unitData.slice(stopByte + 11, stopByte + 15)),
      weapon_texture: "",
      weapon: false,
    }

    let textureLength = intFromBytes(unitData.slice(stopByte + 15, stopByte + 16));
    if (textureLength > 0) {
      slotData.weapon_texture = String.fromCharCode.apply(String, unitData.subarray(stopByte + 16, stopByte + 16 + textureLength))
      slotData.weapon = slotData.weapon_texture !== ""
    }

    stopByte = stopByte + 16 + textureLength
    unit.weapon_slots_data.push(slotData)
  }

  return unit
}

function createUnit(unit, x, y, bodyState) {
  let body = Scene.make.image({
    key: 'sprites', x: x, y: y, frame: unit.body_state.texture + "_skin_1", add: true
  });
  body.setOrigin(0.5);
  body.setDepth(UnitHeight);
  body.setScale(bodyState.scale / 100);

  unit.sprite = body;
  unit.sprite.setAngle(unit.position_data.rotate);

  mouseBodyOver(body, unit);
  return unit
}

function createWeapons(unit) {
  unit.weapons = {};

  for (let weaponSlot of unit.weapon_slots_data) {

    if (weaponSlot && weaponSlot.weapon) {

      let weapon = Scene.make.image({
        x: unit.sprite.x + (weaponSlot.real_x_attach * unit.sprite.scale),
        y: unit.sprite.y + (weaponSlot.real_y_attach * unit.sprite.scale),
        frame: weaponSlot.weapon_texture + "_skin_1",
        add: true,
        key: "sprites",
      });
      weapon.xAttach = weaponSlot.real_x_attach;
      weapon.yAttach = weaponSlot.real_y_attach;
      weapon.setOrigin(weaponSlot.x_anchor, weaponSlot.y_anchor);

      unit.weapons[weaponSlot.number] = {
        weapon: weapon
      }

      if (weaponSlot && weaponSlot.weapon) {
        RotateUnitGun(unit.id, 64, weaponSlot.gun_rotate, weaponSlot.number);
      }
    }
  }
}

export {CreateNewUnit, UnitHeight}
