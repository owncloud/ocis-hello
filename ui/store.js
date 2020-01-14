import { Greet } from './client/hello'

const state = {
  config: null,
  message: ''
}

const getters = {
  config: state => state.config,
  message: state => state.message
}

const actions = {
  loadConfig ({ commit }, config) {
    commit('LOAD_CONFIG', config)
  },

  submitName ({ commit, dispatch, getters }, value) {
    Greet({
      $domain: getters.config.url,
      body: { name: value }
    })
      .then(response => {
        console.log(response)

        if (response.status === 200 || response.status === 201) {
          commit('SET_MESSAGE', response.data.message)
        } else {
          dispatch('showMessage', {
            title: 'Response failed',
            desc: response.statusText,
            status: 'danger'
          }, { root: true })
        }
      })
      .catch(error => {
        console.error(error)

        dispatch('showMessage', {
          title: 'Saving your name failed',
          desc: error.message,
          status: 'danger'
        }, { root: true })
      })
  }
}

const mutations = {
  SET_MESSAGE (state, payload) {
    state.message = payload
  },

  LOAD_CONFIG (state, config) {
    state.config = config
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
