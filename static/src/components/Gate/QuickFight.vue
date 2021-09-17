<template>
  <div>

    <template v-if="quickFightPlayer">
      <div class="createNewUser button" @click="to('/lobby')">Вход</div>
    </template>

    <template v-if="!quickFightPlayer">
      <div class="createNewUser" style="margin-bottom: 4px">
        Введите имя
      </div>
      <input v-model="newName" class="createNewUser input" placeholder="Имя персонажа"/>
      <div class="createNewUser button" @click="create">Создать</div>
    </template>

  </div>
</template>

<script>
export default {
  name: "QuickFight",
  data() {
    return {
      createMenu: false,
      newName: null,
    }
  },
  mounted() {
    if (this.quickFightPlayer) {
      this.to('/lobby')
    }
  },
  methods: {
    create() {

      this.$store.dispatch("sendSocketData", JSON.stringify({
        event: "CreateNewPlayer",
        service: "system",
        user_name: this.newName,
      }))

      this.createMenu = false;
    },
    to(url) {
      let app = this;

      this.$store.dispatch("sendSocketData", JSON.stringify({
        event: "SelectPlayer",
        service: "system",
        id: this.quickFightPlayer.id,
      }))

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
    quickFightPlayer() {
      for (let player of this.$store.getters.getUserPlayers) {
        return player;
      }
      return null;
    },
  }
}
</script>

<style scoped>
.createNewUser {
  height: 20px;
  width: calc(100% - 4px);
  background: rgba(108, 108, 108, 0.51);
  margin: 2px;
  border-radius: 5px;
  color: white;
  line-height: 20px;
  box-shadow: 0 0 2px black;
}

.createNewUser.input {
  width: calc(100% - 14px);
}

.createNewUser.input::placeholder {
  color: white;
}

.createNewUser.button {
  transition: 100ms;
  background: rgba(255, 129, 0, 0.75);
}

.createNewUser.button:hover {
  background: rgba(255, 129, 0, 1);
}

.createNewUser.button:active {
  transform: scale(0.98);
}
</style>
