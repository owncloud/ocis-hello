import { request } from 'https'

const namespaced = true

const state = {
  message: ''
}

const getters = {
  name: state => state.name,
  message: state => state.message
}

const actions = {
  submitHello ({ commit, dispatch }, value) {
    fetch(`http://localhost:8380/api/hello`, {
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
              commit('MESSAGE', json.message)
            })
        } else {
          dispatch('showMessage', {
            title: 'Response failed',
            desc: request.statusText,
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
  MESSAGE (state, payload) {
    state.message = payload
  }
}

export default {
  namespaced,
  state,
  getters,
  actions,
  mutations
}
