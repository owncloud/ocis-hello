import 'core-js/stable'
import 'regenerator-runtime/runtime'
import HelloApp from './components/HelloApp.vue'

const appInfo = {
  name: 'Hello',
  id: 'hello',
  icon: 'folder',
  isFileEditor: false,
  extensions: [],
  config: {
    url: 'http://localhost:9105'
  }
}

const store = require('./store.js')

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

export default define({
  appInfo,
  store,
  routes,
  navItems
})
