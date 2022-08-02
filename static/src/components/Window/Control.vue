<template>
  <div class="windowsHead"
       v-if="checkViewPort"
       @mousedown="mouseDown"
       @mouseup="mouseUp">
    <span class="windowsHeadTitle">{{ head }}</span>

    <div class="closeWindowButton" v-if="close" @click="closeWindow"/>

  </div>
</template>

<script>
import {gameStore} from "../../game/store";

export default {
  name: "Control",
  props: ['head', 'move', 'close', 'refWindow', 'resizeFunc', 'minSize', 'closeFunc', 'noHeight', 'noWidth', 'noPos'],
  data() {
    return {
      block: null,
      state: {
        id: 0,
        left: 0,
        top: 0,
        height: 0,
        width: 0,
        open: false,
      }
    }
  },
  computed: {
    wState() {
      return this.$store.getters.getInterfaceState
    },
    checkViewPort() {

      setTimeout(function () {
        this.getWindowState()
      }.bind(this), 100)

      return this.$store.getters.getCheckViewPort
    }
  },
  destroyed() {
    if (!this.block) return;
    this.setState(this.state.id, this.state.left, this.state.top, this.state.height, this.state.width, false);
  },
  mounted() {
    this.getWindowState()
  },
  methods: {
    playSound(sound, k) {
      if (sound === "button_press.mp3") k = 0.2
      if (sound === "select_sound.mp3") k = 0.1

      this.$store.dispatch('playSound', {
        sound: sound,
        k: k,
      });
    },
    getWindowState() {
      let app = this;
      let block = app.$parent.$refs[app.$props.refWindow];
      app.block = block;

      if (!app.block) return;

      if (this.wState && this.wState[app.$props.refWindow]) {
        if (!app.$props.noPos) block.style.left = this.wState[app.$props.refWindow].left + "px";
        if (!app.$props.noPos) block.style.top = this.wState[app.$props.refWindow].top + "px";
        if (!app.$props.noHeight) block.style.height = this.wState[app.$props.refWindow].height + "px";
        if (!app.$props.noWidth) block.style.width = this.wState[app.$props.refWindow].width + "px";
      }

      this.setState(block.id, $(block).position().left, $(block).position().top, $(block).height(), $(block).width(), true);

      if (app.$props.resizeFunc) {

        app.$props.resizeFunc(null, null, $(block));

        $(block).resizable({
          minHeight: app.$props.minSize.height,
          minWidth: app.$props.minSize.width,
          handles: "all",
          resize: function (event, ui) {
            app.$props.resizeFunc(event, ui, $(this))
          },
          stop: function (e, ui) {
            app.setState(this.id, $(this).position().left, $(this).position().top, $(this).height(), $(this).width(), true);
          }
        });
      }

      app.checkModalInViewPort()
    },
    checkModalInViewPort() {
      let app = this;
      let block = app.block;
      if (!block) {
        return;
      }

      let top = $(block).position().top;
      let left = $(block).position().left;

      if (top < 5) top = 5;
      if ((top + $(block).outerHeight() + 10) - $(window).height() > 0) {
        top = $(window).height() - $(block).outerHeight() - 10
      }

      if (left < 5) left = 5;
      if ((left + $(block).outerWidth() + 5) - $(window).width() > 0) {
        left = $(window).width() - $(block).outerWidth() - 5
      }

      $(block).css({left: left, top: top});

      return {left: left, top: top};
    },
    mouseDown() {
      let app = this;
      let block = app.block;

      gameStore.MouseMoveInterface = true;

      if (this.$props.move) {
        $(block).draggable({
          disabled: false,
          drag: function (e, ui) {
            let top = ui.position.top;
            let left = ui.position.left;

            if (top < 5) ui.position.top = 5;
            if ((top + $(block).outerHeight() + 10) - $(window).height() > 0) {
              ui.position.top = $(window).height() - $(block).outerHeight() - 10
            }

            if (left < 5) ui.position.left = 5;
            if ((left + $(block).outerWidth() + 5) - $(window).width() > 0) {
              ui.position.left = $(window).width() - $(block).outerWidth() - 5
            }
          },
          stop: function (event, ui) {
            gameStore.MouseMoveInterface = false;
            app.checkModalInViewPort();
            app.setState(block.id, $(block).position().left, $(block).position().top, $(block).height(), $(block).width(), true);
            $(block).draggable({
              disabled: true,
            });
          }
        });
      }
    },
    mouseUp() {

      gameStore.MouseMoveInterface = false;

      $(this.block).draggable({
        disabled: true,
      });
    },
    setState(id, left, top, height, width, open) {

      this.state.id = id;
      this.state.left = left;
      this.state.top = top;
      this.state.height = height;
      this.state.width = width;
      this.state.open = open;

      if (this.$store.getters.getWSConnectState.connect) {

        this.$store.dispatch("sendSocketData", JSON.stringify({
          event: "setWindowState",
          service: "system",
          resolution: $(window).width() + ':' + $(window).height(),
          name: id,
          left: Math.round(Number(left)),
          top: Math.round(Number(top)),
          height: Math.round(Number(height)),
          width: Math.round(Number(width)),
          open: open,
        }))

        this.$store.commit({
          type: 'setWindowState',
          id: id,
          state: {
            left: Math.round(Number(left)),
            top: Math.round(Number(top)),
            height: Math.round(Number(height)),
            width: Math.round(Number(width)),
            open: open,
          }
        });
      }
    },
    closeWindow() {
      this.playSound('window_close.mp3', 0.3)

      if (this.$props.closeFunc) {
        this.$props.closeFunc();
      }

      let block = this.block;
      this.$store.commit({
        type: 'toggleWindow',
        id: block.id,
        component: '',
        forceClose: true,
      });
      this.setState(block.id, $(block).position().left, $(block).position().top, $(block).height(), $(block).width(), false);
    }
  },
}
</script>

<style scoped>
.windowsHead {
  top: -1px;
  left: -1px;
  width: calc(100% - 15px);
  display: block;
  position: absolute;
  height: 17px;
  background: #8aaaaa;
  border: 1px solid #25a0e1;
  border-radius: 5px 5px 0 0;
  box-shadow: inset 0 0 3px rgba(0, 0, 0, 1);
  text-align: left;
  text-indent: 5px;
  font-size: 11px;
  line-height: 18px;
  transition: 200ms;
  color: black;
  text-shadow: none;
  user-select: none;
  word-wrap: normal;
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
  padding-right: 15px;
}

.windowsHead:hover {
  background: #a0c3c3;
}

.closeWindowButton {
  position: absolute;
  right: 4px;
  top: 2px;
  border: 1px solid #7a858c;
  width: 12px;
  height: 11px;
  border-radius: 2px;
  background: #fcbb00;
  opacity: 0.7;
  float: right;
  transition: 200ms;
}

.closeWindowButton:hover {
  opacity: 1;
}

.closeWindowButton:active {
  transform: scale(0.80);
}

.closeWindowButton:before, .closeWindowButton:after {
  position: absolute;
  left: 5px;
  content: ' ';
  height: 11px;
  width: 2px;
  background-color: #333;
}

.closeWindowButton:before {
  transform: rotate(45deg);
}

.closeWindowButton:after {
  transform: rotate(-45deg);
}
</style>
