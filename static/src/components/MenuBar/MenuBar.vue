<template>
  <div class="menuBar">
    <div class="menuBarInner">
      <div class="menu">

        <div class="menuPoint" title="Fullscreen"
             @click.stop="toggleFullScreen"
             @mouseover="playSound('select_sound.mp3', 0.2)">
          <div class="image" :style="{backgroundImage: 'url(' + require('../../assets/icons/fullscreen.png') + ')'}"/>
        </div>

        <div class="menuPoint" title="Выход"
             @click="exitGame"
             @mouseover="playSound('select_sound.mp3', 0.2)">
          <div class="image" :style="{backgroundImage: 'url(' + require('../../assets/icons/exit.png') + ')'}"/>
        </div>
      </div>
    </div>

    <div class="bottom_line">
    </div>
  </div>
</template>

<script>
export default {
  name: "MenuBar",
  data() {
    return {
      sub: '',
      audioPlayerVisible: false,
    }
  },
  mounted() {
    this.$store.dispatch("sendSocketData", JSON.stringify({
      event: "GetCredits",
      service: "market",
    }))
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
    openService(service, meta, component = '', forceOpen = false) {
      this.playSound('button_press.mp3', 0.3)
      this.sub = '';
      this.$store.commit({
        type: 'toggleWindow',
        id: service,
        component: component,
        meta: meta,
        forceOpen: forceOpen,
      });
    },
    exitGame() {
      this.playSound('button_press.mp3', 0.3)
      if (this.$route.path === '/lobby') {
        this.$store.commit({
          type: 'closeWS'
        });
        this.$router.push('/');
      } else {
        this.$store.dispatch("sendSocketData", JSON.stringify({
          event: "ExitGame",
          service: "battle",
        }));
      }
    },
    toggleFullScreen() {
      this.playSound('button_press.mp3', 0.3)
      if (!document.fullscreenElement) {
        document.documentElement.requestFullscreen();
      } else {
        if (document.exitFullscreen) {
          document.exitFullscreen();
        }
      }

      this.setResolution();
    },
    setResolution() {
      this.$store.commit({
        type: 'setResolution',
      })

      this.$store.commit({
        type: 'addCheckViewPort',
      })
    },
  },
  computed: {
    openComponents() {
      return this.$store.getters.getNeedOpenComponents
    },
  },
}
</script>

<style scoped>
.menuBar {
  position: fixed;
  display: block;
  width: calc(100% - 2px);
  height: 32px;
  bottom: 0;
  z-index: 900;
  pointer-events: none;
}

.menuBarInner {
  margin-right: 5px;
  height: calc(100% - 4px);
  border-radius: 5px 5px 0 0;
  user-select: none;
  white-space: nowrap;
  text-overflow: ellipsis;
  border: 1px solid #25a0e1;
  background: rgb(8, 138, 210);
  box-shadow: 0 0 2px black;
  padding: 2px 2px 0 2px;
  float: right;
  position: relative;
  z-index: 12;
  pointer-events: auto;
}

.menu, .credits, .inventory_size, .reactor, .credits_up {
  float: right;
  padding: 1px;
  box-shadow: 0 0 2px black;
  border-radius: 3px;
  background: #06679d;
}

.reactor {
  padding: 2px 0 2px 0;
  width: 103px;
  height: calc(100% - 4px);
}

.credits, .inventory_size, .credits_up {
  width: 60px;
  height: calc(100% - 4px);
}

.inventory_size {
  width: 90px;
}

.credits_up {
  width: 20px;
  height: 16px;
  margin-right: 5px;
  background: none;
}

.credits_up:hover {
  background-color: #b4eaff;
}

.credits_up:active {
  transform: scale(.98);
}

.credits_icon, .inventory_size_icon, .reactor_icon, .credits_up_icon {
  height: 16px;
  width: 16px;
  margin-top: -1px;
  float: left;
  filter: drop-shadow(0px 0px 1px black);
  background-size: contain;
}

.credits_count, .inventory_size_wrapper {
  background: #066ea8;
  float: right;
  width: calc(100% - 21px);
  padding-right: 4px;
  height: 100%;
  border-radius: 4px;
  color: white;
  text-align: right;
  line-height: 14px;
  text-shadow: 1px 1px 1px rgba(0, 0, 0, 1);
  box-shadow: inset 0 0 2px black;
  font-size: 9px;
}

/*.inventory_size_icon {*/
/*  background-image: url(https://img.icons8.com/cotton/64/000000/orthogonal-view.png)*/
/*}*/

