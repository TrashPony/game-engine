import Vue from "vue";
import VueRouter from 'vue-router'
import Login from '../components/Auth/Login'
import Lobby from '../components/Lobby/Lobby'
import Global from '../components/Global/Global'
import Registration from '../components/Auth/Registration'
import Index from '../components/Index/Index'
import Gate from '../components/Gate/Gate'

Vue.use(VueRouter);

const router = new VueRouter({
  mode: 'history',
  routes: [
    {
      path: '/', component: Index, meta: {title: ""},
      children: []
    },
    {path: '/login', component: Login, meta: {title: ""}},
    {path: '/registration', component: Registration, meta: {title: ""}},
    {path: '/lobby', component: Lobby, meta: {title: ""}},
    {path: '/global', component: Global, meta: {title: ""}},
    {path: '/gate', component: Gate, meta: {title: ""}},
  ]
});

export default router
