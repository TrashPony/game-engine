import {gameStore} from "../store";
import {ClearBars} from "../interface/status_layer";

function removeUnit(unit) {

  for (let i in unit.weapons) {
    unit.weapons[i].weapon.destroy();
  }

  if (unit.sprite.bodyBottomRight) unit.sprite.bodyBottomRight.destroy();
  if (unit.sprite.bodyBottomLeft) unit.sprite.bodyBottomLeft.destroy();
  if (unit.sprite) unit.sprite.destroy();

  ClearBars('unit', unit.id, 'hp');
  ClearBars('unit', unit.id, 'shield');

  delete gameStore.units[unit.id]
}

export {removeUnit}
