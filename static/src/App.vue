<template>
  <div>

    <app-game/>
    <app-loader/>
    <router-view/>
    <app-alerts/>

    <template v-if="$route.path === '/lobby' || $route.path === '/global'">
      <app-menu-bar/>
    </template>
  </div>
</template>

<script>
import MenuBar from "./components/MenuBar/MenuBar";
import Game from "./components/Game/Game"
import Loader from "./components/Loader/Loader"
import Alerts from "./components/Alerts/Alerts";
import {gameStore} from "./game/store";
import {RemoveOldMap} from "./game/map/remove_old_map";

export default {
  name: 'app',
  data() {
    return {
      game: false,
    }
  },
  mounted() {

    let app = this;
    let playerVolume = 0;

    document.addEventListener("visibilitychange", function () {
      if (document.hidden) {

        RemoveOldMap();
        gameStore.gameReady = false;
        gameStore.gameDataInit = {
          data: false,
          sendRequest: false,
        }
        gameStore.exitTab = true

        playerVolume = app.$store.getters.getSettings.volume;
        app.$store.commit({
          type: 'setVolume',
          volume: 0,
        });

      } else {

        gameStore.exitTab = false
        app.$store.commit({
          type: 'setVolume',
          volume: playerVolume,
        });
      }
    });

    if (this.$route.path === '/lobby' || this.$route.path === '/global') {
      this.$router.push('/gate');
    }

    gameStore.appInit = true
  },
  computed: {
    openComponents() {
      return this.$store.getters.getNeedOpenComponents
    },
  },
  components: {
    AppGame: Game,
    AppLoader: Loader,
    AppMenuBar: MenuBar,
    AppAlerts: Alerts,
  }
}
</script>

<style>
@import url('https://fonts.googleapis.com/css?family=Comfortaa:300,500,700|Noto+Sans+SC:300,500,700|Source+Code+Pro:300,500,700&display=swap');
@import url('https://fonts.googleapis.com/css2?family=Audiowide&display=swap');

body {
  margin: 0;
  padding: 0;
  font-family: 'Comfortaa', cursive;
  overflow: hidden;
}

body * {
  scrollbar-color: #adc6cd rgba(1, 77, 93, 0.4);
  scrollbar-width: thin;
}

input {
  outline: none;
  padding-left: 5px;
  padding-right: 5px;
  line-height: 18px;
  box-shadow: inset 0 0 1px 1px rgb(118, 118, 118);
  border: 0;
  font-size: 11px;
  background: #79a0b4;
  border-radius: 5px;
  color: #ff7800;
  text-shadow: 1px 1px 1px black;
  cursor: pointer;
}

input[type=text], input[type=password], input[type=number] {
  background: #8d9b9e;
}

input:hover {
  background-color: rgba(0, 120, 211, 0.9);
}

input[type=submit]:active {
  transform: scale(0.95);
}

input::placeholder, textarea::placeholder {
  text-shadow: none;
}

input, textarea, select, option {
  outline: none;
}

.Failed {
  color: rgb(255, 95, 97);
}

/* CUSTOM SCROLLBAR */
body ::-webkit-scrollbar,
body ::-webkit-scrollbar:horizontal {
  width: 12px;
}

body ::-webkit-scrollbar {
  background: rgba(165, 213, 239, 0.3);
  box-shadow: inset 0 0 5px black;
}

/* DISABLED STATES */
body ::-webkit-scrollbar:disabled {
  background: rgba(165, 213, 239, 0.31);
  box-shadow: inset 0 0 5px black;
}

/* HANDLE */
body ::-webkit-scrollbar-thumb,
body ::-webkit-scrollbar-thumb:window-inactive,
body ::-webkit-scrollbar-thumb:disabled {
  box-shadow: inset 0 0 5px black;
  background-size: 5px 2px;
  background: rgba(1, 77, 93, 0.5);
}

/*--------*/
body ::-webkit-scrollbar {
  width: 12px;
  height: 12px;
}

input[type=range] {
  -webkit-appearance: none;
  width: 100%;
  margin: 0 0;
  border: 0;
}

input[type=range]:focus {
  outline: none;
}

input[type=range]::-webkit-slider-runnable-track {
  width: 100%;
  height: 16px;
  cursor: pointer;
  box-shadow: 0px 0px 0px #000000, 0px 0px 0px #0d0d0d;
  background: rgba(165, 213, 239, 0.7);
  border-radius: 0px;
  border: 0.1px solid #000000;
}

