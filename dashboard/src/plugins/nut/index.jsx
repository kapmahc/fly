import Install from './Install'

import UsersSignIn from './users/SignIn'
import UsersSignUp from './users/SignUp'

export default[
  {
    path : '/install',
    component : Install
  }, {
    path : '/users/sign-in',
    component : UsersSignIn
  }, {
    path : '/users/sign-up',
    component : UsersSignUp
  }
]
