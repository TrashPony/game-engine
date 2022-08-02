<template>
  <div id="Alerts" ref="Alerts" v-if="Object.keys(notifications).length > 0" @mousedown="toUp">
    <div class="AlertsWrapper">
      <div class="control_wrapper">
        <app-control v-bind:head="'Оповещения'"
                     v-bind:move="true"
                     v-bind:close="false"
                     v-bind:noHeight="true"
                     v-bind:noWidth="true"
                     v-bind:refWindow="'Alerts'"/>
      </div>

      <template v-for="notification in notifications">
        <div class="alert" v-html="notification.html"/>
      </template>
    </div>
  </div>
</template>

<script>
import Control from '../Window/Control';

export default {
  name: "Alerts",
  components: {
    AppControl: Control,
  },
  methods: {
    playSound(sound, k) {
      if (sound === "button_press.mp3") k = 0.2
      if (sound === "select_sound.mp3") k = 0.1

      this.$store.dispatch('playSound', {
        sound: sound,
        k: k,
      });
    },
    toUp() {
      this.$store.commit({
        type: 'setWindowZIndex',
        id: this.$el.id,
      });
    },
  },
  computed: {
    notifications() {
      return this.$store.getters.getNotifications
    }
  }
}
</script>

<style scoped>
#Alerts {
  position: absolute;
  min-height: 0;
  min-width: 150px;
  max-width: 250px;
  top: 5px;
  left: 385px;
  padding: 18px 3px 3px 3px;
  /*pointer-events: none;*/
  border: 1px solid rgba(37, 160, 225, 0.2);
  background: rgba(8, 138, 210, 0.2);
  border-radius: 5px;
  filter: drop-shadow(0 0 2px black);
}

.AlertsWrapper {
  overflow-y: scroll;
  max-height: 100px;
  padding-right: 2px;
}

.alert {
  min-height: 0;
  box-shadow: inset 0 0 2px black;
  padding: 5px;
  width: calc(100% - 12px);
  transition: 0.2s;
  margin-top: 3px;
  font-size: 10px;
  font-weight: 900;
  border: 1px solid rgba(37, 160, 225, 0.5);
  background: rgba(8, 138, 210, 0.5);
  border-radius: 3px;
  color: white;
  text-shadow: 0 -1px 1px #000000, 0 -1px 1px #000000, 0 1px 1px #000000, 0 1px 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000;
}

.inputs {
  margin-top: 3px;
}

.inputs input {
  transition: 100ms;
  background: rgba(255, 129, 0, 0.75);
  color: white;
}

.inputs input:hover {
  background: rgba(255, 129, 0, 1);
}

.inputs input:active {
  transform: scale(0.98);
}

.inputs input:active {
  transform: scale(0.96);
}

.control_wrapper {
  opacity: 0.5;
}

.control_wrapper:hover {
  opacity: 1;
}
</style>
