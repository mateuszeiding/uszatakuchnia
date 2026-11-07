<script lang="ts" setup>
import type { IIngredientForm } from '../ingredientForm';

import UFieldError from '@/components/form/UFieldError.vue';
import type { IngredientDto } from '@/data/dtos/ingredients/IngredientDto';
import type { IDictionaries } from '@/data/dtos/taxonomy/DictionariesDto';

// 01 Podstawowe: nazwa, składnik nadrzędny, typ, alergeny UE.
defineProps<{
    dicts: IDictionaries;
    /** Dozwoleni rodzice (bez samego składnika i jego potomków). */
    parents: IngredientDto[];
    childCount: number;
    errors: Record<string, string>;
}>();
const form = defineModel<IIngredientForm>({ required: true });

function toggleAllergen(code: string) {
    const a = form.value.allergens;
    form.value.allergens = a.includes(code) ? a.filter((x) => x !== code) : [...a, code];
}
</script>

<template>
    <div class="basics">
        <div class="basics__row">
            <div class="field">
                <label
                    for="ing-name"
                    class="field-label"
                >
                    Nazwa
                    <span class="req">*</span>
                </label>
                <input
                    id="ing-name"
                    v-model="form.name"
                    class="input"
                    :class="{ 'is-invalid': errors.name }"
                    :aria-invalid="!!errors.name || undefined"
                    placeholder="np. fusilli z czerwonej soczewicy"
                />
                <UFieldError v-if="errors.name">{{ errors.name }}</UFieldError>
            </div>
            <div class="field">
                <label
                    for="ing-parent"
                    class="field-label"
                >
                    Składnik nadrzędny
                </label>
                <select
                    id="ing-parent"
                    v-model="form.parentId"
                    class="select"
                    :aria-invalid="!!errors.parentId || undefined"
                >
                    <option :value="null">— brak</option>
                    <option
                        v-for="p in parents"
                        :key="p.id"
                        :value="p.id"
                    >
                        {{ p.name }}
                    </option>
                </select>
                <span
                    v-if="childCount"
                    class="field-hint"
                >
                    {{ childCount }} {{ childCount === 1 ? 'podskładnik' : 'podskładniki' }}
                </span>
                <UFieldError v-if="errors.parentId">{{ errors.parentId }}</UFieldError>
            </div>
            <div class="field">
                <label
                    for="ing-type"
                    class="field-label"
                >
                    Typ
                    <span class="req">*</span>
                </label>
                <select
                    id="ing-type"
                    v-model="form.typeCode"
                    class="select"
                    :class="{ 'is-invalid': errors.typeCode }"
                    :aria-invalid="!!errors.typeCode || undefined"
                >
                    <option value="">Wybierz typ</option>
                    <option
                        v-for="t in dicts.ingredientTypes"
                        :key="t.code"
                        :value="t.code"
                    >
                        {{ t.name }}
                    </option>
                </select>
                <UFieldError v-if="errors.typeCode">{{ errors.typeCode }}</UFieldError>
            </div>
        </div>
        <div
            class="field"
            role="group"
            aria-labelledby="ing-allergens"
        >
            <span
                id="ing-allergens"
                class="field-label"
            >
                Alergeny
            </span>
            <div class="basics__chips">
                <button
                    v-for="a in dicts.allergens"
                    :key="a.code"
                    type="button"
                    class="chip"
                    :class="{ 'is-active': form.allergens.includes(a.code) }"
                    :aria-pressed="form.allergens.includes(a.code)"
                    @click="toggleAllergen(a.code)"
                >
                    {{ a.label }}
                </button>
            </div>
            <span class="field-hint">Pokazywane w przepisie. Zostaw puste, jeśli brak.</span>
        </div>
    </div>
</template>

<style scoped>
.basics {
    display: flex;
    flex-direction: column;
    gap: 20px;
}
.basics__row {
    display: grid;
    grid-template-columns: minmax(0, 1.4fr) minmax(0, 1fr) minmax(0, 1fr);
    gap: 16px;
}
.basics__chips {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
}
</style>
