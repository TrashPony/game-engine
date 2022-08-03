import store from "../store/store";
import {gameStore} from "./store";

let loading = 0;

function preload() {

  let scene = this;

  this.load.on('progress', function (value) {
    loading = Math.round(value * 100);
  });

  this.load.on('fileprogress', function (file) {
    if (!gameStore.noLoader && !gameStore.gameReady) {
      store.commit({
        type: 'setVisibleLoader',
        visible: true,
        text: `загрузка текстур: <span style="font-size: 9px; color: yellow">${file.key}</span> (${loading} %)`,
      });
    }
  });

  this.load.on('complete', function () {
    if (!gameStore.noLoader && !gameStore.gameReady) {
      store.commit({
        type: 'setVisibleLoader',
        visible: true,
        text: `<span style="font-size: 9px; color: yellow">тут я обычно зависаю</span>`,
      });
    }
  });

  const atlases = [
    {
      key: 'sprites',
      textureURL: require("../assets/sprites/sprites.png"),
      atlasURL: require("../assets/sprites/sprites_atlas.json"),
    }, {
      key: 'flares',
      textureURL: require("../assets/fire_effects/flares.png"),
      atlasURL: require("../assets/fire_effects/flares_atlas.json"),
    },
  ];

  for (let atlas of atlases) {
    this.load.atlas(atlas);
  }

  // fire effects
  this.load.image('explosion_border', require("../assets/fire_effects/explosion_border.png"));

  //Brush
  this.load.image('brush_256', require("../assets/terrainTextures/brush_256.png"));
  this.load.image('brush_128', require("../assets/terrainTextures/brush_128.png"));
  this.load.image('brush', require("../assets/terrainTextures/brush.png"));
  this.load.image('grass_1', require("../assets/terrainTextures/grass_1.png"));
  this.load.image('grass_2', require("../assets/terrainTextures/grass_2.png"));
  this.load.image('grass_3', require("../assets/terrainTextures/grass_3.png"));

  // Weapon sound
  this.load.audio('firearms_1', require("../assets/audio/sound_effects/fire_weapon/firearms_1.mp3"));
  this.load.audio('firearms_2', require("../assets/audio/sound_effects/fire_weapon/firearms_2.mp3"));

  this.load.audio('explosion_1', require("../assets/audio/sound_effects/explosions/explosion_1.mp3"));
  this.load.audio('explosion_2', require("../assets/audio/sound_effects/explosions/explosion_2.mp3"));
  this.load.audio('explosion_3', require("../assets/audio/sound_effects/explosions/explosion_3.mp3"));
  this.load.audio('explosion_4', require("../assets/audio/sound_effects/explosions/explosion_4.mp3"));
  this.load.audio('explosion_5', require("../assets/audio/sound_effects/explosions/explosion_5.mp3"));

  this.load.audio('explosion_1x2', require("../assets/audio/sound_effects/explosions/explosion_1x2.mp3"));
  this.load.audio('explosion_2x2', require("../assets/audio/sound_effects/explosions/explosion_2x2.mp3"));
  this.load.audio('explosion_3x2', require("../assets/audio/sound_effects/explosions/explosion_3x2.mp3"));
  this.load.audio('explosion_4x2', require("../assets/audio/sound_effects/explosions/explosion_4x2.mp3"));
  this.load.audio('explosion_5x2', require("../assets/audio/sound_effects/explosions/explosion_5x2.mp3"));

  this.load.audio('explosion_3x3', require("../assets/audio/sound_effects/explosions/explosion_3x3.mp3"));
  this.load.audio('explosion_4x3', require("../assets/audio/sound_effects/explosions/explosion_4x3.mp3"));
}

export {preload}
