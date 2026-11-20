<script lang="ts" setup>
import { ref } from 'vue';

import UModal from '@/components/UModal.vue';
import type { IEquipment } from '@/data/dtos/taxonomy/DictionariesDto';

// Deleting equipment in use: replace it with another one (default; details stay) or remove it from recipes.
const props = defineProps<{ item: IEquipment; others: IEquipment[] }>();
const emit = defineEmits<{
    confirm: [opt: { replaceWith?: number; removeFromRecipes?: boolean }];
    close: [];
}>();

const mode = ref<'replace' | 'remove'>(props.others.length ? 'replace' : 'remove');
const target = ref<number | ''>(props.others[0]?.id ?? '');

const usedText = (n: number) => `${n} ${n === 1 ? 'przepisie' : 'przepisach'}`;

function confirm() {
    if (mode.value === 'replace' && target.value !== '')
        emit('confirm', { replaceWith: target.value });
    else emit('confirm', { removeFromRecipes: true });
}
</script>

<template>
    <UModal
        :title="`Usunąć „${item.name}”?`"
        :width="460"
        @close="emit('close')"
    >
        <p class="modal__text">
            Jest używany w {{ usedText(item.usedIn) }}. Co zrobić z nim w tych przepisach?
        </p>
        <div
            class="del-eq"
            role="radiogroup"
            aria-label="Co zrobić w przepisach"
        >
            <div
                v-if="others.length"
                class="del-eq__opt"
            >
                <input
                    id="del-eq-replace"
                    v-model="mode"
                    type="radio"
                    name="del-eq"
                    value="replace"
                />
                <div class="del-eq__body">
                    <label for="del-eq-replace">Zamień na:</label>
                    <select
                        v-model="target"
                        class="select"
                        aria-label="Sprzęt, na który zamienić"
                        @focus="mode = 'replace'"
                    >
                        <option
                            v-for="o in others"
                            :key="o.id"
                            :value="o.id"
                        >
                            {{ o.name }}
                        </option>
                    </select>
                    <span class="del-eq__hint">Doprecyzowania w przepisach zostają.</span>
                </div>
            </div>
            <div class="del-eq__opt">
                <input
                    id="del-eq-remove"
                    v-model="mode"
                    type="radio"
                    name="del-eq"
                    value="remove"
                />
                <label
                    for="del-eq-remove"
                    class="del-eq__body"
                >
                    Usuń z tych przepisów
                </label>
            </div>
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
                @click="confirm"
            >
                Usuń sprzęt
            </button>
        </template>
    </UModal>
</template>

<style scoped>
.del-eq {
    display: flex;
    flex-direction: column;
    gap: 14px;
}
.del-eq__opt {
    display: grid;
    grid-template-columns: 20px minmax(0, 1fr);
    gap: 8px;
    align-items: start;
    font-size: 14px;
}
.del-eq__opt input[type='radio'] {
    margin: 3px 0 0;
    width: 16px;
    height: 16px;
    accent-color: var(--accent-fill);
}
.del-eq__body {
    display: flex;
    flex-direction: column;
    gap: 6px;
    cursor: pointer;
}
.del-eq__hint {
    font-size: 12px;
    color: var(--ink-muted);
}
</style>
