<script lang="ts" setup>
import { blankIng, type IGroup, type IIngItem, uid } from '../recipeForm';

import GroupShell from './GroupShell.vue';
import IngAutosuggest from './IngAutosuggest.vue';

import UFieldError from '@/components/form/UFieldError.vue';
import { confirmAction } from '@/composables/useConfirm';
import { useListDnd } from '@/composables/useListDnd';
import type { IngredientDto } from '@/data/dtos/ingredients/IngredientDto';
import type { IUnit } from '@/data/dtos/taxonomy/DictionariesDto';
import { parseQty, unitLabel } from '@/shared/utils/quantity';

const groups = defineModel<IGroup<IIngItem>[]>('groups', { required: true });
const props = defineProps<{
    ingredients: IngredientDto[];
    units: IUnit[];
    /** Błędy wierszy (id wiersza → komunikat); puste, dopóki nie próbowano zapisać. */
    rowErrors: Record<number, string>;
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
    // Zmiana tekstu odpina wybrany składnik, dopóki nie wybierze się z listy.
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

function removeItem(g: IGroup<IIngItem>, it: IIngItem) {
    g.items.splice(g.items.indexOf(it), 1);
}
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
                </div>
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
    grid-template-columns: 216px minmax(0, 1fr) 30px;
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
.ing-row__msg {
    padding-left: 28px;
}
.ing-row__preview {
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--ink-muted);
}
</style>
