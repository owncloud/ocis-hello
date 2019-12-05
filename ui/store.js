const hello = require('./client/hello')

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
    hello.Greet({
      $domain: getters.config.url,
      name: value
    })
      .then(response => {
        console.log(response)

        if (response.ok) {
          response.json()
            .then(json => {
              commit('SET_MESSAGE', json.message)
            })
        } else {
          dispatch('showMessage', {
            title: 'Response failed',
            desc: response.statusText,
            status: 'danger'
          }, { root: true })
        }
      })
      .catch(error => {
        console.log(error)

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
