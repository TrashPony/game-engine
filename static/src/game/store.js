//** т.к. нам не нужна реактивность вью для игры то у нас свой стор с блекджеком и филдами **//
const gameStore = {
  // флажки состояния игры, TODO некоторые уже не имеют смысла
  appInit: false,
  noLoader: false,
  radarWork: false,
  gameInit: false,
  gameTypes: {},
  gameDataInit: {
    data: false,
    sendRequest: false,
  },
  exitTab: false,
  gameReady: false,
  unitReady: false,
  reload: false,
  mapDrawing: false,

  // --
  StatusLayer: {
    healBoxes: {},
    bars: {},
    barsCache: {},
    barsCacheSprites: {},
  },
  cashTextures: {}, // движовский checkKey() работает через задницу, поэтому это тут

  // игровые обьекты
  objects: {},
  staticObjects: [],
  player: null,
  units: {},
  bullets: {},
  removeObjects: {}, // если нужно получить доступ к обьекту после удаления
  map: null, //** тут лежит карта с глобальной коорджинатной позициией 0:0, говорит о том что игрок живет тут **//
  maps: null,
  mapsState: {},

  // настройки игры
  selectUnits: [],
  inputType: "wasd", // rts, rts_single_unit, wasd
  gameSettings: {},
  mapBinItems: {},
  unitSpriteSize: 128,

  // параметры инпута
  oldTargetPos: {x: 0, y: 0},
  fireState: {
    target: {
      x: 0, y: 0,
    },
    firePosition: {
      x: 0, y: 0,
    }
  },
  mouseInGameChecker: {
    updater: null,
    time: 0
  },
  HoldAttackMouse: false,
  MouseMoveInterface: false,
};

export {gameStore}
