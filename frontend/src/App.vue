<template>
  <div class="roo">
    <clip-input />
  </div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import ClipInput from "./components/ClipInput.vue";

@Options({
  components: {
    ClipInput
  }
})
export default class App extends Vue {}
</script>
