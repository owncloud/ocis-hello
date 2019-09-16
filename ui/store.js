const namespaced = true

const state = {
  config: null,
  message: ''
}

const getters = {
  config: state => state.config,
  message: state => state.message
}

const actions = {
  // Action triggered from within apps store
  loadConfig ({ commit }, config) {
    commit('LOAD_CONFIG', config)
  },

  submitName ({ commit, dispatch, getters }, value) {
    fetch(getters.config.fetchUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        'name': value
      })
    })
      .then(response => {
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
      .catch((error) => {
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
  namespaced,
  state,
  getters,
  actions,
  mutations
}
