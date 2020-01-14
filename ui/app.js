import 'regenerator-runtime/runtime'
import HelloApp from './components/HelloApp.vue'
import store from './store'

const appInfo = {
  name: 'Hello',
  id: 'hello',
  icon: 'info',
  isFileEditor: false,
  extensions: [],
  config: {
    url: 'http://localhost:9105'
  }
}

const routes = [
  {
    name: 'hello',
    path: '/',
    components: {
      app: HelloApp
    }
  }
]

const navItems = [
  {
    name: 'Hello',
    iconMaterial: appInfo.icon,
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
