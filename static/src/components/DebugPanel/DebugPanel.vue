<template>
  <div id="DebugPanel" ref="DebugPanel" @mousedown="toUp">
    <app-control v-bind:head="'Дебах'"
                 v-bind:move="true"
                 v-bind:close="false"
                 v-bind:no-height="true"
                 v-bind:no-width="true"
                 v-bind:no-pos="true"
                 v-bind:refWindow="'DebugPanel'"/>

    <div style="width: 120px; margin: 2px; float: left">
      <h3>Игрок</h3>
      <input type="button" value="создать юнита" @click="CreateUnit">
      <input type="button" value="создать турель" @click="CreateTurret">
    </div>

    <div style="width: 120px; margin: 2px; float: left">
      <h3>Команда</h3>
      <input type="button" value="создать бота" @click="CreateBot(1)">
      <input type="button" value="создать турель" @click="CreateTurret(1)">
    </div>

    <div style="width: 120px; margin: 2px; float: left">
      <h3>Враги</h3>
      <input type="button" value="создать бота" @click="CreateBot(2)">
      <input type="button" value="создать турель" @click="CreateTurret(2)">
    </div>
  </div>
</template>

<script>
import Control from '../Window/Control';

export default {
  name: "DebugPanel",
  methods: {
    toUp() {
      this.$store.commit({
        type: 'setWindowZIndex',
        id: this.$el.id,
      });
    },
    CreateUnit() {
      this.$store.dispatch("sendSocketData", JSON.stringify({
        event: "CreateUnit",
        service: "battle",
      }))
    },
    CreateBot(teamID) {
      this.$store.dispatch("sendSocketData", JSON.stringify({
        event: "CreateBot",
        service: "battle",
        id: teamID,
      }))
    },
    CreateTurret(teamID) {
      this.$store.dispatch("sendSocketData", JSON.stringify({
        event: "CreateObj",
        service: "battle",
        id: teamID,
      }))
    },
  },
  components: {
    AppControl: Control,
  }
}
</script>

<style scoped>
#DebugPanel {
  background: rgb(8, 138, 210);
  box-shadow: 0 1px 2px rgba(0, 0, 0, .2);
  border: 1px solid #25a0e1;
}

#DebugPanel {
  position: absolute;
  height: 75px;
  width: 400px;
  border-radius: 5px;
  top: 5px;
  left: 5px;
  user-select: none;
  padding: 19px 2px 2px;
}

#DebugPanel input[type="button"] {
  display: block;
  margin: 2px auto 0;
  width: 100%;
  pointer-events: auto;
  font-size: 13px;
  text-align: center;
  transition: .1s;
  background: rgba(255, 129, 0, .6);
  height: 18px;
  border-radius: 5px;
  color: #fff;
  line-height: 15px;
  box-shadow: 0 0 2px #000;
  cursor: pointer;
  font-weight: 700;
}

#DebugPanel input[type="button"]:hover {
  background: #ff8100;
}

#DebugPanel input[type="button"]:active {
  transform: scale(.98);
}

#DebugPanel h3 {
  margin: 3px 0;
  background: #dd7034;
  color: hsla(0, 0%, 100%, .8);
  border-radius: 4px;
  font-size: 13px;
  line-height: 17px;
  height: 16px;
  user-select: none;
  text-shadow: 1px 1px 1px #000;
  font-weight: 700;
  box-shadow: 0 0 2px #000;
  padding-left: 10px;
  text-align: left;
  clear: both;
}
</style>
