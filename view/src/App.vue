<template>
  <v-app style="background-color: #ffffff">
    <v-content>
      <Home/>
    </v-content>
  </v-app>
</template>

<script>
import Home from "./views/Home"
export default {
  name: 'App',
  components: {
    Home,
  },
  data: () => ({
    //
  }),
};
</script>

<style>
</style>
