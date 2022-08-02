import {urls} from "../const";
import createSocketPlugin from "./ws/socket";
import store from './store'
import Vue from 'vue'
import {gameStore} from "../game/store";

const mutations = {
  reconnectWS(state) {

    let token = state.token
    if (token === "") token = window.localStorage.getItem(urls.authTokenKey)
    if (urls.siteUrl.includes('yandex')) if (!token) return;

    state.wsConnectState.pending = true
    const WS = new WebSocket(urls.socketURL + "?token=" + token);
    let plugin = createSocketPlugin(WS);
    plugin(store);
  },
  closeWS(state) {
    state.wsConnectState.ws.noRedirect = true
    state.wsConnectState.ws.close()
  },
  setWSConnectState(state, payload) {
    state.wsConnectState.connect = payload.connect;
    state.wsConnectState.error = payload.error;
    state.wsConnectState.pending = payload.pending;
    state.wsConnectState.ws = payload.ws;
  },
  setVisibleLoader(state, payload) {
    Vue.set(state.EndLoad, 'visible', payload.visible);
    Vue.set(state.EndLoad, 'text', payload.text);
  },
  setGameUserSate(state, payload) {
    state.GameUser = payload.gameUser;
    state.UserPlayers = payload.userPlayers;
    state.CurrentPlayer = payload.player;
  },
  setGameSettings(state, payload) {
    Vue.set(state.Settings, 'playMusic', payload.settings.play_music)
    Vue.set(state.Settings, 'volume', payload.settings.volume_music / 100)
    Vue.set(state.Settings, 'SFXVolume', payload.settings.volume_sfx / 100)
    Vue.set(state.Settings, 'FollowCamera', payload.settings.follow_camera)
    Vue.set(state.Settings, 'ZoomCamera', payload.settings.zoom_camera / 100)
    Vue.set(state.Settings, 'UnitTrack', payload.settings.unit_track)
  },

  setLobbyState(state, payload) {
    state.LobbyState = payload.state
  },
  setShortInfoMaps(state, payload) {
    state.ShortInfoMaps = payload.maps
  },
  setPlayerRole(state, payload) {
    Vue.set(state.Game, 'role', payload.role);
  },
  setBattleState(state, payload) {
    Vue.set(state, 'Battle', payload.battle);
  },
  setBattleCommonData(state, payload) {
    if (state.Battle) {
      if (state.Battle.time_out) Vue.set(state.Battle, 'time_out', payload.time_out);
      if (state.Battle.wait_time_out) Vue.set(state.Battle, 'wait_time_out', payload.wait_time_out);
      if (state.Battle.wait_ready) Vue.set(state.Battle, 'wait_ready', payload.wait_ready);
      if (state.Battle) Vue.set(state.Battle, 'end', payload.end);
      Vue.set(state.Battle, 'data', true);
      Vue.set(state.Battle, 'customData', payload.customData);
      Vue.set(state.Battle, 'main_mission', payload.mainMission);
    }
  },
  setBattleTeamData(state, payload) {
    if (state.Battle && state.Battle.teams) {
      for (let i in state.Battle.teams) {
        if (Number(i) === Number(payload.id)) {
          Vue.set(state.Battle.teams[i], 'winner', payload.winner);
          Vue.set(state.Battle.teams[i], 'points', payload.points);
          Vue.set(state.Battle.teams[i], 'alive', payload.alive);
        }
      }
    }
  },
  setBattleBaseData(state, payload) {
    if (state.Battle && state.Battle.bases) {
      for (let i in state.Battle.bases) {
        if (state.Battle.bases[i] && state.Battle.bases[i].id === payload.id) {
          Vue.set(state.Battle.bases[i], 'capture_team', payload.capture_team);
          Vue.set(state.Battle.bases[i], 'capture', payload.capture);
          Vue.set(state.Battle.bases[i], 'capture_fact', payload.capture_fact);
        }
      }
    }
  },
  setBattlePlayersData(state, payload) {
    if (state.Battle && state.Battle.session_players_state) {
      for (let i in state.Battle.session_players_state) {
        if (Number(i) === payload.id) {
          Vue.set(state.Battle.session_players_state[i], 'winner', payload.winner);
          Vue.set(state.Battle.session_players_state[i], 'capture', payload.capture);
          Vue.set(state.Battle.session_players_state[i], 'deaths', payload.deaths);
          Vue.set(state.Battle.session_players_state[i], 'count_respawn', payload.count_respawn);
          Vue.set(state.Battle.session_players_state[i], 'respawn_time', payload.respawn_time);
          Vue.set(state.Battle.session_players_state[i], 'kills', payload.kills);
          Vue.set(state.Battle.session_players_state[i], 'assist', payload.assist);
          Vue.set(state.Battle.session_players_state[i], 'damage', payload.damage);
          Vue.set(state.Battle.session_players_state[i], 'points', payload.points);
          Vue.set(state.Battle.session_players_state[i], 'intelligence_damage', payload.intelligence_damage);
          Vue.set(state.Battle.session_players_state[i], 'live', payload.live);
          Vue.set(state.Battle.session_players_state[i], 'leave', payload.leave);
        }
      }
    }
  },
  setBattleLoses(state, payload) {
    Vue.set(state.Battle, 'loses', payload.loses);
  },
  setBattleReward(state, payload) {
    Vue.set(state.Battle, 'reward', {
      items: payload.reward,
      points: payload.points,
      players_state: payload.players_state,
      teams: payload.teams,
      end_page: payload.end_page,
      avatar: payload.avatar,
    });
  },
  resetBattleReward(state, payload) {
    Vue.set(state.Battle, 'reward', null)
  },
  addPointsNotify(state, payload) {
    state.PointsNotify.push(payload.notify)
  },
  SetWaitGame(state, payload) {
    Vue.set(state.WaitGame, 'playerState', payload.wp)
    Vue.set(state.WaitGame, 'countPlayers', payload.cp)
  },
  SetInventorySlots(state, payload) {
    Vue.set(state.Inventory, 'slots', payload.slots)
  },
  SetHangarSlots(state, payload) {
    Vue.set(state.Hangar, 'slots', payload.slots)
    Vue.set(state.Hangar, 'select_slot', payload.select_slot)
    gameStore.hangar.slots = payload.slots;
    gameStore.hangar.select_slot = payload.select_slot;
  },
  SetSelectHangarSlot(state, payload) {
    Vue.set(state.Hangar, 'select_slot', payload.select_slot)
    gameStore.hangar.select_slot = payload.select_slot;
  },
  setInventoryFilters(state, payload) {
    Vue.set(state.Inventory, 'filters', payload.filters);
  },
  setBlueprints(state, payload) {
    Vue.set(state.WorkBench, 'blueprints', payload.blueprints);
  },
  setWorks(state, payload) {
    Vue.set(state.WorkBench, 'works', payload.works);
  },
  addTakeWork(state, payload) {
    state.WorkBench.takes.push(payload.workID)
  },
  removeTakeWork(state, payload) {
    const index = state.WorkBench.takes.indexOf(payload.workID);
    if (index > -1) {
      state.WorkBench.takes.splice(index, 1);
    }

    let workIndex = -1
    for (let i in state.WorkBench.works) {
      if (state.WorkBench.works[i].id === payload.workID) {
        workIndex = i
        break
      }
    }

    if (workIndex > -1) {
      state.WorkBench.works.splice(workIndex, 1);
    }
  },
  setCredits(state, payload) {
    state.Credits = payload.credits
  },
  setMarketFilter(state, payload) {
    Vue.set(state.Market.filters, 'selectType', payload.main);
    Vue.set(state.Market.filters, 'item', payload.item);
    if (payload.main !== '') {
      Vue.set(state.Market.filters[payload.main], 'type', payload.filterType);
      Vue.set(state.Market.filters[payload.main], 'size', payload.size);
      Vue.set(state.Market.filters[payload.main], 'id', payload.id);
    }
  },
  setMarketAssortment(state, payload) {
    Vue.set(state.Market.assortment, 'body', payload.assortment.body);
    Vue.set(state.Market.assortment, 'detail', payload.assortment.detail);
    Vue.set(state.Market.assortment, 'equip', payload.assortment.equip);
    Vue.set(state.Market.assortment, 'weapon', payload.assortment.weapon);
    Vue.set(state.Market.assortment, 'other', payload.assortment.other);
  },
  setOrders(state, payload) {
    Vue.set(state.Market, 'orders', payload.orders);
  },
  setMyOrders(state, payload) {
    Vue.set(state.Market, 'my_orders', payload.my_orders);
  },
  addMyOrder(state, payload) {
    if (state.Market.my_orders) state.Market.my_orders.push(payload.orderID)
  },
  addOrder(state, payload) {
    if (state.Market.orders) {
      Vue.set(state.Market.orders, payload.order.Id, payload.order);
    }
  },
  removeOrder(state, payload) {
    if (state.Market.orders) {
      Vue.delete(state.Market.orders, payload.order.Id)
    }
  },
  updateOrders(state, payload) {
    if (state.Market.orders) {
      for (let order of payload.orders) {
        if (order.Count === 0) {
          Vue.delete(state.Market.orders, order.Id)
        } else {
          Vue.set(state.Market.orders, order.Id, order);
        }
      }
    }
  },
  /** WINDOWS MANAGER **/
  setWindowState(state, payload) {

    if (!state.Interface.state) {
      state.Interface.state = {}
    }

    let newResolution = $(window).width() + ':' + $(window).height()
    if (!state.Interface.state[newResolution]) {
      if (!state.Interface.state[state.Interface.resolution]) {
        Vue.set(state.Interface.state, newResolution, {});
      } else {
        Vue.set(state.Interface.state, newResolution, state.Interface.state[state.Interface.resolution]);
      }
    }

    Vue.set(state.Interface, "resolution", newResolution);
    Vue.set(state.Interface.state[state.Interface.resolution], payload.id, payload.state);

    if (payload.state.open) {
      this.commit({
        type: 'setWindowZIndex',
        id: payload.id,
      });
    } else {
      if (!state.Interface.allowIDs.hasOwnProperty(payload.id)) {
        delete state.Interface.state[state.Interface.resolution][payload.id]
      }
    }

  },
  setWindowZIndex(state, payload) {
    let windowIndex = state.Interface.openQueue.indexOf(payload.id);
    if (windowIndex > -1) {
      state.Interface.openQueue.splice(windowIndex, 1);
    }

    /** это немного не повьюшному но вот как то так, это позволяет выводить открываемые/выделяемые окна на передний план **/
    if (state.Interface.state && state.Interface.state[state.Interface.resolution]) {
      for (let id in state.Interface.state[state.Interface.resolution]) {
        if (state.Interface.state[state.Interface.resolution].hasOwnProperty(id)) {
          let modal = document.getElementById(id);
          if (modal) {
            let z = state.Interface.openQueue.indexOf(id) + 1;
            if (z < 1) {
              z = 1;
            }

            modal.style.zIndex = z;
          }
        }
      }
    }

    state.Interface.openQueue.push(payload.id)
    if (state.Interface.openQueue.length > 15) {
      state.Interface.openQueue.shift();
    }

    let modal = document.getElementById(payload.id);
    if (modal) {
      modal.style.zIndex = state.Interface.openQueue.indexOf(payload.id) + 1;
    }
  },
  toggleWindow(state, payload) {

    if (state.NeedOpenComponents.hasOwnProperty(payload.id)) {
      payload.open = !state.NeedOpenComponents[payload.id].open

      // если окно на заднем плане то делаем его на передний
      if (state.Interface.openQueue[state.Interface.openQueue.length - 1] !== payload.id) {
        payload.open = true;
      }
    } else {
      payload.open = true;
    }

    if (payload.forceOpen) {
      payload.open = true;
    }

    if (payload.forceClose) {
      payload.open = false;
    }

    let window = {
      id: payload.id,
      component: payload.component,
      open: payload.open,
      meta: payload.meta,
    };

    if (payload.open) {
      this.commit({
        type: 'setWindowZIndex',
        id: payload.id,
      });
    } else {
      state.Interface.openQueue.splice(state.Interface.openQueue.indexOf(payload.id), 1);
    }

    Vue.set(state.NeedOpenComponents, payload.id, window);

    if (!payload.open && !state.Interface.allowIDs.hasOwnProperty(payload.id)) {
      Vue.delete(state.NeedOpenComponents, payload.id)
    }
  },
  removeAllWindowsByComponentName(state, payload) {
    for (let i in state.NeedOpenComponents) {
      if (state.NeedOpenComponents.hasOwnProperty(i) && state.NeedOpenComponents[i].component === payload.component) {
        this.commit({
          type: 'toggleWindow',
          id: state.NeedOpenComponents[i].id,
          component: state.NeedOpenComponents[i].component,
          forceClose: true,
        });
      }
    }
  },
  setInterfaceState(state, payload) {

    Vue.set(state.Interface, 'state', payload.user_interface);
    Vue.set(state.Interface, 'allowIDs', payload.allow_window_save);

    if (payload.user_interface) {

      if (!state.Interface.state[state.Interface.resolution]) {
        state.Interface.state[state.Interface.resolution] = {};
      }

      for (let id in payload.user_interface[state.Interface.resolution]) {

        if (!payload.user_interface[state.Interface.resolution].hasOwnProperty(id)) continue;

        Vue.set(state.NeedOpenComponents, id, {
          id: id,
          open: payload.user_interface[state.Interface.resolution][id].open
        });
      }
    }
  },
  setUnitHP(state, payload) {
    Vue.set(state.unit, 'hp', payload.hp);
  },
  setUnitEnergy(state, payload) {
    Vue.set(state.unit, 'power', payload.energy);
  },
  setEquipPanel(state, payload) {
    Vue.set(state, 'EquipPanel', payload.equipPanel);
  },
  /** Volume, SFXVolume **/
  /** Volume, SFXVolume **/
  setVolume(state, payload) {
    Vue.set(state.Settings, 'volume', payload.volume)
  },
  setSFXVolume(state, payload) {
    Vue.set(state.Settings, 'SFXVolume', payload.SFXVolume)
  },
  setUnitTrack(state, payload) {
    Vue.set(state.Settings, 'UnitTrack', payload.UnitTrack)
  },
  setLanguage(state, payload) {
    Vue.set(state.Settings, 'Language', payload.language)
  },
  setFollowCamera(state, payload) {
    Vue.set(state.Settings, 'FollowCamera', payload.follow)
  },
  setZoomCamera(state, payload) {
    Vue.set(state.Settings, 'ZoomCamera', payload.zoom)
  },
  setHandBook(state, payload) {
    Vue.set(state, 'HandBook', payload.description_items)
  },
  addNotification(state, payload) {
    let app = this;

    app.dispatch('playSound', {
      sound: 'message.mp3',
      k: 1,
    });

    if (!state.Notifications.hasOwnProperty(payload.id)) {
      if (payload.removeSec) {

        payload.newHtml = payload.html.replace(/time_place/gi, payload.removeSec);
        Vue.set(state.Notifications, payload.id, {html: payload.newHtml, time: payload.removeSec});

        let time = setInterval(function () {
          if (state.Notifications.hasOwnProperty(payload.id) && state.Notifications[payload.id].time > 0) {

            let newTime = state.Notifications[payload.id].time - 1;
            payload.newHtml = payload.html.replace(/time_place/gi, newTime);
            Vue.set(state.Notifications, payload.id, {html: payload.newHtml, time: newTime, input: payload.input});
          } else {
            app.commit({
              type: 'removeNotification',
              id: payload.id,
            });
            clearInterval(time)
          }

        }, 1000)
      } else {
        Vue.set(state.Notifications, payload.id, {html: payload.html});
      }

      if (payload.input) {
        Vue.set(state.Notifications[payload.id], 'input', payload.input);
      }
    }
  },
  removeNotification(state, payload) {
    Vue.delete(state.Notifications, payload.id);
  },
  addAvatar(state, payload) {
    Vue.set(state.Chat.avatars, payload.id, payload.avatar);
  },
  removeAvatar(state, payload) {
    Vue.delete(state.Chat.avatars, payload.id);
  },
  setChatHistory(state, payload) {
    Vue.set(state.Chat, 'history', payload.history);
  },
  addChatMsg(state, payload) {
    state.Chat.history.push(payload.msg)
    if (payload.npc_animate) {

      this.dispatch('playSound', {
        sound: 'npc_message.mp3',
        k: 0.6,
      });

      Vue.set(state.Chat, 'npc_animate', true);
      clearTimeout(state.Chat.npc_animate_timeout)
      state.Chat.npc_animate_timeout = setTimeout(function () {
        Vue.set(state.Chat, 'npc_animate', false);
      }, 15000)
    }
  },
  setGroupState(state, payload) {
    Vue.set(state, 'Group', payload.state);
  },
  setPrivateChat(state, payload) {
    Vue.set(state.Chat, 'private', payload.user_name);
    Vue.set(state.Chat, 'private_force', payload.force);
  },
  addUserState(state, payload) {
    Vue.set(state.UsersStat.users, payload.user.user_id, payload.user);
  },
  setFriends(state, payload) {
    Vue.set(state.Chat, 'friends', payload.friends);
  },
  setInitGame(state, payload) {
    Vue.set(state, 'InitGame', payload.init);
  },
  setSuicide(state, payload) {
    Vue.set(state.Suicide, 'current', payload.current);
    Vue.set(state.Suicide, 'deadTime', payload.deadTime);
  },
  setBattlesState(state, payload) {
    Vue.set(state, 'BattlesState', payload.data);
  },
  setSkins(state, payload) {
    Vue.set(state, 'Skins', payload.data);
  },
  setMapsInMapEditor(state, payload) {
    Vue.set(state.MapEditor, 'maps', payload.maps)
  },
  setTypeCoordinates(state, payload) {
    Vue.set(state.MapEditor, 'typeCoordinates', payload.typeCoordinates)
  },
  setSocialMechanics(state, payload) {
    Vue.set(state, 'SocialMechanics', payload.data)
  },
  setSocialMechanicRequest(state, payload) {

    if (!state.SocialMechanics) return
    if (!state.SocialMechanics[payload.typeSM]) {
      state.SocialMechanics[payload.typeSM] = {}
    }

    Vue.set(state.SocialMechanics[payload.typeSM], 'request', payload.request)
  },
  setRewardResource(state, payload) {
    Vue.set(state, 'RewardResource', payload.resource)
  },
  setGameType(state, payload) {
    Vue.set(state, 'GameType', payload.typeGame)
  },
  addCheckViewPort(state, payload) {
    Vue.set(state, 'CheckViewPort', state.CheckViewPort + 1)
  },
  setMissions(state, payload) {
    Vue.set(state, 'Missions', payload.data)
  },
  setResolution(state, payload) {
    Vue.set(state.Interface, 'resolution', $(window).width() + ':' + $(window).height());
  },
  setLogs(state, payload) {
    if (payload.logType === 'market') {
      Vue.set(state.ServerState, 'MarketLogs', payload.data);
    }

    if (payload.logType === 'player') {
      Vue.set(state.ServerState, 'PlayerLogs', payload.data);
    }

    if (payload.logType === 'vk') {
      Vue.set(state.ServerState, 'VkOrders', payload.data);
    }
  },
  setServerState(state, payload) {
    Vue.set(state.ServerState, 'Common', payload.data);
  },
  setPossibleChangeName(state, payload) {
    Vue.set(state, 'PossibleChangeName', payload.data);
  },
  setToken(state, payload) {
    Vue.set(state, 'token', payload.data);
  },
  setServerTime(state, payload) {
    Vue.set(state, 'ServerTime', Number(payload.data));
  },
  setOpenDialog(state, payload) {
    Vue.set(state.OpenDialog, 'page', payload.page);
    Vue.set(state.OpenDialog, 'visited_pages', payload.visited_pages);
  },
  setOperations(state, payload) {
    Vue.set(state, 'Operations', payload.operations);
  },
  setOperation(state, payload) {
    state.SelectOperationID = payload.id
  },
  setBattleEnd: function (state, payload) {
    Vue.set(state.Battle, 'end', payload.end);
  },
  setBattleRespawn: function (state, payload) {
    Vue.set(state.Battle, 'respawn', payload.respawn);
    Vue.set(state.Battle, 'count_respawn', payload.count_respawn);
    Vue.set(state.Battle, 'respawn_time', payload.respawn_time);
  },
  setBattleType: function (state, payload) {
    Vue.set(state.Battle, 'type', payload.typeBattle);
  },
  setBattleMainMission: function (state, payload) {
    Vue.set(state.Battle, 'main_mission', payload.main_mission);
  },
};

export default mutations;
