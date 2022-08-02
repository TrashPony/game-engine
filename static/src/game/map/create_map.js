import {gameStore} from "../store";
import {CreateObject} from "./objects";
import {GetGlobalPos} from "./gep_global_pos";
import {initGunDesignator} from "../interface/server_gun_designator";
import {base64ToArrayBuffer, intFromBytes, mobileAndTabletCheck} from "../../utils";
import {updateMiniMap} from "../interface/mini_map";
import {createDynamicObject} from "../watch/create";

let MapSize = 4096;

function CreateAllMaps(scene) {

  if (Object.keys(gameStore.maps).length === 1) {
    for (let i in gameStore.maps) {
      MapSize = gameStore.maps[i].map.x_size;
    }
  }

  // наращивание размера карты не влияет на производительность
  scene.cameras.main.setBounds(-600, -600, (MapSize * 3) + 1200, (MapSize * 3) + 1200);
  initGunDesignator(scene);
  updateMiniMap();

  if (mobileAndTabletCheck()) {
    scene.cameras.main.setZoom(0.75);
  } else {
    scene.cameras.main.setZoom(1);
  }

  if (!gameStore.mouseInGameChecker.updater) {

    gameStore.mouseInGameChecker.updater = true;
    scene.input.on('pointermove', function (pointer) {
      gameStore.HoldAttackMouse = false;
    });

    scene.input.on('gameout', function (pointer) {
      gameStore.HoldAttackMouse = true;
    });

    scene.input.on('gameover', function (pointer) {
      gameStore.HoldAttackMouse = false;
    });
  }

  for (let i in gameStore.maps) {
    if (gameStore.maps.hasOwnProperty(i) && gameStore.maps[i].map) {

      gameStore.mapsState[i] = {
        bmdTerrain: {},
        staticObjects: [],
      };

      gameStore.map = gameStore.maps[i].map

      Create_map(scene, gameStore.maps[i].map)
    }
  }
}

function Create_map(scene, map) {
  setTimeout(function () {
    CreateFlore(map, scene);
    CreateObjects(map, scene);
  }, 100)
}

function CreateFlore(map, scene) {
  let pos = GetGlobalPos(0, 0, map.id);
  //bitmapData для отрисовки статичного нижнего слоя
  let bmdTerrain = scene.add.renderTexture(pos.x, pos.y, MapSize, MapSize);
  bmdTerrain.setInteractive();
  bmdTerrain.setDepth(-3)

  let bmdObject = scene.add.renderTexture(pos.x, pos.y, MapSize, MapSize);
  bmdObject.setDepth(300)
  gameStore.mapsState[map.id].bmdTerrain = {
    bmd: bmdTerrain,
    bmdObject: bmdObject
  };

  // сортировка по приоритету отрисовки текстур
  let flores = [];

  for (let x in map.flore) {
    for (let y in map.flore[x]) {
      flores.push(map.flore[x][y]);
    }
  }

  flores.sort(function (a, b) {
    return a.texture_priority - b.texture_priority;
  });

  let brush = scene.textures.get('brush').getSourceImage();
  let brush128 = scene.textures.get('brush_128').getSourceImage();

  for (let i in flores) {

    if (!flores.hasOwnProperty(i)) continue;

    let flore = flores[i];

    if (flore.texture_over_flore !== '') {

      let textureKey = "terrain_" + flore.texture_over_flore + "_brush";

      if (!gameStore.cashTextures.hasOwnProperty(textureKey)) {
        gameStore.cashTextures[textureKey] = true;
        let texture = scene.textures.get(flore.texture_over_flore).getSourceImage();

        if (flore.texture_over_flore === "water_1") {
          let bmd = scene.textures.createCanvas(textureKey, 128, 128);
          bmd.draw(0, 0, brush128);
          bmd.context.globalCompositeOperation = 'source-in';
          bmd.draw(0, 0, texture);
        } else {
          let bmd = scene.textures.createCanvas(textureKey, 512, 512);
          bmd.draw(0, 0, brush);
          bmd.context.globalCompositeOperation = 'source-in';
          bmd.draw(0, 0, texture);
        }
      }

      if (flore.texture_over_flore === "water_1") {
        gameStore.mapsState[map.id].bmdTerrain.bmd.drawFrame(textureKey, undefined, flore.x - 64, flore.y - 64);
      } else {
        gameStore.mapsState[map.id].bmdTerrain.bmd.drawFrame(textureKey, undefined, flore.x - 256, flore.y - 256);
      }
    }
  }

  scene.cameras.main.centerOn(pos.x + MapSize / 2, pos.y + MapSize / 2);
}

