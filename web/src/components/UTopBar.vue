<script lang="ts" setup>
import { useAuth0 } from '@auth0/auth0-vue'
import AuthButton from './auth/AuthButton.vue'
import { RouterLink } from 'vue-router'

const { isAuthenticated } = useAuth0()
</script>
<template>
  <header class="px-7 py-2 row border-bottom bg-gray-ghost-white">
    <div class="col-2 d-flex align-items-center justify-content-start">
      <RouterLink to="/type/new" v-if="isAuthenticated">Nowy typ</RouterLink>
      <RouterLink to="/recipes/new" v-if="isAuthenticated">Nowy przepis</RouterLink>
    </div>

    <div class="col d-flex justify-content-center">
      <h1 class="uppercase">Uszata Kuchnia</h1>
    </div>
    <div class="col-2 d-flex align-items-center justify-content-end">
      <AuthButton />
    </div>
  </header>
</template>
