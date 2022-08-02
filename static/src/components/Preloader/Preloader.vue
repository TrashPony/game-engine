<template>
  <div class="lds-default">
    <div></div>
    <div></div>
    <div></div>
    <div></div>
    <div></div>
    <div></div>
    <div></div>
    <div></div>
    <div></div>
    <div></div>
    <div></div>
    <div></div>
  </div>
</template>

<script>
export default {
  name: "Preloader"
}
</script>

<style scoped>
.lds-default {
  display: inline-block;
  position: relative;
  width: 43px;
  height: 43px;
  float: right;
  margin-top: 3px;
  filter: drop-shadow(1px 1px 1px black);
}

.lds-default div {
  position: absolute;
  width: 6px;
  height: 6px;
  background: #fff;
  border-radius: 50%;
  animation: lds-default 1.2s linear infinite;
}

.lds-default div:nth-child(1) {
  animation-delay: 0s;
  top: 18px;
  left: 33px;
}

.lds-default div:nth-child(2) {
  animation-delay: -0.1s;
  top: 11px;
  left: 31px;
}

.lds-default div:nth-child(3) {
  animation-delay: -0.2s;
  top: 6px;
  left: 26px;
}

.lds-default div:nth-child(4) {
  animation-delay: -0.3s;
  top: 3px;
  left: 18px;
}

.lds-default div:nth-child(5) {
  animation-delay: -0.4s;
  top: 5px;
  left: 11px;
}

.lds-default div:nth-child(6) {
  animation-delay: -0.5s;
  top: 11px;
  left: 5px;
}

.lds-default div:nth-child(7) {
  animation-delay: -0.6s;
  top: 18px;
  left: 3px;
}

.lds-default div:nth-child(8) {
  animation-delay: -0.7s;
  top: 26px;
  left: 6px;
}

.lds-default div:nth-child(9) {
  animation-delay: -0.8s;
  top: 31px;
  left: 11px;
}

.lds-default div:nth-child(10) {
  animation-delay: -0.9s;
  top: 33px;
  left: 18px;
}

.lds-default div:nth-child(11) {
  animation-delay: -1s;
  top: 31px;
  left: 26px;
}

.lds-default div:nth-child(12) {
  animation-delay: -1.1s;
  top: 26px;
  left: 31px;
}

@keyframes lds-default {
  0%, 20%, 80%, 100% {
    transform: scale(0.5);
  }
  50% {
    transform: scale(1);
  }
}
</style>