function CreateObjects(map, scene) {
  for (let i in map.static_objects_json) {
    let bytes = base64ToArrayBuffer(map.static_objects_json[i]);

    let obj = ParseObject(bytes);
    if (obj.static) {
      gameStore.mapsState[map.id].staticObjects.push(obj);
    } else {
      createDynamicObject(obj);
    }
  }

  // сортировка по приоритету отрисовки обьектов
  gameStore.mapsState[map.id].staticObjects.sort(function (a, b) {
    return a.priority - b.priority;
  });

  for (let i in gameStore.mapsState[map.id].staticObjects) {

    if (!gameStore.mapsState[map.id].staticObjects.hasOwnProperty(i)) {
      continue
    }

    let obj = gameStore.mapsState[map.id].staticObjects[i];
    CreateObject(obj, obj.x, obj.y, false, scene);
  }
}

function ParseObject(data) {

  let typeLength = intFromBytes(data.slice(69, 70))
  let type = String.fromCharCode.apply(String, data.subarray(70, 70 + typeLength))
  let stopByte = 70 + typeLength;

  let textureLength = intFromBytes(data.slice(stopByte, stopByte + 1))
  let texture = String.fromCharCode.apply(String, data.subarray(stopByte + 1, stopByte + 1 + textureLength))
  stopByte = stopByte + 1 + textureLength;

  let weaponCounts = intFromBytes(data.slice(stopByte, stopByte + 1))
  stopByte = stopByte + 1;

  let weapons = [];

  for (let i = 0; i < weaponCounts; i++) {
    let slotData = {
      gun_rotate: intFromBytes(data.slice(stopByte, stopByte + 4)),
      real_x_attach: intFromBytes(data.slice(stopByte + 4, stopByte + 8)),
      real_y_attach: intFromBytes(data.slice(stopByte + 8, stopByte + 12)),
      number: intFromBytes(data.slice(stopByte + 12, stopByte + 13)),
      x_anchor: intFromBytes(data.slice(stopByte + 13, stopByte + 14)) / 100,
      y_anchor: intFromBytes(data.slice(stopByte + 14, stopByte + 15)) / 100,
      weapon_texture: "",
    }

    let textureLength = intFromBytes(data.slice(stopByte + 15, stopByte + 16));
    if (textureLength > 0) {
      slotData.weapon_texture = String.fromCharCode.apply(String, data.subarray(stopByte + 16, stopByte + 16 + textureLength))
    }

    weapons.push(slotData);
    stopByte = stopByte + 16 + textureLength;
  }

  let geoData = [];
  for (; stopByte < data.length;) {
    geoData.push({
      x: intFromBytes(data.slice(stopByte, stopByte + 4)),
      y: intFromBytes(data.slice(stopByte + 4, stopByte + 8)),
      radius: intFromBytes(data.slice(stopByte + 8, stopByte + 12)),
    })
    stopByte += 12
  }

  let template = {
    id: intFromBytes(data.slice(0, 4)),
    x: intFromBytes(data.slice(4, 8)),
    y: intFromBytes(data.slice(8, 12)),
    rotate: intFromBytes(data.slice(12, 16)),
    height: intFromBytes(data.slice(16, 20)),

    hp: intFromBytes(data.slice(20, 24)),
    max_hp: intFromBytes(data.slice(24, 28)),
    current_energy: intFromBytes(data.slice(28, 32)),
    max_energy: intFromBytes(data.slice(32, 36)),
    view_range: intFromBytes(data.slice(36, 40)),

    x_shadow_offset: intFromBytes(data.slice(40, 44)),
    y_shadow_offset: intFromBytes(data.slice(44, 48)),
    shadow_intensity: intFromBytes(data.slice(48, 52)),
    owner_id: intFromBytes(data.slice(52, 56)),
    priority: intFromBytes(data.slice(56, 60)),

    team_id: intFromBytes(data.slice(60, 61)),
    work: intFromBytes(data.slice(61, 62)) === 1,
    build: intFromBytes(data.slice(62, 63)) === 1,
    scale: intFromBytes(data.slice(63, 64)),
    shadow: intFromBytes(data.slice(64, 65)) === 1,
    animate: intFromBytes(data.slice(65, 66)) === 1,
    animation_speed: intFromBytes(data.slice(66, 67)),
    animate_loop: intFromBytes(data.slice(67, 68)) === 1,
    static: intFromBytes(data.slice(68, 69)) === 1,

    type: type,
    texture: texture,
    weapons: weapons,
    geo_data: geoData,
  }

  return template
}

export {Create_map, CreateAllMaps, MapSize, ParseObject}
