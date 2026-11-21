<script lang="ts" setup>
import { computed, nextTick, ref } from 'vue';

import {
    blankIng,
    type IGroup,
    type IIngItem,
    type ILinkedRef,
    type ISubItem,
    uid,
} from '../recipeForm';

import GroupShell from './GroupShell.vue';
import IngAutosuggest from './IngAutosuggest.vue';
import RecipeLinkPicker from './RecipeLinkPicker.vue';
import RemovedIngNotice from './RemovedIngNotice.vue';
import SubstitutePicker from './SubstitutePicker.vue';

import BookIcon from '@/components/BookIcon.vue';
import UFieldError from '@/components/form/UFieldError.vue';
import SwapIcon from '@/components/SwapIcon.vue';
import { confirmAction } from '@/composables/useConfirm';
import { useListDnd } from '@/composables/useListDnd';
import type { IngredientDto } from '@/data/dtos/ingredients/IngredientDto';
import type { RecipeBaseDto } from '@/data/dtos/recipe/RecipeDto';
import type { IUnit } from '@/data/dtos/taxonomy/DictionariesDto';
import { parseQty, unitLabel } from '@/shared/utils/quantity';
import { stepRefs } from '@/shared/utils/stepRefs';

const groups = defineModel<IGroup<IIngItem>[]>('groups', { required: true });
const props = defineProps<{
    ingredients: IngredientDto[];
    units: IUnit[];
    /** Row errors (row id → message); empty until a save was attempted. */
    rowErrors: Record<number, string>;
    /** Recipes to link to (all of them, drafts too). */
    recipes: RecipeBaseDto[];
    /** The edited recipe (null = new); excluded from the picker. */
    currentId: number | null;
    /** Step texts in editor order (refs to rows, for the notice after removing a row). */
    stepTexts: string[];
}>();
const emit = defineEmits<{ 'add-new': [name: string, item: IIngItem] }>();

const dnd = useListDnd(groups);

const unitOf = (it: IIngItem) => props.units.find((u) => u.id === it.unitId);

function onUnit(it: IIngItem, value: string) {
    it.unitId = value ? Number(value) : null;
    if (unitOf(it)?.noQty) it.qty = '';
}

function onName(it: IIngItem, v: string) {
    it.name = v;
    // Editing the text unlinks the chosen ingredient until one is picked from the list.
    const match = props.ingredients.find((i) => i.id === it.ingredientId);
    if (!match || match.name !== v) it.ingredientId = 0;
}

function preview(it: IIngItem): string | null {
    const unit = unitOf(it);
    const q = parseQty(it.qty);
    if (!unit?.forms || !q.ok || q.empty) return null;
    return `w przepisie: ${it.qty.trim()} ${unitLabel(unit, null, q.max)}${it.name ? ` · ${it.name}` : ''}`;
}

function addGroup() {
    groups.value.push({ id: uid(), label: '', items: [blankIng()] });
}

async function removeGroup(g: IGroup<IIngItem>) {
    const hasContent = g.items.some((i) => i.name.trim() || i.qty.trim());
    if (
        hasContent &&
        !(await confirmAction({
            title: `Usunąć grupę${g.label ? ` „${g.label}”` : ''}?`,
            body: `Razem z nią znikną składniki z tej grupy (${g.items.length}).`,
            confirmLabel: 'Usuń grupę',
            danger: true,
        }))
    )
        return;
    groups.value.splice(groups.value.indexOf(g), 1);
}

interface IRemoved {
    item: IIngItem;
    groupId: number;
    group: string | null;
    index: number;
}
const removed = ref<IRemoved | null>(null);

/** Steps (editor numbering) that refer to a row, from the current steps (they may move or lose the ref). */
const stepsUsing = (key: string) =>
    props.stepTexts.flatMap((t, i) => (stepRefs(t).some((r) => r.key === key) ? [i + 1] : []));
/** The notice lasts while some step still refers to the removed row. */
const removedSteps = computed(() => (removed.value ? stepsUsing(removed.value.item.key) : []));

