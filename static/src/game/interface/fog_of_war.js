import {gameStore} from "../store";
import {MapSize} from "../map/create_map";

function UpdateObjectFog(scene, obj, sprite, type, range_view) {
  if (gameStore.gameReady) {

    let brush = gameStore.FogOfWar.objBrush[type][obj.id];
    if (!brush) {
      brush = createFogBrush(scene);
      gameStore.FogOfWar.objBrush[type][obj.id] = brush;
      gameStore.FogOfWar.brushContainer.add([brush]);
    }

    UpdateFogOfWar(scene, sprite, range_view * 2, brush);
  }
}

function UpdateFogOfWar(scene, sprite, range_view, brush) {
  brush.setPosition(sprite.x, sprite.y);
  brush.displayHeight = range_view;
  brush.displayWidth = range_view;
}

function UpdateFogBack(scene) {
  for (let i in gameStore.FogOfWar.objBrush.object) {
    if (!gameStore.objects[i]) {
      gameStore.FogOfWar.objBrush.object[i].destroy();
      delete gameStore.FogOfWar.objBrush.object[i];
    }
  }

  for (let i in gameStore.FogOfWar.objBrush.unit) {
    if (!gameStore.units[i]) {
      gameStore.FogOfWar.objBrush.unit[i].destroy();
      delete gameStore.FogOfWar.objBrush.unit[i];
    }
  }

  for (let i in gameStore.FogOfWar.objBrush.drone) {
    if (!gameStore.drones[i]) {
      gameStore.FogOfWar.objBrush.drone[i].destroy();
      delete gameStore.FogOfWar.objBrush.drone[i];
    }
  }
}

let initBrush = false;

function createFogBrush(scene) {

  if (!initBrush) {
    initBrush = true;
    let graphics = scene.add.graphics();
    graphics.setDefaultStyles({
      lineStyle: {
        width: 0,
      },
      fillStyle: {
        color: 0xFFFFFF,
        alpha: 1
      }
    });

    let circle = {x: 300, y: 300, radius: 300};
    graphics.fillCircleShape(circle);
    graphics.strokeCircleShape(circle);
    graphics.generateTexture("fog_of_war", 600, 600);
    graphics.destroy();
  }

  return scene.make.image({
    x: 0,
    y: 0,
    key: "fog_of_war",
    add: false
  });
}

function createFogBack(scene) {
  let graphics = scene.add.graphics();
  graphics.setDefaultStyles({
    fillStyle: {
      color: 0x000000,
      alpha: 0.4
    }
  });

  graphics.fillRect(0, 0, 1024, 1024);
  graphics.generateTexture("fog_of_war_back", 1024, 1024);
  graphics.destroy();

  let back = scene.make.image({
    x: 0,
    y: 0,
    key: "fog_of_war_back",
    add: true
  });
  back.setDepth(999);
  return back
}

function initFog(scene) {

  let back = createFogBack(scene);
  let brush = createFogBrush(scene);

  let maskContainer = scene.make.container({
    x: 0,
    y: 0,
    width: MapSize,
    height: MapSize,
  }, false);
  maskContainer.add([brush]);

  back.mask = new Phaser.Display.Masks.BitmapMask(scene, maskContainer);
  back.mask.invertAlpha = true;

  return {
    back: back,
    brush: brush,
    brushContainer: maskContainer,
    objBrush: {
      object: {},
      unit: {},
      drone: {},
    },
  }
}

function updateFogSize() {
  gameStore.FogOfWar.back.x = MapSize + (MapSize / 2) //scene.cameras.main.worldView.centerX;
  gameStore.FogOfWar.back.y = MapSize + (MapSize / 2) //scene.cameras.main.worldView.centerY;
  gameStore.FogOfWar.back.displayHeight = MapSize + 300   //scene.cameras.main.displayHeight + additionalFog * 2;
  gameStore.FogOfWar.back.displayWidth = MapSize + 300//scene.cameras.main.displayWidth + additionalFog * 2;
}

export {usersVisibleZone, initFog, UpdateFogBack, UpdateObjectFog, updateFogSize}
