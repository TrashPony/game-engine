import {Scene} from "../create";
import {GetGlobalPos} from "../map/gep_global_pos";

function damageText(data) {
  let pos = GetGlobalPos(data.x, data.y, data.m);

  function randomInt(min, max) {
    return min + Math.floor((max - min) * Math.random());
  }

  pos.x += randomInt(-15, 15);
  pos.y += randomInt(-15, 15);

  let text = Scene.add.bitmapText(pos.x, pos.y, 'bit_text', data.d, 32);
  if (data.t === "shield") {
    text.setTint(0x00d3ff)
  }

  text.setDepth(1001);
  text.setOrigin(0.5);
  text.setScale(0.45);

  Scene.tweens.add({
    targets: text,
    x: text.x,
    y: text.y - 100,
    ease: 'Linear',
    duration: 1000,
    onComplete: function () {
      text.destroy();
    }
  });
}

export {damageText}
