<template>
  <div>
    <h1>home</h1>
    <router-link :to="{ name: 'users.sign-in'}">Sign in</router-link>
    <router-link :to="{ name: 'users.sign-up'}">Sign Up</router-link>
    <router-link :to="{ name: 'home'}">Home</router-link>
  </div>
</template>

<script>
export default {
  name: 'home',
  data () {
    return {
    }
  }
}
</script>
