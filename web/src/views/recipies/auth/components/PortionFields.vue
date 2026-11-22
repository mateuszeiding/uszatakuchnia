<script lang="ts" setup>
import { computed, useId, watch } from 'vue';

import { blankSize, type IPortionForm, pieceUnits } from '../portionForm';

import UFieldError from '@/components/form/UFieldError.vue';
import type { IUnit } from '@/data/dtos/taxonomy/DictionariesDto';
import { MAX_CM, MIN_CM, type PortionMode } from '@/shared/utils/portion';

// "Czym jest porcja" (Podstawy): the mode and its fields. Servings recipes keep their count in the side card;
// for pieces the same count is "W przepisie: n ×" here; a form has no count.
const portion = defineModel<IPortionForm>({ required: true });
const servings = defineModel<string>('servings', { required: true });
const props = defineProps<{ units: IUnit[]; error?: string }>();

const uid = useId();
const MODES: { k: PortionMode; l: string; hint: string }[] = [
    { k: 'servings', l: 'Porcje', hint: 'Osoba albo sztuka — licznik −/+ jak dotąd.' },
    {
        k: 'sized',
        l: 'Sztuki w rozmiarach',
        hint: 'Sztuki w kilku rozmiarach (pizza, tortille). Czytelnik wybiera rozmiar i liczbę, przeliczamy po masie na sztukę.',
    },
    {
        k: 'form',
        l: 'Forma lub blacha',
        hint: 'Ciasto w formie lub na blasze. Czytelnik wybiera swoją formę, przeliczamy po powierzchni dna.',
    },
];
const hint = computed(() => MODES.find((m) => m.k === portion.value.mode)?.hint);

const units = computed(() => pieceUnits(props.units));
// first switch to pieces: suggest "placek" and one size row to fill
watch(
    () => portion.value.mode,
    (m) => {
        if (m !== 'sized') return;
        if (!portion.value.pieceUnitId)
            portion.value.pieceUnitId =
                units.value.find((u) => u.name === 'placek')?.id ?? units.value[0]?.id ?? null;
        if (!portion.value.sizes.length) portion.value.sizes = [blankSize(true)];
    }
);

const base = computed(() => portion.value.sizes.find((s) => s.base));
const setBase = (id: number) => portion.value.sizes.forEach((s) => (s.base = s.id === id));
function removeSize(id: number) {
    const wasBase = base.value?.id === id;
    portion.value.sizes = portion.value.sizes.filter((s) => s.id !== id);
    if (wasBase && portion.value.sizes[0]) portion.value.sizes[0].base = true;
}
const addSize = () => portion.value.sizes.push(blankSize(!portion.value.sizes.some((s) => s.base)));
</script>

