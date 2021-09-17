import {gameStore} from "../store";
import {Scene} from "../create";

function FocusUnit(id) {
  let unit = gameStore.units[id];

  if (unit && unit.sprite) {
    Scene.cameras.main.startFollow(unit.sprite);
  }
}

function FocusMS() {
  if (!gameStore.gameReady) return;

  FocusUnit(gameStore.user_squad_id)
}

export {FocusUnit, FocusMS}
