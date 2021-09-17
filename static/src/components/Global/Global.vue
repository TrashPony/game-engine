<template>
  <div>
    <app-mini-map/>
    <app-debug-panel/>
  </div>
</template>

<script>
import DebugPanel from "../DebugPanel/DebugPanel";
import MiniMap from "../MiniMap/MiniMap";
import {CreateGame} from "../../game/create";

export default {
  name: "Global",
  data() {
    return {}
  },
  mounted() {
    this.$store.dispatch("sendSocketData", JSON.stringify({
      event: "GetPlayers",
      service: "system",
    }))

    CreateGame();
  },
  methods: {},
  components: {
    AppMiniMap: MiniMap,
    AppDebugPanel: DebugPanel,
  }
}
</script>

<style scoped>

</style>
