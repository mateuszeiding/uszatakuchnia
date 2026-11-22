<script lang="ts" setup>
import { computed, useId } from 'vue';

import CountStepper from './CountStepper.vue';
import SizeChips from './SizeChips.vue';

import {
    type IPortionRecipe,
    MAX_PIECES,
    type PieceChoice,
    pieceForms,
    piecesLabel,
    portionOf,
} from '@/shared/utils/portion';

// "Sztuki w rozmiarach": how many pieces (−/+) and which of the author's sizes.
const props = defineProps<{ recipe: IPortionRecipe; value: PieceChoice; big?: boolean }>();
const emit = defineEmits<{ set: [s: PieceChoice] }>();

const id = useId();
const sizes = computed(() => portionOf(props.recipe).sizes);
const dough = computed(() => sizes.value.find((s) => s.diameterCm === props.value.d)?.doughG);
const label = computed(() => piecesLabel(props.recipe));
</script>

<template>
    <div class="pctl">
        <div class="pctl__row">
            <span
                :id="id"
                class="pctl__label"
            >
                {{ label }}
            </span>
            <CountStepper
                :value="value.n"
                :label="label"
                :max="MAX_PIECES"
                :big="big"
                @set="emit('set', { n: $event, d: value.d })"
            />
        </div>
        <div>
            <div class="pctl__cap">
                Średnica
                <template v-if="dough">
                    · ciasto na {{ pieceForms(recipe)[0] }} ok. {{ dough }} g
                </template>
            </div>
            <SizeChips
                :sizes="sizes"
                :value="value.d"
                :label="`Rozmiar: ${pieceForms(recipe)[0]}`"
                :big="big"
                @pick="emit('set', { n: value.n, d: $event })"
            />
        </div>
    </div>
</template>

<style scoped>
.pctl {
    background: var(--bg);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    padding: 8px 8px 12px 14px;
    display: flex;
    flex-direction: column;
    gap: 10px;
}
.pctl__row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 10px;
}
.pctl__label {
    font-size: 14px;
    font-weight: 500;
    color: var(--ink-dim);
}
.pctl__cap {
    margin-bottom: 6px;
    font-size: 12px;
    color: var(--ink-muted);
}
</style>
