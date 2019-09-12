const namespaced = true

const state = {
  name: 'World'
}

const getters = {
  name: state => state.name
}

export default {
  namespaced,
  state,
  getters
}
