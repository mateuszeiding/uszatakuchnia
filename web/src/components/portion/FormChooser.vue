<script lang="ts" setup>
import { computed, ref, useId, watch } from 'vue';

import FormGlyph from './FormGlyph.vue';
import PortionDialog from './PortionDialog.vue';

import {
    type FormSize,
    MAX_CM,
    MIN_CM,
    RECT_SIZES,
    ROUND_SIZES,
    sameSize,
} from '@/shared/utils/portion';

// "Twoja forma": typical sizes per shape + a custom size (10–60 cm). Picking applies at once.
const props = defineProps<{
    value: FormSize;
    /** The recipe's form: "w przepisie" under its size. */
    recipe: FormSize;
    sheet: boolean;
    anchor: HTMLElement | null;
}>();
const emit = defineEmits<{ select: [f: FormSize]; close: [] }>();

const uid = useId();
const shape = ref<'round' | 'rect'>(props.value.shape);

const options = computed<FormSize[]>(() =>
    shape.value === 'round'
        ? ROUND_SIZES.map((d) => ({ shape: 'round', d }))
        : RECT_SIZES.map(([a, b]) => ({ shape: 'rect', a, b }))
);
// the recipe's own size shows in the grid even when it isn't a typical one
const grid = computed(() =>
    props.recipe.shape === shape.value && !options.value.some((o) => sameSize(o, props.recipe))
        ? [...options.value, props.recipe]
        : options.value
);
const label = (f: FormSize) => (f.shape === 'round' ? `Ø ${f.d}` : `${f.a} × ${f.b}`);

// Custom size: prefilled with the current form of this shape.
const d = ref<string>('');
const a = ref<string>('');
const b = ref<string>('');
function prefill() {
    const v = props.value;
    d.value = v.shape === 'round' ? String(v.d) : '';
    a.value = v.shape === 'rect' ? String(v.a) : '';
    b.value = v.shape === 'rect' ? String(v.b) : '';
}
prefill();
watch(() => props.value, prefill);

const cm = (s: string) => {
    const n = Number(s.trim().replace(',', '.'));
    return Number.isInteger(n) && n >= MIN_CM && n <= MAX_CM ? n : null;
};
const touched = ref(false);
const custom = computed<FormSize | null>(() => {
    if (shape.value === 'round') {
        const v = cm(d.value);
        return v ? { shape: 'round', d: v } : null;
    }
    const x = cm(a.value);
    const y = cm(b.value);
    return x && y ? { shape: 'rect', a: x, b: y } : null;
});
const invalid = computed(
    () =>
        touched.value &&
        !custom.value &&
        (shape.value === 'round' ? !!d.value.trim() : !!(a.value.trim() && b.value.trim()))
);

function applyCustom() {
    touched.value = true;
    if (custom.value && !sameSize(custom.value, props.value)) emit('select', custom.value);
}

function onTabKey(e: KeyboardEvent) {
    if (e.key === 'ArrowRight' || e.key === 'ArrowLeft') {
        e.preventDefault();
        shape.value = shape.value === 'round' ? 'rect' : 'round';
        document.getElementById(`${uid}-t-${shape.value}`)?.focus();
    }
}
</script>

