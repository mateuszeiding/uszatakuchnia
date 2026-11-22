<script lang="ts" setup>
import { computed, ref } from 'vue';

import FormButton from './FormButton.vue';
import FormChooser from './FormChooser.vue';

import { useNarrow } from '@/composables/useNarrow';
import type { FormSize } from '@/shared/utils/portion';

// "Forma: tortownica Ø 24 cm ⌄" → the form chooser (sheet on a phone, popover on desktop).
defineProps<{ value: FormSize; recipe: FormSize }>();
const emit = defineEmits<{ set: [f: FormSize] }>();

const narrow = useNarrow();
const btn = ref<{ $el: HTMLElement } | null>(null);
const anchor = computed(() => btn.value?.$el ?? null);
const open = ref(false);
</script>

<template>
    <div class="fctl">
        <span class="fctl__label">Forma</span>
        <FormButton
            ref="btn"
            :value="value"
            :short="narrow"
            :open="open"
            :big="narrow"
            @click="open = !open"
        />
        <FormChooser
            v-if="open"
            :value="value"
            :recipe="recipe"
            :sheet="narrow"
            :anchor="anchor"
            @select="emit('set', $event)"
            @close="open = false"
        />
    </div>
</template>

<style scoped>
.fctl {
    background: var(--bg);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    padding: 8px 8px 8px 14px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 10px;
}
.fctl__label {
    font-size: 14px;
    font-weight: 500;
    color: var(--ink-dim);
}
</style>
