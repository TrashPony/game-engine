<template>
  <div class="gate" id="gateBlock" ref="gateBlock">
    <div class="gateBlock" v-if="currentUser || currentPlayer">
      <app-control v-bind:head="'Вход'" v-bind:move="false" v-bind:close="false" v-bind:refWindow="'gateBlock'"/>
      <app-quick-fight v-if="selectGameMode==='quick_battle'"/>
    </div>
  </div>
</template>

<script>
import Control from '../Window/Control';
import QuickFight from './QuickFight';
import {urls} from '../../const';

export default {
  name: "Gate",
  data() {
    return {
      selectGameMode: 'quick_battle',
      error: '',
    }
  },
  created() {
    this.$store.commit({
      type: 'setVisibleLoader',
      visible: true,
      text: `<span>Получаем информацию...</span>`,
    });
  },
  mounted() {
    let app = this;

    if (app.$route.query['auth_key'] && app.$route.query['api_id'] && app.$route.query['viewer_id'] && app.$route.query['access_token']) {

      let queryStr = '?auth_key=' + app.$route.query['auth_key'] + '&api_id=' + app.$route.query['api_id'] +
        '&viewer_id=' + app.$route.query['viewer_id'] + '&access_token=' + app.$route.query['access_token']

      app.$api.get(urls.vkAppLogin + queryStr, {
        withCredentials: true,
      }).then(function (response) {
        if (response.data.success) {
          app.getPlayers()
        } else {
          app.error = response.data.error
        }
      });

    } else {
      app.getPlayers()
    }
  },
  methods: {
    getPlayers() {
      this.$store.dispatch("sendSocketData", JSON.stringify({
        event: "GetPlayers",
        service: "system",
      }))
    },
    to(url) {
      let app = this;

      if (url === '/lobby') {
        app.$store.commit({
          type: 'setVisibleLoader',
          visible: true,
          text: 'Пытаемся понять что происходит...'
        });
      }

      setTimeout(function () {
        app.$router.push({path: url});
      }, 1000);
    },
  },
  computed: {
    currentUser() {
      let currentUser = this.$store.getters.getGameUser

      if (currentUser && currentUser.hasOwnProperty('id')) {
        this.$store.commit({
          type: 'setVisibleLoader',
          visible: false,
        });
      }

      return currentUser
    },
    currentPlayer() {
      let currentPlayer = this.$store.getters.getCurrentPlayer
      if (currentPlayer && currentPlayer.hasOwnProperty('id')) {
        this.to('/lobby')
      }
      return ''
    }
  },
  components: {
    AppControl: Control,
    AppQuickFight: QuickFight,
  }
}
</script>

<style scoped>
.gate {
  height: 100vh;
  width: 100%;
  text-align: center;
  background-color: #7f7f7f;
  background-image: url('../../assets/bases/base.jpg');
  background-size: cover;
  background-attachment: fixed;
  background-position: center;
}

.gateBlock {
  position: absolute;
  left: calc(50% - 107px);
  top: 20%;
  display: block;
  border-radius: 5px;
  width: 246px;
  min-height: 40px;
  border: 1px solid #25a0e1;
  background: rgb(8, 138, 210);
  z-index: 11;
  padding: 20px 0 2px 0;
  box-shadow: 0 0 2px black;
}

.entryButton {
  height: 100px;
  width: 100px;
  box-shadow: 0 0 3px 1px rgba(0, 0, 0, 0.6);
  float: left;
  margin: 15px 0 5px 15px;
  border-radius: 7px;
  background: rgba(255, 255, 255, 0.2);
  position: relative;
  transition: 100ms;
  background-size: cover;
}

.entryButton span {
  position: absolute;
  width: 77px;
  padding: 4px 4px 0 4px;
  left: calc(50% - 43px);
  bottom: 5px;
  font-size: 13px;
  color: white;
  background: rgba(0, 0, 0, 0.25);
  border-radius: 7px;
  text-shadow: 0 -1px 1px #000000, 0 -1px 1px #000000, 0 1px 1px #000000, 0 1px 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000;
}

.entryButton:hover {
  box-shadow: 0 0 3px 1px white;
}

.entryButton:active {
  transform: scale(0.98);
}

.disable {
  filter: grayscale(1);
  pointer-events: none;
}

.language {
  position: absolute;
  left: 100%;
  border: 1px solid #25a0e1;
  border-left: 0;
  background-color: rgb(8, 138, 210);
  box-shadow: 0 0 2px 0 black;
  border-radius: 0 5px 5px 0;
}

.language div {
  height: 25px;
  width: 25px;
  background-size: cover;
  filter: drop-shadow(0 0 1px black);
}

.language .disable_language {
  filter: grayscale(75%);
  opacity: 0.75;
}

.language .disable_language:hover {
  filter: drop-shadow(0 0 1px white) grayscale(0%);
  opacity: 1;
}
</style>

<style>
.gateBlock .chatUserLine {
  background: rgba(250, 235, 215, 0.28);
  text-align: left;
}

.gateBlock .chatUserLine .chatUserIcon {
  height: 40px;
  width: 40px;
}

.gateBlock .chatUserLine .chatUserName {
  line-height: 30px;
  font-size: 20px;
}
</style>