/*.reactor_icon {*/
/*  background-image: url('../../assets/resource/enriched_thorium.png')*/
/*}*/

.inventory_size_wrapper {
  width: calc(100% - 16px);
  padding-right: 0;
  overflow: hidden;
  height: 100%;
}

.menuPoint {
  height: 22px;
  width: 22px;
  background: #06679d;
  border-radius: 5px;
  margin: 1px 2px 1px 2px;
  float: left;
  transition: 0.1s;
  background-size: contain;
  position: relative;
  z-index: 100;
  box-shadow: 0 0 2px 0 black;
}

.menuPoint:hover {
  background-color: #b4eaff;
}

.menuPoint:active {
  transform: scale(0.95);
}

.sub {
  position: absolute;
  min-height: 20px;
  min-width: 200px;
  border: 1px solid rgba(37, 160, 225, 0.2);
  background: rgba(8, 138, 210, 0.2);
  border-radius: 5px;
  bottom: calc(100% + 5px);
  right: 0;
}

.sub > div {
  height: 35px;
  clear: both;
  background: #0cc2fb;
  border: 1px solid rgba(37, 160, 225, 0.5);
  background: rgba(8, 138, 210, 0.5);
  border-radius: 3px;
  margin: 2px;
  line-height: 35px;
  color: white;
  text-shadow: 1px 1px 1px black;
}

.sub > div .image_wrapper {
  height: 30px;
  width: 30px;
  /*filter: drop-shadow(0 0 2px black);*/
  box-shadow: inset 0 0 2px black;
  border: 1px solid grey;
  background: #8cb3c7;
  border-radius: 5px;
  margin: 2px 5px 2px 5px;
  float: left;
  transition: 0.1s;
  background-size: contain;
}

.sub > div h4 {
  margin: 0 2px;
  float: left;
  opacity: 0.8;
  font-size: 18px;
}

.sub > div:hover {
  border: 1px solid rgba(37, 160, 225, 0.8);
  background: rgba(8, 138, 210, 0.8);
}

.sub > div:hover h4 {
  opacity: 1;
  cursor: pointer;
  -webkit-touch-callout: none; /* iOS Safari */
  -webkit-user-select: none; /* Safari */
  -khtml-user-select: none; /* Konqueror HTML */
  -moz-user-select: none; /* Old versions of Firefox */
  -ms-user-select: none; /* Internet Explorer/Edge */
  user-select: none;
}

.bottom_line {
  border: 1px solid rgba(37, 160, 225, 0.3);
  background: rgba(8, 138, 210, 0.3);
  height: 5px;
  box-shadow: 0 0 2px black;
  position: absolute;
  bottom: 0;
  width: 100%;
  z-index: 0;
}

.credits_wrapper, .inventory_size_out_wrapper, .reactor_out_wrapper {
  float: right;
  z-index: 999;
  position: relative;
  margin-right: 5px;
  pointer-events: auto;
  border-radius: 5px 5px 0 0;
  padding: 2px 2px 0 2px;
  border: 1px solid #25a0e1;
  background: rgb(8, 138, 210);
  box-shadow: 0 0 2px black;
  height: calc(100% - 15px);
  top: 11px;
}

.credits_up {

}

.menuPoint_audio {
  height: 22px;
  width: 22px;
  filter: none;
  background-size: contain;
  border-radius: 5px;
  padding: 2px;
}

.menuPoint_audio:hover {
  background-color: #b4eaff;
}

.menuPoint_audio:active {
  transform: scale(0.95);
}

.audio_player_wrapper {
  position: absolute;
  bottom: calc(100% + 5px);
  right: calc(100% - 35px);
  transition: 100ms;
}

.on {
  background: #00ff22;
}

.image {
  height: 100%;
  width: 100%;
  background-size: contain;
  background-position: center;
  background-repeat: no-repeat;
  filter: contrast(50%) sepia(100%) hue-rotate(346deg) brightness(0.8) saturate(800%) drop-shadow(0px 1px 0px black);
}

</style>

<style>
.menuBar .sizeInventoryInfo {
  box-shadow: inset 0 0 2px black !important;
  background: none !important;
  border: none !important;
  border-radius: 2px !important;
  width: 100% !important;
}

.menuBar .sizeInventoryInfo span {
  background-image: none !important;
  font-size: 9px !important;
  line-height: 13px !important;
}

.menuBar .sizeInventoryInfo .realSize {
  background: rgba(255, 96, 0, 1) !important;
  box-shadow: inset 0 0 2px black !important;
}
</style>
