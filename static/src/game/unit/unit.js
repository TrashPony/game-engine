import {gameStore} from "../store";
import {Scene} from "../create";
import {RotateUnitGun, SetShadowAngle} from "./rotate"
import {CreateMapBar} from "../interface/status_layer";
import {GetGlobalPos} from "../map/gep_global_pos";
import {mouseBodyOver} from "./select_units";
import {SelectSprite} from "../utils/select_sprite";
import {GetMapLvl, GetSpriteSizeByMapLvl} from "../map/height_offset";

let UnitHeight = 35

/*{"id":-1,"owner_id":4,"map_id":-1,
"position_data":{"real_x":0,"real_y":0,"x":1592,"y":2026,"z":0,"rotate":0,"power_move":0,"reverse":0,"angular_velocity":0,"x_velocity":0,"y_velocity":0,"need_z":0,"speed":90,"reverse_speed":22.5,"power_factor":15,"reverse_factor":1.5,"turn_speed":0.3,"height":18,"width":26}}*/
function CreateNewUnit(newUnit) {
  if (!newUnit.hasOwnProperty('id')) {
    newUnit = JSON.parse(newUnit);
  }

  if (!gameStore.gameReady) return;

  let unit = gameStore.units[newUnit.id];
  if (!unit || !unit.sprite) {

    let pos = GetGlobalPos(newUnit.position_data.x, newUnit.position_data.y, newUnit.map_id);

    unit = createUnit(
      newUnit,
      pos.x,
      pos.y,
      newUnit.position_data.rotate,
    );

    gameStore.units[newUnit.id] = unit;

    createWeapons(unit)
    CreateMapBar(unit.sprite, unit.body.max_hp, unit.hp, 10, null, Scene, 'unit', unit.id, 'hp', 50);
  }
}

function createUnit(unit, x, y, rotate) {

  let unitBox = Scene.add.container(x, y);
  unitBox.setDepth(UnitHeight);
  unitBox.setScale(GetSpriteSizeByMapLvl(GetMapLvl(unitBox.x, unitBox.y, unit.map_id), unit.body.scale / 100, 0.02).x);

  bodyAnimate(unit)

  let bodyBottomShadow = Scene.make.sprite({
    key: unit.body.name,
    x: Scene.shadowXOffset / 2,
    y: Scene.shadowYOffset / 2,
    frame: unit.body.name + "_bottom",
    add: true
  });

  bodyBottomShadow.setOrigin(0.5);
  bodyBottomShadow.setAlpha(0.2);
  bodyBottomShadow.setTint(0x000000);

  let bodyBottomLeft = Scene.make.sprite({
    x: 0,
    y: 0,
    key: unit.body.name + "_bottom_animate_left",
    add: true
  });
  bodyBottomLeft.setOrigin(0.5);
  bodyBottomLeft.play(unit.body.name + "_bottom_animate_left");
  bodyBottomLeft.anims.pause();

  let bodyBottomRight = Scene.make.sprite({
    x: 0,
    y: 0,
    key: unit.body.name + "_bottom_animate_right",
    add: true
  });
  bodyBottomRight.setOrigin(0.5);
  bodyBottomRight.play(unit.body.name + "_bottom_animate_right");
  bodyBottomRight.anims.pause();

  let bodyBox = Scene.add.container(0, 0);
  let bodyShadow = Scene.make.sprite({
    key: unit.body.name,
    x: Scene.shadowXOffset / 2,
    y: Scene.shadowYOffset / 2,
    frame: unit.body.texture,
    add: true
  });
  bodyShadow.setOrigin(0.5);
  bodyShadow.setAlpha(0.2);
  bodyShadow.setTint(0x000000);

  let body = Scene.make.sprite({
    key: unit.body.name,
    x: 0,
    y: 0,
    frame: unit.body.texture,
    add: true
  });
  body.setOrigin(0.5);

  bodyBox.add(bodyBottomShadow);
  bodyBox.add(bodyBottomLeft);
  bodyBox.add(bodyBottomRight);
  bodyBox.add(bodyShadow);
  bodyBox.add(body);

  unit.selectSprite = SelectSprite(0, 0, null, 0xffffff, 0xffffff, body.displayHeight + 5);
  unit.selectSprite.visible = false
  bodyBox.add(unit.selectSprite);
  bodyBox.add(bodyShadow);
  bodyBox.add(body);

  unit.sprite = unitBox;
  unit.sprite.unitBody = body;
  unit.sprite.bodyShadow = bodyShadow;
  unit.sprite.bodyBottomLeft = bodyBottomLeft;
  unit.sprite.bodyBottomRight = bodyBottomRight;
  unit.sprite.bodyBottomShadow = bodyBottomShadow;

  unit.sprite.setAngle(rotate);
  unitBox.add(bodyBox);
  mouseBodyOver(body, unit, unitBox);

  SetShadowAngle(unit, rotate, Scene.shadowXOffset * 3);
  return unit
}

