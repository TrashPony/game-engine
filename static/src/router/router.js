import Vue from "vue";
import VueRouter from 'vue-router'
import Lobby from '../components/Lobby/Lobby'
import Global from '../components/Global/Global'
import Index from '../components/Index/Index'

Vue.use(VueRouter);

const router = new VueRouter({
  mode: 'history',
  routes: [
    {
      path: '/', component: Index, meta: {title: ""},
      children: []
    },
    {path: '/lobby', component: Lobby, meta: {title: ""}},
    {path: '/global', component: Global, meta: {title: ""}},
    {
      path: '/index.html',
      name: 'roomPage',
      component: Index
    },
    {
      path: '/*',
      name: 'roomPage',
      component: Index
    },
  ]
});

export default router
