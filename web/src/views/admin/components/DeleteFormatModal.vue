<script lang="ts" setup>
import { ref } from 'vue';

import UModal from '@/components/UModal.vue';
import type { IChallengeFormat } from '@/data/dtos/challenge/ChallengeDto';
import { plural } from '@/shared/utils/challenges';

// Usunięcie używanego formatu: format jest wymagany w wyzwaniu, więc wyzwania trzeba przenieść do innego.
const props = defineProps<{ item: IChallengeFormat; others: IChallengeFormat[] }>();
const emit = defineEmits<{ confirm: [moveTo: number]; close: [] }>();

const target = ref<number | ''>(props.others[0]?.id ?? '');
</script>

<template>
    <UModal
        :title="`Usunąć „${item.name}”?`"
        :width="460"
        @close="emit('close')"
    >
        <p class="modal__text">
            Ma go {{ item.usedIn }} {{ plural(item.usedIn, 'wyzwanie', 'wyzwania', 'wyzwań') }}.
            {{
                others.length
                    ? 'Format jest wymagany, więc przenieś je do innego:'
                    : 'Format jest wymagany, a nie ma innego — najpierw dodaj nowy format.'
            }}
        </p>
        <div
            v-if="others.length"
            class="field"
        >
            <label
                for="format-move-to"
                class="field-label"
            >
                Przenieś do
            </label>
            <select
                id="format-move-to"
                v-model="target"
                class="select"
                autofocus
            >
                <option
                    v-for="o in others"
                    :key="o.id"
                    :value="o.id"
                >
                    {{ o.name }}
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
                :disabled="target === ''"
                @click="target !== '' && emit('confirm', target)"
            >
                Przenieś i usuń
            </button>
        </template>
    </UModal>
</template>
