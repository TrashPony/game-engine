const getDefaultState = () => {
  return {
    token: "",
    wsConnectState: {
      connect: false,
      error: false,
      ws: null,
    },
    EndLoad: {
      visible: false,
      text: '',
    },
    OpenDialog: {
      page: null,
      visited_pages: null,
    },
    InitGame: false,
    Interface: {
      resolution: $(window).width() + ':' + $(window).height(),
      state: {},
      allowIDs: {},
      openQueue: [],
    },
    NeedOpenComponents: {
      /** в общем т.к. тут один поток данных, а в моем интерфейсы мелкие компоненты должны уметь открыть
       более высокий уровень (модальники), то вот эта штука заставляет открывать компоненты
       тут хранятся обьекты с полями id: component:, component определяет тип окна **/
    },
    ServerTime: 0,
    GameUser: {},
    EquipPanel: null,
    CurrentPlayer: {},
    UserPlayers: [],
    PossibleChangeName: -1,
    Settings: {
      playMusic: true,
      FollowCamera: true,
      ZoomCamera: 1,
      volume: 0.33,
      SFXVolume: 0.33,
      Language: "RU",
      UnitTrack: 1,
      ObjectShadows: true,
      Clouds: true,
      MoveAnimate: true,
      EffectsCut: 1,
      MoveAndRotateTween: true,
    },
    Chat: {
      history: [],
      avatars: {},
      private: "",
      private_force: false,
      friends: [],
      npc_animate: false,
      npc_animate_timeout: null,
    },
    UsersStat: {
      users: {},
    },
    Group: null,
    LobbyState: null,
    ShortInfoMaps: {},
    Game: {
      role: "",
    },
    unit: {
      hp: 0,
      power: 0,
    },
    Battle: {},
    PointsNotify: [],
    WaitGame: {},
    Inventory: {
      filters: {category: "all"},
      cellSize: 50,
    },
    WorkBench: {
      blueprints: {},
      works: [],
      takes: [],
    },
    Hangar: {},
    HandBook: {},
    Credits: 0,
    Notifications: {},
    Market: {
      assortment: {
        weapon: null,
        body: null,
        detail: null,
        equip: null,
        other: null,
      },
      orders: null,
      my_orders: null,
      filters: {
        // указывает на фильтр по которому будет ариентироватся таблицы заказов
        selectType: '', // ['ammo', 'body', 'weapon', 'equip', 'resource', 'blueprints', 'boxes', 'trash']
        item: null,
        equip: {
          type: '', // ['laser', 'missile', 'firearms']
          size: 0,  // [1, 2, 3]
          id: 0,
        },
        weapon: {
          type: '', // ['laser', 'missile', 'firearms']
          size: 0,  // [1, 2, 3]
          id: 0,
        },
        body: {
          type: '',
          size: 0,  // [1, 2, 3]
          id: 0,
        },
        detail: {
          type: '',
          size: 0,
          id: 0,
        },
        other: {
          type: '',
          size: 0,
          id: 0,
        }
      }
    },
    Suicide: {
      current: -1,
      deadTime: -1,
    },
    BattlesState: {},
    Skins: {},
    MapEditor: {
      maps: null,
      typeCoordinates: null,
    },
    RewardResource: "steel",
    GameType: "all_pvp",
    SocialMechanics: null,
    CheckViewPort: 1,
    Missions: [],
    ServerState: {
      Common: {},
      PlayerLogs: [],
      MarketLogs: [],
      VkOrders: [],
    },
    Operations: [],
    SelectOperationID: "",
  }
};

export {getDefaultState}
