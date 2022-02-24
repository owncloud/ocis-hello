import 'regenerator-runtime/runtime'
import App from './components/App.vue'
import store from './store'

const appInfo = {
  name: 'Hello',
  id: 'hello',
  icon: 'chat-smile',
  isFileEditor: false
}

const routes = [
  {
    name: 'hello',
    path: '/',
    component: App
  }
]

const navItems = [
  {
    name: 'Hello',
    icon: appInfo.icon,
    route: {
      name: 'hello',
      path: `/${appInfo.id}/`
    }
  }
]

export default {
  appInfo,
  store,
  routes,
  navItems
}
