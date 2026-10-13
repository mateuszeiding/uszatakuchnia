<script lang="ts" setup>
import { computed, ref } from 'vue';

import UModal from '@/components/UModal.vue';

// Usuwanie używanej wartości słownika = przeniesienie przepisów do innej wartości.
// Dla osi wymaganej (Typ dania) przeniesienie jest obowiązkowe.
const props = defineProps<{
    label: string;
    /** Ile przepisów / składników używa wartości. */
    count: number;
    countNoun: string;
    others: { id: number; label: string }[];
    required: boolean;
}>();
const emit = defineEmits<{ confirm: [moveTo: number | null]; close: [] }>();

const target = ref<number | ''>('');
const canDelete = computed(() => !props.required || target.value !== '');
</script>

<template>
    <UModal
        :title="`Usunąć „${label}”?`"
        @close="emit('close')"
    >
        <p class="modal__text">
            Używa jej {{ count }} {{ countNoun }}.
            {{
                required
                    ? 'Ta wartość jest wymagana w przepisie, więc wybierz, dokąd je przenieść.'
                    : 'Możesz przenieść je do innej wartości albo po prostu usunąć wartość.'
            }}
        </p>
        <div class="field">
            <label
                for="move-to"
                class="field-label"
            >
                Przenieś do
                <span
                    v-if="required"
                    class="req"
                >
                    *
                </span>
            </label>
            <select
                id="move-to"
                v-model="target"
                class="select"
                autofocus
            >
                <option value="">{{ required ? 'Wybierz' : 'Nie przenoś' }}</option>
                <option
                    v-for="o in others"
                    :key="o.id"
                    :value="o.id"
                >
                    {{ o.label }}
                </option>
            </select>
        </div>
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
                class="btn btn--danger-fill"
                :disabled="!canDelete"
                @click="emit('confirm', target === '' ? null : target)"
            >
                {{ target === '' ? 'Usuń' : 'Przenieś i usuń' }}
            </button>
        </template>
    </UModal>
</template>
