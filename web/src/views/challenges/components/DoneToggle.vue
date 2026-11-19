<script lang="ts" setup>
// Okrągły przełącznik „zrobione” na karcie: cel 44 px, kółko 24 px, aria-pressed.
const props = defineProps<{ on: boolean; title: string }>();
defineEmits<{ toggle: [] }>();
</script>

<template>
    <button
        type="button"
        class="done"
        :class="{ 'is-on': props.on }"
        :aria-pressed="props.on"
        :aria-label="`${props.on ? 'Zrobione' : 'Oznacz jako zrobione'}: ${props.title}`"
        :title="props.on ? 'Zrobione — kliknij, żeby cofnąć' : 'Oznacz jako zrobione'"
        @click="$emit('toggle')"
    >
        <span class="done__dot">
            <svg
                width="12"
                height="12"
                viewBox="0 0 12 12"
                fill="none"
                aria-hidden="true"
            >
                <path
                    d="M2.5 6.3L5 8.8l4.6-5.3"
                    stroke="currentColor"
                    stroke-width="1.8"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                />
            </svg>
        </span>
    </button>
</template>

<style scoped>
.done {
    width: 44px;
    height: 44px;
    margin: -8px;
    padding: 0;
    border: none;
    background: transparent;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    border-radius: var(--r-pill);
}
.done__dot {
    width: 24px;
    height: 24px;
    border-radius: var(--r-pill);
    border: 1.5px solid var(--rule-strong);
    color: var(--ink-faint);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    transition:
        background var(--dur-fast),
        border-color var(--dur-fast);
}
.done:hover .done__dot {
    border-color: var(--accent);
}
.done.is-on .done__dot {
    background: var(--accent-fill);
    border-color: var(--accent-fill);
    color: var(--on-accent);
}
</style>
