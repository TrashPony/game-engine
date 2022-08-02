import Phaser from "phaser";
import {preload} from "./preload"
import {update} from "./update"
import {gameStore} from "./store";
import store from "../store/store";
import {cacheBars} from "./interface/status_layer";
import {updateMiniMap} from "./interface/mini_map";

let Game = {};
let Scene = null;

function CreateGame() {
  if (!gameStore.gameInit) {

    gameStore.gameInit = true;

    let config = {
      type: Phaser.WEBGL, scene: {
        preload: preload, create: create, update: update
      }, scale: {
        width: '100', height: '100', parent: 'game-container',
      }, render: {
        antialias: true, antialiasGL: false, batchSize: 512, powerPreference: "high-performance",
      }, fps: {
        target: 30, min: 20,
      }, physics: {
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

  this.scale.on('resize', function (gameSize, baseSize, displaySize, previousWidth, previousHeight) {
    if (previousWidth !== window.innerWidth || previousHeight !== window.innerHeight) {
      Scene.scale.setGameSize(window.innerWidth, window.innerHeight);
      updateMiniMap();
      store.commit({
        type: 'addCheckViewPort',
      })
    }
  });

  this.wasd = this.input.keyboard.addKeys({
    up: Phaser.Input.Keyboard.KeyCodes.W,
    down: Phaser.Input.Keyboard.KeyCodes.S,
    left: Phaser.Input.Keyboard.KeyCodes.A,
    right: Phaser.Input.Keyboard.KeyCodes.D,
  }, false, false);

  // параметры смещения тени игры
  this.shadowXOffset = 6;
  this.shadowYOffset = 8;

  // воспроизведение звука даже в неактивной вкладке
  this.sound.pauseOnBlur = false;

  this.cameras.main.setBackgroundColor('#000000');
  this.input.topOnly = false; // нажатия по всем обьектам под мышкой

  Scene.time.addEvent({
    delay: 1000, callback: function () {
      // подчищаем список т.к. метод destroy() почему отуда их не удаляет и в итоге разрастается до бесконечности
      clear()
    }, loop: true
  });

  cacheBars()
}

function clear() {
  for (let obj of Scene.children.systems.updateList._active) {
    if (!obj.active) Scene.children.systems.updateList.remove(obj)
  }
}

export {CreateGame, Game, Scene}
