import {gameStore} from "../store";
import {Scene} from "../create";
import {MapSize} from "./createMap";

function GetSpriteSizeByMapLvl(z, normalSize, K) {
  // пуля имеет размер(normalSize) 1 на высоте 0(lvl)
  // каждый 1 уровень высоты размер пули растет на 0.05 (на высоте 1 размер 1.05) (на высоте 1.5 размер 1.075)
  // size = normalSize + lvl * 0.05
  let size = normalSize + (z / 100) * K;
  return {x: size, y: size}
}

function GetOffsetShadowByMapLvl(z, x, y, normalSize, K, mpID) {

  let mapX = x - MapSize;
  let mapY = y - MapSize;

  let zoneLvl = GetLvl(mapX, mapY, mpID);

  let offset = normalSize + ((z - zoneLvl) / 100) * K;
  let currPos = {
    x: x + Scene.shadowXOffset * offset,
    y: y + Scene.shadowYOffset * offset,
    z: z - zoneLvl,
    zoneLvl: zoneLvl
  };

  return currPos
  // TODO искать зону пока не найдется так которая подходит но пока тут мои любимые костыли)
  // for (let i = 0; i < 2; i++) {
  //
  //   let zoneLvl = GetLvl(currPos.x, currPos.y, mpID);
  //   let offset = normalSize + ((z - zoneLvl) / 100) * K;
  //
  //   if (zoneLvl === currPos.zoneLvl) {
  //     return currPos
  //   } else {
  //     return currPos = {
  //       x: x + Scene.shadowXOffset * offset,
  //       y: y + Scene.shadowYOffset * offset,
  //       z: z - zoneLvl,
  //       zoneLvl: zoneLvl,
  //     };
  //   }
  // }
}

function GetLvl(x, y, mpID) {
  let zoneLvl = GetMapLvl(x, y, mpID);
  //zoneLvl += GetObjectLvl(x, y); //todo дыра в производительности
  return zoneLvl
}

function GetMapLvl(x, y, mpID) {

  let lvlX = Math.round((x) / 16);
  let lvlY = Math.round((y) / 16);

  if (gameStore.maps[mpID].map.level_map.hasOwnProperty(lvlX + ":" + lvlY)) {
    let lvlMap = gameStore.maps[mpID].map.level_map[lvlX + ":" + lvlY];
    return lvlMap.level
  } else {
    return gameStore.maps[mpID].map.default_level
  }
}

/**
 * @return {number}
 */
function GetObjectLvl(x, y) {

  function checkObj(obj) {
    if (!obj || !obj.geo_data || obj.geo_data.length === 0 || obj.height === 0) return -1;

    if (Phaser.Math.Distance.Between(x, y, obj.x, obj.y) < 200) {
      for (let i = 0; i < obj.geo_data.length; i++) {
        if (obj.geo_data[i].radius > Phaser.Math.Distance.Between(x, y, obj.geo_data[i].x, obj.geo_data[i].y)) {
          return obj.height;
        }
      }
    }

    return -1
  }

  for (let i in gameStore.objects) {
    let lvl = checkObj(gameStore.objects[i]);
    if (lvl >= 0) {
      return lvl
    }
  }

  for (let i in gameStore.statickObjects) {
    let lvl = checkObj(gameStore.statickObjects[i]);
    if (lvl >= 0) {
      return lvl
    }
  }

  return 0
}

export {GetSpriteSizeByMapLvl, GetOffsetShadowByMapLvl, GetLvl, GetMapLvl, GetObjectLvl}
