<template>
  <div class="uk-flex uk-flex-column uk-flex-middle uk-margin-xlarge-top">
    <h1>
      Hello, tell us your name
    </h1>

    <form v-on:submit.prevent="submitName(name)" action="#" method="post">
      <oc-text-input v-model="name" placeholder="Your name" />

      <oc-button variation="primary" class="uk-width-1-1 uk-margin-top">
        Submit
      </oc-button>
    </form>

    <p class="uk-text-lead" v-if="message" v-text="message" />
  </div>
</template>

<script>
export default {
  name: 'App',
  data: function () {
    return {
      name: 'World'
    }
  },
  computed: {
    message () {
      return this.$store.getters['Hello/message']
    }
  },
  methods: {
    submitName (name) {
      this.$store.dispatch('Hello/submitName', name)
    }
  }
}
</script>
