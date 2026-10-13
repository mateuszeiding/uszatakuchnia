<script lang="ts" setup>
import { blankStep, type IGroup, type IStepItem, PREP_LEADS, uid } from '../recipeForm';

import GroupShell from './GroupShell.vue';

import { confirmAction } from '@/composables/useConfirm';
import { useListDnd } from '@/composables/useListDnd';

const groups = defineModel<IGroup<IStepItem>[]>('groups', { required: true });
const dnd = useListDnd(groups);

/** Numer kroku w całym przepisie (ciągła numeracja przez etapy). */
function stepNo(gi: number, ii: number): number {
    let n = 0;
    for (let g = 0; g < gi; g++) n += groups.value[g]!.items.length;
    return n + ii + 1;
}

function togglePrep(s: IStepItem) {
    s.prep = !s.prep;
    if (!s.prep) s.prepLead = '';
}

async function removeGroup(g: IGroup<IStepItem>) {
    const hasContent = g.items.some((s) => s.text.trim());
    if (
        hasContent &&
        !(await confirmAction({
            title: `Usunąć etap${g.label ? ` „${g.label}”` : ''}?`,
            body: `Razem z nim znikną kroki z tego etapu (${g.items.length}).`,
            confirmLabel: 'Usuń etap',
            danger: true,
        }))
    )
        return;
    groups.value.splice(groups.value.indexOf(g), 1);
}
</script>

<template>
    <div class="steps-editor">
        <GroupShell
            v-for="(g, gi) in groups"
            :key="g.id"
            v-model:label="g.label"
            :index="gi"
            :group-count="groups.length"
            :placeholder="`Etap ${gi + 1}, np. „Ciasto”`"
            noun="krok"
            :empty="!g.items.length"
            :end-line="dnd.endLine(g.id)"
            :end-props="dnd.groupEnd(g.id)"
            @remove="removeGroup(g)"
        >
            <div
                v-for="(s, ii) in g.items"
                :key="s.id"
                v-bind="dnd.row(g.id, s.id)"
                class="step-row"
                :class="{ 'is-prep': s.prep, 'is-dragged': dnd.isDragged(s.id) }"
            >
                <span
                    v-if="dnd.lineAt(g.id, s.id)"
                    class="drop-line"
                    :class="`drop-line--${dnd.lineAt(g.id, s.id)}`"
                    aria-hidden="true"
                />
                <button
                    type="button"
                    class="drag-handle step-row__handle"
                    aria-label="Zmień kolejność: przeciągnij albo użyj strzałek ↑ ↓"
                    title="Przeciągnij"
                    v-bind="dnd.handle(g.id, s.id)"
                >
                    <svg
                        width="10"
                        height="14"
                        viewBox="0 0 10 14"
                        fill="currentColor"
                        aria-hidden="true"
                    >
                        <circle
                            v-for="(c, ci) in [
                                [3, 3],
                                [7, 3],
                                [3, 7],
                                [7, 7],
                                [3, 11],
                                [7, 11],
                            ]"
                            :key="ci"
                            :cx="c[0]"
                            :cy="c[1]"
                            r="1.2"
                        />
                    </svg>
                </button>
                <div
                    class="step-row__no"
                    aria-hidden="true"
                >
                    {{ stepNo(gi, ii) }}
                </div>
                <div class="step-row__body">
                    <textarea
                        v-model="s.text"
                        class="textarea step-row__text"
                        rows="2"
                        :aria-label="`Krok ${stepNo(gi, ii)}`"
                        :placeholder="
                            stepNo(gi, ii) === 1 ? 'Rozgrzej piekarnik do 200°C…' : 'Następny krok…'
                        "
                    />
                    <div class="step-row__prep">
                        <button
                            type="button"
                            class="prep-toggle"
                            :class="{ 'is-active': s.prep }"
                            :aria-pressed="s.prep"
                            @click="togglePrep(s)"
                        >
                            <svg
                                width="12"
                                height="12"
                                viewBox="0 0 12 12"
                                aria-hidden="true"
                            >
                                <circle
                                    cx="6"
                                    cy="6"
                                    r="4.8"
                                    stroke="currentColor"
                                    stroke-width="1.3"
                                />
                                <path
                                    d="M6 3.5v2.7l1.7 1"
                                    stroke="currentColor"
                                    stroke-width="1.3"
                                    stroke-linecap="round"
                                />
                            </svg>
                            {{ s.prep ? 'Przygotuj wcześniej' : 'Przygotuj wcześniej?' }}
                        </button>
                        <select
                            v-if="s.prep"
                            v-model="s.prepLead"
                            class="select prep-lead"
                            aria-label="Ile wcześniej"
                        >
                            <option value="">ile wcześniej?</option>
                            <option
                                v-for="l in PREP_LEADS"
                                :key="l"
                                :value="l"
                            >
                                {{ l === 'dzień wcześniej' ? l : `${l} wcześniej` }}
                            </option>
                        </select>
                    </div>
                </div>
                <button
                    type="button"
                    class="row-icon-btn step-row__del"
                    :aria-label="`Usuń krok ${stepNo(gi, ii)}`"
                    :title="`Usuń krok ${stepNo(gi, ii)}`"
                    @click="g.items.splice(ii, 1)"
                >
                    <svg
                        width="11"
                        height="11"
                        viewBox="0 0 11 11"
                        aria-hidden="true"
                    >
                        <path
                            d="M2 2l7 7M9 2L2 9"
                            stroke="currentColor"
                            stroke-width="1.5"
                            stroke-linecap="round"
                        />
                    </svg>
                </button>
            </div>
            <template #footer>
                <button
                    type="button"
                    class="add-link steps-editor__add"
                    @click="g.items.push(blankStep())"
                >
                    + krok
                </button>
            </template>
        </GroupShell>
        <button
            type="button"
            class="add-group-btn"
            @click="groups.push({ id: uid(), label: '', items: [blankStep()] })"
        >
            + Dodaj etap
        </button>
    </div>
</template>

<style scoped>
.steps-editor {
    display: flex;
    flex-direction: column;
    gap: 16px;
}
.steps-editor__add {
    margin-left: 60px;
}
.step-row {
    position: relative;
    display: grid;
    grid-template-columns: 20px 32px minmax(0, 1fr) 30px;
    gap: 8px;
    align-items: start;
    border-radius: var(--r-lg);
}
.step-row.is-prep {
    padding: 10px 10px 10px 0;
    margin-right: -10px;
    background: var(--prep-soft);
}
.step-row.is-dragged {
    opacity: 0.4;
}
.step-row__handle,
.step-row__del {
    margin-top: 6px;
}
.step-row__no {
    font-size: 22px;
    font-weight: 600;
    color: var(--accent);
    line-height: 1;
    padding-top: 10px;
    text-align: center;
}
.step-row__body {
    display: flex;
    flex-direction: column;
    gap: 8px;
}
.step-row__text {
    font-size: 15px;
}
.step-row__prep {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
}
.prep-toggle {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 5px 10px;
    border-radius: var(--r-pill);
    font-size: 12px;
    font-weight: 600;
    cursor: pointer;
    color: var(--ink-muted);
    background: transparent;
    border: 1px solid var(--rule);
}
.prep-toggle.is-active {
    color: var(--prep);
    background: var(--bg-alt);
    border-color: var(--prep-line);
}
.prep-lead {
    width: auto;
    padding: 5px 30px 5px 8px;
    font-size: 12px;
    background-color: var(--bg-alt);
}
</style>
