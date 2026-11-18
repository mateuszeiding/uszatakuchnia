<script lang="ts" setup>
// Wiersze wskazówek kroku w formularzu: ikona, pole, ✕ (bez potwierdzenia — nic nie jest zapisane do „Zapisz”).
import { MAX_HINT_LEN } from '../recipeForm';

import BulbIcon from '@/components/BulbIcon.vue';

defineProps<{ stepId: number; stepNo: number }>();
const hints = defineModel<string[]>('hints', { required: true });
</script>

<template>
    <div
        v-for="(_, hi) in hints"
        :key="hi"
        class="hint-row"
    >
        <span class="hint-row__icon"><BulbIcon /></span>
        <textarea
            :id="`hint-${stepId}-${hi}`"
            v-model="hints[hi]"
            v-autosize
            class="textarea hint-row__input"
            rows="1"
            :maxlength="MAX_HINT_LEN"
            :aria-label="`Krok ${stepNo}, wskazówka ${hi + 1}`"
            placeholder="np. Jeśli ciasto się klei, odstaw je na 10 minut"
        />
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
    /* pole rośnie w dół — ikona i ✕ zostają przy pierwszej linii */
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
/* jedna linia jak .input (40 px), dalej rośnie z treścią */
.hint-row__input {
    min-height: 40px;
    font-size: 13px;
    line-height: 18px;
}
</style>