input[type=range]::-webkit-slider-thumb {
  box-shadow: 0px 0px 0px rgba(0, 0, 0, 0), 0px 0px 0px rgba(13, 13, 13, 0);
  border: 0px solid rgba(0, 0, 0, 0);
  height: 16px;
  width: 16px;
  border-radius: 0px;
  background: rgba(0, 73, 93, 0.7);
  cursor: pointer;
  -webkit-appearance: none;
  margin-top: -0.1px;
}

input[type=range]:focus::-webkit-slider-runnable-track {
  background: rgba(NaN, NaN, NaN, 0.7);
}

input[type=range]::-moz-range-track {
  width: 100%;
  height: 16px;
  cursor: pointer;
  box-shadow: 0px 0px 0px #000000, 0px 0px 0px #0d0d0d;
  background: rgba(165, 213, 239, 0.7);
  border-radius: 0px;
  border: 0.1px solid #000000;
}

input[type=range]::-moz-range-thumb {
  box-shadow: 0px 0px 0px rgba(0, 0, 0, 0), 0px 0px 0px rgba(13, 13, 13, 0);
  border: 0px solid rgba(0, 0, 0, 0);
  height: 16px;
  width: 16px;
  border-radius: 0px;
  background: rgba(0, 73, 93, 0.7);
  cursor: pointer;
}

input[type=range]::-ms-track {
  width: 100%;
  height: 16px;
  cursor: pointer;
  background: transparent;
  border-color: transparent;
  color: transparent;
}

input[type=range]::-ms-fill-lower {
  background: rgba(NaN, NaN, NaN, 0.7);
  border: 0.1px solid #000000;
  border-radius: 0px;
  box-shadow: 0px 0px 0px #000000, 0px 0px 0px #0d0d0d;
}

input[type=range]::-ms-fill-upper {
  background: rgba(165, 213, 239, 0.7);
  border: 0.1px solid #000000;
  border-radius: 0px;
  box-shadow: 0px 0px 0px #000000, 0px 0px 0px #0d0d0d;
}

input[type=range]::-ms-thumb {
  box-shadow: 0px 0px 0px rgba(0, 0, 0, 0), 0px 0px 0px rgba(13, 13, 13, 0);
  border: 0px solid rgba(0, 0, 0, 0);
  width: 16px;
  border-radius: 0px;
  background: rgba(0, 73, 93, 0.7);
  cursor: pointer;
  height: 16px;
}

input[type=range]:focus::-ms-fill-lower {
  background: rgba(165, 213, 239, 0.7);
}

input[type=range]:focus::-ms-fill-upper {
  background: rgba(NaN, NaN, NaN, 0.7);
}

input[type=range] {
  -webkit-appearance: none;
  width: 100%;
  margin: 0px 0;
}

input[type=range]:focus {
  outline: none;
}

input[type=range]::-webkit-slider-runnable-track {
  width: 100%;
  height: 16px;
  cursor: pointer;
  box-shadow: 0px 0px 0px #000000, 0px 0px 0px #0d0d0d;
  background: rgba(165, 213, 239, 0.7);
  border-radius: 0px;
  border: 0.1px solid #000000;
}

input[type=range]::-webkit-slider-thumb {
  box-shadow: 0px 0px 0px rgba(0, 0, 0, 0), 0px 0px 0px rgba(13, 13, 13, 0);
  border: 0px solid rgba(0, 0, 0, 0);
  height: 16px;
  width: 16px;
  border-radius: 0px;
  background: rgba(0, 73, 93, 0.7);
  cursor: pointer;
  -webkit-appearance: none;
  margin-top: -0.1px;
}

input[type=range]:focus::-webkit-slider-runnable-track {
  background: rgba(NaN, NaN, NaN, 0.7);
}

input[type=range]::-moz-range-track {
  width: 100%;
  height: 16px;
  cursor: pointer;
  box-shadow: 0px 0px 0px #000000, 0px 0px 0px #0d0d0d;
  background: rgba(165, 213, 239, 0.7);
  border-radius: 0px;
  border: 0.1px solid #000000;
}

input[type=range]::-moz-range-thumb {
  box-shadow: 0px 0px 0px rgba(0, 0, 0, 0), 0px 0px 0px rgba(13, 13, 13, 0);
  border: 0px solid rgba(0, 0, 0, 0);
  height: 16px;
  width: 16px;
  border-radius: 0px;
  background: rgba(0, 73, 93, 0.7);
  cursor: pointer;
}

input[type=range]::-ms-track {
  width: 100%;
  height: 16px;
  cursor: pointer;
  background: transparent;
  border-color: transparent;
  color: transparent;
}

input[type=range]::-ms-fill-lower {
  background: rgba(NaN, NaN, NaN, 0.7);
  border: 0.1px solid #000000;
  border-radius: 0px;
  box-shadow: 0px 0px 0px #000000, 0px 0px 0px #0d0d0d;
}

