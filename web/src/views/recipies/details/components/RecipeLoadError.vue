<script lang="ts" setup>
import ErrorState from '@/views/errors/ErrorState.vue';

// Sieć / 5xx / timeout (20.11): „Spróbuj ponownie” ponawia pobranie w miejscu, bez przeładowania strony.
defineProps<{ retrying: boolean }>();
const emit = defineEmits<{ retry: [] }>();
</script>

<template>
    <ErrorState
        kicker="Brak połączenia"
        title="Nie udało się wczytać przepisu"
        warn
    >
        To raczej internet niż przepis. Sprawdź połączenie i spróbuj jeszcze raz.
        <template #actions>
            <button
                type="button"
                class="btn btn--accent"
                :disabled="retrying"
                @click="emit('retry')"
            >
                <svg
                    width="15"
                    height="15"
                    viewBox="0 0 16 16"
                    fill="none"
                    aria-hidden="true"
                >
                    <path
                        d="M13 8a5 5 0 11-1.5-3.6M13 2.5v3h-3"
                        stroke="currentColor"
                        stroke-width="1.5"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    />
                </svg>
                {{ retrying ? 'Wczytuję…' : 'Spróbuj ponownie' }}
            </button>
            <RouterLink
                to="/recipes"
                class="btn btn--ghost"
            >
                Wszystkie przepisy
            </RouterLink>
        </template>
    </ErrorState>
</template>
