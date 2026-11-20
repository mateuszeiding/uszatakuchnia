<script lang="ts" setup>
import { useId } from 'vue';

// Licznik porcji (1–20). `compact` — wersja w nagłówku listy zakupów (bez ramki, mniejsze przyciski);
// `big` — przyciski 44 px (arkusz składników w trybie gotowania).
defineProps<{ value: number; compact?: boolean; big?: boolean }>();
const emit = defineEmits<{ set: [n: number] }>();

const id = useId();
</script>

<template>
    <div
        class="servings"
        :class="{ 'servings--compact': compact, 'servings--big': big }"
    >
        <span
            :id="id"
            class="servings__label"
        >
            Porcje
        </span>
        <div
            class="servings__ctrl"
            role="group"
            :aria-labelledby="id"
        >
            <button
                type="button"
                class="servings__btn"
                aria-label="Mniej porcji"
                :disabled="value <= 1"
                @click="emit('set', value - 1)"
            >
                −
            </button>
            <output
                class="servings__value"
                aria-live="polite"
            >
                <span
                    v-if="!value"
                    class="sk servings__sk"
                />
                <template v-else>{{ value }}</template>
            </output>
            <button
                type="button"
                class="servings__btn"
                aria-label="Więcej porcji"
                :disabled="value >= 20"
                @click="emit('set', value + 1)"
            >
                +
            </button>
        </div>
    </div>
</template>

<style scoped>
.servings {
    background: var(--bg);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    padding: 8px 8px 8px 14px;
    display: flex;
    align-items: center;
    justify-content: space-between;
}
.servings__label {
    font-size: 14px;
    font-weight: 500;
    color: var(--ink-dim);
}
.servings__ctrl {
    display: flex;
    align-items: center;
    gap: 4px;
}
.servings__btn {
    width: 36px;
    height: 36px;
    border-radius: var(--r-md);
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    color: var(--ink);
    font-size: 18px;
    cursor: pointer;
}
.servings__btn:disabled {
    background: transparent;
    color: var(--ink-faint);
    cursor: default;
}
.servings__sk {
    display: inline-block;
    width: 18px;
    height: 18px;
    vertical-align: middle;
}
.servings__value {
    min-width: 40px;
    text-align: center;
    font-size: 22px;
    font-weight: 600;
}

.servings--compact {
    background: none;
    border: none;
    padding: 0;
    gap: 8px;
}
.servings--compact .servings__label {
    font-size: 13px;
    font-weight: 400;
}
.servings--compact .servings__btn {
    width: 32px;
    height: 32px;
}
.servings--compact .servings__value {
    min-width: 30px;
    font-size: 18px;
}
.servings--big .servings__btn {
    width: 44px;
    height: 44px;
}
</style>
