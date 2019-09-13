<template>
  <div>
    <h1>
      Hello
    </h1>

    <form v-on:submit.prevent="submitHello" action="#" method="post">
      <input v-model="name" placeholder="What is your name?">

      <button type="submit">
        Submit
      </button>
    </form>

    <p v-if="message">
      {{ message }}
    </p>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'

export default {
  name: 'HelloApp',
  methods: {
    ...mapActions('Hello', ['submitHello'])
  },
  computed: {
    ...mapGetters('Hello', ['name', 'message'])
  }
}
</script>
