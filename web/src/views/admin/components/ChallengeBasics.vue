<script lang="ts" setup>
import type { IChallengeForm } from '../challengeForm';

import ChallengeDifficultyPicker from './ChallengeDifficultyPicker.vue';

import UFieldError from '@/components/form/UFieldError.vue';
import type { IChallengeFormat } from '@/data/dtos/challenge/ChallengeDto';

// Formularz wyzwania, sekcja 01 Podstawowe: tytuł, format, limit czasu, trudność, źródło.
defineProps<{ formats: IChallengeFormat[]; errors: Record<string, string> }>();
const form = defineModel<IChallengeForm>({ required: true });
</script>

<template>
    <div class="cb">
        <div class="field">
            <label
                for="ch-title"
                class="field-label"
            >
                Tytuł
                <span
                    class="req"
                    aria-hidden="true"
                >
                    *
                </span>
            </label>
            <input
                id="ch-title"
                v-model="form.title"
                class="input"
                :class="{ 'is-invalid': errors.title }"
                maxlength="255"
                required
            />
            <UFieldError v-if="errors.title">{{ errors.title }}</UFieldError>
        </div>
        <div class="cb__row">
            <div class="field">
                <label
                    for="ch-format"
                    class="field-label"
                >
                    Format
                    <span
                        class="req"
                        aria-hidden="true"
                    >
                        *
                    </span>
                </label>
                <select
                    id="ch-format"
                    v-model="form.formatId"
                    class="select"
                    :class="{ 'is-invalid': errors.formatId }"
                    required
                >
                    <option
                        value=""
                        disabled
                    >
                        Wybierz
                    </option>
                    <option
                        v-for="f in formats"
                        :key="f.id"
                        :value="f.id"
                    >
                        {{ f.name }}
                    </option>
                </select>
                <UFieldError v-if="errors.formatId">{{ errors.formatId }}</UFieldError>
                <span class="field-hint">Formaty: Panel → Klasyfikacja → Formaty wyzwań</span>
            </div>
            <div class="field">
                <label
                    for="ch-time"
                    class="field-label"
                >
                    Limit czasu
                </label>
                <span class="cb__time">
                    <input
                        id="ch-time"
                        v-model="form.timeLimit"
                        class="input input--mono cb__time-input"
                        :class="{ 'is-invalid': errors.timeLimitMin }"
                        inputmode="numeric"
                        aria-describedby="ch-time-hint"
                    />
                    <span
                        class="cb__unit"
                        aria-hidden="true"
                    >
                        min
                    </span>
                </span>
                <UFieldError v-if="errors.timeLimitMin">{{ errors.timeLimitMin }}</UFieldError>
                <span
                    id="ch-time-hint"
                    class="field-hint"
                >
                    Puste = bez limitu. Tylko informacja, bez zegara.
                </span>
            </div>
        </div>
        <div class="field">
            <span
                id="ch-diff"
                class="field-label"
            >
                Trudność
                <span
                    class="req"
                    aria-hidden="true"
                >
                    *
                </span>
            </span>
            <ChallengeDifficultyPicker
                v-model="form.difficulty"
                labelledby="ch-diff"
                :invalid="!!errors.difficulty"
            />
            <UFieldError v-if="errors.difficulty">{{ errors.difficulty }}</UFieldError>
        </div>
        <div class="field">
            <label
                for="ch-source"
                class="field-label"
            >
                Źródło
            </label>
            <input
                id="ch-source"
                v-model="form.source"
                class="input"
                :class="{ 'is-invalid': errors.source }"
                placeholder="np. program kulinarny, sezon 3 · wymyślone"
                maxlength="255"
                aria-describedby="ch-source-hint"
            />
            <UFieldError v-if="errors.source">{{ errors.source }}</UFieldError>
            <span
                id="ch-source-hint"
                class="field-hint"
            >
                Krótko, skąd zadanie. Puste = nie pokazujemy.
            </span>
        </div>
    </div>
</template>

<style scoped>
.cb {
    display: flex;
    flex-direction: column;
    gap: 16px;
}
.cb__row {
    display: grid;
    grid-template-columns: minmax(0, 1.3fr) minmax(0, 1fr);
    gap: 16px;
    align-items: start;
}
.cb__time {
    position: relative;
    display: block;
}
.cb__time-input {
    padding-right: 44px;
}
.cb__unit {
    position: absolute;
    right: 10px;
    top: 50%;
    transform: translateY(-50%);
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--ink-muted);
}
</style>
