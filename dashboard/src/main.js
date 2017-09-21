import Vue from 'vue'

import App from './App'
import router from './router'
import store from './store'
import {loadLocales, i18n} from './intl'
import './layouts'

Vue.config.productionTip = false

loadLocales()

/* eslint-disable no-new */
new Vue({
  router,
  i18n,
  store,
  el: '#app',
  template: '<App/>',
  components: { App }
})
