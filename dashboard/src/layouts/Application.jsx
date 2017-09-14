import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {Layout, Menu, Breadcrumb} from 'antd'

import {signIn, signOut, refresh} from '../actions'

const {Header, Content, Footer} = Layout

class Widget extends Component {
  render() {
    const {children} = this.props
    return (
      <Layout>
        <Header style={{
          position: 'fixed',
          width: '100%'
        }}>
          <div className="logo"/>
          <Menu theme="dark" mode="horizontal" defaultSelectedKeys={[]} style={{
            lineHeight: '64px'
          }}>
            <Menu.Item key="1">nav 1</Menu.Item>
            <Menu.Item key="2">nav 2</Menu.Item>
            <Menu.Item key="3">nav 3</Menu.Item>
          </Menu>
        </Header>
        <Content style={{
          padding: '0 50px',
          marginTop: 64
        }}>
          <Breadcrumb style={{
            margin: '12px 0'
          }}>
            <Breadcrumb.Item>Home</Breadcrumb.Item>
            <Breadcrumb.Item>List</Breadcrumb.Item>
            <Breadcrumb.Item>App</Breadcrumb.Item>
          </Breadcrumb>
          <div style={{
            background: '#fff',
            padding: 24,
            minHeight: 380
          }}>{children}</div>
        </Content>
        <Footer style={{
          textAlign: 'center'
        }}>
          Ant Design ©2016 Created by Ant UED
        </Footer>
      </Layout>
    )
  }
}

Widget.propTypes = {
  children: PropTypes.node.isRequired,
  push: PropTypes.func.isRequired,
  refresh: PropTypes.func.isRequired,
  signIn: PropTypes.func.isRequired,
  signOut: PropTypes.func.isRequired,
  user: PropTypes.object.isRequired,
  info: PropTypes.object.isRequired,
  breads: PropTypes.array.isRequired,
  intl: intlShape.isRequired
}

const WidgetT = injectIntl(Widget)

export default connect(state => ({user: state.currentUser, info: state.siteInfo}), {
  push,
  signIn,
  refresh,
  signOut
},)(WidgetT)
