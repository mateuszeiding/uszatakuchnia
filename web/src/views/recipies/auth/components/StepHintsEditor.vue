<script lang="ts" setup>
// ✕ deletes without confirmation: nothing is saved until "Zapisz".
import { MAX_HINT_LEN } from '../recipeForm';
import { HINT_LIMIT, HINTS_LIMIT } from '../stepRefsForm';

import BulbIcon from '@/components/BulbIcon.vue';

defineProps<{ stepId: number; stepNo: number }>();
const hints = defineModel<string[]>('hints', { required: true });

// Cooking mode: a warning only, saving and publishing stay possible.
const tooLong = (hi: number) => [...(hints.value[hi] ?? '').trim()].length > HINT_LIMIT;
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
            <label
                :for="`hint-${stepId}-${hi}`"
                class="sr-only"
            >
                Krok {{ stepNo }}, wskazówka {{ hi + 1 }}
            </label>
            <textarea
                :id="`hint-${stepId}-${hi}`"
                v-model="hints[hi]"
                v-autosize
                class="textarea hint-row__input"
                :class="{ 'is-warn': warnAt(hi) }"
                rows="1"
                :maxlength="MAX_HINT_LEN"
                :aria-describedby="warnAt(hi) ? `hint-warn-${stepId}-${hi}` : undefined"
                placeholder="np. Jeśli ciasto się klei, odstaw je na 10 minut"
            />
            <p
                v-if="warnAt(hi)"
                :id="`hint-warn-${stepId}-${hi}`"
                class="hint-row__warn"
            >
                <template v-if="tooLong(hi)">
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
.hint-row__input.is-warn {
    border-color: var(--prep-line);
}
.hint-row__warn {
    margin: 0;
    font-size: 12px;
    line-height: 1.45;
    color: var(--prep);
}
.hint-row__input {
    min-height: 40px;
    font-size: 13px;
    line-height: 18px;
}
</style>
