<script lang="ts" setup>
import { ref } from 'vue';

// One line instead of the recipe name + counter: the name ate a line of the step on a phone.
defineProps<{
    /** Step number from 1; 0 = the "Zanim zaczniesz" screen. */
    n: number;
    total: number;
}>();
const emit = defineEmits<{ ingredients: []; close: [] }>();

const ingBtn = ref<HTMLButtonElement>();
defineExpose({ focusIngredients: () => ingBtn.value?.focus() });
</script>

<template>
    <div class="cook-bar">
        <div class="cook-bar__count">
            <template v-if="n">
                Krok {{ n }}
                <span class="cook-bar__of">z {{ total }}</span>
            </template>
            <template v-else>Zanim zaczniesz</template>
        </div>
        <button
            ref="ingBtn"
            type="button"
            class="cook-bar__pill"
            aria-haspopup="dialog"
            @click="emit('ingredients')"
        >
            Składniki
        </button>
        <button
            type="button"
            class="cook-bar__round"
            aria-label="Zamknij tryb gotowania"
            @click="emit('close')"
        >
            <svg
                width="11"
                height="11"
                viewBox="0 0 11 11"
                aria-hidden="true"
            >
                <path
                    d="M2 2l7 7M9 2L2 9"
                    stroke="currentColor"
                    stroke-width="1.5"
                    stroke-linecap="round"
                />
            </svg>
        </button>
    </div>
    <div
        class="cook-bar__progress"
        aria-hidden="true"
    >
        <div
            class="cook-bar__fill"
            :style="{ width: `${(n / total) * 100}%` }"
        />
    </div>
</template>

<style scoped>
.cook-bar {
    grid-area: bar;
    height: 56px;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 0 8px 0 20px;
}
.cook-bar__count {
    flex: 1;
    min-width: 0;
    font-size: 17px;
    font-weight: 600;
    white-space: nowrap;
}
.cook-bar__of {
    color: var(--ink-muted);
    font-weight: 500;
}
.cook-bar__pill {
    height: 44px;
    padding: 0 14px;
    border-radius: var(--r-pill);
    border: 1px solid var(--rule-strong);
    background: var(--bg-alt);
    color: var(--ink);
    font-size: 15px;
    font-weight: 600;
    cursor: pointer;
}
.cook-bar__round {
    width: 44px;
    height: 44px;
    flex-shrink: 0;
    border-radius: var(--r-pill);
    border: 1px solid var(--rule-strong);
    background: var(--bg-alt);
    color: var(--ink);
    cursor: pointer;
    display: inline-flex;
    align-items: center;
    justify-content: center;
}
.cook-bar__progress {
    grid-area: progress;
    height: 3px;
    flex-shrink: 0;
    background: var(--rule);
}
.cook-bar__fill {
    height: 100%;
    background: var(--accent-fill);
    transition: width 0.25s;
}

@media (orientation: landscape) and (max-height: 500px) {
    .cook-bar {
        height: 48px;
    }
    .cook-bar__count {
        font-size: 16px;
    }
    .cook-bar__pill {
        height: 40px;
        padding: 0 12px;
        font-size: 14px;
    }
    .cook-bar__round {
        width: 40px;
        height: 40px;
    }
}
</style>