input[type=range]::-ms-fill-upper {
  background: rgba(165, 213, 239, 0.7);
  border: 0.1px solid #000000;
  border-radius: 0px;
  box-shadow: 0px 0px 0px #000000, 0px 0px 0px #0d0d0d;
}

input[type=range]::-ms-thumb {
  box-shadow: 0px 0px 0px rgba(0, 0, 0, 0), 0px 0px 0px rgba(13, 13, 13, 0);
  border: 0px solid rgba(0, 0, 0, 0);
  width: 16px;
  border-radius: 0px;
  background: rgba(0, 73, 93, 0.7);
  cursor: pointer;
  height: 14px;
}

input[type=range]:focus::-ms-fill-lower {
  background: rgba(165, 213, 239, 0.7);
}

input[type=range]:focus::-ms-fill-upper {
  background: rgba(NaN, NaN, NaN, 0.7);
}

input[type=range] {
  -webkit-appearance: none;
  width: 100%;
  margin: 0px 0;
}

input[type=range]:focus {
  outline: none;
}

input[type=range]::-webkit-slider-runnable-track {
  width: 100%;
  height: 16px;
  cursor: pointer;
  box-shadow: 0px 0px 0px #000000, 0px 0px 0px #0d0d0d;
  background: rgba(165, 213, 239, 0.7);
  border-radius: 0px;
  border: 0.1px solid #000000;
}

input[type=range]::-webkit-slider-thumb {
  box-shadow: 0px 0px 0px rgba(0, 0, 0, 0), 0px 0px 0px rgba(13, 13, 13, 0);
  border: 0px solid rgba(0, 0, 0, 0);
  height: 14px;
  width: 16px;
  border-radius: 0px;
  background: rgba(0, 73, 93, 0.7);
  cursor: pointer;
  -webkit-appearance: none;
  margin-top: -0.1px;
}

input[type=range]:focus::-webkit-slider-runnable-track {
  background: rgba(NaN, NaN, NaN, 0.7);
}

input[type=range]::-moz-range-track {
  width: 100%;
  height: 16px;
  cursor: pointer;
  box-shadow: 0px 0px 0px #000000, 0px 0px 0px #0d0d0d;
  background: rgba(165, 213, 239, 0.7);
  border-radius: 0px;
  border: 0.1px solid #000000;
}

input[type=range]::-moz-range-thumb {
  box-shadow: 0px 0px 0px rgba(0, 0, 0, 0), 0px 0px 0px rgba(13, 13, 13, 0);
  border: 0px solid rgba(0, 0, 0, 0);
  height: 14px;
  width: 16px;
  border-radius: 0px;
  background: rgba(0, 73, 93, 0.7);
  cursor: pointer;
}

input[type=range]::-ms-track {
  width: 100%;
  height: 16px;
  cursor: pointer;
  background: transparent;
  border-color: transparent;
  color: transparent;
}

input[type=range]::-ms-fill-lower {
  background: rgba(NaN, NaN, NaN, 0.7);
  border: 0.1px solid #000000;
  border-radius: 0px;
  box-shadow: 0px 0px 0px #000000, 0px 0px 0px #0d0d0d;
}

input[type=range]::-ms-fill-upper {
  background: rgba(165, 213, 239, 0.7);
  border: 0.1px solid #000000;
  border-radius: 0px;
  box-shadow: 0px 0px 0px #000000, 0px 0px 0px #0d0d0d;
}

input[type=range]::-ms-thumb {
  box-shadow: 0px 0px 0px rgba(0, 0, 0, 0), 0px 0px 0px rgba(13, 13, 13, 0);
  border: 0px solid rgba(0, 0, 0, 0);
  width: 16px;
  border-radius: 0px;
  background: rgba(0, 73, 93, 0.7);
  cursor: pointer;
  height: 16px;
}

input[type=range]:focus::-ms-fill-lower {
  background: rgba(165, 213, 239, 0.7);
}

input[type=range]:focus::-ms-fill-upper {
  background: rgba(NaN, NaN, NaN, 0.7);
}

.importantly {
  text-decoration: none;
  color: #faff00;
  text-shadow: 1px 1px 1px rgb(0 0 0);
  font-weight: bold;
}

@keyframes selectMenu {
  0% {
    box-shadow: 0 0 15px 0 rgb(252, 255, 0);
  }
  95% {
    box-shadow: 0 0 30px 15px rgba(204, 169, 44, 0);
  }
  100% {
    box-shadow: 0 0 0 20px rgba(204, 169, 44, 0);
  }
}

