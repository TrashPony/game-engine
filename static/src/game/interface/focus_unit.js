import {gameStore} from "../store";
import {Scene} from "../create";

function FocusUnit(id, force) {
  if (gameStore.gameSettings.follow_camera && !force) return;

  let unit = gameStore.units[id];

  if (unit && unit.sprite) {
    Scene.cameras.main.startFollow(unit.sprite);
    return true
  }

  return false
}

function FocusMS() {
  if (!gameStore.gameReady) return;

  for (let i in gameStore.units) {
    if (gameStore.player.id === gameStore.units[i].owner_id) {
      return FocusUnit(gameStore.units[i].id, true)
    }
  }

  return false
}

export {FocusUnit, FocusMS}
