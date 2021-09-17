<template>
  <div id="miniMap" ref="miniMap" @mousedown="toUp">

    <app-control v-bind:head="'Карта'"
                 v-bind:move="true"
                 v-bind:close="false"
                 v-bind:refWindow="'miniMap'"
                 v-bind:resizeFunc="resize"
                 v-bind:minSize="{height: 100, width: 100}"/>

    <app-mini-map-canvas/>

    <div class="zoomButton" title="на транспорт" @mousedown="mouseZoomPress(-0.0025)" @mouseup="mouseZoomUp"
         style="background-image: url('https://img.icons8.com/officel/16/000000/zoom-out.png'); top: 44px;">
    </div>

    <div class="zoomButton" @mousedown="mouseZoomPress(+0.0025)" @mouseup="mouseZoomUp"
         style="background-image: url('https://img.icons8.com/officel/16/000000/zoom-in.png'); top: 19px;">
    </div>

  </div>
</template>

<script>
import Control from '../Window/Control';
import MiniMapCanvas from './MiniMapCanvas';

import {Scene} from '../../game/create';
import {gameStore} from "../../game/store";

export default {
  name: "MiniMap",
  data() {
    return {
      zoomChange: false
    }
  },
  methods: {
    toUp() {
      this.$store.commit({
        type: 'setWindowZIndex',
        id: this.$el.id,
      });
    },
    resize(event, ui, el) {

      let sizeCanvas = function (id) {
        el.find(id).css("height", el.height());
        el.find(id).css("width", el.width());
        el.find(id).prop('width', el.height())
        el.find(id).prop('height', el.width())
      }

      sizeCanvas('#canvasFog')
      sizeCanvas('#canvasMap')
      sizeCanvas('#canvasMapStatic')
    },
    mouseZoomPress(size) {
      let app = this;

      app.zoomChange = true;
      let zoomer = setInterval(function () {
        if (app.zoomChange) {
          Scene.cameras.main.setZoom(Scene.cameras.main.zoom + size);

          if (gameStore.mapEditor) return;

          let minSize = 0.5
          if (app.currentUser.user_role === 'admin') {
            minSize = 0.2
          }

          if (Scene.cameras.main.zoom < minSize) {
            Scene.cameras.main.setZoom(minSize);
          } else if (Scene.cameras.main.zoom > 2) {
            Scene.cameras.main.setZoom(2);
          }
        } else {
          clearInterval(zoomer);
        }
      });
    },
    mouseZoomUp() {
      this.zoomChange = false;
    },
  },
  computed: {
    currentUser() {
      return this.$store.getters.getGameUser
    },
  },
  components: {
    AppControl: Control,
    AppMiniMapCanvas: MiniMapCanvas,
  }
}
</script>

<style scoped>
#miniMap, #miniMap .zoomButton {
  background: rgb(8, 138, 210);
  box-shadow: 0 1px 2px rgba(0, 0, 0, .2);
  border: 1px solid #25a0e1;
}

#miniMap {
  position: absolute;
  height: 200px;
  width: 200px;
  border-radius: 5px;
  top: 5px;
  right: 5px;
  user-select: none;
  padding: 19px 2px 2px;
}

#miniMap .zoomButton {
  height: 20px;
  width: 20px;
  background-color: rgb(19, 76, 105);
  border-right: transparent;
  border-bottom-left-radius: 10px;
  border-top-left-radius: 10px;
  color: #f9ff00;
  font-size: 29px;
  line-height: 20px;
  font-weight: 900;
  transition: 1s;
  box-shadow: inset 0 0 2px black;
  background-position: center;
  background-size: 20px;
  background-repeat: no-repeat;
  position: absolute;
  right: calc(100% + 1px);
  top: 22px;
}

#miniMap .zoomButton:hover {
  cursor: pointer;
  background-color: rgb(33, 176, 255);
}

#miniMap .zoomButton:nth-child(5) {
  line-height: 15px;
}

#miniMap .zoomButton:nth-child(6) {
  font-weight: 100;
}

#miniMap .topButton:nth-child(1) {
  font-size: 22px;
  line-height: 0;
}
</style>
