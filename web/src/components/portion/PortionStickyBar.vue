<script lang="ts" setup>
import { computed, ref } from 'vue';

import CountStepper from './CountStepper.vue';
import FormButton from './FormButton.vue';
import FormChooser from './FormChooser.vue';
import PortionDialog from './PortionDialog.vue';
import SizeChips from './SizeChips.vue';

import {
    effectiveMode,
    fmtScale,
    type IPortionRecipe,
    isForm,
    isPieces,
    MAX_PIECES,
    pieceForms,
    piecesLabel,
    portionOf,
    recipeForm,
    type SizeChoice,
} from '@/shared/utils/portion';

// Phone: the short control in the bar that sticks over the ingredient list ("Placki [−] 2 [+] [30 cm ⌄]",
// "Forma ×1,17 [Ø 26 cm ⌄]"); the size button opens a sheet.
// plain: without the bar chrome (shopping list header on a phone).
const props = defineProps<{
    recipe: IPortionRecipe;
    size: SizeChoice | null;
    scale: number;
    plain?: boolean;
}>();
const emit = defineEmits<{ size: [s: SizeChoice] }>();

const mode = computed(() => effectiveMode(props.recipe));
const own = computed(() => recipeForm(portionOf(props.recipe)));
const label = computed(() => piecesLabel(props.recipe));

const btn = ref<{ $el: HTMLElement } | HTMLElement | null>(null);
const anchor = computed(() => {
    const b = btn.value;
    return b instanceof HTMLElement ? b : (b?.$el ?? null);
});
const open = ref(false);
</script>

<template>
    <div
        class="pbar"
        :class="{ 'pbar--plain': plain }"
    >
        <template v-if="mode === 'sized' && isPieces(size)">
            <span class="pbar__label">{{ label }}</span>
            <CountStepper
                :value="size.n"
                :label="label"
                :max="MAX_PIECES"
                big
                @set="emit('size', { n: $event, d: size.d })"
            />
            <button
                ref="btn"
                type="button"
                class="pbar__size"
                aria-haspopup="dialog"
                :aria-expanded="open"
                @click="open = !open"
            >
                <span class="sr-only">Rozmiar:</span>
                {{ size.d }} cm
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
            <PortionDialog
                v-if="open"
                :title="`Rozmiar: ${pieceForms(recipe)[0]}`"
                sheet
                :anchor="anchor"
                @close="open = false"
            >
                <SizeChips
                    :sizes="portionOf(recipe).sizes"
                    :value="size.d"
                    :label="`Rozmiar: ${pieceForms(recipe)[0]}`"
                    big
                    @pick="emit('size', { n: size.n, d: $event })"
                />
            </PortionDialog>
        </template>
        <template v-else-if="mode === 'form' && isForm(size) && own">
            <span class="pbar__label">Forma</span>
            <span class="pbar__form">
                <span class="pbar__scale">{{ fmtScale(scale) }}</span>
                <FormButton
                    ref="btn"
                    :value="size"
                    short
                    big
                    :open="open"
                    @click="open = !open"
                />
            </span>
            <FormChooser
                v-if="open"
                :value="size"
                :recipe="own"
                sheet
                :anchor="anchor"
                @select="emit('size', $event)"
                @close="open = false"
            />
        </template>
    </div>
</template>

<style scoped>
.pbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 10px;
    padding: 8px 16px;
    background: var(--bg-alt);
    border-bottom: 1px solid var(--rule-strong);
    box-shadow: 0 6px 16px rgba(10, 10, 15, 0.08);
}
.pbar--plain {
    padding: 0;
    background: none;
    border: none;
    box-shadow: none;
}
.pbar__label {
    font-size: 14px;
    font-weight: 500;
    color: var(--ink-dim);
}
.pbar__size {
    min-height: 44px;
    padding: 0 12px;
    border-radius: var(--r-lg);
    border: 1px solid var(--rule-strong);
    background: var(--bg-alt);
    color: var(--ink);
    font-size: 15px;
    font-weight: 600;
    display: inline-flex;
    align-items: center;
    gap: 6px;
    cursor: pointer;
}
.pbar__form {
    display: inline-flex;
    align-items: center;
    gap: 8px;
}
.pbar__scale {
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--ink-muted);
}
</style>
