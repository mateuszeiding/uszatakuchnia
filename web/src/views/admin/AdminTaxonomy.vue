<script lang="ts" setup>
import { useQueryClient } from '@tanstack/vue-query';
import { computed, onMounted, ref } from 'vue';

import ChallengeFormatsTable from './components/ChallengeFormatsTable.vue';
import EquipmentTable from './components/EquipmentTable.vue';
import UnitsTable from './components/UnitsTable.vue';
import ValuesTable from './components/ValuesTable.vue';

import USpinner from '@/components/USpinner.vue';
import { invalidateChallenges } from '@/data/api/challenges/fetch';
import { invalidateRecipes } from '@/data/api/recipe/fetch';
import { fetchDictionariesWithClient, invalidateDictionaries } from '@/data/api/taxonomy/fetch';
import type { IDictionaries } from '@/data/dtos/taxonomy/DictionariesDto';

const qc = useQueryClient();
const dicts = ref<IDictionaries>();
const tab = ref<string>('dish');

const DESCRIPTIONS: Record<string, string> = {
    dish: 'Wymagane w przepisie. Kolor pojawia się na kartach.',
    base: 'Główne składniki, po których się szuka.',
    cuisine: 'Opcjonalna, jedna na przepis.',
    diet: 'Filtr działa jako „spełnia wszystkie”.',
    practical: 'Lekkie tagi pomocnicze.',
    units: 'Formy odmiany: 1 / 2–4 / 5+. Puste dla skrótów (g, ml).',
};

async function load() {
    dicts.value = await fetchDictionariesWithClient(qc);
}
async function reload() {
    // labels / colors / moved values also show on recipe list cards (60 s cache)
    invalidateRecipes(qc);
    // format names also show on challenge cards
    invalidateChallenges(qc);
    await invalidateDictionaries(qc);
    await load();
}
onMounted(load);

const tabs = computed(() => [
    ...(dicts.value?.axes ?? []).map((a) => ({ id: a.code, label: a.label, n: a.values.length })),
    { id: 'units', label: 'Jednostki', n: dicts.value?.units.length ?? 0 },
    { id: 'equipment', label: 'Sprzęt', n: dicts.value?.equipment.length ?? 0 },
    { id: 'formats', label: 'Formaty wyzwań', n: dicts.value?.challengeFormats.length ?? 0 },
]);
const axis = computed(() => dicts.value?.axes.find((a) => a.code === tab.value));
</script>

<template>
    <h1 class="admin-title">Klasyfikacja</h1>
    <template v-if="!dicts">
        <div
            class="tabs tax-tabs"
            aria-hidden="true"
        >
            <span
                v-for="n in 8"
                :key="n"
                class="tab"
                :class="{ 'is-active': n === 1 }"
            >
                <span class="sk tax-sk__label" />
                <span class="count"><span class="tax-sk__count">00</span></span>
            </span>
        </div>
        <p class="tax-desc">{{ DESCRIPTIONS.dish }}</p>
        <div class="tax-wrap">
            <div
                class="tax-loading"
                role="status"
                aria-label="Wczytywanie klasyfikacji"
            >
                <USpinner :size="22" />
            </div>
            <ValuesTable :axis="null" />
        </div>
    </template>
    <template v-else>
        <div
            class="tabs tax-tabs"
            role="tablist"
        >
            <button
                v-for="t in tabs"
                :key="t.id"
                type="button"
                role="tab"
                class="tab"
                :class="{ 'is-active': tab === t.id }"
                :aria-selected="tab === t.id"
                @click="tab = t.id"
            >
                {{ t.label }}
                <span class="count">{{ t.n }}</span>
            </button>
        </div>
        <p
            v-if="DESCRIPTIONS[tab]"
            class="tax-desc"
        >
            {{ DESCRIPTIONS[tab] }}
        </p>
        <UnitsTable
            v-if="tab === 'units'"
            :units="dicts.units"
            @changed="reload"
        />
        <EquipmentTable
            v-else-if="tab === 'equipment'"
            :items="dicts.equipment"
            @changed="reload"
        />
        <ChallengeFormatsTable
            v-else-if="tab === 'formats'"
            :items="dicts.challengeFormats"
            @changed="reload"
        />
        <ValuesTable
            v-else-if="axis"
            :key="axis.id"
            :axis="axis"
            @changed="reload"
        />
    </template>
</template>

<style scoped>
.admin-title {
    font-size: 28px;
    font-weight: 600;
    letter-spacing: -0.8px;
    margin-bottom: 22px;
}
.tax-wrap {
    position: relative;
}
.tax-loading {
    position: absolute;
    left: 0;
    right: 0;
    top: 56px;
    z-index: 1;
    display: flex;
    justify-content: center;
    pointer-events: none;
}
.tax-sk__label {
    display: inline-block;
    width: 64px;
    height: 12px;
    vertical-align: middle;
}
.tax-sk__count {
    visibility: hidden;
}
.tax-tabs {
    overflow-x: auto;
}
/* 8 tabs (with "Formaty wyzwań") don't fit at ~1120 px: labels stay on one line and the bar scrolls */
.tax-tabs .tab {
    white-space: nowrap;
}
.tax-desc {
    margin: 16px 0;
    font-size: 13px;
    color: var(--ink-muted);
}
</style>
