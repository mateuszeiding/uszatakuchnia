<script setup lang="ts">
import { useAuth0 } from '@auth0/auth0-vue'
import { computed } from 'vue'

const { isAuthenticated, loginWithRedirect, logout, isLoading } = useAuth0()

const handleLogin = () => {
  loginWithRedirect()
}

const handleLogout = () => {
  logout({
    logoutParams: {
      returnTo: window.location.origin,
    },
  })
}

const btnText = computed(() => (isAuthenticated.value ? 'Wyloguj' : 'Zaloguj'))
</script>
<template>
  <button
    @click="() => (isAuthenticated ? handleLogout() : handleLogin())"
    class="button login"
    :disabled="isLoading"
  >
    {{ isLoading ? 'Loading...' : btnText }}
  </button>
</template>
