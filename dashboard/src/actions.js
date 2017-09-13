export const REFRESH = 'refresh'
export const USERS_SIGN_IN = "users.sign-in"
export const USERS_SIGN_OUT = "users.sign-out"

export const signIn = (token) => {
  return {type: USERS_SIGN_IN, token}
}

export const signOut = () => {
  return {type: USERS_SIGN_OUT}
}

export const refresh = (info) => {
  return {type: REFRESH, info}
}
