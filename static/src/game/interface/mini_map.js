import {Scene} from "../create";
import {MapSize} from "../map/create_map";

let minimap = {
  backRenderTexture: null,
  nominalSize: 2048,
  mainMapRectangle: null,
  background: null,
  size: 200,
  init: false,
  mapPoints: {},
}

function initMiniMap() {
  if (minimap.backRenderTexture) return;

  if ($(window).width() < 1000) {
    minimap.size = 118
  }

  if ($(window).width() === 1000) {
    minimap.size = 150
  }

  minimap.init = true;
  minimap.backRenderTexture = Scene.cameras.add(Scene.cameras.main.width - minimap.size - 10, 10, minimap.size, minimap.size).setName('mini');
  updateMiniMap()
}

function setPositionMapRectangle() {
  if (minimap.mainMapRectangle) {
    minimap.mainMapRectangle.setPosition(Scene.cameras.main.worldView.centerX, Scene.cameras.main.worldView.centerY)
    minimap.mainMapRectangle.setDisplaySize(Scene.cameras.main.displayWidth, Scene.cameras.main.displayHeight)
  }
}

function updateMiniMap() {

  if ($(window).width() < 1000) {
    minimap.size = 118
  } else if ($(window).width() === 1000) {
    minimap.size = 150
  } else {
    minimap.size = 200
  }


  let miniMapWrapper = document.getElementById("miniMap")
  if (miniMapWrapper) {
    miniMapWrapper.style.width = (minimap.size - 5) + 'px'
    miniMapWrapper.style.height = (minimap.size - 22) + 'px'
    miniMapWrapper.style.display = 'block'
  }

  if (minimap.backRenderTexture) {
    minimap.backRenderTexture.setZoom(1 / (MapSize / minimap.size))
    minimap.backRenderTexture.setBounds(MapSize, MapSize, MapSize, MapSize);
    minimap.backRenderTexture.setBackgroundColor(0x002244);
    minimap.backRenderTexture.scrollX = MapSize;
    minimap.backRenderTexture.scrollY = MapSize;
    minimap.backRenderTexture.disableCull = false
    minimap.backRenderTexture.inputEnabled = false
    minimap.backRenderTexture.setPosition(Scene.cameras.main.width - minimap.size - 10, 10)
    minimap.backRenderTexture.setSize(minimap.size, minimap.size)
  }

  for (let id in minimap.mapPoints) {
    minimap.mapPoints[id].sprite.tween.remove();
    minimap.mapPoints[id].sprite.destroy();
    delete minimap.mapPoints[id]
  }

  if (!minimap.mainMapRectangle) {
    let mainMapRectangle = Scene.add.graphics({x: 0, y: 0, add: true});
    mainMapRectangle.lineStyle(15, 0xaaaa00);
    mainMapRectangle.strokeRect(0, 0, Scene.cameras.main.displayWidth / 2, Scene.cameras.main.displayHeight / 2);
    mainMapRectangle.generateTexture("main_map_rectangle", Scene.cameras.main.displayWidth / 2, Scene.cameras.main.displayHeight / 2);
    mainMapRectangle.destroy()

    minimap.mainMapRectangle = Scene.make.image({
      x: 0, y: 0, key: "main_map_rectangle", add: true
    });
    minimap.mainMapRectangle.setDepth(1000)
    Scene.cameras.main.ignore(minimap.mainMapRectangle)
  }
}

export {
  initMiniMap, updateMiniMap, setPositionMapRectangle, minimap
}
