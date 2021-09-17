import Vue from "vue";
import Vuex from 'vuex';
import mutations from './mutations';
import getters from './getters';
import actions from './actions';
import {getDefaultState} from "./default_state";

Vue.use(Vuex);

let state = getDefaultState();
// выносим соеденения что бы не терять их при перегрузке стореджа
state.connections = {
  chat: {
    socket: null,
    connect: false,
    error: false,
  }
};

let store = new Vuex.Store({
  state: state,
  mutations: mutations,
  getters: getters,
  actions: actions,
});

export default store
