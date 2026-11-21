<script lang="ts" setup>
// Old recipe hint (recipes.tip) above the steps: text + "Przenieś do kroku…".
// Choosing a step appends it to that step's hints and clears tip (in the database after "Zapisz").
import { computed } from 'vue';

import { type IRecipeForm, moveTipToStep } from '../recipeForm';

import BulbIcon from '@/components/BulbIcon.vue';

const form = defineModel<IRecipeForm>('form', { required: true });

const options = computed(() =>
    form.value.stepGroups
        .flatMap((g) => g.items)
        .map((s, i) => ({ s, no: i + 1 }))
        .filter(({ s }) => s.text.trim())
        .map(({ s, no }) => {
            const t = s.text.trim();
            return { id: s.id, label: `Krok ${no} — ${t.length > 40 ? `${t.slice(0, 40)}…` : t}` };
        })
);

function onPick(e: Event) {
    const id = Number((e.target as HTMLSelectElement).value);
    if (id) moveTipToStep(form.value, id);
}
</script>

<template>
    <div
        v-if="form.tip.trim()"
        class="tip-move"
    >
        <div class="field-label">Stara wskazówka</div>
        <p class="tip-move__text">
            <span class="tip-move__icon"><BulbIcon /></span>
            <span>{{ form.tip.trim() }}</span>
        </p>
        <label
            class="tip-move__pick"
            for="tip-move"
        >
            <span class="tip-move__hint">Wskazówki są teraz przy krokach.</span>
            <select
                id="tip-move"
                class="select tip-move__select"
                :disabled="!options.length"
                @change="onPick"
            >
                <option value="">Przenieś do kroku…</option>
                <option
                    v-for="o in options"
                    :key="o.id"
                    :value="o.id"
                >
                    {{ o.label }}
                </option>
            </select>
        </label>
    </div>
</template>

<style scoped>
.tip-move {
    margin-bottom: 16px;
    padding: 14px 16px;
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    background: var(--bg);
    display: flex;
    flex-direction: column;
    gap: 10px;
}
.tip-move__text {
    display: grid;
    grid-template-columns: 16px minmax(0, 1fr);
    gap: 8px;
    align-items: start;
    font-size: 14px;
    line-height: 1.5;
    color: var(--ink-dim);
}
.tip-move__icon {
    display: inline-flex;
    padding-top: 2px;
    color: var(--accent);
}
.tip-move__pick {
    display: flex;
    align-items: center;
    gap: 12px;
    flex-wrap: wrap;
}
.tip-move__hint {
    font-size: 12px;
    color: var(--ink-muted);
}
.tip-move__select {
    width: auto;
    min-width: 260px;
    font-size: 13px;
}
</style>