function removeItem(g: IGroup<IIngItem>, it: IIngItem) {
    const index = g.items.indexOf(it);
    g.items.splice(index, 1);
    removed.value =
        it.ingredientId && stepsUsing(it.key).length
            ? { item: it, groupId: g.id, group: g.label.trim() || null, index }
            : null;
}
function undoRemove() {
    const r = removed.value;
    if (!r) return;
    const g = groups.value.find((x) => x.id === r.groupId) ?? groups.value[0];
    g?.items.splice(Math.min(r.index, g.items.length), 0, r.item);
    removed.value = null;
}

const picking = ref<number | null>(null);

/** After the picker closes, focus returns to "Przepis…" / "Zmień". */
async function focusLinkBtn(it: IIngItem) {
    await nextTick();
    document.getElementById(`lk-btn-${it.id}`)?.focus();
}
function closePicker(it: IIngItem) {
    picking.value = null;
    focusLinkBtn(it);
}
function setLink(it: IIngItem, r: ILinkedRef) {
    it.linked = r;
    closePicker(it);
}
function unlink(it: IIngItem) {
    it.linked = null;
    focusLinkBtn(it);
}
const togglePicker = (it: IIngItem) => {
    subPicking.value = null;
    picking.value = picking.value === it.id ? null : it.id;
};

const subPicking = ref<number | null>(null);

async function focusSubBtn(it: IIngItem) {
    await nextTick();
    document.getElementById(`sub-btn-${it.id}`)?.focus();
}
const toggleSubPicker = (it: IIngItem) => {
    picking.value = null;
    subPicking.value = subPicking.value === it.id ? null : it.id;
};
function closeSubPicker(it: IIngItem) {
    subPicking.value = null;
    focusSubBtn(it);
}
function addSub(it: IIngItem, s: ISubItem) {
    it.substitutes.push(s);
    closeSubPicker(it);
}
function removeSub(it: IIngItem, s: ISubItem) {
    it.substitutes.splice(it.substitutes.indexOf(s), 1);
    focusSubBtn(it);
}
const subLabel = (s: ISubItem) => (s.mix ? `pół na pół z: ${s.name}` : s.name);
</script>

