<template>
  <div id="CreateGameBlock" ref="CreateGameBlock" @mousedown="toUp">
    <app-control v-bind:head="'Создание игры'"
                 v-bind:move="true"
                 v-bind:close="true"
                 v-bind:close-func="removeLobbySession"
                 v-bind:noHeight="true"
                 v-bind:noWidth="true"
                 v-bind:refWindow="'CreateGameBlock'"/>

    <div v-if="!lobbyState">Что то пошло не так</div>

    <template v-if="lobbyState">
      <div class="players">
        <h3 class="players_head">Игроки:</h3>

        <app-user-line v-for="p in lobbyState.players" v-bind:user="p"/>
      </div>

      <div class="settings">
        <div class="map_list">

          <div v-if="lobbyState" class="map_line" v-for="game_map in maps" @click="selectGameMap(game_map.id)">
            {{ game_map.name }}
          </div>
        </div>

        <div>
          <h3 class="map_head">{{ selectMap }}</h3>
          <div class="map_icon"></div>
          <input type="button" value="Начать" @click="StartGame">
          <input type="button" value="Закрыть">
        </div>
      </div>
    </template>
  </div>
</template>

<script>
import Control from '../Window/Control';
import UserLine from "../UserLine/UserLine";

export default {
  name: "CreateGame",
  mounted() {
    this.$store.dispatch("sendSocketData", JSON.stringify({
      event: "CreateLobbySession",
      service: "lobby",
    }))
  },
  methods: {
    toUp() {
      this.$store.commit({
        type: 'setWindowZIndex',
        id: this.$el.id,
      });
    },
    removeLobbySession() {
      // TODO
    },
    selectGameMap(id) {
      this.$store.dispatch("sendSocketData", JSON.stringify({
        event: "SelectLobbyMap",
        service: "lobby",
        id: id,
        uuid: this.lobbyState.uuid,
      }))
    },
    StartGame() {
      this.$store.dispatch("sendSocketData", JSON.stringify({
        event: "StartGame",
        service: "lobby",
        uuid: this.lobbyState.uuid,
      }))
    }
  },
  computed: {
    lobbyState() {
      return this.$store.getters.getLobbyState;
    },
    maps() {
      return this.$store.getters.getShortInfoMaps;
    },
    selectMap() {
      for (let id in this.maps) {
        if (Number(id) === Number(this.lobbyState.map_id)) {
          return this.maps[id].name
        }
      }

      return "Не выбрано"
    }
  },
  components: {
    AppControl: Control,
    AppUserLine: UserLine,
  }
}
</script>

<style scoped>
#CreateGameBlock {
  position: absolute;
  display: block;
  border-radius: 5px;
  z-index: 11;
  width: 402px;
  height: 165px;
  top: 50px;
  left: calc(50% - 200px);
  padding: 19px 1px 1px;
  color: #0cc2fb;
  text-shadow: 0 -1px 1px #000000, 0 -1px 1px #000000, 0 1px 1px #000000, 0 1px 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000;
  border: 1px solid #25a0e1;
  background: rgb(8, 138, 210);
  box-shadow: 0 0 2px black;
}

.players, .map_list {
  float: left;
  height: 100%;
  width: 150px;
  box-shadow: inset 0 0 2px #000;
  border-radius: 5px;
  background: rgba(0, 0, 0, .2);
  margin-right: 2px;
  overflow-y: scroll;
  overflow-x: hidden;
  position: relative;
}

.map_list {
  width: 135px;
}

.settings {
  margin: 0 auto;
  overflow-y: scroll;
  height: calc(100% - 2px);
  box-shadow: inset 0 0 1px #000;
  background: #89969c;
  pointer-events: auto;
  overflow-x: hidden;
  width: 245px;
  text-align: left;
  font-size: 13px;
  border: 1px solid rgba(0, 0, 0, .4);
  border-radius: 5px;
}

.players_head, .map_head {
  margin: 2px 2px 3px;
  background: #dd7034;
  color: hsla(0, 0%, 100%, .8);
  border-radius: 4px;
  font-size: 13px;
  line-height: 17px;
  height: 17px;
  user-select: none;
  text-shadow: 1px 1px 1px #000;
  font-weight: 700;
  box-shadow: 0 0 2px #000;
  text-align: left;
  text-indent: 10px;
}

.map_head {
  float: left;
  width: 90px;
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
}

.map_icon {
  float: left;
  width: 90px;
  height: 90px;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 4px;
  box-shadow: 0 0 2px 0 black;
  margin: 3px 3px 3px 3px;
}

.map_line {
  transition: .1s;
  height: 15px;
  font-size: 10px;
  line-height: 15px;
  text-align: left;
  background: #5d5d5d57;
  padding-left: 10px;
  color: #ffe510;
  cursor: pointer;
  text-shadow: 1px 1px 1px #000;
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
}

.map_line:active {
  background: hsla(0, 0%, 100%, .73);
  color: #ff4208;
}

.map_line:hover {
  background: hsla(0, 0%, 92%, .34);
  color: #ff7800;
}


input[type=button]:hover {
  background: #ff8100;
}

input[type=button]:active {
  transform: scale(.98);
}

input[type=button] {
  display: block;
  width: 94px;
  margin: 4px 0 0 0;
  pointer-events: auto;
  font-size: 9px;
  text-align: center;
  transition: .1s;
  background: rgba(255, 129, 0, .6);
  height: 16px;
  border-radius: 5px;
  color: #fff;
  line-height: 15px;
  box-shadow: 0 0 2px #000;
  float: left;
}
</style>
