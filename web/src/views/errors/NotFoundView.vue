<script lang="ts" setup>
import { useRoute } from 'vue-router';

import ErrorState from './ErrorState.vue';

// Nieznany adres (catch-all w routerze): „Nie ma takiej strony” z wpisaną ścieżką. Bez przekierowania —
// adres w pasku zostaje taki, jak wpisany (20.11).
const route = useRoute();
</script>

<template>
    <ErrorState
        kicker="404"
        title="Nie ma takiej strony"
    >
        Pod adresem
        <code>{{ route.path }}</code>
        nic tu nie ma. Może link uciął się przy kopiowaniu.
        <template #actions>
            <RouterLink
                to="/recipes"
                class="btn btn--accent"
            >
                Przepisy
            </RouterLink>
            <RouterLink
                to="/challenges"
                class="btn btn--secondary"
            >
                Wyzwania
            </RouterLink>
        </template>
    </ErrorState>
</template>
