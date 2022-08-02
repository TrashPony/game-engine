import {gameStore} from "../store";
import {removeUnit} from "../unit/remove";
import {deleteBullet} from "../bullet/remove";
import {Scene} from "../create";
import {getEffectsParticles} from "../weapon/explosion";
import {ClearBars} from "../interface/status_layer";

function RemoveRadarObject(mark) {
  if (mark.to === "unit") {
    let unit = gameStore.units[mark.id];
    if (unit) removeUnit(unit);
  }

  if (mark.to === "object" || mark.to === "dynamic_objects") {
    removeDynamicObject(gameStore.objects[mark.id])
    gameStore.objects[mark.id] = null;
  }

  if (mark.to === "bullet") {
    deleteBullet(mark.id);
  }
}

function removeAllObj() {
  removeDynamicObjects();
}

function removeDynamicObjects() {
  for (let i in gameStore.objects) {
    if (gameStore.objects[i] && gameStore.objects[i].objectSprite) {
      removeDynamicObject(gameStore.objects[i])
      gameStore.removeObjects[i] = gameStore.objects[i]
      gameStore.objects[i] = null;
    }
  }
}

function removeDynamicObject(obj) {

  if (!obj) return

  Scene.tweens.add({
    targets: [obj.objectSprite, obj.objectSprite.weapon, obj.objectSprite.top, obj.objectSprite.equip,
      obj.objectSprite.weaponShadow, obj.objectSprite.weaponBox, obj.objectSprite.equip, obj.objectSprite.equipShadow,
      obj.objectSprite.equipBox, obj.objectSprite.shadowTop, obj.objectSprite.shadow, obj.objectSprite.light],
    alpha: 0,
    ease: 'Linear',
    duration: 150,
    onComplete: function () {
      if (obj.objectSprite.weapon) obj.objectSprite.weapon.destroy();
      if (obj.objectSprite.weaponShadow) obj.objectSprite.weaponShadow.destroy();
      if (obj.objectSprite.equip) obj.objectSprite.equip.destroy();
      if (obj.objectSprite.equipShadow) obj.objectSprite.equipShadow.destroy();
      if (obj.objectSprite.equipBox) obj.objectSprite.equipBox.destroy();
      if (obj.objectSprite.top) obj.objectSprite.top.destroy();
      if (obj.objectSprite.shadowTop) obj.objectSprite.shadowTop.destroy();
      if (obj.objectSprite.shadow) obj.objectSprite.shadow.destroy();
      if (obj.objectSprite.passBuildSelectSprite) obj.objectSprite.passBuildSelectSprite.destroy();
      if (obj.objectSprite.noPassbuildSelectSprite) obj.objectSprite.noPassbuildSelectSprite.destroy();
      if (obj.objectSprite.light) obj.objectSprite.light.destroy();
      if (obj.objectSprite.RadarMark) obj.objectSprite.RadarMark.destroy();
      if (obj.objectSprite.emitter) {
        obj.objectSprite.emitter.stop();
        getEffectsParticles().emitters.remove(obj.objectSprite.emitter)
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

  if (obj.type === "rope_trap") {
    for (let i in gameStore.ropes) {
      if (i.includes(id + '_')) {
        gameStore.ropes[i].sprite.destroy()
        gameStore.ropes[i].shadow.destroy()
        delete gameStore.ropes[i]
      }
    }
  }
}

export {
  RemoveRadarObject,
  removeAllObj,
  removeDynamicObject,
  removeDynamicObjects
}