<template>
    <div class="ing-editor">
        <GroupShell
            v-for="(g, gi) in groups"
            :key="g.id"
            v-model:label="g.label"
            :index="gi"
            :group-count="groups.length"
            :placeholder="`Grupa ${gi + 1}, np. „Sos”`"
            noun="składnik"
            :empty="!g.items.length"
            :end-line="dnd.endLine(g.id)"
            :end-props="dnd.groupEnd(g.id)"
            @remove="removeGroup(g)"
        >
            <div
                v-for="it in g.items"
                :key="it.id"
                v-bind="dnd.row(g.id, it.id)"
                class="ing-row"
                :class="{ 'is-dragged': dnd.isDragged(it.id) }"
            >
                <span
                    v-if="dnd.lineAt(g.id, it.id)"
                    class="drop-line"
                    :class="`drop-line--${dnd.lineAt(g.id, it.id)}`"
                    aria-hidden="true"
                />
                <div class="ing-row__main">
                    <button
                        type="button"
                        class="drag-handle"
                        aria-label="Zmień kolejność: przeciągnij albo użyj strzałek ↑ ↓"
                        title="Przeciągnij"
                        v-bind="dnd.handle(g.id, it.id)"
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
                    <input
                        v-model="it.qty"
                        class="input input--mono ing-row__qty"
                        :class="{ 'is-invalid': rowErrors[it.id] && it.ingredientId }"
                        aria-label="Ilość"
                        :disabled="unitOf(it)?.noQty"
                        :placeholder="unitOf(it)?.noQty ? '—' : '2'"
                    />
                    <select
                        :value="it.unitId ?? ''"
                        class="select ing-row__unit"
                        aria-label="Jednostka"
                        @change="onUnit(it, ($event.target as HTMLSelectElement).value)"
                    >
                        <option value="">bez jednostki</option>
                        <option
                            v-for="u in units"
                            :key="u.id"
                            :value="u.id"
                        >
                            {{ u.name }}
                        </option>
                    </select>
                    <IngAutosuggest
                        :id="`ing-${it.id}`"
                        :name="it.name"
                        :ingredients="ingredients"
                        :invalid="!!rowErrors[it.id] && !it.ingredientId"
                        @update:name="onName(it, $event)"
                        @select="it.ingredientId = $event.id"
                        @add-new="emit('add-new', $event, it)"
                    />
                    <button
                        type="button"
                        class="row-icon-btn"
                        aria-label="Usuń składnik"
                        title="Usuń składnik"
                        @click="removeItem(g, it)"
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
                <div class="ing-row__extra">
                    <button
                        type="button"
                        class="ing-row__optional"
                        :class="{ 'is-active': it.optional }"
                        :aria-pressed="it.optional"
                        @click="it.optional = !it.optional"
                    >
                        {{ it.optional ? '✓ ' : '' }}Składnik opcjonalny
                    </button>
                    <input
                        v-model="it.note"
                        class="input ing-row__note"
                        aria-label="Notatka do składnika"
                        placeholder="notatka, np. ~150 g"
                    />
                    <button
                        :id="`sub-btn-${it.id}`"
                        type="button"
                        class="ing-row__link-btn ing-row__sub-btn"
                        :class="{ 'is-open': subPicking === it.id }"
                        aria-haspopup="listbox"
                        :aria-expanded="subPicking === it.id"
                        :disabled="!it.ingredientId"
                        :title="it.ingredientId ? undefined : 'Najpierw wybierz składnik'"
                        @click="toggleSubPicker(it)"
                    >
                        <SwapIcon :size="13" />
                        Zamiennik…
                    </button>
                    <button
                        v-if="!it.linked"
                        :id="`lk-btn-${it.id}`"
                        type="button"
                        class="ing-row__link-btn"
                        :class="{ 'is-open': picking === it.id }"
                        aria-haspopup="listbox"
                        :aria-expanded="picking === it.id"
                        @click="togglePicker(it)"
                    >
                        <BookIcon :size="13" />
                        Przepis…
                    </button>
                </div>
                <div
                    v-if="it.substitutes.length || it.linked"
                    class="ing-row__lines"
                >
                    <template v-if="it.substitutes.length">
                        <span class="ing-row__line-label">Zamienniki:</span>
                        <div class="ing-row__line ing-row__subs">
                            <span
                                v-for="s in it.substitutes"
                                :key="`${s.ingredientId}-${s.mix}`"
                                class="ing-row__chip"
                            >
                                <span class="ing-row__chip-title">{{ subLabel(s) }}</span>
                                <button
                                    type="button"
                                    class="ing-row__chip-x"
                                    :aria-label="`Usuń zamiennik ${subLabel(s)}`"
                                    @click="removeSub(it, s)"
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
                            </span>
                        </div>
                    </template>
                    <template v-if="it.linked">
                        <span class="ing-row__line-label">Prowadzi do:</span>
                        <div class="ing-row__line">
                            <span class="ing-row__chip">
                                <BookIcon :size="13" />
                                <span class="ing-row__chip-title">{{ it.linked.title }}</span>
                                <button
                                    type="button"
                                    class="ing-row__chip-x"
                                    :aria-label="`Usuń link do przepisu „${it.linked.title}”`"
                                    @click="unlink(it)"
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
                            </span>
                            <span
                                v-if="it.linked.status === 'draft'"
                                class="ing-row__draft"
                            >
                                szkic — link niewidoczny na stronie, dopóki go nie opublikujesz
                            </span>
                            <button
                                :id="`lk-btn-${it.id}`"
                                type="button"
                                class="ing-row__change"
                                aria-haspopup="listbox"
                                :aria-expanded="picking === it.id"
                                @click="togglePicker(it)"
                            >
                                Zmień
                            </button>
                        </div>
                    </template>
                </div>
                <SubstitutePicker
                    v-if="subPicking === it.id"
                    :id="`sp-${it.id}`"
                    :ingredient-name="it.name"
                    :main-id="it.ingredientId"
                    :current="it.substitutes"
                    :ingredients="ingredients"
                    @add="addSub(it, $event)"
                    @close="closeSubPicker(it)"
                />
                <RecipeLinkPicker
                    v-if="picking === it.id"
                    :id="`lk-${it.id}`"
                    :ingredient-name="it.name"
                    :recipes="recipes"
                    :current-id="currentId"
                    @select="setLink(it, $event)"
                    @close="closePicker(it)"
                />
                <div
                    v-if="rowErrors[it.id] || preview(it)"
                    class="ing-row__msg"
                >
                    <UFieldError v-if="rowErrors[it.id]">{{ rowErrors[it.id] }}</UFieldError>
                    <span
                        v-else
                        class="ing-row__preview"
                    >
                        {{ preview(it) }}
                    </span>
                </div>
            </div>
            <RemovedIngNotice
                v-if="removed?.groupId === g.id && removedSteps.length"
                :name="removed.item.name"
                :group="removed.group"
                :qty="removed.item.qty"
                :unit="unitOf(removed.item)?.name ?? ''"
                :steps="removedSteps"
                @undo="undoRemove"
            />
            <template #footer>
                <button
                    type="button"
                    class="add-link"
                    @click="g.items.push(blankIng())"
                >
                    + składnik
                </button>
            </template>
        </GroupShell>
        <button
            type="button"
            class="add-group-btn"
            @click="addGroup"
        >
            + Dodaj grupę składników
        </button>
    </div>
