import {gameStore} from "../store";

function mouseBodyOver(body, unit, unitBox) {

  body.setInteractive({
    pixelPerfect: true,
    alphaTolerance: 1
  });

  body.on('pointerover', function () {
    unit.selectSprite.visible = true
  }, this);

  body.on('pointerout', function () {
    unit.selectSprite.visible = false
  }, this);

  body.on('pointerdown', function (pointer, localX, localY, event) {
    event.stopPropagation()
    removeAllSelect();
    gameStore.selectUnits.push(unit.id);
  }, this);
}

function unitInfo(unit, unitBox, body) {

}

function unitRemoveInfo(unit, unitBox, force) {

}

function removeAllSelect() {
  for (let id of gameStore.selectUnits) {
    let unit = gameStore.units[id];
    if (unit) {
      unit.selectSprite.setTint(0xffffff);
      unit.selectSprite.visible = false
    }
  }
  gameStore.selectUnits = []
}

export {mouseBodyOver}
