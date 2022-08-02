<template>
  <div>
    <app-mini-map/>
    <app-debug-panel/>
  </div>
</template>

<script>
import DebugPanel from "../DebugPanel/DebugPanel";
import {CreateGame} from "../../game/create";
import {gameStore} from "../../game/store";
import MiniMap from "../MiniMap/MiniMap"

export default {
  name: "Global",
  data() {
    return {}
  },
  mounted() {
    if (gameStore.appInit)
      this.$store.dispatch("sendSocketData", JSON.stringify({
        event: "GetPlayers",
        service: "system",
      }))

    if (gameStore.appInit) {
      CreateGame();
    } else {
      this.$router.push('/gate');
    }
  },
  methods: {
    respawn() {
      this.$store.dispatch("sendSocketData", JSON.stringify({
        event: "Respawn",
        service: "battle",
      }));
    },
    openService(service, meta, component = '', forceOpen = false) {
      this.sub = '';
      this.$store.commit({
        type: 'toggleWindow',
        id: service,
        component: component,
        meta: meta,
        forceOpen: forceOpen,
      });
    },
  },
  computed: {},
  components: {
    AppDebugPanel: DebugPanel,
    AppMiniMap: MiniMap,
  }
}
</script>

<style scoped>

#displayNotify {
  position: absolute;
  width: 100%;
  pointer-events: none;
  top: 15%;
}

#dead {
  font-family: 'Audiowide', cursive;
  width: 100%;
  text-align: center;
  height: 50px;
  color: #bdbd00;
  font-size: 50px;
  text-shadow: 0 -1px 1px #000000, 0 -1px 1px #000000, 0 1px 1px #000000, 0 1px 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000;
}

.displayNotifyText {
  animation: LowGravity 2s infinite;
}

.respawnBlock {
  font-size: 16px;
  color: white;
  border-radius: 5px;
  padding: 3px 4px;
  width: 150px;
  margin: 6px auto;
  -webkit-touch-callout: none; /* iOS Safari */
  -webkit-user-select: none; /* Safari */
  -khtml-user-select: none; /* Konqueror HTML */
  -moz-user-select: none; /* Old versions of Firefox */
  -ms-user-select: none; /* Internet Explorer/Edge */
  user-select: none;
  border: 1px solid rgba(37, 160, 225, 0.8);
  background: rgba(8, 138, 210, 0.8);
}

.respawnBlock.disable {
  opacity: 0.5;
}

.respawnButton {
  background: rgba(255, 129, 0, 0.75);
  box-shadow: 0 0 2px black;
  pointer-events: all;
  transition: 100ms;
  border-radius: 5px;
}

.respawnButton.disable {
  background: rgba(84, 78, 72, 0.75);
  color: rgba(145, 145, 145, 0.75);
}

.respawnButton:hover {
  background: rgba(255, 129, 0, 1);
}

.respawnButton.disable:hover {
  background: rgba(84, 78, 72, 0.75);
}

.respawnButton:active {
  transform: scale(0.98);
}

.respawnButton.disable:active {
  transform: scale(1);
}

.countRespawn {
  color: rgba(255, 255, 255, 0.8);
  margin-top: 5px;
  font-size: 14px;
}

.hp_bar_wrapper {
  width: 300px;
  height: 14px;
  border: 1px solid #4c4c4c;
  text-align: left;
  display: block;
  box-shadow: inset 0 0 2px black, 0 0 2px black;
  border-radius: 10px;
  background-size: 12%;
  overflow: hidden;
  background-color: #959595;
  margin: 100px auto 0;
}

.hp_bar_wrapper span {
  display: block;
  width: 100%;
  text-align: center;
  text-shadow: 1px 1px 1px rgba(0, 0, 0, 1);
  font-weight: bold;
  margin: auto;
  float: left;
  color: white;
  font-size: 11px;
  line-height: 14px;
  font-family: 'Comfortaa', cursive;
}

.hp_bar_inner {
  overflow: visible;
  height: 100%;
  box-shadow: inset 0 0 2px black;
  background: rgba(255, 129, 0, 1);
}

@media (max-width: 1000px) {
  .displayNotifyText {
    font-size: 25px !important;
  }
}

.training_helpers, .training_helpers_2 {
  pointer-events: all;
  float: left;
  clear: both;
  width: 260px;
  top: 35px;
  left: calc(50% - 130px);
  position: absolute;
}

.training_helpers_2 {
  top: 75px;
  left: 7px;
  color: white;
  background: rgba(0, 0, 0, 0.35);
  padding: 3px;
  border-radius: 5px;
  border: 1px solid rgba(128, 128, 128, 0.35);
  text-shadow: 1px 1px 1px black;
  /*pointer-events: none;*/
  font-size: 14px;
}

.training_helpers h3 {
  margin: 0 0 5px 0;
  background: rgb(221, 112, 52);
  color: rgba(255, 255, 255, 0.8);
  border-radius: 4px;
  font-size: 8px;
  line-height: 12px;
  height: 10px;
  user-select: none;
  text-shadow: 1px 1px 1px rgb(0 0 0);
  font-weight: bold;
  box-shadow: 0 0 2px rgb(0 0 0);
  text-align: center;
}

.buttons {
  display: block;
  width: 100%;
  margin: 2px auto 0;
  pointer-events: auto;
  font-size: 11px;
  text-align: center;
  transition: 100ms;
  background: rgba(255, 129, 0, 0.6);
  height: 12px;
  border-radius: 5px;
  color: #fff;
  line-height: 13px;
  box-shadow: 0 0 2px #000;
  cursor: pointer;
}

.buttons:hover {
  background: rgba(255, 129, 0, 1);
}

.buttons:active {
  transform: scale(0.98);
}

.training_helpers_head {
  margin: 0 0 3px 0;
  background: rgb(221, 112, 52);
  color: rgba(255, 255, 255, 0.8);
  border-radius: 4px;
  font-size: 13px;
  line-height: 17px;
  height: 16px;
  user-select: none;
  text-shadow: 1px 1px 1px rgb(0 0 0);
  font-weight: bold;
  box-shadow: 0 0 2px rgb(0 0 0);
  padding-left: 10px;
  text-align: left;
  clear: both;
}

.training_helpers_content {

}

.training_helpers_text {

}

.training_helpers_img {
  width: 256px;
  height: 54px;
  background-size: contain;
}
</style>