</template>

<style scoped>
.ing-editor {
    display: flex;
    flex-direction: column;
    gap: 16px;
}
.ing-row {
    position: relative;
    display: flex;
    flex-direction: column;
    gap: 6px;
    padding-bottom: 10px;
    border-bottom: 1px solid var(--rule);
}
.ing-row.is-dragged {
    opacity: 0.4;
}
.ing-row__main {
    display: grid;
    grid-template-columns: 20px 84px 124px minmax(0, 1fr) 30px;
    gap: 8px;
    align-items: center;
}
.ing-row__qty,
.ing-row__unit {
    padding: 8px 10px;
    font-size: 13px;
}
.ing-row__extra {
    display: grid;
    grid-template-columns: 216px minmax(0, 1fr) auto auto 30px;
    gap: 8px;
    align-items: center;
    padding-left: 28px;
}
.ing-row__optional {
    padding: 6px 10px;
    border-radius: var(--r-md);
    font-size: 12px;
    font-weight: 500;
    color: var(--ink-muted);
    border: 1px solid var(--rule);
    background: transparent;
    cursor: pointer;
}
.ing-row__optional.is-active {
    color: var(--accent);
    border-color: var(--accent);
    background: var(--accent-soft);
}
.ing-row__note {
    padding: 6px 10px;
    font-size: 13px;
}
.ing-row__link-btn.ing-row__sub-btn {
    grid-column: 3;
}
.ing-row__sub-btn:disabled {
    opacity: 0.45;
    cursor: default;
}
.ing-row__link-btn {
    grid-column: 4;
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 6px 10px;
    border-radius: var(--r-md);
    font-size: 12px;
    font-weight: 500;
    white-space: nowrap;
    color: var(--ink-dim);
    border: 1px solid var(--rule);
    background: transparent;
    cursor: pointer;
}
.ing-row__link-btn.is-open {
    color: var(--accent);
    border-color: var(--accent);
}
.ing-row__lines {
    padding-left: 28px;
    display: grid;
    grid-template-columns: 96px minmax(0, 1fr);
    gap: 8px;
    align-items: center;
}
.ing-row__line-label {
    font-size: 12px;
    color: var(--ink-muted);
}
.ing-row__line {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 6px 10px;
}
.ing-row__subs {
    gap: 6px;
}
.ing-row__chip {
    max-width: 420px;
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 4px 6px 4px 10px;
    border-radius: var(--r-pill);
    background: var(--accent-soft);
    color: var(--accent);
    font-size: 13px;
    font-weight: 600;
}
.ing-row__chip svg {
    flex-shrink: 0;
}
.ing-row__chip-title {
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}
.ing-row__chip-x {
    width: 24px;
    height: 24px;
    flex-shrink: 0;
    border-radius: var(--r-pill);
    border: none;
    background: transparent;
    color: var(--accent);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
}
.ing-row__chip-x:hover {
    background: color-mix(in srgb, var(--accent) 14%, transparent);
}
.ing-row__draft {
    padding: 2px 9px;
    border-radius: var(--r-pill);
    font-size: 12px;
    font-weight: 600;
    background: var(--prep-soft);
    color: var(--prep);
}
.ing-row__change {
    background: none;
    border: none;
    padding: 0;
    font-size: 13px;
    color: var(--ink-dim);
    text-decoration: underline;
    text-underline-offset: 3px;
    cursor: pointer;
}
.ing-row__msg {
    padding-left: 28px;
}
.ing-row__preview {
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--ink-muted);
}
</style>
