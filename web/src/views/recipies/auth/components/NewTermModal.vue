<script lang="ts" setup>
import { AxiosError } from 'axios';
import { computed, ref, useId } from 'vue';

import UFieldError from '@/components/form/UFieldError.vue';
import UModal from '@/components/UModal.vue';
import { createTerm, type IGlossaryTerm } from '@/data/api/glossary/fetch';

// "+ Nowe pojęcie" from the recipe form: saved to the glossary right away; the recipe form stays open with its
// unsaved changes. The name comes from the selection — usually inflected, so the hint asks for the base form.
const props = defineProps<{ word: string }>();
const emit = defineEmits<{ created: [term: IGlossaryTerm]; close: [] }>();

const MAX = 280;
const uid = useId();
const name = ref(props.word);
const explanation = ref('');
const saving = ref(false);
const errors = ref<{ name?: string; explanation?: string; form?: string }>({});
const len = computed(() => [...explanation.value.trim()].length);

async function save() {
    errors.value = {};
    if (!name.value.trim()) errors.value.name = 'Podaj nazwę pojęcia';
    if (!explanation.value.trim()) errors.value.explanation = 'Napisz wyjaśnienie';
    else if (len.value > MAX) errors.value.explanation = `Najwyżej ${MAX} znaków`;
    if (errors.value.name || errors.value.explanation) return;
    saving.value = true;
    try {
        emit('created', await createTerm({ name: name.value, explanation: explanation.value }));
    } catch (e) {
        const data = (e as AxiosError<{ fields?: { field: string; message: string }[] }>).response
            ?.data;
        for (const f of data?.fields ?? [])
            if (f.field === 'name' || f.field === 'explanation') errors.value[f.field] = f.message;
        if (!data?.fields?.length) errors.value.form = 'Nie udało się zapisać pojęcia.';
    } finally {
        saving.value = false;
    }
}
</script>

<template>
    <UModal
        title="Nowe pojęcie"
        @close="emit('close')"
    >
        <form
            :id="`${uid}-form`"
            class="new-term"
            @submit.prevent="save"
        >
            <div class="new-term__field">
                <label
                    :for="`${uid}-name`"
                    class="field-label"
                >
                    Nazwa
                </label>
                <input
                    :id="`${uid}-name`"
                    v-model="name"
                    class="input"
                    maxlength="60"
                    autofocus
                    :aria-describedby="`${uid}-name-hint`"
                />
                <UFieldError v-if="errors.name">{{ errors.name }}</UFieldError>
                <span
                    :id="`${uid}-name-hint`"
                    class="new-term__hint"
                >
                    Wpisane z zaznaczenia „{{ word }}” — popraw na formę podstawową.
                </span>
            </div>
            <div class="new-term__field">
                <label
                    :for="`${uid}-expl`"
                    class="field-label"
                >
                    Wyjaśnienie
                </label>
                <textarea
                    :id="`${uid}-expl`"
                    v-model="explanation"
                    v-autosize
                    class="textarea"
                    rows="3"
                    placeholder="1–3 zdania"
                />
                <UFieldError v-if="errors.explanation">{{ errors.explanation }}</UFieldError>
                <span
                    class="new-term__count"
                    :class="{ 'is-over': len > MAX }"
                >
                    {{ len }} / {{ MAX }}
                </span>
            </div>
            <UFieldError v-if="errors.form">{{ errors.form }}</UFieldError>
            <p class="new-term__note">
                Zapis od razu trafia do słowniczka; formularz przepisu zostaje otwarty, niezapisane
                zmiany w nim nie giną.
            </p>
        </form>
        <template #footer>
            <button
                type="button"
                class="btn btn--ghost"
                @click="emit('close')"
            >
                Anuluj
            </button>
            <button
                type="submit"
                class="btn btn--accent"
                :form="`${uid}-form`"
                :disabled="saving"
            >
                {{ saving ? 'Zapisywanie…' : 'Dodaj i oznacz' }}
            </button>
        </template>
    </UModal>
</template>

<style scoped>
.new-term {
    display: flex;
    flex-direction: column;
    gap: 14px;
}
.new-term__field {
    display: flex;
    flex-direction: column;
    gap: 6px;
}
.new-term__hint,
.new-term__note {
    margin: 0;
    font-size: 12px;
    color: var(--ink-muted);
}
.new-term__count {
    align-self: flex-end;
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--ink-muted);
}
.new-term__count.is-over {
    color: var(--danger);
}
</style>
