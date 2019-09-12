<template>
  <div class="uk-flex uk-flex-center uk-flex-middle">
    <h1>
      Hello {{ name }}
    </h1>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'HelloApp',
  computed: {
    ...mapGetters('Hello', ['name'])
  }
}
</script>
