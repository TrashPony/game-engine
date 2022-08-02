import {Scene} from "../create";
import {GetGlobalPos} from "../map/gep_global_pos";
import {gameStore} from "../store";
import {MoveSprite} from "../utils/move_sprite";

let Designators = {
  targetCircle: null,
}

let ammoState = null;

function createCircleTexture(scene, styles, key) {
  let graphics = scene.add.graphics();
  graphics.setDefaultStyles(styles);
  let circle = {x: 312, y: 312, radius: 300};
  graphics.fillCircleShape(circle);
  graphics.strokeCircleShape(circle);
  graphics.generateTexture(key, 624, 624);
  graphics.destroy();
}

function initGunDesignator(scene) {

  if (!Designators.targetCircle) {
    createCircleTexture(scene, {
      lineStyle: {
        width: 12, color: 0xFF0000, alpha: 0.5
      }, fillStyle: {
        color: 0xFF4444, alpha: 0.2
      }
    }, "GunDesignator")

    Designators.targetCircle = scene.make.image({
      x: 0, y: 0, key: "GunDesignator", add: true
    });
    Designators.targetCircle.setOrigin(0.5);
    Designators.targetCircle.setVisible(false);
    Designators.targetCircle.setDepth(900);
  }
}

function ServerGunDesignator(data) {
  if (!Designators.targetCircle) {
    return
  }

  let key = "targetCircle";
  if (data.hide) {
    if (Designators.targetCircle.visible) Designators.targetCircle.setVisible(false);
    return
  } else {
    if (!Designators.targetCircle.visible) Designators.targetCircle.setVisible(true);
  }

  let pos = GetGlobalPos(data.x, data.y, gameStore.map.id);
  if (Designators[key].displayHeight !== data.rd || Designators[key].displayWidth !== data.rd) {
    Scene.tweens.add({
      targets: Designators[key], props: {
        displayHeight: {value: data.rd, duration: 48, ease: 'Linear'},
        displayWidth: {value: data.rd, duration: 48, ease: 'Linear'},
      }, repeat: 0,
    });
  }

  MoveSprite(Designators[key], pos.x, pos.y, 48)
}

export {initGunDesignator, ServerGunDesignator, ammoState, createCircleTexture}
