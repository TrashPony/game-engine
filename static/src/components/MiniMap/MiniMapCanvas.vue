<template>
  <div class="can">
    <div class="back"></div>
    <canvas width="200" height="200" id="canvasMapStatic" ref="canvasMapStatic"/>
    <canvas width="200" height="200" id="canvasMap" ref="canvasMap"
            @mousemove="fastMove($event)"
            @mouseup="moveCamera = false"
            @mouseout="moveCamera = false"
            @mousedown="fastMove($event, true)"/>
  </div>
</template>

<script>
import {gameStore} from '../../game/store'
import {Scene} from '../../game/create'
import {GetGlobalPos} from "../../game/map/gep_global_pos";
import {MapSize} from "../../game/map/createMap";

export default {
  name: "MiniMapCanvas",
  data() {
    return {
      offsetX: null,
      offsetY: null,
      radarMarks: {},
      mapScale: 2048, // размер карты при котором увелечение х1, следовательно 4096 будет х2
      scale: 1,
      map: 0,
      staticObjCount: 0,
      moveCamera: false,
      updaters: {
        static: null,
        other: null,
      }
    }
  },
  destroyed() {
    this.map = 0;
    if (this.updaters.static) clearInterval(this.updaters.static)
    if (this.updaters.other) clearInterval(this.updaters.other)
  },
  created() {
    let app = this;
    app.updaters.static = setInterval(function () {
      if (gameStore.map && gameStore.mapsState[gameStore.map.id].staticObjects.length > 0 && Scene && app.$refs['canvasMapStatic']) {
        if (app.map !== gameStore.map.id || gameStore.mapsState[gameStore.map.id].staticObjects.length !== app.staticObjCount) {
          app.createStaticMap()
        }
      }
    }, 1000);

    app.updaters.other = setInterval(function () {
      if (gameStore.map && Scene && Scene.cameras && app.$refs['canvasMap']) {
        let canvas = app.$refs['canvasMap'];
        let ctx = canvas.getContext("2d");

        let oldOffsetX = app.offsetX;
        let oldOffsetY = app.offsetY;

        app.offsetX = gameStore.map.x_size / canvas.width;
        app.offsetY = gameStore.map.y_size / canvas.height;

        if (oldOffsetX !== app.offsetX || oldOffsetY !== app.offsetY) {
          app.createStaticMap()
        }

        ctx.clearRect(0, 0, canvas.width, canvas.height);

        app.createMap(ctx);
      }
    }, 100)
  },
  methods: {
    fastMove(e, on) {

      if (on) this.moveCamera = true;

      if (this.moveCamera) {
        let x = e.offsetX * this.offsetX;
        let y = e.offsetY * this.offsetY;

        let pos = GetGlobalPos(x, y, gameStore.map.id);

        Scene.cameras.main.stopFollow();
        Scene.cameras.main.centerOn(pos.x, pos.y);
      }
    },
    getObjectColorByUserID(userID, fraction, battle_group_uuid) {

      if (battle_group_uuid) {
        if (gameStore.user.battle_group_uuid === battle_group_uuid) {
          return {strokeColor: "rgba(0, 0, 0, 0.7)", fillColor: "#659DFF"}
        } else {
          return {strokeColor: "rgba(0, 0, 0, 0.7)", fillColor: "#ff0000"}
        }
      }

      if (gameStore.user && userID === gameStore.user.id) {
        return {strokeColor: "rgb(0,0,0)", fillColor: "#19ff00"}
      } else {
        if (this.groupState.members && this.groupState.members.hasOwnProperty(userID)) {
          if (this.battleState) {
            return {strokeColor: "rgba(0, 0, 0, 0.7)", fillColor: "#659DFF"}
          } else {
            return {strokeColor: "rgba(0, 0, 0, 0.7)", fillColor: "#00ff1c"}
          }
        } else {
          if (this.battleState) {
            return {strokeColor: "rgba(0, 0, 0, 0.7)", fillColor: "#ff0000"}
          } else {
            if (fraction) {
              if (fraction === 'Replics') return {strokeColor: "rgba(0, 0, 0, 0.7)", fillColor: "#ee7015"}
              if (fraction === 'Reverses') return {strokeColor: "rgba(0, 0, 0, 0.7)", fillColor: "#659DFF"}
              if (fraction === 'Explores') return {strokeColor: "rgba(0, 0, 0, 0.7)", fillColor: "#7aba00"}
              if (fraction === 'APD') return {strokeColor: "rgba(0, 0, 0, 0.7)", fillColor: "#d7bc09"}
            } else {
              return {strokeColor: "rgba(0, 0, 0, 0.7)", fillColor: "#ffffff"}
            }
          }
        }
      }

      return {strokeColor: "rgba(0, 0, 0, 0.7)", fillColor: "#ffffff"}
    },
    fillGeoData(ctx, geoPoint) {
      let app = this;
      app.createMapRing(ctx, geoPoint.x, geoPoint.y, geoPoint.radius, "rgba(255, 255, 255, 0)", "#000000")
    },
    createStaticMap() {
      let app = this;
      app.map = gameStore.map.id;
      app.staticObjCount = gameStore.mapsState[gameStore.map.id].staticObjects.length;
      app.scale = gameStore.map.x_size / app.mapScale

      let canvas = app.$refs['canvasMapStatic'];
      let ctx = canvas.getContext("2d");
      ctx.clearRect(0, 0, canvas.width, canvas.height);

      ctx.fillStyle = "rgba(0,0,0,0.4)";
      ctx.fillRect(0, 0, canvas.width, canvas.height);
      app.staticMap(ctx);
    },
    createMapText(ctx, x, y, strokeColor, fillColor, text) {
      let app = this;

      if (fillColor) ctx.fillStyle = fillColor;
      if (strokeColor) ctx.strokeStyle = strokeColor;

      let textSize = 18
      x = (x / app.offsetX);
      y = (y / app.offsetY);

      ctx.textBaseline = 'middle';
      ctx.textAlign = "center";

      ctx.font = "bold " + textSize + "px Arial";
      ctx.lineWidth = 2;

      if (strokeColor) ctx.strokeText(text, x, y);
      if (fillColor) ctx.fillText(text, x, y);
    },
    createMapRing(ctx, x, y, radius, strokeColor, fillColor, needScale) {

      let app = this;

      ctx.beginPath();
      if (fillColor) ctx.fillStyle = fillColor;
      if (strokeColor) ctx.strokeStyle = strokeColor;

      if (needScale) {
        radius = radius * app.scale;
      }

      ctx.ellipse(
        (x / app.offsetX),
        (y / app.offsetY),
        radius / app.offsetX,
        radius / app.offsetY,
        0, 0, 2 * Math.PI, true);

      if (fillColor) ctx.fill();
      if (strokeColor) ctx.stroke();
    },
    createMapRect(ctx, x, y, width, height, strokeColor, fillColor, center, needScale) {

      let app = this;

      if (fillColor) ctx.fillStyle = fillColor;
      if (strokeColor) ctx.strokeStyle = strokeColor;

      if (needScale) {
        width = width * app.scale;
        height = height * app.scale;
      }

      width = width / app.offsetX;
      height = height / app.offsetY;

      if (center) {
        if (fillColor) ctx.fillRect((x / app.offsetX) - width / 2, (y / app.offsetY) - height / 2, width, height);
        if (strokeColor) ctx.strokeRect((x / app.offsetX) - width / 2, (y / app.offsetY) - height / 2, width, height);
      } else {
        if (fillColor) ctx.fillRect((x / app.offsetX), (y / app.offsetY), width, height);
        if (strokeColor) ctx.strokeRect((x / app.offsetX), (y / app.offsetY), width, height);
      }
    },
    createMapRhombus(ctx, x, y, width, height, strokeColor, fillColor, needScale) {
      let app = this;

      if (fillColor) ctx.fillStyle = fillColor;
      if (strokeColor) ctx.strokeStyle = strokeColor;

      if (needScale) {
        width = width * app.scale;
        height = height * app.scale;
      }

      width = width / app.offsetX;
      height = height / app.offsetY;

      x = x / app.offsetX;
      y = y / app.offsetY;

      ctx.beginPath();
      ctx.moveTo(x - width / 2, y);
      ctx.lineTo(x, y + height / 2);
      ctx.lineTo(x + width / 2, y);
      ctx.lineTo(x, y - height / 2);
      ctx.lineTo(x - width / 2, y);
      ctx.closePath();

      if (fillColor) ctx.fill();
      if (strokeColor) ctx.stroke();
    },
    createMapHexagon(ctx, x, y, radius, strokeColor, fillColor, needScale) {
      let app = this;

      if (fillColor) ctx.fillStyle = fillColor;
      if (strokeColor) ctx.strokeStyle = strokeColor;

      if (needScale) {
        radius = radius * app.scale;
      }

      radius = radius / app.offsetX;
      const a = 2 * Math.PI / 6;

      x = x / app.offsetX;
      y = y / app.offsetY;

      ctx.beginPath();
      for (let i = 0; i < 6; i++) {
        ctx.lineTo(x + radius * Math.cos(a * i), y + radius * Math.sin(a * i));
      }
      ctx.closePath();

      if (fillColor) ctx.fill();
      if (strokeColor) ctx.stroke();
    },
    createMapTriangle(ctx, x, y, width, height, strokeColor, fillColor, needScale) {
      let app = this;

      if (fillColor) ctx.fillStyle = fillColor;
      if (strokeColor) ctx.strokeStyle = strokeColor;

      if (needScale) {
        width = width * app.scale;
        height = height * app.scale;
      }

      width = width / app.offsetX;
      height = height / app.offsetY;

      x = x / app.offsetX;
      y = y / app.offsetY;

      ctx.beginPath();
      ctx.moveTo(x, y + height / 2);
      ctx.lineTo(x - width / 2, y + height / 2);
      ctx.lineTo(x, y - height / 2);
      ctx.lineTo(x + width / 2, y + height / 2);
      ctx.lineTo(x, y + height / 2);
      ctx.closePath();

      if (fillColor) ctx.fill();
      if (strokeColor) ctx.stroke();
    },
    staticMap(ctx) {
      let app = this;

      for (let i in gameStore.map.geo_data) {
        let obstacle = gameStore.map.geo_data[i];
        if (obstacle) {
          app.fillGeoData(ctx, obstacle)
        }
      }

      for (let i in gameStore.mapsState[gameStore.map.id].staticObjects) {
        let obj = gameStore.mapsState[gameStore.map.id].staticObjects[i];
        if (obj && obj.geo_data && obj.geo_data.length > 0) {
          for (let j in obj.geo_data) {
            let obstacle = obj.geo_data[j];
            app.fillGeoData(ctx, obstacle)
          }
        }
      }
    },
    createMap(ctx) {
      let app = this;

      for (let i in gameStore.objects) {
        let obj = gameStore.objects[i];
        if (obj && obj.geo_data && obj.geo_data.length > 0) {

          if (obj.build) {
            let {strokeColor, fillColor} = app.getObjectColorByUserID(obj.owner_id, obj.fraction)
            app.createMapRect(ctx, obj.x, obj.y, 50, 50, strokeColor, fillColor, true, true)
          } else {
            for (let j in obj.geo_data) {
              app.fillGeoData(ctx, obj.geo_data[j])
            }
          }
        }
      }

      for (let id in gameStore.units) {
        // TODO
      }

      // летящие пули
      for (let i in gameStore.bullets) {
        let bullet = gameStore.bullets[i];
        if (bullet && bullet.sprite && gameStore.bullets.hasOwnProperty(i)) {
          app.createMapRing(ctx, bullet.sprite.x - MapSize, bullet.sprite.y - MapSize, 15, "#7c6d17", "#ffd700", true)
        }
      }

      if (Scene && Scene.cameras.main) {
        app.createMapRect(
          ctx,
          Scene.cameras.main.worldView.x - MapSize,
          Scene.cameras.main.worldView.y - MapSize,
          Scene.cameras.main.worldView.width,
          Scene.cameras.main.worldView.height,
          "#fffc1f",
          null,
        )
      }
    }
  }
}
</script>

<style scoped>
#canvasMap {
  height: 100%;
  width: 100%;
  position: relative;
}

#canvasMapStatic {
  height: 100%;
  width: 100%;
  position: absolute;
}

.back {
  height: 100%;
  width: 100%;
  position: absolute;
  left: 0;
  top: 0;
  background: grey;
}

.can {
  position: relative;
  height: 100%;
  width: 100%;
}
</style>
