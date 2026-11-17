<script lang="ts" setup>
import ServingsStepper from '@/views/recipies/details/components/ServingsStepper.vue';
import { servingsLabel } from '@/views/recipies/details/recipeView';

// Nagłówek listy zakupów: powrót, Drukuj (desktop), tytuł, licznik porcji, postęp „n z m”.
// Na telefonie przyklejony pod paskiem nawigacji.
defineProps<{
    recipeId: number;
    recipeName: string;
    servings: number;
    done: number;
    total: number;
    /** Ładowanie: ten sam nagłówek, szkielet w miejscu nazwy przepisu, porcji i licznika. */
    loading?: boolean;
}>();
const emit = defineEmits<{ servings: [n: number] }>();

const print = () => window.print();
</script>

<template>
    <header class="head">
        <div class="head__inner">
            <div class="head__bar">
                <RouterLink
                    :to="`/recipes/${recipeId}`"
                    class="head__back"
                >
                    <svg
                        viewBox="0 0 14 14"
                        fill="none"
                        aria-hidden="true"
                    >
                        <path
                            d="M8.5 2.5L4 7l4.5 4.5"
                            stroke="currentColor"
                            stroke-width="1.6"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        />
                    </svg>
                    Przepis
                </RouterLink>
                <button
                    type="button"
                    class="btn btn--secondary btn--sm head__print"
                    @click="print"
                >
                    <svg
                        viewBox="0 0 16 16"
                        fill="none"
                        aria-hidden="true"
                    >
                        <path
                            d="M4 6V2h8v4M4 12H2.5A1 1 0 011.5 11V7a1 1 0 011-1h11a1 1 0 011 1v4a1 1 0 01-1 1H12"
                            stroke="currentColor"
                            stroke-width="1.3"
                            stroke-linejoin="round"
                        />
                        <path
                            d="M4 9.5h8V14H4z"
                            stroke="currentColor"
                            stroke-width="1.3"
                            stroke-linejoin="round"
                        />
                    </svg>
                    Drukuj
                </button>
            </div>
            <div class="head__titles">
                <div class="head__titlebox">
                    <h1 class="head__title">Lista zakupów</h1>
                    <p class="head__recipe">
                        <span
                            v-if="loading"
                            class="sk head__sk-name"
                        />
                        {{ recipeName }}
                        <span class="head__print-only">
                            · {{ servings }} {{ servingsLabel(servings) }}
                        </span>
                    </p>
                </div>
                <span
                    class="head__brand head__print-only"
                    aria-hidden="true"
                >
                    uszatakuchnia.
                </span>
                <ServingsStepper
                    class="head__servings"
                    compact
                    :value="servings"
                    @set="emit('servings', $event)"
                />
            </div>
            <div class="head__progress">
                <div
                    class="head__track"
                    role="progressbar"
                    aria-label="Odhaczone pozycje"
                    :aria-valuemin="0"
                    :aria-valuemax="total"
                    :aria-valuenow="done"
                >
                    <div
                        class="head__fill"
                        :style="{ width: `${total ? (done / total) * 100 : 0}%` }"
                    />
                </div>
                <span class="head__count">
                    <span
                        v-if="loading"
                        class="sk head__sk-count"
                    />
                    <template v-else>{{ done }} z {{ total }}</template>
                </span>
            </div>
        </div>
    </header>
</template>

<style lang="scss" scoped>
@use '@/scss/mixins/media' as m;

/* szkielet nazwy przepisu i licznika (ładowanie) */
.head__sk-name {
    display: inline-block;
    width: 180px;
    height: 12px;
    vertical-align: middle;
}
.head__sk-count {
    display: inline-block;
    width: 34px;
    height: 10px;
    vertical-align: middle;
}

.head {
    background: var(--bg);
    border-bottom: 1px solid var(--rule);
    padding: 18px 0 16px;
}
.head__inner {
    max-width: 640px;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    gap: 14px;
}
.head__bar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
}
.head__back {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    min-height: 40px;
    font-size: 14px;
    color: var(--ink-dim);
    text-decoration: none;
}
.head__back svg {
    width: 14px;
    height: 14px;
}
.head__titles {
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
    gap: 12px;
    flex-wrap: wrap;
}
.head__titlebox {
    min-width: 0;
}
.head__title {
    font-size: 32px;
    font-weight: 500;
    letter-spacing: -1px;
    line-height: 1.1;
}
.head__recipe {
    margin-top: 4px;
    font-size: 14px;
    color: var(--ink-muted);
}
.head__print-only {
    display: none;
}
.head__progress {
    display: flex;
    align-items: center;
    gap: 10px;
}
.head__track {
    flex: 1;
    height: 4px;
    background: var(--rule);
    border-radius: var(--r-pill);
    overflow: hidden;
}
.head__fill {
    height: 100%;
    background: var(--accent);
    transition: width var(--dur);
}
.head__count {
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--ink-muted);
    white-space: nowrap;
}

/* telefon: nagłówek przyklejony pod paskiem nawigacji (57 px), bez Drukuj */
@media (max-width: 900px) {
    .head {
        position: sticky;
        top: 57px;
        z-index: calc(var(--z-sticky) - 1);
        padding: 10px 16px 12px;
    }
    .head__inner {
        max-width: none;
        gap: 10px;
    }
    .head__print {
        display: none;
    }
    .head__title {
        font-size: 26px;
        letter-spacing: -0.8px;
    }
}

@include m.print {
    .head {
        position: static;
        padding: 0 0 10px;
        border-bottom: 1.5pt solid #000;
    }
    .head__inner {
        max-width: none;
    }
    .head__bar,
    .head__servings,
    .head__progress {
        display: none;
    }
    .head__titles {
        align-items: baseline;
    }
    .head__title {
        font-size: 22px;
        font-weight: 600;
        letter-spacing: -0.5px;
    }
    .head__recipe {
        margin-top: 2px;
        font-size: 13px;
        color: #000;
    }
    .head__print-only {
        display: inline;
    }
    .head__brand {
        font-size: 13px;
        font-weight: 600;
    }
}
</style>
