<script lang="ts" setup>
// ✕ deletes without confirmation: nothing is saved until "Zapisz".
// A hint uses the step's editor without "@": it can mark glossary terms, but has no quantities.
import { MAX_HINT_LEN } from '../recipeForm';
import { HINT_LIMIT, HINTS_LIMIT } from '../stepRefsForm';

import StepTextEditor from './StepTextEditor.vue';

import BulbIcon from '@/components/BulbIcon.vue';
import { type ITerm, termPlainText } from '@/shared/utils/stepRefs';

defineProps<{ stepId: number; stepNo: number; terms: ITerm[] }>();
const emit = defineEmits<{ 'term-created': [term: ITerm] }>();
const hints = defineModel<string[]>('hints', { required: true });

// Length as the reader sees it: a marked term counts as its word.
const len = (hi: number) => [...termPlainText(hints.value[hi] ?? '').trim()].length;
// Over MAX_HINT_LEN the API refuses to save; over HINT_LIMIT / HINTS_LIMIT it's a cooking-mode warning only.
const tooLongToSave = (hi: number) => len(hi) > MAX_HINT_LEN;
const tooLong = (hi: number) => len(hi) > HINT_LIMIT;
const tooMany = (hi: number) => hi >= HINTS_LIMIT;
const warnAt = (hi: number) => tooLong(hi) || tooMany(hi);
</script>

<template>
    <div
        v-for="(_, hi) in hints"
        :key="hi"
        class="hint-row"
    >
        <span class="hint-row__icon"><BulbIcon /></span>
        <div class="hint-row__field">
            <StepTextEditor
                :field-id="`hint-${stepId}-${hi}`"
                v-model="hints[hi]!"
                :rows="[]"
                :terms="terms"
                :ing-refs="false"
                compact
                :label="`Krok ${stepNo}, wskazówka ${hi + 1}`"
                :warn="warnAt(hi)"
                :broken="tooLongToSave(hi)"
                :describedby="warnAt(hi) ? `hint-warn-${stepId}-${hi}` : undefined"
                placeholder="np. Jeśli ciasto się klei, odstaw je na 10 minut"
                @term-created="emit('term-created', $event)"
            />
            <p
                v-if="warnAt(hi)"
                :id="`hint-warn-${stepId}-${hi}`"
                class="hint-row__warn"
                :class="{ 'is-error': tooLongToSave(hi) }"
            >
                <template v-if="tooLongToSave(hi)">
                    Wskazówka może mieć najwyżej {{ MAX_HINT_LEN }} znaków (teraz {{ len(hi) }}).
                </template>
                <template v-else-if="tooLong(hi)">
                    Za długa do trybu gotowania — na telefonie nie zmieści się bez przewijania.
                    Skróć ją do ok. {{ HINT_LIMIT }} znaków.
                </template>
                <template v-else>
                    W trybie gotowania mieszczą się najwyżej {{ HINTS_LIMIT }} wskazówki — połącz je
                    albo usuń którąś.
                </template>
            </p>
        </div>
        <button
            type="button"
            class="row-icon-btn"
            :aria-label="`Usuń wskazówkę ${hi + 1} z kroku ${stepNo}`"
            :title="`Usuń wskazówkę ${hi + 1}`"
            @click="hints.splice(hi, 1)"
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
</template>

<style scoped>
.hint-row {
    display: grid;
    grid-template-columns: 16px minmax(0, 1fr) 30px;
    gap: 8px;
    /* the field grows downward; the icon and ✕ stay by the first line */
    align-items: start;
}
.hint-row__icon {
    display: inline-flex;
    margin-top: 12px;
    color: var(--accent);
}
.hint-row .row-icon-btn {
    margin-top: 5px;
}
.hint-row__field {
    display: flex;
    flex-direction: column;
    gap: 4px;
}
.hint-row__warn {
    margin: 0;
    font-size: 12px;
    line-height: 1.45;
    color: var(--prep);
}
.hint-row__warn.is-error {
    color: var(--danger);
}
</style>
