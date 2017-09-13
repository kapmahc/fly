import jwtDecode from 'jwt-decode'

import {USERS_SIGN_IN, USERS_SIGN_OUT, REFRESH} from './actions'

const currentUser = (state = {}, action) => {
  switch (action.type) {
    case USERS_SIGN_IN:
      try {
        return jwtDecode(action.token)
      } catch (e) {
        console.error(e)
      }
      return {}
    case USERS_SIGN_OUT:
      return {}
    default:
      return state
  }
}

const siteInfo = (state = {
  languages: []
}, action) => {
  switch (action.type) {
    case REFRESH:
      return Object.assign({}, action.info)
    default:
      return state;
  }
}

export default {
  currentUser,
  siteInfo
}
