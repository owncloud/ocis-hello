const namespaced = true

const state = {
  name: 'World',
  message: ''
}

const getters = {
  name: state => state.name,
  message: state => state.message
}

const actions = {
  submitHello (context, value) {
    fetch(`http://localhost:8380/api/hello`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        'name': this.state.Hello.name
      })
    })
      .then(response => {
        if (response.ok) {
          response.json()
            .then(json => {
              context.commit('MESSAGE', json.message)
            })
        } else {
          console.error('response', response)
        }
      })
      .catch((error) => {
        console.error('catch', error)
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
