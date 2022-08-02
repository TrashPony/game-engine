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
  getGroup: function (state) {
    return state.Group
  },
  getGameUser: function (state) {
    return state.GameUser
  },
  getUserPlayers: function (state) {
    return state.UserPlayers
  },
  getCheckViewPort: function (state) {
    return state.CheckViewPort
  },
  getNeedOpenComponents: function (state) {
    return state.NeedOpenComponents
  },
  getChatState: function (state) {
    return state.Chat
  },
  getAvatars: function (state) {
    return state.Chat.avatars
  },
  getUsers: function (state) {
    return state.UsersStat
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
  getUnitState: function (state) {
    return state.unit
  },
  getBattleState: function (state) {
    return state.Battle
  },
  getPointsNotify: function (state) {
    return state.PointsNotify
  },
  getWaitGame: function (state) {
    return state.WaitGame
  },
  getInventory: function (state) {
    return state.Inventory
  },
  getHandBook: function (state) {
    return state.HandBook
  },
  getHangar: function (state) {
    return state.Hangar
  },
  getBlueprints: function (state) {
    return state.WorkBench.blueprints
  },
  getWorks: function (state) {
    return state.WorkBench.works
  },
  getTakeWorks: function (state) {
    return state.WorkBench.takes
  },
  getCredits: function (state) {
    return state.Credits
  },
  getMarket: function (state) {
    return state.Market
  },
  getNotifications: function (state) {
    return state.Notifications
  },
  getMarketFilter: function (state) {
    return state.Market.filters
  },
  GetColorDamage: state => percentHP => {
    if (percentHP >= 80) {
      return "#5ece4d"
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
  GetEquipPanel: function (state) {
    return state.EquipPanel
  },
  getInterfaceState: function (state) {
    if (!state.Interface.state) {
      return null
    }
    return state.Interface.state[state.Interface.resolution]
  },
  getInitGame: function (state) {
    return state.InitGame
  },
  getSuicide: function (state) {
    return state.Suicide
  },
  getBattlesState: function (state) {
    return state.BattlesState
  },
  getBattleEnd: function (state) {
    return state.Battle.end
  },
  getSkins: function (state) {
    return state.Skins
  },
  getMapEditorData: function (state) {
    return state.MapEditor
  },
  getRewardResource: function (state) {
    return state.RewardResource
  },
  getGameType: function (state) {
    return state.GameType
  },
  getSocialMechanics: function (state) {
    return state.SocialMechanics
  },
  getMissions: function (state) {
    return state.Missions
  },
  getServerState: function (state) {
    return state.ServerState
  },
  getPossibleChangeName: function (state) {
    return state.PossibleChangeName
  },
  getServerTime: function (state) {
    return state.ServerTime
  },
  getOpenDialog: function (state) {
    return state.OpenDialog
  },
  getOperations: function (state) {
    return state.Operations
  },
  getSelectOperation: function (state) {
    return state.SelectOperationID
  },
};

export default getters;
