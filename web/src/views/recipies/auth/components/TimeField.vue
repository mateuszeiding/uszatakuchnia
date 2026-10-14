<script lang="ts" setup>
import UFieldError from '@/components/form/UFieldError.vue';
import { fromMinutes, type TimeUnit, toMinutes } from '@/shared/utils/time';

// Czas przepisu: „od” (wymagane), „do” (opcjonalnie, zakres) + jednostka min/h.
// Przełączenie jednostki przelicza wpisane wartości; zapis w minutach robi toRequest().
const time = defineModel<string>('time', { required: true });
const timeMax = defineModel<string>('timeMax', { required: true });
const unit = defineModel<TimeUnit>('unit', { required: true });
defineProps<{ error?: string }>();

const UNITS: TimeUnit[] = ['min', 'h'];

function setUnit(next: TimeUnit) {
    if (next === unit.value) return;
    // niepoprawny tekst zostaje bez zmian — walidacja go pokaże
    const conv = (v: string) => {
        const m = toMinutes(v, unit.value);
        return m ? fromMinutes(m, next) : v;
    };
    time.value = conv(time.value);
    timeMax.value = conv(timeMax.value);
    unit.value = next;
}
</script>

<template>
    <div
        id="f-time"
        class="field"
    >
        <label
            for="in-time"
            class="field-label"
        >
            Czas
            <span class="req">*</span>
        </label>
        <div class="time-row">
            <input
                id="in-time"
                v-model="time"
                class="input input--mono"
                :class="{ 'is-invalid': error }"
                :aria-invalid="!!error || undefined"
                inputmode="decimal"
                :placeholder="unit === 'h' ? '2,5' : '40'"
            />
            <span
                class="time-row__dash"
                aria-hidden="true"
            >
                –
            </span>
            <input
                v-model="timeMax"
                class="input input--mono"
                :class="{ 'is-invalid': error }"
                inputmode="decimal"
                aria-label="Czas maksymalny (opcjonalnie)"
                placeholder="do"
            />
            <div
                class="time-unit"
                role="radiogroup"
                aria-label="Jednostka czasu"
            >
                <button
                    v-for="u in UNITS"
                    :key="u"
                    type="button"
                    role="radio"
                    class="time-unit__btn"
                    :class="{ 'is-active': unit === u }"
                    :aria-checked="unit === u"
                    @click="setUnit(u)"
                >
                    {{ u }}
                </button>
            </div>
        </div>
        <UFieldError v-if="error">{{ error }}</UFieldError>
        <span
            v-else
            class="field-hint"
        >
            Drugie pole opcjonalnie, dla zakresu (np. 2,5–4,5 h).
        </span>
    </div>
</template>

<style scoped>
.time-row {
    display: grid;
    grid-template-columns: minmax(0, 1fr) auto minmax(0, 1fr) auto;
    gap: 6px;
    align-items: center;
}
.time-row__dash {
    color: var(--ink-muted);
}
.time-unit {
    display: flex;
    border: 1px solid var(--rule-strong);
    border-radius: var(--r-md);
    overflow: hidden;
}
.time-unit__btn {
    min-width: 40px;
    min-height: 38px;
    padding: 0 10px;
    border: none;
    background: transparent;
    color: var(--ink);
    font-size: 13px;
    font-weight: 500;
    cursor: pointer;
}
.time-unit__btn.is-active {
    background: var(--ink);
    color: var(--bg);
}
</style>
