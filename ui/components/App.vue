<template>
  <div class="oc-flex oc-flex-column oc-flex-middle oc-p-l">
    <h1>Hello, tell us your name</h1>
    <form v-on:submit.prevent="submitName(name)" action="#" method="post">
      <oc-text-input v-model="name" placeholder="Your name" />
      <oc-button variation="primary" class="oc-width-1-1 oc-mt-m">
        Submit
      </oc-button>
    </form>
    <p class="oc-text-lead" v-if="message" v-text="message" />
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
