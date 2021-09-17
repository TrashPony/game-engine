//** т.к. нам не нужна реактивность вью для игры то у нас свой стор с блекджеком и филдами **//
const gameStore = {
  noLoader: false,
  gameInit: false,
  gameTypes: {},
  gameDataInit: {
    data: false,
    sendRequest: false,
  },
  gameReady: false,
  unitReady: false,
  reload: false,
  mapDrawing: false,

  HoldAttackMouse: false,
  map: null, //** тут лежит карта с глобальной коорджинатной позициией 0:0, говорит о том что игрок живет тут **//
  maps: null,
  mapsState: {},
  StatusLayer: {
    healBoxes: {},
    bars: {},
    barsCache: {},
  },
  bmdTerrain: null,
  flore: null,
  cashTextures: {}, // движовский checkKey() работает через задницу, поэтому это тут

  // игровые обьекты
  objects: {},
  staticObjects: [],
  removeObjects: {}, // если нужно получить доступ к обьекту после удаления

  user_squad_id: null,
  user: null,

  units: {},
  bullets: {},
  clouds: {},
  lights: [],
  cacheAnimate: {},
  unitSpriteSize: 128,
  mouseInGameChecker: {
    updater: null,
    time: 0
  },
  /** что бы не создавать постоянно новые спрайты после смерти пули попадают сюда **/
  /** [bullet_sprite_name] []bullets **/
  cacheSpriteBullets: {},
  spawns: {},
  selectUnits: [],
  inputType: "wasd", // rts_single_unit, wasd
};

export {gameStore}
