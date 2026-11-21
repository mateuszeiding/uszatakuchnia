<script lang="ts" setup>
import { AxiosError } from 'axios';
import { computed, ref, useId } from 'vue';

import UFieldError from '@/components/form/UFieldError.vue';
import TermExplainCard from '@/components/TermExplainCard.vue';
import UModal from '@/components/UModal.vue';
import {
    createTerm,
    type IGlossaryTerm,
    type IUpsertTerm,
    updateTerm,
} from '@/data/api/glossary/fetch';

// Adding / editing a term, with a live preview of the explanation popover.
const props = defineProps<{ term: IGlossaryTerm | null }>();
const emit = defineEmits<{ saved: [term: IGlossaryTerm]; close: [] }>();

const MAX = 280;
const uid = useId();
const name = ref(props.term?.name ?? '');
const explanation = ref(props.term?.explanation ?? '');
const saving = ref(false);
const errors = ref<{ name?: string; explanation?: string; form?: string }>({});
const len = computed(() => [...explanation.value.trim()].length);
const usedText = computed(() => {
    const n = props.term?.usedIn ?? 0;
    return `Użyte w ${n} ${n === 1 ? 'przepisie' : 'przepisach'}.`;
});

async function save() {
    errors.value = {};
    if (!name.value.trim()) errors.value.name = 'Podaj nazwę pojęcia';
    if (!explanation.value.trim()) errors.value.explanation = 'Napisz wyjaśnienie';
    else if (len.value > MAX) errors.value.explanation = `Najwyżej ${MAX} znaków`;
    if (errors.value.name || errors.value.explanation) return;
    saving.value = true;
    const body: IUpsertTerm = { name: name.value, explanation: explanation.value };
    try {
        const t = props.term ? await updateTerm(props.term.id, body) : await createTerm(body);
        emit('saved', { ...t, usedIn: props.term?.usedIn ?? 0 });
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
        :title="term ? 'Edytuj pojęcie' : 'Nowe pojęcie'"
        :width="480"
        @close="emit('close')"
    >
        <form
            :id="`${uid}-form`"
            class="term-form"
            @submit.prevent="save"
        >
            <div class="term-form__field">
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
                />
                <UFieldError v-if="errors.name">{{ errors.name }}</UFieldError>
            </div>
            <div class="term-form__field">
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
                    rows="4"
                    :aria-describedby="`${uid}-expl-hint`"
                />
                <UFieldError v-if="errors.explanation">{{ errors.explanation }}</UFieldError>
                <span
                    :id="`${uid}-expl-hint`"
                    class="term-form__hint"
                >
                    <span>1–3 zdania, zwykły tekst.</span>
                    <span
                        class="term-form__count"
                        :class="{ 'is-over': len > MAX }"
                    >
                        {{ len }} / {{ MAX }}
                    </span>
                </span>
            </div>
            <div>
                <div class="term-form__label">Podgląd</div>
                <div class="term-form__preview">
                    <TermExplainCard
                        class="term-form__card"
                        :name="name.trim() || 'Nazwa pojęcia'"
                        :explanation="explanation.trim() || 'Tu pojawi się wyjaśnienie.'"
                        preview
                    />
                </div>
            </div>
            <p
                v-if="term"
                class="term-form__used"
            >
                {{ usedText }}
            </p>
            <UFieldError v-if="errors.form">{{ errors.form }}</UFieldError>
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
                {{ saving ? 'Zapisywanie…' : 'Zapisz' }}
            </button>
        </template>
    </UModal>
</template>

<style scoped>
.term-form {
    display: flex;
    flex-direction: column;
    gap: 14px;
}
.term-form__field {
    display: flex;
    flex-direction: column;
    gap: 6px;
}
.term-form__hint {
    display: flex;
    justify-content: space-between;
    gap: 12px;
    font-size: 12px;
    color: var(--ink-muted);
}
.term-form__count {
    font-family: var(--font-mono);
}
.term-form__count.is-over {
    color: var(--danger);
}
.term-form__label {
    font-family: var(--font-mono);
    font-size: 11px;
    font-weight: 600;
    letter-spacing: 0.5px;
    text-transform: uppercase;
    color: var(--ink-muted);
    margin-bottom: 6px;
}
.term-form__preview {
    padding: 12px;
    border: 1px dashed var(--rule-strong);
    border-radius: var(--r-lg);
}
.term-form__card {
    max-width: 340px;
    padding: 14px 16px;
    background: var(--bg-alt);
    border: 1px solid var(--rule-strong);
    border-radius: var(--r-xl);
    box-shadow: 0 12px 32px rgba(10, 10, 15, 0.18);
}
.term-form__used {
    margin: 0;
    font-size: 12.5px;
    color: var(--ink-muted);
}
</style>
