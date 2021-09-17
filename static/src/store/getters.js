const getters = {
  getWSConnectState: function (state) {
    return state.wsConnectState
  },
  getEndLoad: function (state) {
    return state.EndLoad
  },
  getCurrentPlayer: function (state) {
    return state.CurrentPlayer
  },
  getGameUser: function (state) {
    return state.GameUser
  },
  getUserPlayers: function (state) {
    return state.UserPlayers
  },
  getNeedOpenComponents: function (state) {
    return state.NeedOpenComponents
  },
  getChatState: function (state) {
    return state.Chat
  },
  getLobbyState: function (state) {
    return state.LobbyState
  },
  getShortInfoMaps: function (state) {
    return state.ShortInfoMaps
  },
  getSettings: function (state) {
    return state.Settings
  },
  GetColorDamage: state => percentHP => {
    if (percentHP >= 80) {
      return "#48FF28"
    } else if (percentHP < 80 && percentHP >= 75) {
      return "#fff326"
    } else if (percentHP < 75 && percentHP >= 50) {
      return "#fac227"
    } else if (percentHP < 50 && percentHP >= 25) {
      return "#fa7b31"
    } else if (percentHP < 25) {
      return "#ff2615"
    }
  },
};

export default getters;
