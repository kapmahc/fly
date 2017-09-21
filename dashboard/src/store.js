import Vue from 'vue'
import Vuex from 'vuex'
import jwtDecode from 'jwt-decode'

Vue.use(Vuex)

const store = new Vuex.Store({
  state: {
    currentUser: null,
    siteInfo: {
      languages: []
    }
  },
  mutations: {
    refresh (state, info) {
      state.siteInfo = info
    },
    signIn (state, token) {
      try {
        state.currentUser = jwtDecode(token)
      } catch (e) {
        console.error(e)
      }
    },
    signOut (state) {
      state.currentUser = null
    }
  }
})

export default store
