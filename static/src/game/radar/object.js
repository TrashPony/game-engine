import {CreateAnimate, CreateObject} from "../map/objects";
import {gameStore} from "../store";
import {Scene} from "../create";
import {deleteBullet} from "../bullet/remove";
import store from "../../store/store";
import {ClearBars, createObjectBars} from "../interface/status_layer";
import {ParseObject} from "../map/createMap";
import {GetCacheBulletSprite} from "../bullet/create";
import {GetGlobalPos} from "../map/gep_global_pos";
import {CreateNewUnit} from "../unit/unit";

function CreateRadarObject(mark, object) {
  if (mark.to === "unit") {
    CreateNewUnit(object)
  }

  if (mark.to === "dynamic_objects") {
    createDynamicObject(ParseObject(object));
  }

  if (mark.to === "bullet") {
    object = JSON.parse(object);

    let pos = GetGlobalPos(object.x, object.y, object.m);
    object.x = pos.x;
    object.y = pos.y;

    let infoBullet = gameStore.gameTypes.ammo[object.type_id];
    GetCacheBulletSprite(object, infoBullet, object.m);
  }
}

function RemoveRadarObject(mark) {
  // if (mark.to === "unit") {
  //   let unit = gameStore.units[mark.id];
  //   if (unit) removeUnit(unit);
  // }

  if (mark.to === "dynamic_objects") {
    removeDynamicObject(mark.id)
  }

  if (mark.to === "bullet") {
    deleteBullet(mark.id);
  }
}

function removeAllObj() {

  // for (let i in gameStore.units) {
  //   removeUnit(gameStore.units[i]);
  // }

  removeDynamicObjects();
}

function removeDynamicObjects() {
  for (let i in gameStore.objects) {
    if (gameStore.objects[i] && gameStore.objects[i].objectSprite) {
      removeDynamicObject(i)
    }
  }
}

function createDynamicObject(object) {
  if (gameStore.objects[object.id]) return;

  if (object.texture !== '') {
    CreateObject(object, object.position_data.x, object.position_data.y, true, Scene);
  }

  if (object.animate_sprite_sheets !== '') {
    CreateAnimate(object, object.position_data.x, object.position_data.y, Scene);
  }

  gameStore.objects[object.id] = object;
  createObjectBars(object.id)
}

function removeDynamicObject(id) {
  if (gameStore.objects[id]) {
    let obj = gameStore.objects[id];

    store.commit({
      type: 'toggleWindow',
      id: 'reservoirTip' + obj.x + "" + obj.y,
      component: 'reservoirTip',
      forceClose: true,
    });

    store.commit({
      type: 'toggleWindow',
      id: 'ObjectDialog' + obj.id,
      component: 'ObjectDialog',
      forceClose: true,
    });

    Scene.tweens.add({
      targets: [obj.objectSprite, obj.objectSprite.weapon, obj.objectSprite.top, obj.objectSprite.equip,
        obj.objectSprite.weaponShadow, obj.objectSprite.weaponBox, obj.objectSprite.equip, obj.objectSprite.equipShadow,
        obj.objectSprite.equipBox, obj.objectSprite.shadowTop, obj.objectSprite.shadow, obj.objectSprite.light,
        obj.objectSprite.emitter],
      alpha: 0,
      ease: 'Linear',
      duration: 150,
      onComplete: function () {

        if (obj.objectSprite.weapon) obj.objectSprite.weapon.destroy();
        if (obj.objectSprite.weaponShadow) obj.objectSprite.weaponShadow.destroy();
        if (obj.objectSprite.weaponBox) obj.objectSprite.weaponBox.destroy();
        if (obj.objectSprite.equip) obj.objectSprite.equip.destroy();
        if (obj.objectSprite.equipShadow) obj.objectSprite.equipShadow.destroy();
        if (obj.objectSprite.equipBox) obj.objectSprite.equipBox.destroy();
        if (obj.objectSprite.top) obj.objectSprite.top.destroy();
        if (obj.objectSprite.shadowTop) obj.objectSprite.shadowTop.destroy();
        if (obj.objectSprite.shadow) obj.objectSprite.shadow.destroy();
        if (obj.objectSprite.passBuildSelectSprite) obj.objectSprite.passBuildSelectSprite.destroy();
        if (obj.objectSprite.noPassbuildSelectSprite) obj.objectSprite.noPassbuildSelectSprite.destroy();
        if (obj.objectSprite.light) obj.objectSprite.light.destroy();
        if (obj.objectSprite.emitter) {
          obj.objectSprite.emitter.emitter.stop();
          obj.objectSprite.emitter.destroy();
        }

        if (obj.placeSprite) obj.placeSprite.destroy();
        if (obj.objectSprite.lights) {
          for (let i of obj.objectSprite.lights) {
            i.green.destroy();
            i.red.destroy();
          }
        }

        obj.objectSprite.destroy();
      }
    });

    ClearBars('object', obj.id, 'hp');
    ClearBars('object', obj.id, 'build');
    ClearBars('object', obj.id, 'energy');
    ClearBars('object', obj.id, 'shield');
    ClearBars('object', obj.id, 'progress');

    // TODO костыль, убрать когда радарные сообщения будут бинарными
    gameStore.removeObjects[id] = gameStore.objects[id]
    gameStore.objects[id] = null;
  }
}

function CreateDynamicObjects(dynamicObjects) {
  removeDynamicObjects();

  for (let id in dynamicObjects) {
    createDynamicObject(ParseObject(dynamicObjects[id]));
  }
}

export {CreateRadarObject, RemoveRadarObject, removeAllObj, CreateDynamicObjects}
