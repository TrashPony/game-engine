import {gameStore} from "../store";
import {Scene} from "../create";

function createBox(scene, sizeBox, boxKey, x, y, color) {
  if (!gameStore.StatusLayer.healBoxes.hasOwnProperty(boxKey)) {
    let graphics = scene.add.graphics();
    graphics.setDefaultStyles({
      lineStyle: {
        width: 4,
        color: 0x000000,
        alpha: 1
      },
      fillStyle: {
        color: color,
        alpha: 1
      }
    });

    graphics.fillRect(0, 0, sizeBox * 2, sizeBox * 2);
    graphics.strokeRect(0, 0, sizeBox * 2, sizeBox * 2);
    graphics.generateTexture(boxKey, sizeBox * 2, sizeBox * 2);
    graphics.destroy();
  }
}

function ClearBars(type, id, typeBar) {
  let oldBar = gameStore.StatusLayer.bars[type + id + typeBar];
  if (oldBar) {
    oldBar.bar.destroy();
    delete gameStore.StatusLayer.bars[type + id + typeBar];
  }
}

function cacheBars() {
  for (let i = 0; i < 100; i++) {
    CreateMapBar(null, 1000, i * 10, 10, null, Scene, 'unit', 0, 'hp', 50);
    CreateMapBar(null, 400, i * 10, 10, null, Scene, 'unit', 1, 'hp', 50);
    CreateMapBar(null, 1000, i * 10, 10, null, Scene, 'object', 0, 'hp', 50);
    CreateMapBar(null, 100, i, 10, null, Scene, 'object', 1, 'hp', 50);
  }
  ClearBars('unit', 0, 'hp')
  ClearBars('unit', 1, 'hp')
  ClearBars('object', 0, 'hp')
  ClearBars('object', 1, 'hp')

}

function CreateMapBar(sprite, maxHP, hp, offsetY, color, scene, type, id, typeBar, hpInBox) {
  
  let sizeBox = 6;
  let interval = 1; // промеж уток между квадратиками

  let countBoxes = Math.ceil(maxHP / hpInBox);
  // для особо жирных
  if (countBoxes > 10) {
    countBoxes = 10;
    hpInBox = Math.ceil(maxHP / countBoxes);
  }

  if (countBoxes < 1) {
    countBoxes = 1;
  }

  if (!sprite) {
    sprite = {displayHeight: 0, originY: 0, x: 0, y: 0}
  }
  let displayHeight = sprite.displayHeight;
  if (displayHeight === 0 && sprite.hasOwnProperty('unitBody')) {
    displayHeight = sprite.unitBody.displayHeight * sprite.scale;
  }

  let boxY = Math.round(offsetY + displayHeight * sprite.originY);
  let sizeBar = (sizeBox * countBoxes * 2) + (interval * countBoxes)
  let startX = 0;
  let percentHP = Math.round(100 / (maxHP / hp));

  let barKey = getBarKey(countBoxes, color, hp, percentHP, hpInBox, startX, sizeBox, interval)

  let oldBar = gameStore.StatusLayer.bars[type + id + typeBar];
  if (oldBar) {
    if (oldBar.key === barKey) {
      return
    }
    ClearBars(type, id, typeBar);
  }

  if (!gameStore.StatusLayer.barsCache[barKey]) {
    //null 400 338 85 50 "unit" -307
    // console.log(color, maxHP, hp, percentHP, hpInBox, type, id)
    let bar = Scene.add.renderTexture(0, 0, sizeBar, sizeBox * 2);

    getBarKey(countBoxes, color, hp, percentHP, hpInBox, startX, sizeBox, interval, bar, scene, boxY)

    bar.saveTexture(barKey)
    bar.destroy();
    gameStore.StatusLayer.barsCache[barKey] = true
  }

  let barSprite = Scene.make.sprite({
    x: sprite.x,
    y: sprite.y + boxY,
    key: barKey,
    add: true,
  });
  gameStore.StatusLayer.bars[type + id + typeBar] = {
    bar: barSprite,
    key: barKey,
  }

  barSprite.setOrigin(0.5);
  barSprite.setScale(0.5);
  barSprite.setDepth(900);
}

function getBarKey(countBoxes, color, hp, percentHP, hpInBox, startX, sizeBox, interval, bar, scene, boxY) {

  let BarKey = countBoxes

  for (let i = 0; i < countBoxes; i++) {

    if (hp > 0) {
      if (!color) {
        color = Phaser.Display.Color.HexStringToColor(GetColorDamage(percentHP).color).color;
      }
    } else {
      color = 0x999b9f;
    }

    BarKey += ":" + color

    if (bar) {
      let boxKey = 'box_' + color + "" + sizeBox;
      createBox(scene, sizeBox, boxKey, startX, boxY, color)
      bar.drawFrame(boxKey, undefined, startX, 0);
    }

    hp -= hpInBox;
    startX += (sizeBox * 2) + interval
  }

  return BarKey
}

function UpdatePosBars(sprite, maxHP, hp, offsetY, color, scene, type, id, typeBar, hpInBox) {
  let bar = gameStore.StatusLayer.bars[type + id + typeBar];
  if (bar) {
    let displayHeight = sprite.displayHeight;
    if (displayHeight === 0 && sprite.hasOwnProperty('unitBody')) {
      displayHeight = sprite.unitBody.displayHeight * sprite.scale;
    }

    let boxY = Math.round(offsetY + displayHeight * sprite.originY);
    bar.bar.setPosition(sprite.x, sprite.y + boxY);
  } else {
    CreateMapBar(sprite, maxHP, hp, offsetY, color, scene, type, id, typeBar, hpInBox);
  }
}

function createObjectBars(id) {
  let obj = gameStore.objects[id];
  if (obj && obj.hp > -2) {
    if ((obj.build && obj.complete >= 100) || !obj.build) {

      ClearBars('object', obj.id, 'build');

      if (obj.hp >= 0) {
        CreateMapBar(obj.objectSprite, obj.max_hp, obj.hp, 0, null, Scene, 'object', obj.id, 'hp', 50);
      } else {
        CreateMapBar(obj.objectSprite, 250, 250, 0, null, Scene, 'object', obj.id, 'hp', 50);
      }

      if (obj.max_energy > 0) {
        CreateMapBar(obj.objectSprite, obj.max_energy / 100, obj.current_energy / 100, 7,
          0x00ffd6, Scene, 'object', obj.id, 'energy', 5);
      }

    } else if (obj.build && obj.complete < 100) {
      CreateMapBar(obj.objectSprite, obj.max_hp, obj.hp, 0, null, Scene, 'object', obj.id, 'hp', 50);
      CreateMapBar(obj.objectSprite, 100, obj.complete, 7, 0xff00e1, Scene, 'object', obj.id, 'build', 5);
    }
  }
}

function GetColorDamage(percentHP) {
  if (percentHP >= 80) {
    return {key: 'g', color: "#00ff0f"}
  } else if (percentHP < 80 && percentHP >= 75) {
    return {key: 'u', color: "#fff326"}
  } else if (percentHP < 75 && percentHP >= 50) {
    return {key: 'y', color: "#fac227"}
  } else if (percentHP < 50 && percentHP >= 25) {
    return {key: 'o', color: "#fa7b31"}
  } else if (percentHP < 25) {
    return {key: 'r', color: "#ff2615"}
  }
}

export {CreateMapBar, ClearBars, createObjectBars, UpdatePosBars, cacheBars}
