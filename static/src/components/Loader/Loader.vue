<template>
  <div class="wrapper" v-if="loader.visible">
    <div class="mask"></div>

    <div id="loader" ref="loader">
      <app-control v-bind:head="'Загрузка... '"
                   v-bind:move="true"
                   v-bind:close="false"
                   v-bind:refWindow="'loader'"
                   v-bind:noHeight="true"
                   v-bind:noPos="true"
                   v-bind:noWidth="true"/>

      <div class="text" v-if="loader.text !== ''" v-html="loader.text"/>
      <div class="text" v-else>
        Происходит что то неведомое... ¯\\_(ツ)_/¯
      </div>

      <div id="OutDialog">
        <div>
          <div class="arrow">
            <span/>
            <span/>
            <span/>
          </div>
        </div>
      </div>


    </div>
  </div>
</template>

<script>
  import Control from '../Window/Control';

  export default {
    name: "Loader",
    computed: {
      loader() {
        return this.$store.getters.getEndLoad
      }
    },
    components: {
      AppControl: Control,
    }
  }
</script>

<style scoped>

  .wrapper {
    height: 100vh;
    width: 100%;
    text-align: center;
    background-color: #7f7f7f;
    background-image: url('../../assets/bases/base.jpg');
    background-size: cover;
    background-attachment: fixed;
    background-position: center;
    z-index: 999;
    position: absolute;
    left: 0;
    top: 0;
  }

  .text {
    border-radius: 5px;
    background: #89969c;
    pointer-events: auto;
    padding: 5px;
    box-shadow: inset 0 0 2px black;
  }

  .mask {
    position: absolute;
    background-color: rgba(0, 0, 0, 0.7);
    height: 100%;
    width: 100%;
    top: 0;
    left: 0;
  }

  #loader {
    position: absolute;
    left: calc(50% - 140px);
    top: 20%;
    display: block;
    border-radius: 5px;
    width: 280px;
    height: 24px;
    border: 1px solid #25a0e1;
    background: rgb(8, 138, 210);
    padding: 20px 2px 2px 2px;
    box-shadow: 0 0 2px black;
    font-size: 12px;
    text-align: left;
  }

  #OutDialog {
    position: absolute;
    top: 56px;
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
    background: rgba(255, 255, 255, 0.2);
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
</style>
