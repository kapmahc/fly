import React, {Component} from 'react'
import {FormattedMessage} from 'react-intl'
class Widget extends Component {
  render() {
    return (
      <div>home
        <FormattedMessage id="buttons.submit"/>
      </div>
    )
  }
}

export default Widget
