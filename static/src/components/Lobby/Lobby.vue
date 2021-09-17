<template>
  <div id="lobbyWrapper">
    <div class="RankWrapper">
      <app-session-battle-rank/>
    </div>

    <div id="ServiceTable">
      <div class="sub" @click.stop="">
        <div @click.stop="openService('CreateGameBlock')">
          <img src="https://img.icons8.com/material/40/000000/exit.png"/>
          <h4>Создать</h4>
        </div>
      </div>

      <div class="sub" @click.stop="">
        <div @click.stop="openService('')">
          <img src="https://img.icons8.com/material/40/000000/exit.png"/>
          <h4>Поиск игорей</h4>
        </div>
      </div>

      <h2 class="head_category">Какаято хуйня: </h2>
      <div class="sub" @click.stop="">
        <div id="hangarButton">
          <img src="https://img.icons8.com/ios-glyphs/30/000000/hangar.png"/>
          <h4>Какаято хуйня</h4>
          <div class="help_point_menu"></div>
        </div>
      </div>
    </div>

    <app-create-game v-if="openComponents['CreateGameBlock'] && openComponents['CreateGameBlock'].open"/>
  </div>
</template>

<script>
import SessionBattleRank from './SessionBattleRank'
import CreateGame from "../CreateGame/CreateGame";

export default {
  name: "Lobby",
  data() {
    return {}
  },
  mounted() {
    this.$store.dispatch("sendSocketData", JSON.stringify({
      event: "GetPlayers",
      service: "system",
    }))
  },
  methods: {
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
  computed: {
    openComponents() {
      return this.$store.getters.getNeedOpenComponents
    }
  },
  components: {
    AppSessionBattleRank: SessionBattleRank,
    AppCreateGame: CreateGame,
  }
}
</script>

<style scoped>
#lobbyWrapper {
  height: 100vh;
  width: 100%;
  text-align: center;
  background-color: #7f7f7f;
  background-image: url('../../assets/bases/base.jpg');
  background-size: cover;
  background-attachment: fixed;
  background-position: center;
}

#ServiceTable {
  top: 60px;
  left: 5px;
  text-align: center;
  float: left;
  border-radius: 5px;
  position: relative;
  padding: 2px;
  width: 220px;
  border: 1px solid #25a0e1;
  background: rgb(8, 138, 210);
  box-shadow: 0 0 2px black;
}

#dialogBase {
  pointer-events: auto;
  height: 80px;
  width: 267px;
  border: 1px solid rgb(37, 160, 225);
  background-size: 10px 2px;
  background-image: linear-gradient(1deg, rgba(33, 176, 255, 0.6), rgba(37, 160, 225, 0.6) 6px);
  position: absolute;
  border-radius: 5px;
  padding: 20px 2px 2px;
  left: calc(50% - 133px) !important;
  top: 40% !important;
}

#dialogBase #textWrapper {
  box-shadow: inset 0 0 5px black;
  background: #8cb3c7;
  border-radius: 5px;
  height: calc(100% - 42px);
  position: relative;
  color: #ff7800;
  text-shadow: 0 -1px 1px #000000, 0 -1px 1px #000000, 0 1px 1px #000000, 0 1px 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000;
  padding: 10px;
  font-size: 15px;
}

#dialogBase input {
  margin: 2px auto;
  width: 100%;
  background: rgb(221, 112, 52);
  box-shadow: inset 0 0 4px 0 white;
  color: rgba(255, 255, 255, 0.8);
  display: block;
}

#dialogBase input:hover {
  cursor: pointer;
  box-shadow: inset 0 0 4px 0 #20fffd;
}

#dialogBase input:active {
  transform: scale(0.98);
}

#OutDialog {
  position: absolute;
  top: 50px;
  left: calc(50% - 150px);
  display: block;
  width: 300px;
  height: 100px;
  z-index: 999;
}

#OutDialog h3 {
  color: #ff8100;
  font-size: 11px;
  text-shadow: 0 -1px 1px #000000, 0 -1px 1px #000000, 0 1px 1px #000000, 0 1px 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000;
  margin: 8px 4px;
  padding: 5px 0;
  box-shadow: inset 0 0 2px 0 black;
  border-radius: 9px;
  border: 1px solid #25a0e1;
  background: rgb(8, 138, 210);
}

#OutDialog > div {
  height: 52px;
  width: 175px;
  margin: 0 auto;
  background: rgba(2, 2, 2, 0.2);
  padding-left: 0px;
  border-radius: 0 0 40% 40%;
  margin-top: -9px;
}

.arrow {
  width: 60px;
  height: 32px;
  margin-left: 14px;
  padding-top: 6px;
  transform: rotate(270deg) scale(0.9);
}

.arrow span {
  display: block;
  width: 30px;
  height: 30px;
  border-bottom: 5px solid rgb(222, 156, 0);
  border-right: 5px solid rgb(255, 204, 0);
  transform: rotate(45deg);
  margin: 5px;
  animation: animate 2s infinite;
  filter: drop-shadow(0 0 5px rgba(255, 250, 0, 1));
}

.arrow span:nth-child(2) {
  animation-delay: -0.2s;
}

.arrow span:nth-child(3) {
  animation-delay: -0.4s;
}

