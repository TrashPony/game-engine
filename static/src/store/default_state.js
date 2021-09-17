const getDefaultState = () => {
  return {
    wsConnectState: {
      connect: false,
      error: false,
      ws: null,
    },
    EndLoad: {
      visible: false,
      text: '',
    },
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
    GameUser: {},
    CurrentPlayer: {},
    UserPlayers: [],
    Settings: {
      volume: 0.1,
      SFXVolume: 0.1,
      Language: "RU",
    },
    Chat: {
      currentChatID: 0,
      groups: [],
      users: {},
      avatars: {},
    },
    LobbyState: null,
    ShortInfoMaps: {},
    Game: {
      role: "",
    }
  }
};

export {getDefaultState}
