import {gameStore} from "../store";
import {GetGlobalPos} from "./gep_global_pos";
import {CreateSpriteTurret} from "./structures/turret";
import store from "../../store/store";


let template = {
  id: 0,
  x: 0,
  y: 0,
  rotate: 0,
  height: 0,

  hp: 0,
  max_hp: 0,
  current_energy: 0,
  max_energy: 0,
  view_range: 0,

  x_shadow_offset: 0,
  y_shadow_offset: 0,
  shadow_intensity: 0,
  owner_id: 0,
  priority: 0,

  team_id: 0,          // byte
  work: false,         // byte
  build: false,        // byte
  scale: 0,            // byte
  shadow: false,       // byte
  animate: false,      // byte
  animation_speed: 0,  // byte
  animate_loop: false, // byte

  type: "",
  texture: "",

  weapons: [{
    gun_rotate: 0,
    real_x_attach: 0,
    real_y_attach: 0,
    number: 0,   // byte
    x_anchor: 0, // byte
    y_anchor: 0, // byte
  }],

  geo_data: [{
    x: 0,
    y: 0,
    radius: 0,
  }],
}

function CreateObject(coordinate, x, y, noBmd, scene) {

  coordinate.objectSprite = gameObjectCreate(x, y, coordinate.texture, coordinate.scale,
    coordinate.shadow && coordinate.shadow_intensity > 0 && store.getters.getSettings.ObjectShadows,
    coordinate.rotate, coordinate.x_shadow_offset, coordinate.y_shadow_offset, coordinate.shadow_intensity,
    noBmd || coordinate.build, 'sprites', scene, gameStore.map.id, coordinate);

  coordinate.objectSprite.setDepth(coordinate.height);
  if (coordinate.objectSprite.shadow) {
    coordinate.objectSprite.shadow.setDepth(coordinate.height - 1);
  }

  if (coordinate.type === "turret") {
    CreateSpriteTurret(coordinate, coordinate.objectSprite, scene)
  }

  return coordinate
}

function gameObjectCreate(x, y, texture, scale, needShadow, rotate, xShadowOffset, yShadowOffset, shadowIntensity, noBmd, atlasName, scene, mapID, coordinate) {
  if (!gameStore.gameReady) return;

  let shadow;

  let pos = GetGlobalPos(x, y, mapID);
  if (needShadow) {
    shadow = scene.make.image({
      x: pos.x + xShadowOffset, y: pos.y + yShadowOffset, frame: texture, add: true, key: atlasName,
    });

    shadow.setOrigin(0.5);
    shadow.setScale(scale / 100);
    shadow.setAngle(rotate);
    shadow.setAlpha(shadowIntensity / 100);
    shadow.setTint(0x000000);
  }

  let object = scene.make.image({
    x: pos.x, y: pos.y, frame: texture, add: true, key: atlasName,
  });
  object.setOrigin(0.5);
  object.setScale(scale / 100);
  object.setAngle(rotate);

  object.shadow = shadow;

  return object
}

export {CreateObject, gameObjectCreate}
