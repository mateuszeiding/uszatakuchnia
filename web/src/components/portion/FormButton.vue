<script lang="ts" setup>
import FormGlyph from './FormGlyph.vue';

import { formName, formShort, type FormSize } from '@/shared/utils/portion';

// Opens the form chooser (panel control and the phone sticky bar).
defineProps<{ value: FormSize; short?: boolean; open?: boolean; big?: boolean }>();
const emit = defineEmits<{ click: [] }>();
</script>

<template>
    <button
        type="button"
        class="fbtn"
        :class="{ 'fbtn--big': big }"
        aria-haspopup="dialog"
        :aria-expanded="open"
        @click="emit('click')"
    >
        <FormGlyph :shape="value.shape" />
        <span class="sr-only">Forma:</span>
        {{ short ? formShort(value) : formName(value) }}
        <svg
            width="11"
            height="11"
            viewBox="0 0 12 12"
            fill="none"
            aria-hidden="true"
        >
            <path
                d="M3 4.5L6 7.5 9 4.5"
                stroke="currentColor"
                stroke-width="1.6"
                stroke-linecap="round"
                stroke-linejoin="round"
            />
        </svg>
    </button>
</template>

<style scoped>
.fbtn {
    min-height: 36px;
    padding: 0 12px;
    border-radius: var(--r-lg);
    border: 1px solid var(--rule-strong);
    background: var(--bg-alt);
    color: var(--ink);
    font-size: 15px;
    font-weight: 600;
    display: inline-flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    white-space: nowrap;
}
.fbtn:hover {
    border-color: var(--ink-muted);
}
.fbtn--big {
    min-height: 44px;
}
</style>