<template>
    <div
        id="f-portion"
        class="pf"
    >
        <fieldset class="pf__group">
            <legend class="field-label">Czym jest porcja</legend>
            <div class="pf__seg">
                <label
                    v-for="m in MODES"
                    :key="m.k"
                    class="pf__seg-opt"
                    :class="{ 'is-on': portion.mode === m.k }"
                >
                    <input
                        v-model="portion.mode"
                        type="radio"
                        class="sr-only"
                        :name="`${uid}-mode`"
                        :value="m.k"
                    />
                    {{ m.l }}
                </label>
            </div>
            <p class="field-hint pf__hint">{{ hint }}</p>
        </fieldset>

        <template v-if="portion.mode === 'sized'">
            <div class="pf__row">
                <div class="field">
                    <label
                        :for="`${uid}-unit`"
                        class="field-label"
                    >
                        Nazwa sztuki
                        <span class="req">*</span>
                    </label>
                    <select
                        :id="`${uid}-unit`"
                        v-model="portion.pieceUnitId"
                        class="select pf__unit"
                    >
                        <option :value="null">— wybierz —</option>
                        <option
                            v-for="u in units"
                            :key="u.id"
                            :value="u.id"
                        >
                            {{ u.forms!.join(' / ') }}
                        </option>
                    </select>
                    <span class="field-hint">Odmiana jak w słowniku jednostek.</span>
                </div>
                <div class="field">
                    <label
                        :for="`${uid}-n`"
                        class="field-label"
                    >
                        W przepisie
                        <span class="req">*</span>
                    </label>
                    <span class="pf__inline">
                        <input
                            :id="`${uid}-n`"
                            v-model="servings"
                            class="input input--mono pf__n"
                            type="number"
                            min="1"
                            max="20"
                        />
                        <span class="pf__times">× rozmiar zaznaczony niżej</span>
                    </span>
                </div>
            </div>

            <fieldset class="pf__group">
                <legend class="field-label">Rozmiary</legend>
                <div
                    v-if="portion.sizes.length"
                    class="pf__sizes"
                    role="table"
                    aria-label="Rozmiary"
                >
                    <div
                        class="pf__size pf__size--head"
                        role="row"
                    >
                        <span role="columnheader">Średnica</span>
                        <span role="columnheader">Ciasto na sztukę</span>
                        <span role="columnheader">W przepisie</span>
                        <span role="columnheader"><span class="sr-only">Usuń</span></span>
                    </div>
                    <div
                        v-for="(s, i) in portion.sizes"
                        :key="s.id"
                        class="pf__size"
                        role="row"
                    >
                        <label
                            class="pf__cell"
                            role="cell"
                        >
                            <span class="sr-only">Średnica, rozmiar {{ i + 1 }}</span>
                            <input
                                v-model="s.d"
                                class="input input--mono"
                                inputmode="numeric"
                                :placeholder="`${MIN_CM}–${MAX_CM}`"
                            />
                            <span class="pf__unit-cm">cm</span>
                        </label>
                        <label
                            class="pf__cell"
                            role="cell"
                        >
                            <span class="sr-only">Ciasto na sztukę, rozmiar {{ i + 1 }}</span>
                            <input
                                v-model="s.g"
                                class="input input--mono"
                                inputmode="numeric"
                            />
                            <span class="pf__unit-cm">g</span>
                        </label>
                        <label
                            class="pf__cell pf__base"
                            role="cell"
                        >
                            <input
                                type="radio"
                                :name="`${uid}-base`"
                                :checked="s.base"
                                @change="setBase(s.id)"
                            />
                            <span>{{ s.base ? 'ten' : '' }}</span>
                            <span class="sr-only">w przepisie, rozmiar {{ i + 1 }}</span>
                        </label>
                        <span role="cell">
                            <button
                                type="button"
                                class="pf__del"
                                :aria-label="`Usuń rozmiar ${s.d ? `${s.d} cm` : i + 1}`"
                                @click="removeSize(s.id)"
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
                </div>
                <button
                    type="button"
                    class="pf__add"
                    @click="addSize"
                >
                    + rozmiar
                </button>
                <p class="field-hint">
                    Czytelnik wybiera tylko z tych rozmiarów — dodaj wszystkie, które mają sens.
                    Ilości składników wpisujesz dla
                    {{ servings || 'n' }} × {{ base?.d ? `${base.d} cm` : 'rozmiar w przepisie' }}.
                </p>
            </fieldset>
        </template>

        <template v-if="portion.mode === 'form'">
            <div class="pf__row pf__row--end">
                <fieldset class="pf__group">
                    <legend class="field-label">Kształt</legend>
                    <div class="pf__seg">
                        <label
                            v-for="s in [
                                { k: 'round', l: 'Okrągła' },
                                { k: 'rect', l: 'Prostokątna' },
                            ] as const"
                            :key="s.k"
                            class="pf__seg-opt"
                            :class="{ 'is-on': portion.shape === s.k }"
                        >
                            <input
                                v-model="portion.shape"
                                type="radio"
                                class="sr-only"
                                :name="`${uid}-shape`"
                                :value="s.k"
                            />
                            {{ s.l }}
                        </label>
                    </div>
                </fieldset>
                <div
                    v-if="portion.shape === 'round'"
                    class="field"
                >
                    <label
                        :for="`${uid}-d`"
                        class="field-label"
                    >
                        Średnica (cm)
                        <span class="req">*</span>
                    </label>
                    <input
                        :id="`${uid}-d`"
                        v-model="portion.d"
                        class="input input--mono pf__dim"
                        inputmode="numeric"
                    />
                </div>
                <template v-else>
                    <div class="field">
                        <label
                            :for="`${uid}-a`"
                            class="field-label"
                        >
                            Długość (cm)
                            <span class="req">*</span>
                        </label>
                        <input
                            :id="`${uid}-a`"
                            v-model="portion.a"
                            class="input input--mono pf__dim"
                            inputmode="numeric"
                        />
                    </div>
                    <div class="field">
                        <label
                            :for="`${uid}-b`"
                            class="field-label"
                        >
                            Szerokość (cm)
                            <span class="req">*</span>
                        </label>
                        <input
                            :id="`${uid}-b`"
                            v-model="portion.b"
                            class="input input--mono pf__dim"
                            inputmode="numeric"
                        />
                    </div>
                </template>
                <div class="field">
                    <label
                        :for="`${uid}-pieces`"
                        class="field-label"
                    >
                        Kawałków z formy
                    </label>
                    <input
                        :id="`${uid}-pieces`"
                        v-model="portion.pieces"
                        class="input input--mono pf__dim"
                        inputmode="numeric"
                        :aria-describedby="`${uid}-pieces-h`"
                    />
                    <span
                        :id="`${uid}-pieces-h`"
                        class="field-hint"
                    >
                        Opcjonalne — do wartości odżywczych.
                    </span>
                </div>
            </div>
            <p class="field-hint">
                Ilości składników wpisujesz dla tej formy. Wymiary w pełnych centymetrach,
                {{ MIN_CM }}–{{ MAX_CM }} cm.
            </p>
        </template>

        <UFieldError v-if="error">{{ error }}</UFieldError>
    </div>
