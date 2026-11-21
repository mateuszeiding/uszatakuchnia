<script lang="ts" setup>
import { computed } from 'vue';

import CautionsBlock, { type ICautionPrep } from '@/components/CautionsBlock.vue';
import { pluralPL } from '@/shared/utils/quantity';

// Notes and "prepare ahead" read before cooking starts, so step 1 gets the whole screen.
const props = defineProps<{
    name: string;
    steps: number;
    cautions: string[];
    prep: ICautionPrep[];
}>();
const emit = defineEmits<{ step: [no: number] }>();

const count = computed(() => `${props.steps} ${pluralPL(props.steps, 'krok', 'kroki', 'kroków')}`);
</script>

<template>
    <div class="cook-before">
        <div class="cook-before__meta">{{ name }} · {{ count }}</div>
        <CautionsBlock
            :items="cautions"
            :prep="prep"
            variant="large"
            @step="emit('step', $event)"
        />
        <p class="cook-before__note">„Zaczynam” niczego nie odhacza.</p>
    </div>
</template>

<style scoped>
.cook-before__meta {
    font-size: 15px;
    color: var(--ink-muted);
    margin-bottom: 14px;
}
.cook-before__note {
    font-size: 15px;
    line-height: 1.45;
    color: var(--ink-muted);
    margin: 16px 0 0;
}
</style>
