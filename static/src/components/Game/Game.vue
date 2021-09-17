<template>
  <div id="game_wrap">
    <div id="game-container"></div>
  </div>
</template>

<script>
  export default {
    name: "Game",
    data() {
      return {}
    },
    mounted() {
      //launch();
    }
  }
</script>

<style>
  #game_wrap {
    position: absolute;
    z-index: -1;
    left: 0;
    top: 0;
    height: 100vh;
    user-select: none;
  }

  #game-container {
    height: 100%;
    width: 100%;
  }
</style>
