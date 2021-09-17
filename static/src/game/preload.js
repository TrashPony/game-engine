import store from "../store/store";
import {gameStore} from "./store";

let loading = 0;

function preload() {

  this.load.on('progress', function (value) {
    loading = Math.round(value * 100);
  });

  this.load.on('fileprogress', function (file) {
    if (!gameStore.noLoader) {
      store.commit({
        type: 'setVisibleLoader',
        visible: true,
        text: `загрузка текстур: <span style="font-size: 9px; color: yellow">${file.key}</span> (${loading} %)`,
      });
    }
  });

  this.load.on('complete', function () {
    if (!gameStore.noLoader) {
      store.commit({
        type: 'setVisibleLoader',
        visible: true,
        text: `<span style="font-size: 9px; color: yellow">тут я обычно зависаю</span>`,
      });
    }
  });

  this.load.bitmapFont('bit_text',
    require("../assets/bit_map_text/bit_map_text.png"),
    require("../assets/bit_map_text/bit_map_text.xml"),
  );

  const atlases = [
    {
      key: 'weapons',
      textureURL: require("../assets/units/weapon/atlas/weapons.png"),
      atlasURL: require("../assets/units/weapon/atlas/weapons_atlas.json"),
    }, {
      key: 'mountains',
      textureURL: require("../assets/map/objects/mountains/atlas/mountains.png"),
      atlasURL: require("../assets/map/objects/mountains/atlas/mountains_atlas.json"),
    }, {
      key: 'plants',
      textureURL: require("../assets/map/objects/plants/atlas/plants.png"),
      atlasURL: require("../assets/map/objects/plants/atlas/plants_atlas.json"),
    }, {
      key: 'ravines',
      textureURL: require("../assets/map/objects/ravines/atlas/ravines.png"),
      atlasURL: require("../assets/map/objects/ravines/atlas/ravines_atlas.json"),
    }, {
      key: 'roads',
      textureURL: require("../assets/map/objects/roads/atlas/roads.png"),
      atlasURL: require("../assets/map/objects/roads/atlas/roads_atlas.json"),
    }, {
      key: 'unknown_civilization',
      textureURL: require("../assets/map/objects/unknown_civilization/atlas/unknown_civilization.png"),
      atlasURL: require("../assets/map/objects/unknown_civilization/atlas/unknown_civilization_atlas.json"),
    }, {
      key: 'unit_wreckage',
      textureURL: require("../assets/map/objects/unit_wreckage/atlas/unit_wreckage.png"),
      atlasURL: require("../assets/map/objects/unit_wreckage/atlas/unit_wreckage_atlas.json"),
    }, {
      key: 'other',
      textureURL: require("../assets/map/objects/other/atlas/other.png"),
      atlasURL: require("../assets/map/objects/other/atlas/other_atlas.json"),
    }, {
      key: 'elevations',
      textureURL: require("../assets/map/objects/elevations/atlas/elevations.png"),
      atlasURL: require("../assets/map/objects/elevations/atlas/elevations_atlas.json"),
    }, {
      key: 'craters',
      textureURL: require("../assets/map/objects/craters/atlas/craters.png"),
      atlasURL: require("../assets/map/objects/craters/atlas/craters_atlas.json"),
    }, {
      key: 'flares',
      textureURL: require("../assets/fire_effects/flares.png"),
      atlasURL: require("../assets/fire_effects/flares_atlas.json"),
    }, {
      key: 'bullets',
      textureURL: require("../assets/units/gameAmmo/atlas/bullets.png"),
      atlasURL: require("../assets/units/gameAmmo/atlas/bullets_atlas.json"),
    }, {
      key: 'apd_light',
      textureURL: require("../assets/units/body/apd_light/atlas/apd_light.png"),
      atlasURL: require("../assets/units/body/apd_light/atlas/apd_light_atlas.json"),
    },
  ];

  for (let atlas of atlases) {
    this.load.atlas(atlas);
  }

  // fire effects
  this.load.image('explosion_border', require("../assets/fire_effects/explosion_border.png"));

  //Brush
  this.load.image('water_1', require("../assets//terrainTextures/water_1.jpg"));
  this.load.image('brush_256', require("../assets/terrainTextures/brush_256.png"));
  this.load.image('brush_128', require("../assets/terrainTextures/brush_128.png"));

  this.load.image('brush', require("../assets/terrainTextures/brush.png"));
  this.load.image('desertDunes', require("../assets/terrainTextures/desertDunes.png"));
  this.load.image('desertDunes_2', require("../assets/terrainTextures/desertDunes_2.png"));
  this.load.image('xenos', require("../assets/terrainTextures/xenos.png"));
  this.load.image('xenos_2', require("../assets/terrainTextures/xenos_2.png"));
  this.load.image('arctic', require("../assets/terrainTextures/arctic.png"));
  this.load.image('arctic_2', require("../assets/terrainTextures/arctic_2.png"));
  this.load.image('tundra', require("../assets/terrainTextures/tundra.png"));
  this.load.image('tundra_2', require("../assets/terrainTextures/tundra_2.png"));
  this.load.image('grass_1', require("../assets/terrainTextures/grass_1.png"));
  this.load.image('grass_2', require("../assets/terrainTextures/grass_2.png"));
  this.load.image('grass_3', require("../assets/terrainTextures/grass_3.png"));
  this.load.image('soil', require("../assets/terrainTextures/soil.png"));
  this.load.image('soil_2', require("../assets/terrainTextures/soil_2.png"));
  this.load.image('ravine_1', require("../assets/terrainTextures/ravine_1.png"));
  this.load.image('ravine_2', require("../assets/terrainTextures/ravine_2.png"));
  this.load.image('paving_stone_1', require("../assets/terrainTextures/paving_stone_1.png"));
  this.load.image('clay_1', require("../assets/terrainTextures/clay_1.png"));
  this.load.image('clay_2', require("../assets/terrainTextures/clay_2.png"));

  //Clouds
  this.load.image('cloud0', require("../assets/map/clouds/cloud13.png"));
  this.load.image('cloud1', require("../assets/map/clouds/cloud1.png"));
  this.load.image('cloud2', require("../assets/map/clouds/cloud2.png"));
  this.load.image('cloud3', require("../assets/map/clouds/cloud3.png"));
  this.load.image('cloud4', require("../assets/map/clouds/cloud4.png"));
  this.load.image('cloud5', require("../assets/map/clouds/cloud5.png"));
  this.load.image('cloud6', require("../assets/map/clouds/cloud6.png"));
  this.load.image('cloud7', require("../assets/map/clouds/cloud7.png"));
  this.load.image('cloud8', require("../assets/map/clouds/cloud8.png"));
  this.load.image('cloud9', require("../assets/map/clouds/cloud9.png"));
  this.load.image('cloud10', require("../assets/map/clouds/cloud10.png"));
  this.load.image('cloud11', require("../assets/map/clouds/cloud11.png"));
  this.load.image('cloud12', require("../assets/map/clouds/cloud12.png"));

  this.load.audio('structure_4', require("../assets/audio/sound_effects/fire_weapon/structure_4.mp3"));
  this.load.audio('structure_5', require("../assets/audio/sound_effects/fire_weapon/structure_5.mp3"));

  // Weapon sound
  this.load.audio('big_laser_1_fire', require("../assets/audio/sound_effects/fire_weapon/big_laser_1_fire.mp3"));
  this.load.audio('big_laser_1_reload', require("../assets/audio/sound_effects/fire_weapon/big_laser_1_reload.mp3"));
  this.load.audio('big_laser_2_fire', require("../assets/audio/sound_effects/fire_weapon/big_laser_2_fire.mp3"));
  this.load.audio('big_laser_2_reload', require("../assets/audio/sound_effects/fire_weapon/big_laser_2_reload.mp3"));
  this.load.audio('big_laser_3_fire', require("../assets/audio/sound_effects/fire_weapon/big_laser_3_fire.mp3"));
  this.load.audio('big_laser_3_reload', require("../assets/audio/sound_effects/fire_weapon/big_laser_3_reload.mp3"));

  this.load.audio('laser_1_fire_1', require("../assets/audio/sound_effects/fire_weapon/laser_1_fire_1.mp3"));
  this.load.audio('laser_1_fire_2', require("../assets/audio/sound_effects/fire_weapon/laser_1_fire_2.mp3"));
  this.load.audio('laser_1_fire_3', require("../assets/audio/sound_effects/fire_weapon/laser_1_fire_3.mp3"));
  this.load.audio('laser_1_reload', require("../assets/audio/sound_effects/fire_weapon/laser_1_reload.mp3"));

  this.load.audio('laser_2_fire_1', require("../assets/audio/sound_effects/fire_weapon/laser_2_fire_1.mp3"));
  this.load.audio('laser_2_fire_2', require("../assets/audio/sound_effects/fire_weapon/laser_2_fire_2.mp3"));
  this.load.audio('laser_2_fire_3', require("../assets/audio/sound_effects/fire_weapon/laser_2_fire_3.mp3"));
  this.load.audio('laser_2_fire_4', require("../assets/audio/sound_effects/fire_weapon/laser_2_fire_4.mp3"));
  this.load.audio('laser_2_reload', require("../assets/audio/sound_effects/fire_weapon/laser_1_reload.mp3"));

  this.load.audio('firearms_1', require("../assets/audio/sound_effects/fire_weapon/firearms_1.mp3"));
  this.load.audio('firearms_2', require("../assets/audio/sound_effects/fire_weapon/firearms_2.mp3"));
  this.load.audio('firearms_3', require("../assets/audio/sound_effects/fire_weapon/firearms_3.mp3"));
  this.load.audio('firearms_4', require("../assets/audio/sound_effects/fire_weapon/firearms_4.mp3"));
  this.load.audio('firearms_5', require("../assets/audio/sound_effects/fire_weapon/firearms_5.mp3"));
  this.load.audio('firearms_6', require("../assets/audio/sound_effects/fire_weapon/firearms_6.mp3"));

  this.load.audio('missile_1', require("../assets/audio/sound_effects/fire_weapon/missile_1.mp3"));
  this.load.audio('missile_2', require("../assets/audio/sound_effects/fire_weapon/missile_2.mp3"));
  this.load.audio('missile_3', require("../assets/audio/sound_effects/fire_weapon/missile_3.mp3"));

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

  this.load.audio('engine_4', require("../assets/audio/sound_effects/engiges/engine_4.mp3"));
  this.load.audio('engine_6', require("../assets/audio/sound_effects/engiges/engine_6.mp3"));
  this.load.audio('engine_7', require("../assets/audio/sound_effects/engiges/engine_7.mp3"));
  this.load.audio('engine_12', require("../assets/audio/sound_effects/engiges/engine_12.mp3"));
  this.load.audio('engine_14', require("../assets/audio/sound_effects/engiges/engine_14.mp3"));

}

export {preload}
