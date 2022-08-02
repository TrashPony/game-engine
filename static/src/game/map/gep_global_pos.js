import {MapSize} from "./create_map";

function GetGlobalPos(x, y, mapID) {
  return {x: x + MapSize, y: y + MapSize}
}

export {GetGlobalPos}
