<script lang="ts" setup>
import { nextTick } from 'vue';

import { blankStep, type IGroup, type IStepItem, PREP_LEADS, uid } from '../recipeForm';
import { COUNT_FROM, type IRefRow, STEP_LIMIT } from '../stepRefsForm';

import GroupShell from './GroupShell.vue';
import StepHintsEditor from './StepHintsEditor.vue';
import StepTextEditor from './StepTextEditor.vue';

import { confirmAction } from '@/composables/useConfirm';
import { useListDnd } from '@/composables/useListDnd';
import { type IRefContext, type ITerm, stepRefs, visibleLength } from '@/shared/utils/stepRefs';

const groups = defineModel<IGroup<IStepItem>[]>('groups', { required: true });
const props = defineProps<{
    /** Ingredient rows a ref can point to. */
    rows: IRefRow[];
    /** For the length counter: refs counted as shown at the default servings. */
    refCtx: IRefContext;
    /** Per step (continuous numbering): refers to a removed row. */
    broken: boolean[];
    /** Rows whose fractions add up to more than the whole. */
    overused: Set<string>;
    /** The glossary (marking terms in steps and hints). */
    terms: ITerm[];
    /** Per step: marks a term that isn't in the glossary. */
    missingTerms: boolean[];
}>();
const emit = defineEmits<{ 'term-created': [term: ITerm] }>();
const dnd = useListDnd(groups);

const len = (s: IStepItem) => visibleLength(s.text, props.refCtx);
const overusedIn = (s: IStepItem) =>
    [...new Set(stepRefs(s.text).map((r) => r.key))]
        .filter((k) => props.overused.has(k))
        .map((k) => props.rows.find((r) => r.key === k)?.name)
        .filter(Boolean);

/** Step number across the whole recipe (continuous numbering through stages). */
function stepNo(gi: number, ii: number): number {
    let n = 0;
    for (let g = 0; g < gi; g++) n += groups.value[g]!.items.length;
    return n + ii + 1;
}

function togglePrep(s: IStepItem) {
    s.prep = !s.prep;
    if (!s.prep) s.prepLead = '';
}

async function addHint(s: IStepItem) {
    s.hints.push('');
    await nextTick();
    document.getElementById(`hint-${s.id}-${s.hints.length - 1}`)?.focus();
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
        <p class="steps-editor__limit">
            Krok w trybie gotowania mieści się na telefonie do ok. {{ STEP_LIMIT }} znaków (ilości
            liczą się tak, jak je widać) i z najwyżej 3 krótkimi wskazówkami.
        </p>
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
                    <StepTextEditor
                        v-model="s.text"
                        :rows="rows"
                        :terms="terms"
                        :label="`Krok ${stepNo(gi, ii)}`"
                        :describedby="`step-msg-${s.id}`"
                        :warn="len(s) > STEP_LIMIT"
                        :broken="broken[stepNo(gi, ii) - 1] || missingTerms[stepNo(gi, ii) - 1]"
                        @term-created="emit('term-created', $event)"
                        :placeholder="
                            stepNo(gi, ii) === 1 ? 'Rozgrzej piekarnik do 200°C…' : 'Następny krok…'
                        "
                    />
                    <div class="step-row__meta">
                        <span
                            :id="`step-msg-${s.id}`"
                            class="step-row__msg"
                            :class="{
                                'is-broken':
                                    broken[stepNo(gi, ii) - 1] || missingTerms[stepNo(gi, ii) - 1],
                                'is-warn':
                                    !broken[stepNo(gi, ii) - 1] &&
                                    !missingTerms[stepNo(gi, ii) - 1] &&
                                    len(s) > STEP_LIMIT,
                            }"
                        >
                            <template v-if="broken[stepNo(gi, ii) - 1]">
                                Składnik z odwołania został usunięty. Wybierz inny albo usuń
                                odwołanie — do tego czasu nie da się opublikować.
                            </template>
                            <template v-else-if="missingTerms[stepNo(gi, ii) - 1]">
                                Oznaczonego pojęcia nie ma już w słowniczku. Wybierz inne albo usuń
                                oznaczenie — do tego czasu nie da się opublikować.
                            </template>
                            <template v-else-if="len(s) > STEP_LIMIT">
                                Za długi do trybu gotowania — na telefonie nie zmieści się bez
                                przewijania. Podziel krok na dwa albo przenieś rady do wskazówek.
                            </template>
                            <template v-else>
                                Wpisz
                                <code>@</code>
                                , żeby wstawić ilość. Zaznacz słowo, żeby oznaczyć pojęcie.
                            </template>
                        </span>
                        <span
                            v-if="len(s) >= COUNT_FROM"
                            class="step-row__count"
                            :class="{ 'is-over': len(s) > STEP_LIMIT }"
                        >
                            {{ len(s) }} / {{ STEP_LIMIT }}
                        </span>
                    </div>
                    <p
                        v-if="overusedIn(s).length"
                        class="step-row__msg is-warn"
                    >
                        Części składnika „{{ overusedIn(s).join('”, „') }}” w krokach dają więcej
                        niż całość albo nie zostawiają nic na „resztę”.
                    </p>
                    <StepHintsEditor
                        v-model:hints="s.hints"
                        :step-id="s.id"
                        :step-no="stepNo(gi, ii)"
                        :terms="terms"
                        @term-created="emit('term-created', $event)"
                    />
                    <div class="step-row__prep">
                        <button
                            type="button"
                            class="add-link hint-add"
                            @click="addHint(s)"
                        >
                            + wskazówka
                        </button>
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
.steps-editor__limit {
    margin: -6px 0 0;
    font-size: 13px;
    line-height: 1.5;
    color: var(--ink-muted);
}
.step-row__meta {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 12px;
}
.step-row__msg {
    margin: 0;
    font-size: 12.5px;
    line-height: 1.45;
    color: var(--ink-muted);
}
.step-row__msg code {
    font-family: var(--font-mono);
}
.step-row__msg.is-warn {
    color: var(--prep);
}
.step-row__msg.is-broken {
    color: var(--danger);
}
.step-row__count {
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--ink-muted);
    white-space: nowrap;
}
.step-row__count.is-over {
    color: var(--prep);
    font-weight: 600;
}
.step-row__prep {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-wrap: wrap;
}
.hint-add {
    margin-right: 8px;
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
