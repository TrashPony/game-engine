import Phaser from "phaser";
import {preload} from "./preload"
import {update} from "./update"
import {gameStore} from "./store";
import store from "../store/store";
import {cacheBars} from "./interface/status_layer";

let Game = {};
let Scene = null;

function CreateGame() {
  if (!gameStore.gameInit) {

    if (!gameStore.mapEditor) {
      store.commit({
        type: 'setVisibleLoader',
        visible: true,
        text: 'Загружаем картинки...'
      });
    }

    gameStore.gameInit = true;

    let config = {
      type: Phaser.WEBGL,
      scene: {
        preload: preload,
        create: create,
        update: update
      },
      scale: {
        autoRound: true,
        width: '100',
        height: '100',
        parent: 'game-container',
      },
      render: {
        //pixelArt: true,
        antialias: true,
        antialiasGL: true,
        //desynchronized: false,
        //batchSize: 16384,
        //roundPixels: true,
        clearBeforeRender: false,
        //failIfMajorPerformanceCaveat: true,
        powerPreference: "high-performance",
        //mipmapFilter: 'LINEAR_MIPMAP_LINEAR',
      },
      fps: {
        //target: 60,
        //min: 30,
        //forceSetTimeOut: true
        smoothStep: false,
      },
      physics: {
        default: false  // no physics system enabled
      },
    };

    Game = new Phaser.Game(config);
  } else {
    gameStore.noLoader = false;
    gameStore.gameReady = false;
    gameStore.unitReady = false;
  }

  gameStore.reload = false
}

function create() {
  Scene = this;

  this.wasd = this.input.keyboard.addKeys({ // todo только для тестов
    up: 'W',
    down: 'S',
    left: 'A',
    right: 'D',
  }, false, false);

  // параметры смещения тени игры
  this.shadowXOffset = 6;
  this.shadowYOffset = 8;

  // воспроизведение звука даже в неактивной вкладке
  this.sound.pauseOnBlur = false;

  this.cameras.main.setBackgroundColor('#000000');
  this.input.topOnly = false; // нажатия по всем обьектам под мышкой

  createSelectSprite(this);

  // подчищаем список т.к. метод destroy() почему отуда их не удаляет и в итоге разрастается до бесконечности
  setInterval(function () {
    for (let obj of Scene.children.systems.updateList._active) {
      if (!obj.active) Scene.children.systems.updateList.remove(obj)
    }
  }, 1000)

  cacheBars()
}

function createSelectSprite(scene) {

  let graphics = scene.add.graphics();
  graphics.setDefaultStyles({
    lineStyle: {
      width: 12,
      color: 0xFFFFFF,
      alpha: 0.5
    },
    fillStyle: {
      color: 0xFFFFFF,
      alpha: 0.1
    }
  });

  let circle = {x: 312, y: 312, radius: 300};
  graphics.fillCircleShape(circle);
  graphics.strokeCircleShape(circle);
  graphics.generateTexture("select_sprite", 624, 624);
  graphics.destroy();

  let select_rect = scene.add.graphics();
  select_rect.setDefaultStyles({
    lineStyle: {
      width: 6,
      color: 0xFFFFFF,
      alpha: 0.5
    },
    fillStyle: {
      color: 0xFFFFFF,
      alpha: 0.1
    }
  });

  select_rect.fillRoundedRect(0, 0, 128, 128, 16);
  select_rect.strokeRoundedRect(3, 3, 122, 122, 16);
  select_rect.generateTexture("select_rect", 128, 128);
  select_rect.destroy();

  let deniedSprite = scene.add.graphics();
  deniedSprite.setDefaultStyles({
    lineStyle: {
      width: 60,
      color: 0xFF0000,
      alpha: 0.8
    },
    fillStyle: {
      color: 0x000000,
      alpha: 0.3
    }
  });

  let deniedCircle = {x: 300, y: 300, radius: 300 - 60};
  deniedSprite.fillCircleShape(deniedCircle);
  deniedSprite.strokeCircleShape(deniedCircle);
  deniedSprite.lineBetween(100, 400, 500, 200);
  deniedSprite.generateTexture("denied_rect", 600, 600);
  deniedSprite.destroy();
}

export {CreateGame, Game, Scene}
