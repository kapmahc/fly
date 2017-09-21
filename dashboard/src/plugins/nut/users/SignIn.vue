<template>
  <application-layout>
    <h1>sign in {{ $t('nut.users.sign-in.title') }}</h1>
    <router-link :to="{ name: 'users.sign-in'}">Sign in</router-link>
    <router-link :to="{ name: 'users.sign-up'}">Sign Up</router-link>
    <router-link :to="{ name: 'home'}">Home</router-link>
  </application-layout>
</template>

<script>
export default {
  name: 'users-sign-in',
  data () {
    return {
    }
  }
}
</script>
