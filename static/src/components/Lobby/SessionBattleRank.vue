<template>
  <div class="sessionRankWrapper">

    <div class="sessionRankWrapper2">

      <div class="userAvatarWrapper"/>

      <div v-if="currentUser">

        <h3 class="userName">
          {{ currentUser.login }}
        </h3>

        <div class="state">
          <div class="Rank">
            0
          </div>

          <div>
            <h5>Боевой ранг 0 / 0</h5>
            <div class="RankBar">
              <div class="RankBarInner" :style="{width: '100%'}"/>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>

export default {
  name: "SessionBattleRank",
  data() {
    return {
      avatar: null,
      logo: null,
    }
  },
  created() {

  },
  methods: {},
  computed: {
    currentUser() {
      return this.$store.getters.getCurrentPlayer
    },
  }
}
</script>

<style scoped>
.sessionRankWrapper {
  height: 50px;
  width: 100%;
  border: 1px solid rgba(37, 160, 225, 0.2);
  background: rgba(8, 138, 210, 0.2);
  box-shadow: none;
  border-radius: 0 0 50px 5px;
}

.sessionRankWrapper2 {
  height: 44px;
  clear: both;
  border: 1px solid rgba(37, 160, 225, 0.5);
  background: rgba(8, 138, 210, 0.5);
  border-radius: 0 0 50px 5px;
  margin: 2px;
  line-height: 35px;
  color: white;
  text-shadow: 1px 1px 1px black;
}

.userAvatarWrapper {
  background-color: rgb(128, 128, 128);
  height: 42px;
  width: 42px;
  float: left;
  margin: 0 5px 5px 1px;
  border: 1px solid rgba(0, 0, 0, 0.4);
  box-shadow: inset 0 0 1px 0 black;
  background-size: cover;
  border-radius: 5px;
  transition: 100ms;
}

.userAvatarWrapper:hover {
  box-shadow: inset 0 0 2px 0 white, 0 0 2px 0 white;
  border: 1px solid white;
}

.userAvatarWrapper:active {
  transform: scale(0.97);
}

.userName {
  margin: 1px 0 0 2px;
  background: rgb(221, 112, 52);
  color: rgba(255, 255, 255, 0.8);
  border-radius: 0 0 50px 5px;
  font-size: 13px;
  line-height: 17px;
  height: 17px;
  user-select: none;
  text-shadow: 1px 1px 1px rgba(0, 0, 0, 1);
  font-weight: bold;
  box-shadow: 0 0 2px rgba(0, 0, 0, 1);
  padding-left: 10px;
  text-align: left;
  width: calc(100% - 16px);
}

.userName span {
  color: #e4ff00;
  font-size: 10px;
}

.Rank {
  height: 21px;
  width: 31px;
  float: left;
  background-color: rgb(128, 128, 128);
  border: 1px solid rgba(0, 0, 0, 0.4);
  box-shadow: inset 0 0 1px 0 black;
  margin: -1px 0 0 -1px;
  font-family: 'Audiowide', cursive;
  line-height: 22px;
  font-size: 17px;
  text-align: center;
  border-radius: 3px;
}

.state {
  margin-top: 3px;
}

.state h5 {
  margin: 0;
  width: calc(100% - 100px);
  float: left;
  height: 12px;
  font-size: 9px;
  line-height: 10px;
  text-align: left;
  text-indent: 5px;
}

.RankBar {
  position: relative;
  width: calc(100% - 125px);
  height: 7px;
  margin: 1px 2px;
  background: rgb(145, 145, 145);
  float: left;
  border-radius: 5px;
  box-shadow: inset 0 0 2px black;
  overflow: hidden;
}

.RankBarInner {
  top: 0;
  left: 0;
  position: absolute;
  height: 100%;
  width: 50%;
  box-shadow: inset 0 0 2px 0 black;
  background: #ff8100;
}

.exitTime, .lawTime {
  position: absolute;
  height: 35px;
  width: 35px;
  background-size: cover;
  background-image: url(https://img.icons8.com/fluent/48/000000/medium-risk.png);
  left: calc(100% + 12px);
  top: 1px;
  filter: drop-shadow(0px 0px 2px black);
}

.lawTime {
  background-image: url(https://img.icons8.com/emoji/48/000000/red-triangle-pointed-up-emoji.png);
  left: calc(100% + 55px);
  top: 2px;
}

.exitTime .time, .lawTime .time {
  color: yellow;
  font-weight: 900;
  font-size: 10px;
  position: absolute;
  top: 100%;
  width: 100%;
  text-align: center;
  text-shadow: 0 -1px 1px #000000, 0 -1px 1px #000000, 0 1px 1px #000000, 0 1px 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000;
}

.vio {
  font-size: 8px;
  line-height: 45px;
  color: white;
  width: 100%;
  text-align: center;
  text-shadow: 0 -1px 1px #000000, 0 -1px 1px #000000, 0 1px 1px #000000, 0 1px 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000;
}
</style>
