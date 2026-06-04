<script setup lang="ts">
import { useAuth0 } from '@auth0/auth0-vue';
import { computed } from 'vue';

const { isAuthenticated, loginWithRedirect, logout, isLoading } = useAuth0();

const handleLogin = () => loginWithRedirect();
const handleLogout = () => logout({ logoutParams: { returnTo: globalThis.location.origin } });

const btnText = computed(() => (isAuthenticated.value ? 'Wyloguj' : 'Zaloguj'));
</script>

<template>
    <button
        class="btn btn--ghost btn--sm"
        :disabled="isLoading"
        @click="() => (isAuthenticated ? handleLogout() : handleLogin())"
    >
        {{ isLoading ? '...' : btnText }}
    </button>
</template>
