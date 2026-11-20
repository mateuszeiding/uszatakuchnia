<script lang="ts" setup>
import { useRoute } from 'vue-router';

import ErrorState from './ErrorState.vue';

// Unknown URL (catch-all in the router): "Nie ma takiej strony" with the typed path. No redirect:
// the URL in the address bar stays as typed.
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
