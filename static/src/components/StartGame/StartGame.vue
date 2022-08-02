<template>
  <div id="StartGameMenu" ref="StartGameMenu" class="noSelect">
    <div style="clear: both;">
      <div class="StartButton" @click="StartGame">В бой!</div>
    </div>
  </div>
</template>

<script>
import Preloader from "../Preloader/Preloader";

export default {
  name: "StartGame",
  data() {
    return {
      checker: null,
    }
  },
  destroyed() {
    clearInterval(this.checker)
  },
  mounted() {
    this.checker = setInterval(function () {
      this.$store.dispatch("sendSocketData", JSON.stringify({
        event: "gsgs",
        service: "lobby",
      }));
    }.bind(this), 1000)
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
    StartGame() {
      this.playSound('button_press.mp3', 0.3)
      this.$store.dispatch("sendSocketData", JSON.stringify({
        event: "StartGame",
        service: "lobby",
      }));
    },
  },
  computed: {},
  components: {
    AppPreloader: Preloader,
  }
}
</script>

<style scoped>
#StartGameMenu {
  position: absolute;
  border-radius: 5px;
  background: rgb(8, 138, 210);
  border: 1px solid #25a0e1;
  left: calc(50% - 132px);
  top: 5px;
  padding: 1px 1px 0;
  height: 30px;
  width: 264px;
  box-shadow: 1px 1px 3px black;
}

.StartButton {
  height: 25px;
  width: calc(100% - 4px);
  margin: 2px;
  border-radius: 5px;
  color: white;
  line-height: 28px;
  box-shadow: 0 0 2px black;
  transition: 100ms;
  background: rgba(255, 129, 0, 0.75);
  font-weight: 900;
  font-size: 22px;
  text-shadow: 1px 1px 1px black;
  cursor: pointer;
}

.StartButton:hover {
  background: rgba(255, 129, 0, 1);
}

.StartButton:active {
  transform: scale(.98);
}

.reward_block_title {
  float: left;
  height: 22px;
  width: 70px;
  border-radius: 5px 0 0 5px;
  overflow: hidden;
  background: #3486dd;
  color: #fff;
  line-height: 24px;
  font-size: 12px;
  font-weight: bold;
  box-shadow: -1px 1px 2px #000;
  background-size: contain;
  text-shadow: 1px 1px 1px black;
  cursor: pointer;
  margin: 2px 10px 0 3px;
  white-space: nowrap;
  text-overflow: ellipsis;
  outline: none;
  border: 0;
}

.wait_status {
  position: absolute;
  top: calc(100% + 10px);
  background: rgb(8, 138, 210);
  border: 1px solid #25a0e1;
  box-shadow: 1px 1px 3px black;
  border-radius: 5px;
  width: 146px;
  left: calc(50% - 75px);
  padding: 2px;
}

.wait_status_inner {
  background: #89969c;
  padding: 4px;
  border-radius: 5px;
  box-shadow: inset 0 0 2px black;
  float: left;
}

.search_table {
  width: 90px;
  margin: auto;
  border-spacing: 0;
  float: left;
}

.search_table tr:first-child {
  box-shadow: 0 0 1px rgb(0 0 0);
}

.search_table tr:first-child th {
  background: #216a8e;
  color: rgba(255, 255, 255, 0.8);
  font-size: 8px;
  user-select: none;
  margin: 2px auto;
  position: sticky;
  top: 0;
  padding-top: 2px;
  z-index: 2;
}

.search_table td:first-child {
  border-left: 1px solid rgba(0, 0, 0, 0.2);
}

.search_table td {
  border-right: 1px solid rgba(0, 0, 0, 0.2);
  border-bottom: 1px solid rgba(0, 0, 0, 0.2);
  color: rgba(255, 255, 255, 0.8);
  font-size: 10px;
  text-align: left;
  text-shadow: 1px 1px 1px black;
}

.search_table td:last-child {
  text-align: center;
}

.disable {
  opacity: 0.5 !important;
}

