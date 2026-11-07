<script lang="ts" setup>
import { computed, ref } from 'vue';

import UFieldError from './form/UFieldError.vue';
import UModal from './UModal.vue';

import { createIngredient } from '@/data/api/ingredients/fetch';
import type { IngredientDto } from '@/data/dtos/ingredients/IngredientDto';
import type { IDictionaries } from '@/data/dtos/taxonomy/DictionariesDto';

// Szybkie dodanie składnika z formularza przepisu: nazwa, typ, alergeny UE.
// Edycja (z makro i przelicznikami) jest w panelu — AdminIngredientEdit.
const props = defineProps<{
    dicts: IDictionaries;
    initialName?: string;
    /** Tekst przycisku dla nowego składnika (z formularza: „Dodaj i wstaw”). */
    createLabel?: string;
}>();
const emit = defineEmits<{ saved: [ing: IngredientDto]; close: [] }>();

const name = ref(props.initialName ?? '');
const typeCode = ref('');
const allergens = ref<string[]>([]);
const saving = ref(false);
const error = ref('');

const valid = computed(() => !!name.value.trim() && !!typeCode.value);

function toggle(code: string) {
    allergens.value = allergens.value.includes(code)
        ? allergens.value.filter((a) => a !== code)
        : [...allergens.value, code];
}

async function save() {
    if (!valid.value || saving.value) return;
    saving.value = true;
    error.value = '';
    try {
        const ing = await createIngredient({
            name: name.value.trim(),
            typeCode: typeCode.value,
            allergens: allergens.value,
        });
        emit('saved', ing);
    } catch (e) {
        const status = (e as { response?: { status?: number } }).response?.status;
        error.value =
            status === 409 ? 'Taki składnik już jest w bazie.' : 'Nie udało się zapisać składnika.';
    } finally {
        saving.value = false;
    }
}
</script>

<template>
    <UModal
        eyebrow="Baza składników"
        title="Nowy składnik"
        :width="480"
        @close="emit('close')"
    >
        <p
            v-if="createLabel"
            class="modal__text"
        >
            Trafi do bazy i od razu do tego wiersza.
        </p>
        <div class="field">
            <label
                for="ing-modal-name"
                class="field-label"
            >
                Nazwa
                <span class="req">*</span>
            </label>
            <input
                id="ing-modal-name"
                v-model="name"
                class="input"
                autofocus
                placeholder="np. masło klarowane"
                @keydown.enter.prevent="save"
            />
        </div>
        <div class="field">
            <label
                for="ing-modal-type"
                class="field-label"
            >
                Typ
                <span class="req">*</span>
            </label>
            <select
                id="ing-modal-type"
                v-model="typeCode"
                class="select"
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
        </div>
        <div
            class="field"
            role="group"
            aria-labelledby="ing-modal-allergens"
        >
            <span
                id="ing-modal-allergens"
                class="field-label"
            >
                Alergeny
            </span>
            <div class="ing-modal__chips">
                <button
                    v-for="a in dicts.allergens"
                    :key="a.code"
                    type="button"
                    class="chip"
                    :class="{ 'is-active': allergens.includes(a.code) }"
                    :aria-pressed="allergens.includes(a.code)"
                    @click="toggle(a.code)"
                >
                    {{ a.label }}
                </button>
            </div>
            <span class="field-hint">Pokazywane w przepisie. Zostaw puste, jeśli brak.</span>
        </div>
        <UFieldError v-if="error">{{ error }}</UFieldError>
        <template #footer>
            <button
                type="button"
                class="btn btn--ghost"
                @click="emit('close')"
            >
                Anuluj
            </button>
            <button
                type="button"
                class="btn btn--accent"
                :disabled="!valid || saving"
                @click="save"
            >
                {{ saving ? 'Zapisywanie…' : (createLabel ?? 'Dodaj składnik') }}
            </button>
        </template>
    </UModal>
</template>

<style scoped>
.ing-modal__chips {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
}
</style>