function createWeapons(unit) {
  unit.weapons = {};

  for (let i in unit.weapon_slots_data) {
    let weaponSlot = unit.weapon_slots_data[i];

    if (weaponSlot && weaponSlot.weapon) {

      let weaponBox = Scene.add.container(0, 0);

      let weapon = Scene.make.sprite({
        x: weaponSlot.real_x_attach,
        y: weaponSlot.real_y_attach,
        frame: weaponSlot.weapon.weapon_texture,
        add: true,
        key: "weapons",
      });
      weapon.xAttach = weaponSlot.real_x_attach;
      weapon.yAttach = weaponSlot.real_y_attach;
      weapon.setOrigin(weaponSlot.x_anchor, weaponSlot.y_anchor);

      let weaponShadow = Scene.make.sprite({
        x: weaponSlot.real_x_attach + Scene.shadowXOffset / 2,
        y: weaponSlot.real_y_attach + Scene.shadowYOffset / 2,
        frame: weaponSlot.weapon.weapon_texture,
        add: true,
        key: "weapons",
      });
      weaponShadow.setAlpha(0.5);
      weaponShadow.setTint(0x000000);
      weaponShadow.setOrigin(weaponSlot.x_anchor, weaponSlot.y_anchor);

      weaponBox.add(weaponShadow);
      weaponBox.add(weapon);
      unit.sprite.add(weaponBox);

      unit.weapons[weaponSlot.number] = {
        weapon: weapon,
        shadow: weaponShadow,
      }

      if (weaponSlot && weaponSlot.weapon) {
        RotateUnitGun(unit.id, weaponSlot.number, weaponSlot.gun_rotate, 64);
      }
    }
  }
}

function bodyAnimate(unit) {
  if (!gameStore.cacheAnimate.hasOwnProperty(unit.body.name)) {
    gameStore.cacheAnimate[unit.body.name] = true;

    Scene.textures.addSpriteSheetFromAtlas(unit.body.name + "_bottom_animate_left", {
      atlas: unit.body.name,
      frame: unit.body.name + "_bottom_animate_left",
      frameWidth: gameStore.unitSpriteSize,
      frameHeight: gameStore.unitSpriteSize,
    });

    Scene.textures.addSpriteSheetFromAtlas(unit.body.name + "_bottom_animate_right", {
      atlas: unit.body.name,
      frame: unit.body.name + "_bottom_animate_right",
      frameWidth: gameStore.unitSpriteSize,
      frameHeight: gameStore.unitSpriteSize,
    });


    Scene.anims.create({
      key: unit.body.name + "_bottom_animate_left",
      frames: Scene.anims.generateFrameNumbers(unit.body.name + "_bottom_animate_left"),
      frameRate: 16,
      repeat: -1
    });

    Scene.anims.create({
      key: unit.body.name + "_bottom_animate_right",
      frames: Scene.anims.generateFrameNumbers(unit.body.name + "_bottom_animate_right"),
      frameRate: 16,
      repeat: -1
    });
  }
}

export {CreateNewUnit, UnitHeight}
