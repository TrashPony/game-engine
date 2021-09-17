import {gameStore} from "../store";
import {GetGlobalPos} from "./gep_global_pos";

function CreateObject(coordinate, x, y, noBmd, scene) {

  if (coordinate.type === "mission") {
    noBmd = true
  }

  let structures = ['unknown_civilization_jammer', 'repair_station', 'turret', 'beacon', 'storage', 'shield',
    'generator', 'missile_defense', 'turret', 'turret', 'jammer', 'radar', 'meteorite_defense', 'extractor',
    'explores_antenna', 'expensive_tower'];

  let atlasName = coordinate.type;
  if (structures.includes(coordinate.type)) {
    atlasName = 'structures';
  }

  coordinate.objectSprite = gameObjectCreate(x, y, coordinate.texture, coordinate.scale, coordinate.shadow, coordinate.position_data.rotate,
    coordinate.x_shadow_offset, coordinate.y_shadow_offset, coordinate.shadow_intensity,
    noBmd || coordinate.build, atlasName, scene, coordinate.map_id);

  if (coordinate.objectSprite) {
    coordinate.objectSprite.setDepth(coordinate.position_data.height);
    if (coordinate.objectSprite.shadow) {
      coordinate.objectSprite.shadow.setDepth(coordinate.position_data.height - 1);
    }
  }

  return coordinate
}

function gameObjectCreate(x, y, texture, scale, needShadow, rotate, xShadowOffset, yShadowOffset, shadowIntensity, noBmd, atlasName, scene, mapID) {
  let shadow;

  let pos = GetGlobalPos(x, y, mapID);
  if (needShadow) {
    shadow = scene.make.image({
      x: pos.x + xShadowOffset,
      y: pos.y + yShadowOffset,
      frame: texture,
      add: true,
      key: atlasName,
    });

    shadow.setOrigin(0.5);
    shadow.setScale(scale / 100);
    shadow.setAngle(rotate);
    shadow.setAlpha(shadowIntensity / 100);
    shadow.setTint(0x000000);
  }

  let object = scene.make.image({
    x: pos.x,
    y: pos.y,
    frame: texture,
    add: true,
    key: atlasName,
  });
  object.setOrigin(0.5);
  object.setScale(scale / 100);
  object.setAngle(rotate);

  object.shadow = shadow;

  if (!needShadow && !noBmd && gameStore.mapsState[mapID].bmdTerrain && !texture.includes('geyser')) {
    object.setPosition(x, y);
    gameStore.mapsState[mapID].bmdTerrain.bmd.draw(object);
    object.destroy();
    object = null;
  }

  return object
}

function CreateAnimate(coordinate, x, y, scene) {

  if (coordinate.unit_overlap) {
    coordinate.objectSprite = gameAnimateObjectCreate(x, y, coordinate.animate_sprite_sheets, coordinate.scale, coordinate.shadow,
      coordinate.position_data.rotate, coordinate.animation_speed, coordinate.animate_loop, scene, coordinate.map_id);
  } else {
    coordinate.objectSprite = gameAnimateObjectCreate(x, y, coordinate.animate_sprite_sheets, coordinate.scale, coordinate.shadow,
      coordinate.position_data.rotate, coordinate.animation_speed, coordinate.animate_loop, scene, coordinate.map_id);
  }

  if (coordinate.objectSprite) {
    coordinate.objectSprite.setDepth(coordinate.height);
    if (coordinate.objectSprite.shadow) coordinate.objectSprite.shadow.setDepth(coordinate.height - 1);
  }
}

function gameAnimateObjectCreate(x, y, texture, scale, needShadow, rotate, speed, needAnimate, scene, mapID) {
  let shadow;

  let pos = GetGlobalPos(x, y, mapID);

  scene.anims.create({
    key: texture,
    frames: scene.anims.generateFrameNumbers(texture),
    frameRate: speed,
    repeat: -1
  });

  if (needShadow) {

    shadow = scene.make.sprite({
      x: pos.x + scene.shadowXOffset,
      y: pos.y + scene.shadowYOffset,
      key: texture,
      add: true
    });

    shadow.setOrigin(0.5);
    shadow.setScale(scale / 100);
    shadow.setAngle(rotate);
    shadow.setTint(0x000000);

    if (needAnimate) {
      shadow.anims.play(texture);
    }
  }

  let object = scene.make.sprite({
    x: pos.x,
    y: pos.y,
    key: texture,
    add: true
  });
  object.setOrigin(0.5);
  object.setScale(scale / 100);
  object.setAngle(rotate);

  if (needAnimate) {
    object.anims.play(texture);
  }

  object.shadow = shadow;

  return object
}

export {CreateObject, CreateAnimate, gameObjectCreate, gameAnimateObjectCreate}
