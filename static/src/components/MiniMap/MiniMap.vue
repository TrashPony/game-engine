<template>
  <div id="miniMap" ref="miniMap"
       @mouseup="moveCamera = false"
       @mouseout="moveCamera = false"
       @mousedown="fastMove($event, true)"
       @mousemove="fastMove($event)">

    <div class="zoomButton" @mousedown="mouseZoomPress(-0.0025)" @mouseup="mouseZoomUp" title="зум"
         :style="{backgroundImage: 'url(' + require('../../assets/icons/zoom_minus.png') + ')'}"
         style="top: calc(5% + 25px);">

    </div>

    <div class="zoomButton" @mousedown="mouseZoomPress(+0.0025)" @mouseup="mouseZoomUp" title="зум"
         :style="{backgroundImage: 'url(' + require('../../assets/icons/zoom_plus.png') + ')'}"
         style="top: calc(5%);">
    </div>

  </div>
</template>

<script>
import Control from '../Window/Control';
import {Scene} from '../../game/create';
import {gameStore} from "../../game/store";
import {userUnit} from "../../game/update";
import {GetGlobalPos} from "../../game/map/gep_global_pos";
import {minimap} from "../../game/interface/mini_map";

export default {
  name: "MiniMap",
  data() {
    return {
      zoomChange: false,
      moveCamera: false,
    }
  },
  created() {
    window.addEventListener('wheel', this.wheelZoom);
  },
  mounted() {
    let wait = setInterval(function () {
      if (minimap.init) {
        let miniMap = document.getElementById("miniMap")
        if (miniMap) {
          miniMap.style.width = (minimap.size - 5) + 'px'
          miniMap.style.height = (minimap.size - 22) + 'px'
          miniMap.style.display = 'block'
          clearInterval(wait)
        }
      }
    }, 100)
  },
  destroyed() {
    window.removeEventListener('wheel', this.wheelZoom);
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
          app.zoom(null, size)
        } else {
          clearInterval(zoomer);
        }
      });
    },
    wheelZoom(event, size) {
      if (!gameStore.HoldAttackMouse) {
        this.zoom(event, size)
      }
    },
    zoom(event, size) {

      if (event) {
        size = (event.deltaY * 0.001) * -1
      }

      let zoom = Scene.cameras.main.zoom + size
      let minSize = 0.5

      if (zoom < minSize) {
        zoom = minSize
      } else if (zoom > 2) {
        zoom = 2
      }

      Scene.cameras.main.setZoom(zoom);
      this.$store.dispatch('changeSettings', {
        name: "ZoomCamera",
        count: zoom * 100
      });

      this.$store.commit({
        type: 'setZoomCamera',
        zoom: zoom,
      });
    },
    mouseZoomUp() {
      this.zoomChange = false;
    },
    fastMove(e, on) {
      if (e.path[0].id !== 'miniMap') return;
      if (!gameStore.gameSettings.follow_camera || !userUnit) {
        if (on) this.moveCamera = true;

        let offsetX = gameStore.map.x_size / minimap.size;
        let offsetY = gameStore.map.y_size / minimap.size;

        if (this.moveCamera) {
          let x = e.offsetX * offsetX;
          let y = e.offsetY * offsetY;

          let pos = GetGlobalPos(x, y, gameStore.map.id);

          Scene.cameras.main.stopFollow();
          Scene.cameras.main.centerOn(pos.x, pos.y);
        }
      }
    },
  },
  computed: {
    settings() {
      return this.$store.getters.getSettings
    }
  },
  components: {
    AppControl: Control,
  }
}
</script>

<style scoped>
#miniMap, #miniMap .zoomButton {
  /*background: rgb(8, 138, 210);*/
  box-shadow: 0 0 2px rgba(0, 0, 0, 1), inset 0 0 2px rgba(0, 0, 0, 1);
  border: 3px solid #25a0e1;
}

#miniMap {
  position: absolute;
  height: 178px;
  width: 195px;
  border-radius: 5px;
  top: 7px;
  right: 7px;
  user-select: none;
  padding: 19px 2px 2px;
  display: none;
}

#miniMap .zoomButton {
  height: 20px;
  width: 20px;
  background-color: rgba(0, 0, 0, 0.2);
  border-bottom-left-radius: 5px;
  border-top-left-radius: 5px;
  line-height: 20px;
  background-position: center;
  background-size: contain;
  position: absolute;
  right: calc(100% + 1px);
  top: 22px;
  border: 1px solid rgba(0, 0, 0, 0.2);
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

#miniMap .zoomButton .label {
  font-size: 11px;
  color: #ffed93;
  font-weight: bold;
  text-align: center;
  text-shadow: 0 -1px 1px #000000, 0 -1px 1px #000000, 0 1px 1px #000000, 0 1px 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px 0 1px #000000, 1px 0 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000, -1px -1px 1px #000000, 1px -1px 1px #000000, -1px 1px 1px #000000, 1px 1px 1px #000000;
  margin-top: 2px;
}

.image {
  height: 100%;
  width: 100%;
  background-size: contain;
  background-position: center;
  /*filter: contrast(45%) sepia(100%) hue-rotate(11deg) brightness(0.8) saturate(800%) drop-shadow(1px 1px 0px black);*/
}
</style>
