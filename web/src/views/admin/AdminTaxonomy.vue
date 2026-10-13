<script lang="ts" setup>
import { useQueryClient } from '@tanstack/vue-query';
import { computed, onMounted, ref } from 'vue';

import UnitsTable from './components/UnitsTable.vue';
import ValuesTable from './components/ValuesTable.vue';

import USpinner from '@/components/USpinner.vue';
import { fetchDictionariesWithClient, invalidateDictionaries } from '@/data/api/taxonomy/fetch';
import type { IDictionaries } from '@/data/dtos/taxonomy/DictionariesDto';

// Klasyfikacja: osie (Typ dania, Baza, Kuchnia, Diety, Praktyczne, …) + Jednostki.
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
    await invalidateDictionaries(qc);
    await load();
}
onMounted(load);

const tabs = computed(() => [
    ...(dicts.value?.axes ?? []).map((a) => ({ id: a.code, label: a.label, n: a.values.length })),
    { id: 'units', label: 'Jednostki', n: dicts.value?.units.length ?? 0 },
]);
const axis = computed(() => dicts.value?.axes.find((a) => a.code === tab.value));
</script>

<template>
    <h1 class="admin-title">Klasyfikacja</h1>
    <div
        v-if="!dicts"
        class="admin-loading"
    >
        <USpinner :size="22" />
    </div>
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
.admin-loading {
    display: flex;
    justify-content: center;
    padding: 60px 0;
}
.tax-tabs {
    overflow-x: auto;
}
.tax-desc {
    margin: 16px 0;
    font-size: 13px;
    color: var(--ink-muted);
}
</style>