@keyframes animate {
  0% {
    opacity: 0;
    transform: rotate(45deg) translate(-20px, -20px);
  }
  50% {
    opacity: 1;
  }
  100% {
    opacity: 0;
    transform: rotate(45deg) translate(20px, 20px);
  }
}

.sub {
  min-height: 20px;
  min-width: 200px;
  border-radius: 5px;
  top: calc(100% + 8px);
  filter: drop-shadow(0px 0px 6px black);
  box-shadow: 0 0 2px black;
}

.sub > div {
  height: 20px;
  clear: both;
  background: #0cc2fb;
  border: 1px solid rgba(37, 160, 225, 0.5);
  background: rgba(8, 138, 210, 0.5);
  border-radius: 3px;
  margin: 1px;
  line-height: 35px;
  color: white;
  text-shadow: 1px 1px 1px black;
  cursor: pointer;
  -webkit-touch-callout: none; /* iOS Safari */
  -webkit-user-select: none; /* Safari */
  -khtml-user-select: none; /* Konqueror HTML */
  -moz-user-select: none; /* Old versions of Firefox */
  -ms-user-select: none; /* Internet Explorer/Edge */
  user-select: none;
  overflow: hidden;
  position: relative;
}

.sub > div img {
  height: 20px;
  width: 20px;
  background: #8cb3c7;
  float: left;
  transition: 0.1s;
  background-size: contain;
}

.sub > div h4 {
  margin: 0 5px;
  float: left;
  opacity: 0.8;
  font-size: 12px;
  line-height: 21px;
}

.sub > div:hover {
  border: 1px solid rgba(37, 160, 225, 0.8);
  background: rgba(8, 138, 210, 0.8);
}

.sub > div:active {
  transform: scale(0.97);
}

.sub > div:hover h4 {
  opacity: 1;
}

.RankWrapper {
  position: absolute;
  top: 5px;
  left: 5px;
  width: 270px;
  background: none;
  box-shadow: none;
  z-index: 0;
}

.base_icon {
  height: 50px;
  width: 50px;
  float: left;
  margin: 1px 5px 5px 5px;
  background-size: contain;
  background-color: rgba(6, 110, 168, 0.5);
  border-radius: 5px;
  border: 1px solid rgba(6, 110, 168, 0.6);
  box-shadow: inset 0 0 2px black;
}

.base_name {
  margin: 0 0 7px 0;
  background: rgb(221, 112, 52);
  color: rgba(255, 255, 255, 0.8);
  border-radius: 4px;
  font-size: 13px;
  line-height: 18px;
  height: 17px;
  user-select: none;
  text-shadow: 1px 1px 1px rgba(0, 0, 0, 1);
  font-weight: bold;
  box-shadow: 0 0 2px rgba(0, 0, 0, 1);
  padding-left: 10px;
  float: right;
  width: calc(100% - 74px);
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
}

.base_type {
  float: left;
  color: #e4ff00;
  font-size: 10px;
  margin-right: 4px;
}

.section {
  line-height: 8px;
  background: rgba(77, 77, 84, 0.3);
  width: 140px;
  height: 28px;
  border-radius: 0 5px 5px 0;
  color: white;
  text-shadow: 1px 1px 1px rgb(0 0 0);
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
  margin-bottom: 5px;
  font-size: 11px;
  text-align: left;
}

.fractionLogo {
  height: 28px;
  width: 28px;
  float: left;
  margin: -1px 5px 5px -1px;
  background-size: cover;
  background-color: rgba(6, 110, 168, 0.2);
  border-radius: 5px;
  border: 1px solid rgba(6, 110, 168, 0.3);
  box-shadow: inset 0 0 2px black;
}

.importantly {
  color: #f1bd00;
  font-weight: 100;
}

.head_category {
  margin: 3px 0 3px 0;
  background: rgb(221, 112, 52);
  color: rgba(255, 255, 255, 0.8);
  border-radius: 4px;
  font-size: 13px;
  line-height: 17px;
  height: 16px;
  user-select: none;
  text-shadow: 1px 1px 1px rgba(0, 0, 0, 1);
  font-weight: bold;
  box-shadow: 0 0 2px rgba(0, 0, 0, 1);
  padding-left: 10px;
  text-align: left;
  clear: both;
}

.help_point_menu {
  height: 16px;
  width: 16px;
  /* background: black; */
  position: absolute;
  right: 3px;
  bottom: 2px;
  background-image: url(https://img.icons8.com/flat_round/64/000000/question-mark.png);
  background-size: contain;
  border-radius: 50%;
  box-shadow: 0 0 2px black;
  opacity: 0.5;
}

.help_point_menu:hover {
  opacity: 1;
  box-shadow: 0 0 4px 2px orange;
}

@media (max-width: 1000px) {
  .sub > div h4 {
    font-size: 10px;
    line-height: 17px;
  }

  .sub > div {
    height: 16px;
  }

  .sub > div img {
    height: 17px;
    width: 17px;
  }

  #ServiceTable {
    width: 202px;
  }

  .sub {
    min-height: 17px;
    filter: drop-shadow(0px 0px 2px black);
    box-shadow: 0 0 2px black;
  }

  .help_point_menu {
    height: 13px;
    width: 13px;
  }
}
</style>
