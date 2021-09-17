import Vue from 'vue'
import App from './App.vue'
import axios from 'axios'
import store from './store/store'
import router from './router/router'

require('webpack-jquery-ui');
require('webpack-jquery-ui/css');

// инструмент для http запросов
Vue.use({
  install(Vue) {
    Vue.prototype.$api = axios.create({
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      }
    });
  }
});

new Vue({
  el: '#app',
  router: router,
  store,
  render: h => h(App),
});
