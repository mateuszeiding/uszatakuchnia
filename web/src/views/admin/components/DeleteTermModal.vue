<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue';

import UModal from '@/components/UModal.vue';
import { fetchTermUsage, type IGlossaryTerm, type ITermUsage } from '@/data/api/glossary/fetch';

// Deleting a term marked in recipes: remove the marks (default; the words stay as plain text) or point them at
// another term. The recipes are named, so the owner knows what changes.
const props = defineProps<{ term: IGlossaryTerm; others: IGlossaryTerm[] }>();
const emit = defineEmits<{
    confirm: [opt: { replaceWith?: number; removeMarks?: boolean }];
    close: [];
}>();

const mode = ref<'remove' | 'replace'>('remove');
const target = ref<number | ''>('');
const usage = ref<ITermUsage[] | null>(null);
/** The recipe list didn't load: keep the known count instead of showing "0 recipes". */
const usageFailed = ref(false);
async function loadUsage() {
    usageFailed.value = false;
    try {
        usage.value = await fetchTermUsage(props.term.id);
    } catch {
        usageFailed.value = true;
    }
}
onMounted(loadUsage);

const n = computed(() => usage.value?.length ?? props.term.usedIn ?? 0);
const names = computed(() =>
    (usage.value ?? []).map((u) => (u.status === 'draft' ? `${u.name} (szkic)` : u.name)).join(', ')
);
const canConfirm = computed(() => mode.value === 'remove' || target.value !== '');

function confirm() {
    if (mode.value === 'replace' && target.value !== '')
        emit('confirm', { replaceWith: target.value });
    else emit('confirm', { removeMarks: true });
}
</script>

<template>
    <UModal
        :title="`Usunąć „${term.name}”?`"
        :width="460"
        @close="emit('close')"
    >
        <p class="modal__text">
            Pojęcie jest oznaczone w {{ n }} {{ n === 1 ? 'przepisie' : 'przepisach'
            }}{{ names ? `: ${names}` : '' }}.
        </p>
        <p
            v-if="usageFailed"
            class="modal__text del-term__err"
        >
            Nie udało się wczytać listy przepisów.
            <button
                type="button"
                class="del-term__retry"
                @click="loadUsage"
            >
                Spróbuj ponownie
            </button>
        </p>
        <div
            class="del-term"
            role="radiogroup"
            aria-label="Co zrobić z oznaczeniami"
        >
            <div class="del-term__opt">
                <input
                    id="del-term-remove"
                    v-model="mode"
                    type="radio"
                    name="del-term"
                    value="remove"
                />
                <label
                    for="del-term-remove"
                    class="del-term__body"
                >
                    <b>Usuń oznaczenia</b>
                    <span class="del-term__hint">Słowa w krokach zostaną, bez podkreślenia.</span>
                </label>
            </div>
            <div
                v-if="others.length"
                class="del-term__opt"
            >
                <input
                    id="del-term-replace"
                    v-model="mode"
                    type="radio"
                    name="del-term"
                    value="replace"
                />
                <div class="del-term__body">
                    <label for="del-term-replace"><b>Zamień na…</b></label>
                    <select
                        v-model="target"
                        class="select"
                        aria-label="Pojęcie, na które zamienić"
                        @focus="mode = 'replace'"
                    >
                        <option
                            value=""
                            disabled
                        >
                            wybierz pojęcie
                        </option>
                        <option
                            v-for="o in others"
                            :key="o.id"
                            :value="o.id"
                        >
                            {{ o.name }}
                        </option>
                    </select>
                    <span class="del-term__hint">
                        Słowa zostają, wyjaśnienie będzie z wybranego.
                    </span>
                </div>
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
                :disabled="!canConfirm"
                @click="confirm"
            >
                Usuń pojęcie
            </button>
        </template>
    </UModal>
</template>

<style scoped>
.del-term {
    display: flex;
    flex-direction: column;
    gap: 14px;
}
.del-term__opt {
    display: grid;
    grid-template-columns: 20px minmax(0, 1fr);
    gap: 8px;
    align-items: start;
    font-size: 14px;
}
.del-term__opt input[type='radio'] {
    margin: 3px 0 0;
    width: 16px;
    height: 16px;
    accent-color: var(--accent-fill);
}
.del-term__body {
    display: flex;
    flex-direction: column;
    gap: 6px;
    cursor: pointer;
}
.del-term__err {
    color: var(--danger);
}
.del-term__retry {
    background: none;
    border: none;
    padding: 0;
    font: inherit;
    font-weight: 600;
    color: var(--accent);
    text-decoration: underline;
    cursor: pointer;
}
.del-term__hint {
    font-size: 13px;
    color: var(--ink-muted);
}
</style>
