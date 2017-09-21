import Vue from 'vue'

import App from './App'
import router from './router'
import {loadLocales, i18n} from './intl'
import './layouts'

Vue.config.productionTip = false

loadLocales()

/* eslint-disable no-new */
new Vue({
  router,
  i18n,
  el: '#app',
  template: '<App/>',
  components: { App }
})
