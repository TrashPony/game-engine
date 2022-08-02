<template>
  <div class="index_background">
    <div class="scroll_div">

      <header>
        <div class="WrapLeftTopBlock">
          <div class="leftTopBlock">
            Server online
          </div>
        </div>

        <div class="play_button_wrapper">
          <div class="play_button" @click="to('/lobby')">
            PLAY
          </div>
        </div>

        <div class="WrapRightTopBlock">
          <div class="rightTopBlock" style="color: white">
            Search...
          </div>
        </div>
      </header>

      <div class="logo"/>

      <router-view/>

      <footer>
        <div id="about">
            <span
              style="font-size: 14px;text-align: center;width: 100%;color: #a0a0a0;">TrashPony13@gmail.com 2020</span>
        </div>
      </footer>
    </div>
  </div>
</template>

<script>
import Control from '../Window/Control';

export default {

  name: "Index",
  methods: {
    to(url) {
      this.$router.push(url)
    }
  },
  components: {
    AppControl: Control,
  }
}
</script>

<style scoped>
body {
  background: #000;
  scroll-behavior: smooth;
  overflow-x: hidden;
  padding: 0;
  margin: 0;
  text-align: justify;
  text-shadow: 0 0 2px #343434, 2px 1px 0 hsl(0, 0%, 27%), 2px 1px 1px hsl(0, 0%, 23%), 1px 1px 2px hsl(80, 10%, 15%);
  color: white;
}

.index_background {
  background: url('../../assets/bases/base.png') center no-repeat;
  background-attachment: fixed;
  background-size: cover;
  height: 100vh;
  width: 100vw;
}
.logo {
  background-size: contain;
  background-position: 5% center;
  height: 100px;
  width: 1100px;
  margin: auto;
  filter: drop-shadow(2px 2px 0px black);
}

.scroll_div {
  overflow-y: scroll;
  height: 100vh;
  width: 100vw;
  position: absolute;
  top: 0;
  left: 0;
  overflow-x: hidden;
}

section {
  margin: auto;
}

a {
  color: #0cfdff;
  text-decoration: none;
}

a:hover {
  text-decoration: underline;
  color: #ffc588;
}

main {
  padding-bottom: 25px;
  /*flex: 1 0 auto;*/
}

footer {
  border: 1px solid #25a0e1;
  background: #088ad2;
  box-shadow: 0 0 2px #000;
  /*flex: 0 0 auto;*/
  position: fixed;
  bottom: 0;
  width: 100%;
}

header {
  border-bottom: 1px solid rgba(37, 160, 225, 0.75);
  background: rgba(8, 138, 210, 0.75);
  box-shadow: 0 0 2px #000;
  height: 35px;
  width: 100%;
  z-index: 100;
  position: relative;
}

#features #features-container .feature {
  box-shadow: inset 0 0 8px rgb(0, 0, 0);
  background: rgba(0, 0, 0, 0.2);
  width: 20%;
  border-radius: 15px;
  color: white;
}

#features #features-container .feature span {
  padding: 10px;
}

#features #features-container .feature div {
  margin: 15px auto 10px;
  height: 100px;
  width: 100px;
  background: grey;
  border-radius: 5px;
  border: 1px solid grey;
  box-shadow: 0 0 8px rgb(0, 0, 0);
}

footer #about {
  margin: auto;
  display: -webkit-box;
  display: -ms-flexbox;
  display: flex;
  -webkit-box-pack: justify;
  -ms-flex-pack: justify;
  justify-content: space-between;
  -webkit-box-align: center;
  -ms-flex-align: center;
  align-items: center;
  height: 30px;
}

.play_button {
  display: block;
  pointer-events: auto;
  text-align: center;
  transition: .1s;
  box-shadow: 0 0 2px #000;
  /* text-indent: 0; */
  height: 26px;
  width: 100px;
  background: rgba(255, 129, 0, .75);
  margin: 2px;
  border-radius: 5px;
  color: #fff;
  line-height: 28px;
  cursor: pointer;
  font-size: 20px;
  font-family: 'Audiowide', cursive;
  text-shadow: 0 -1px 1px #000, 0 -1px 1px #000, 0 1px 1px #000, 0 1px 1px #000, -1px 0 1px #000, 1px 0 1px #000, -1px 0 1px #000, 1px 0 1px #000, -1px -1px 1px #000, 1px -1px 1px #000, -1px 1px 1px #000, 1px 1px 1px #000, -1px -1px 1px #000, 1px -1px 1px #000, -1px 1px 1px #000, 1px 1px 1px #000;
}

.play_button:active {
  transform: scale(.98);
}

.play_button:hover {
  background: #ff8100;
}

.play_button_wrapper {
  position: absolute;
  top: 3px;
  left: calc(50% - 50px);
}

.WrapLeftTopBlock, .WrapRightTopBlock {
  height: 30px;
  width: 15%;
  border: 1px solid rgba(37, 160, 225, .2);
  background: rgba(17, 73, 105, 0.60);
  box-shadow: inset 0 0 2px #00d7ff;
  border-radius: 0 0 50px 50px;
  float: left;
  margin-left: 20%;
  filter: drop-shadow(0 0 2px black);
}

.WrapRightTopBlock {
  float: right;
  margin-right: 20%;
}

.leftTopBlock, .rightTopBlock {
  font-family: 'Audiowide', cursive;
  height: 27px;
  clear: both;
  background: rgba(8, 138, 210, .5);
  border-radius: 0 0 50px 50px;
  margin: 2px;
  line-height: 27px;
  color: #7ebc3d;
  box-shadow: inset 0 0 2px black;
  text-align: center;
  text-shadow: 0 -1px 1px #000, 0 -1px 1px #000, 0 1px 1px #000, 0 1px 1px #000, -1px 0 1px #000, 1px 0 1px #000, -1px 0 1px #000, 1px 0 1px #000, -1px -1px 1px #000, 1px -1px 1px #000, -1px 1px 1px #000, 1px 1px 1px #000, -1px -1px 1px #000, 1px -1px 1px #000, -1px 1px 1px #000, 1px 1px 1px #000;
}

@media (max-width: 1122px) {
  .WrapLeftTopBlock, .WrapRightTopBlock {
    width: 25%;
    margin-left: 10%;
  }

  .WrapRightTopBlock {
    margin-right: 10%;
  }
}

@media (max-width: 600px) {
  .WrapLeftTopBlock, .WrapRightTopBlock {
    width: 32%;
    margin-left: 3%;
  }

  .WrapRightTopBlock {
    margin-right: 3%;
  }
}

@media (max-width: 470px) {
  .WrapLeftTopBlock, .WrapRightTopBlock {
    display: none;
  }
}

</style>