<template>
    <PortionDialog
        title="Twoja forma"
        :sheet="sheet"
        :anchor="anchor"
        @close="emit('close')"
    >
        <div
            class="fch__tabs"
            role="tablist"
            aria-label="Kształt formy"
        >
            <button
                v-for="t in [
                    { k: 'round', l: 'Okrągła' },
                    { k: 'rect', l: 'Prostokątna' },
                ] as const"
                :id="`${uid}-t-${t.k}`"
                :key="t.k"
                type="button"
                role="tab"
                class="fch__tab"
                :class="{ 'is-on': shape === t.k }"
                :aria-selected="shape === t.k"
                :aria-controls="`${uid}-p`"
                :tabindex="shape === t.k ? 0 : -1"
                @click="shape = t.k"
                @keydown="onTabKey"
            >
                <FormGlyph :shape="t.k" />
                {{ t.l }}
            </button>
        </div>

        <div
            :id="`${uid}-p`"
            role="tabpanel"
            class="fch__panel"
            :aria-labelledby="`${uid}-t-${shape}`"
        >
            <div
                class="fch__grid"
                role="radiogroup"
                aria-label="Typowe rozmiary"
            >
                <button
                    v-for="o in grid"
                    :key="label(o)"
                    type="button"
                    role="radio"
                    class="fch__opt"
                    :class="{ 'is-on': sameSize(o, value) }"
                    :aria-checked="sameSize(o, value)"
                    @click="emit('select', o)"
                >
                    {{ label(o) }}
                    <span
                        v-if="sameSize(o, recipe)"
                        class="fch__sub"
                    >
                        w przepisie
                    </span>
                </button>
            </div>

            <fieldset class="fch__custom">
                <legend class="fch__legend">Własny wymiar</legend>
                <div class="fch__inputs">
                    <template v-if="shape === 'round'">
                        <label class="fch__field">
                            <span class="fch__flabel">Średnica</span>
                            <span class="fch__box">
                                <input
                                    v-model="d"
                                    class="fch__input"
                                    inputmode="numeric"
                                    :aria-invalid="invalid || undefined"
                                    :aria-describedby="`${uid}-hint`"
                                    @change="applyCustom"
                                    @keydown.enter.prevent="applyCustom"
                                />
                                <span class="fch__unit">cm</span>
                            </span>
                        </label>
                    </template>
                    <template v-else>
                        <label class="fch__field">
                            <span class="fch__flabel">Długość</span>
                            <span class="fch__box">
                                <input
                                    v-model="a"
                                    class="fch__input"
                                    inputmode="numeric"
                                    :aria-invalid="invalid || undefined"
                                    :aria-describedby="`${uid}-hint`"
                                    @change="applyCustom"
                                    @keydown.enter.prevent="applyCustom"
                                />
                                <span class="fch__unit">cm</span>
                            </span>
                        </label>
                        <label class="fch__field">
                            <span class="fch__flabel">Szerokość</span>
                            <span class="fch__box">
                                <input
                                    v-model="b"
                                    class="fch__input"
                                    inputmode="numeric"
                                    :aria-invalid="invalid || undefined"
                                    :aria-describedby="`${uid}-hint`"
                                    @change="applyCustom"
                                    @keydown.enter.prevent="applyCustom"
                                />
                                <span class="fch__unit">cm</span>
                            </span>
                        </label>
                    </template>
                </div>
                <p
                    v-if="invalid"
                    class="fch__err"
                    role="alert"
                >
                    Podaj pełne centymetry od {{ MIN_CM }} do {{ MAX_CM }}.
                </p>
            </fieldset>
        </div>
        <p
            :id="`${uid}-hint`"
            class="fch__hint"
        >
            Przeliczamy po powierzchni dna, przy tej samej wysokości ciasta. Zakres {{ MIN_CM }}–{{
                MAX_CM
            }}
            cm.
        </p>
    </PortionDialog>
</template>

<style scoped>
.fch__tabs {
    display: flex;
    background: var(--bg);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    padding: 3px;
}
.fch__tab {
    flex: 1;
    min-height: 32px;
    border: none;
    border-radius: var(--r-md);
    background: transparent;
    color: var(--ink);
    font-size: 14px;
    font-weight: 500;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    cursor: pointer;
}
.fch__tab.is-on {
    background: var(--bg-alt);
    box-shadow: 0 0 0 1px var(--rule-strong);
    font-weight: 600;
}
.fch__panel {
    display: flex;
    flex-direction: column;
    gap: 14px;
}
.fch__grid {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 6px;
}
.fch__opt {
    min-height: 38px;
    padding: 4px 8px;
    border-radius: var(--r-lg);
    border: 1px solid var(--rule-strong);
    background: var(--bg-alt);
    color: var(--ink);
    font-size: 14px;
    font-weight: 500;
    line-height: 1.2;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    cursor: pointer;
}
.fch__opt:hover {
    border-color: var(--ink-muted);
}
.fch__opt.is-on {
    border-color: var(--accent);
    background: var(--accent-soft);
    color: var(--accent);
    font-weight: 600;
}
.fch__sub {
    font-size: 11px;
    font-weight: 500;
    color: var(--ink-muted);
}
.fch__opt.is-on .fch__sub {
    color: var(--accent);
}
.fch__custom {
    border: none;
    margin: 0;
    padding: 0;
    min-width: 0;
}
.fch__legend {
    padding: 0;
    margin-bottom: 6px;
    font-size: 12px;
    color: var(--ink-muted);
}
.fch__inputs {
    display: flex;
    gap: 8px;
}
.fch__field {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 4px;
    font-size: 12px;
    color: var(--ink-muted);
}
.fch__box {
    display: flex;
    align-items: center;
    gap: 6px;
    height: 44px;
    padding: 0 10px;
    border: 1px solid var(--rule-strong);
    border-radius: var(--r-md);
    background: var(--bg);
}
.fch__box:focus-within {
    border-color: var(--accent);
    box-shadow: var(--focus-ring);
}
.fch__input {
    flex: 1;
    min-width: 0;
    border: none;
    background: none;
    color: var(--ink);
    font-family: var(--font-mono);
    font-size: 15px;
    outline: none;
}
.fch__input[aria-invalid='true'] {
    color: var(--danger);
}
.fch__unit {
    font-size: 13px;
    color: var(--ink-muted);
}
.fch__err {
    margin: 6px 0 0;
    font-size: 12.5px;
    color: var(--danger);
}
.fch__hint {
    margin: 0;
    font-size: 12.5px;
    line-height: 1.45;
    color: var(--ink-muted);
}

:global(.pdlg--sheet) .fch__tab {
    min-height: 40px;
}
:global(.pdlg--sheet) .fch__opt {
    min-height: 48px;
}
</style>
