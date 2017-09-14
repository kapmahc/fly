import React, {Component} from 'react'
import {FormattedMessage} from 'react-intl'

import Layout from '../../layouts/Application'

class Widget extends Component {
  render() {
    return (
      <Layout breads={[]}>home
        <FormattedMessage id="buttons.submit"/>
      </Layout>
    )
  }
}

export default Widget