@keyframes selectMenu2 {
  0% {
    box-shadow: 0 0 15px 0 rgb(0, 255, 46);
  }
  95% {
    box-shadow: 0 0 30px 15px rgba(204, 169, 44, 0);
  }
  100% {
    box-shadow: 0 0 0 20px rgba(204, 169, 44, 0);
  }
}

@keyframes selectMenu3 {
  0% {
    box-shadow: 0 0 15px 0 rgb(255, 0, 8);
  }
  95% {
    box-shadow: 0 0 30px 15px rgba(204, 169, 44, 0);
  }
  100% {
    box-shadow: 0 0 0 20px rgba(204, 169, 44, 0);
  }
}

@keyframes selectMenu4 {
  0% {
    box-shadow: 0 0 15px 0 rgb(0, 37, 255);
  }
  95% {
    box-shadow: 0 0 30px 15px rgba(204, 169, 44, 0);
  }
  100% {
    box-shadow: 0 0 0 20px rgba(204, 169, 44, 0);
  }
}

::-webkit-scrollbar-corner {
  background: rgb(8, 138, 210);
  box-shadow: inset 0 0 4px black;
}

.InventoryCell.bodyEquip span {
  font-weight: 900;
  font-size: 15px;
  line-height: 38px;
  text-align: center;
  width: 100%;
  display: block;
  color: #ff7800;
  text-shadow: 0 -1px 1px #000000, 0 -1px 1px #000000, 0 1px 1px #000000, 0 1px 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000;
}

.InventoryCell.bodyThorium h6 {
  margin: 4px 0 0 3px;
  padding: 0;
  z-index: 1;
  font-weight: 900;
  font-size: 8px;
  line-height: 8px;
  width: 100%;
  display: block;
  color: #ff7800;
  text-shadow: 0 -1px 1px #000000, 0 -1px 1px #000000, 0 1px 1px #000000, 0 1px 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000;
  text-align: left;
}

.InventoryCell.bodyThorium p {
  margin: 0 2px 0 0;
  font-size: 8px;
  padding: 0;
  text-align: right;
  color: #c5c3c1;
  background: none;
  text-shadow: 0 -1px 1px #000000, 0 -1px 1px #000000, 0 1px 1px #000000, 0 1px 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000;
}

.InventoryCell.bodyThorium .sep {
  text-align: center;
}

.InventoryCell.bodyThorium .itemIconInventoryCell {
  z-index: -1;
}

#hp_bar_wrapper {
  height: 6px;
  border: 1px solid #4c4c4c;
  text-align: left;
  display: block;
  box-shadow: inset 0 0 2px black, 0 0 2px black;
  border-radius: 10px;
  background-size: 12%;
  overflow: hidden;
  background-color: #959595;
  margin: 0 auto 0;
  position: absolute;
  top: 79px;
  z-index: 2;
  width: 80px;
  left: calc(50% - 40px);
  opacity: .75;
}

#hp_bar_wrapper span {
  display: block;
  font-size: 5px;
  line-height: 7px;
  width: 100%;
  text-align: center;
  text-shadow: 1px 1px 1px rgba(0, 0, 0, 1);
  font-weight: bold;
  margin: auto;
  float: left;
  color: white;
}

#hp_bar {
  overflow: visible;
  background: rgba(72, 255, 40, 0.7);
  height: 100%;
}

#bodyFine {
  background-image: url("https://img.icons8.com/fluent/48/000000/high-priority.png");
  background-size: contain;
  height: 50px;
  width: 50px;
  filter: drop-shadow(0px 0px 3px black);
  position: absolute;
  left: -23px;
  top: -2px;
}

.tipBodyFine {
  position: absolute;
  z-index: 15;
  left: 92px;
  top: 117px;
  display: block;
  font-size: 12px;
  color: white;
  width: 138px;
  padding: 5px 5px 5px 20px;
  background: rgba(0, 0, 0, 0.9);
  border-radius: 5px;
  box-shadow: 0 0 2px white;
  text-shadow: 0 -1px 1px #000000, 0 -1px 1px #000000, 0 1px 1px #000000, 0 1px 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000;
}

.disable {
  background: grey !important;
  pointer-events: none !important;
}

.noSelect {
  -webkit-touch-callout: none; /* iOS Safari */
  -webkit-user-select: none; /* Safari */
  -khtml-user-select: none; /* Konqueror HTML */
  -moz-user-select: none; /* Old versions of Firefox */
  -ms-user-select: none; /* Internet Explorer/Edge */
}

.version {
  position: fixed;
  top: 20px;
  right: 20px;
  font-size: 12px;
  font-weight: bold;
  color: white;
  text-shadow: 1px 1px 1px black;
}
</style>
