<template>
  <div id="DebugPanel" ref="DebugPanel" @mousedown="toUp">
    <app-control v-bind:head="'Дебах'"
                 v-bind:move="true"
                 v-bind:close="false"
                 v-bind:refWindow="'DebugPanel'"/>

    <input type="button" value="создать юнита" @click="CreateUnit">
    <input type="button" value="создать бота" @click="CreateBot">

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
    CreateBot() {
      this.$store.dispatch("sendSocketData", JSON.stringify({
        event: "CreateBot",
        service: "battle",
      }))
    }
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
  height: 200px;
  width: 200px;
  border-radius: 5px;
  top: 5px;
  right: 5px;
  user-select: none;
  padding: 19px 2px 2px;
}
</style>