.select {
  opacity: 1 !important;
}

.wait_tip {
  color: wheat;
  text-shadow: 1px 1px 1px black;
  font-size: 8px;
  margin-top: 3px;
  float: left;
  text-align: justify;
  border-top: 1px solid rgba(0, 0, 0, 0.2);
  padding-top: 1px;
  width: 100%;
}

.select_type_battle {
  margin: 1px 0 0 3px;
}

.select_type_battle .cat {
  float: left;
  font-size: 13px;
  text-align: center;
  line-height: 18px;
  width: calc(50% - 6px);
  white-space: nowrap;
  background-color: rgba(76, 76, 76, 0.4);
  border-radius: 5px 5px 0 0;
  border: 1px solid rgba(76, 76, 76, 0.4);
  margin: 0 1px 0 1px;
  color: rgba(255, 255, 255, 0.8);
  text-shadow: 1px 1px 1px black;
  cursor: pointer;
  font-weight: bold;
}

.select_type_battle .cat:hover {
  background-color: rgba(255, 129, 0, 0.3);
}

.cat.select {
  /*border: 1px solid #FF9520;*/
  background-color: rgba(255, 129, 0, 0.5) !important;
  color: rgba(255, 255, 255, 1);
}

.operations {
  position: absolute;
  top: calc(100% + 10px);
  background: rgb(8, 138, 210);
  border: 1px solid #25a0e1;
  box-shadow: 1px 1px 3px black;
  border-radius: 5px;
  width: 400px;
  left: calc(50% - 200px);
  padding: 2px;
  height: 200px;
}

.operations_cat {
  box-shadow: inset 0 0 5px black;
  background: rgb(110 120 125);
  float: left;
  width: 137px;
  margin-left: 3px;
  border-radius: 0 0 5px 0;
  height: 100%;
  font-size: 13px;
  color: white;
  text-shadow: 1px 1px 1px black;
}

.operations_description {
  height: calc(100% - 10px);
  width: 250px;
  overflow-y: scroll;
  overflow-x: hidden;
  float: left;
  box-shadow: inset 0 0 3px black;
  background: linear-gradient(0deg, transparent 60%, rgba(0, 0, 0, 0.05) 21%), #adc6cd;
  background-size: 10px 3px;
  color: #000000;
  font-size: 13px;
  text-indent: 1.5em;
  text-shadow: none;
  font-weight: bold;
  padding: 5px;
  word-wrap: break-word;
  border-radius: 5px 0 0 5px;
  text-align: justify;
  margin: auto auto 2px 0;
}

.operations_icon {
  float: left;
  height: 68px;
  width: 68px;
  margin: 0 4px 0 0;
  border: 1px solid rgba(0, 0, 0, 0.4);
  box-shadow: inset 0 0 1px 0 black, 0 0 1px 0 black;
  background-size: cover;
  border-radius: 5px;
}

.operation_name {
  background: rgb(6 103 157);
  padding: 3px 2px;
  margin: 2px 1px;
  line-height: 14px;
  box-shadow: 0 1px 3px black;
  width: calc(100% - 8px);
  border: 1px solid #1684bf;
}

.operation_name:hover {
  background: rgb(213, 146, 110);
}

.operation_name.active {
  background: rgb(221, 112, 52);
}

.operations_head {
  margin: 0 0 5px 0;
  background: rgb(221, 112, 52);
  color: rgba(255, 255, 255, 0.8);
  border-radius: 4px;
  font-size: 12px;
  line-height: 17px;
  height: 16px;
  user-select: none;
  text-shadow: 1px 1px 1px rgb(0 0 0);
  font-weight: bold;
  box-shadow: 0 0 2px rgb(0 0 0);
  padding-left: 6px;
  text-align: left;
  clear: both;
  text-indent: 0;
}
</style>

<style>
.reward_block_resources .InventoryCell.select {
  background-color: rgba(255, 129, 0, 0.75) !important;
  box-shadow: none !important;
}

.operations_description p {
  margin-top: 0;
}

</style>
