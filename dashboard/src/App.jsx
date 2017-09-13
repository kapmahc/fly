import React, {Component} from 'react'
import {Route} from 'react-router'
import {Switch} from 'react-router-dom'
import {createStore, combineReducers, applyMiddleware} from 'redux'
import {Provider} from 'react-redux'
import createHistory from 'history/createBrowserHistory'
import {ConnectedRouter, routerReducer, routerMiddleware} from 'react-router-redux'
import moment from 'moment'
import {LocaleProvider} from 'antd'
import {IntlProvider, addLocaleData} from 'react-intl'

import {detectLocale} from './locales'
import reducers from './reducers'
import routes from './plugins'
import Home from './plugins/nut/Home'
import NoMatch from './plugins/nut/NoMatch'

const history = createHistory()
const middleware = routerMiddleware(history)

const store = createStore(combineReducers({
  ...reducers,
  router: routerReducer
}), applyMiddleware(middleware))

const intl = detectLocale()
moment.locale(intl.moment)
addLocaleData(intl.data)

class App extends Component {
  render() {
    return (
      <Provider store={store}>
        <LocaleProvider locale={intl.antd}>
          <IntlProvider locale={intl.locale} messages={intl.messages}>
            <ConnectedRouter history={history}>
              <Switch>
                <Route exact path="/" component={Home}/> {routes.map((r, i) => <Route key={i} path={r.path} component={r.component}/>)}
                <Route component={NoMatch}/>
              </Switch>
            </ConnectedRouter>
          </IntlProvider>
        </LocaleProvider>
      </Provider>
    )
  }
}

export default App