</template>

<style scoped>
.pf {
    display: flex;
    flex-direction: column;
    gap: 18px;
    padding-top: 18px;
    border-top: 1px solid var(--rule);
}
.pf__group {
    border: none;
    margin: 0;
    padding: 0;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 8px;
}
.pf__group > legend {
    padding: 0;
    margin-bottom: 8px;
}
.pf__seg {
    display: inline-flex;
    flex-wrap: wrap;
    align-self: flex-start;
    background: var(--bg);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    padding: 3px;
}
.pf__seg-opt {
    min-height: 34px;
    padding: 0 14px;
    border-radius: var(--r-md);
    display: inline-flex;
    align-items: center;
    font-size: 14px;
    font-weight: 500;
    color: var(--ink);
    cursor: pointer;
}
.pf__seg-opt.is-on {
    background: var(--bg-alt);
    box-shadow: 0 0 0 1px var(--rule-strong);
    font-weight: 600;
}
.pf__seg-opt:has(input:focus-visible) {
    box-shadow: var(--focus-ring);
}
.pf__hint {
    margin: 0;
}
.pf__row {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
}
.pf__row--end {
    align-items: flex-start;
}
.pf__unit {
    width: 220px;
}
.pf__inline {
    display: flex;
    align-items: center;
    gap: 8px;
}
.pf__n {
    width: 72px;
}
.pf__times {
    font-size: 14px;
    color: var(--ink-dim);
}
.pf__dim {
    width: 110px;
}
.pf__sizes {
    display: flex;
    flex-direction: column;
    gap: 6px;
}
.pf__size {
    display: grid;
    grid-template-columns: 120px 150px 110px 32px;
    gap: 8px;
    align-items: center;
}
.pf__size--head {
    font-size: 12px;
    color: var(--ink-muted);
}
.pf__cell {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 13px;
    color: var(--ink-muted);
}
.pf__base {
    gap: 8px;
    color: var(--ink);
    font-size: 14px;
    cursor: pointer;
}
.pf__base input {
    accent-color: var(--accent-fill);
}
.pf__del {
    width: 32px;
    height: 32px;
    border: none;
    background: none;
    color: var(--ink-muted);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
}
.pf__del:hover {
    color: var(--danger);
}
.pf__add {
    align-self: flex-start;
    padding: 4px 0;
    border: none;
    background: none;
    color: var(--accent);
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
}
.pf__unit-cm {
    flex-shrink: 0;
}

@media (max-width: 600px) {
    .pf__size {
        grid-template-columns: minmax(0, 1fr) minmax(0, 1fr) 64px 32px;
    }
}
</style>
