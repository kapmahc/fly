import Vue from 'vue'
import VueI18n from 'vue-i18n'
import ElementUI from 'element-ui'
import enUSE from 'element-ui/lib/locale/lang/en'
import zhHansE from 'element-ui/lib/locale/lang/zh-CN'
import zhHantE from 'element-ui/lib/locale/lang/zh-TW'

import messages from './locales'

Vue.use(VueI18n)

const LOCALE = 'locale'
const lang = localStorage.getItem(LOCALE) || 'en-US'

export const i18n = new VueI18n({
  locale: lang,
  messages
})

export const setLocale = (lng) => {
  localStorage.setItem(LOCALE, lng)
  window.location.reload()
}

export const loadLocales = () => {
  switch (lang) {
    case 'zh-Hans':
      Vue.use(ElementUI, {locale: zhHansE})
      break
    case 'zh-Hant':
      Vue.use(ElementUI, {locale: zhHantE})
      break
    default:
      Vue.use(ElementUI, {locale: enUSE})
      break
  }
}
