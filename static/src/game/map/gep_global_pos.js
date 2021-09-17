import {MapSize} from "./createMap";
import {gameStore} from "../store";

function GetGlobalPos(x, y, mapID) {
  if (gameStore.maps[mapID]) {
    return {x: x + (gameStore.maps[mapID].x + 1) * MapSize, y: y + (gameStore.maps[mapID].y + 1) * MapSize}
  }

  return {x: 0, y: 0}
}

export {GetGlobalPos}
