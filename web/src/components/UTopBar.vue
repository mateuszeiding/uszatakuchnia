<script lang="ts" setup>
import { useAuth0 } from '@auth0/auth0-vue';
import { RouterLink } from 'vue-router';

import AuthButton from './auth/AuthButton.vue';

const { isAuthenticated } = useAuth0();
</script>
<template>
    <header class="px-7 py-2 row border-bottom bg-gray-ghost-white">
        <nav
            aria-label="page navigation"
            class="col-4 d-flex gap-4 align-items-center justify-content-start"
        >
            <RouterLink
                to="recipes"
                class="btn"
            >
                Recipes
            </RouterLink>
            <RouterLink
                to="ingredients"
                class="btn"
            >
                Ingredients
            </RouterLink>
        </nav>
        <div class="col d-flex justify-content-center">
            <h1 class="uppercase">Uszata Kuchnia</h1>
        </div>
        <nav
            aria-label="admin page navigation"
            class="col-4 d-flex gap-4 align-items-center justify-content-end"
        >
            <RouterLink
                v-if="isAuthenticated"
                to="/type/new"
                class="btn"
            >
                Nowy typ
            </RouterLink>
            <RouterLink
                v-if="isAuthenticated"
                to="/recipes/new"
                class="btn"
            >
                Nowy przepis
            </RouterLink>
            <AuthButton />
        </nav>
    </header>
</template>
