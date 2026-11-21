// Ingredient refs inside step text: {{ing:<rowKey>}}, {{ing:<rowKey>|1/2}}, {{ing:<rowKey>|rest}}.
// The grammar is generic ({{<kind>:<id>|<arg>}}) so a glossary ref ({{term:…}}) can come later; the API
// (services/steprefs.go) rejects kinds it doesn't know, so here an unknown kind just renders nothing.
// A ref only shows the quantity — the shopping list, macros and allergens don't read it.

import type { RecipeIngredientDto } from '@/data/dtos/recipe/RecipeIngredientDto';
import type { IUnit } from '@/data/dtos/taxonomy/DictionariesDto';
import { scaleQuantity } from '@/shared/utils/quantity';
import { chosenSub, type SubsMap } from '@/shared/utils/substitutes';

export type RefPart = '1/2' | '1/3' | '1/4' | '2/3' | '3/4' | 'rest';

export const PART_VALUE: Record<Exclude<RefPart, 'rest'>, number> = {
    '1/2': 1 / 2,
    '1/3': 1 / 3,
    '1/4': 1 / 4,
    '2/3': 2 / 3,
    '3/4': 3 / 4,
};

/** Parts in picker order, with the glyph shown in the form. */
export const PARTS: { part: RefPart | null; label: string }[] = [
    { part: null, label: 'całość' },
    { part: '1/2', label: '½' },
    { part: '1/3', label: '⅓' },
    { part: '1/4', label: '¼' },
    { part: '2/3', label: '⅔' },
    { part: '3/4', label: '¾' },
    { part: 'rest', label: 'reszta' },
];

export const partLabel = (p: RefPart | null) => PARTS.find((x) => x.part === p)?.label ?? 'całość';

export interface IIngRef {
    kind: 'ing';
    key: string;
    part: RefPart | null;
}
export interface IOtherRef {
    kind: 'other';
    raw: string;
}
export type StepSeg = string | IIngRef | IOtherRef;

const REF_RE = /\{\{([^{}:|]*):([^{}|]*)(?:\|([^{}]*))?\}\}/g;

const isPart = (x: string): x is RefPart => x === 'rest' || x in PART_VALUE;

export function parseStep(text: string): StepSeg[] {
    const out: StepSeg[] = [];
    let last = 0;
    for (const m of text.matchAll(REF_RE)) {
        if (m.index > last) out.push(text.slice(last, m.index));
        const [raw, kind, key, arg] = m;
        if (kind === 'ing' && key && (!arg || isPart(arg)))
            out.push({ kind: 'ing', key, part: (arg as RefPart | undefined) ?? null });
        else out.push({ kind: 'other', raw });
        last = m.index + raw.length;
    }
    if (last < text.length) out.push(text.slice(last));
    return out;
}

export function refToken(key: string, part: RefPart | null): string {
    return part ? `{{ing:${key}|${part}}}` : `{{ing:${key}}}`;
}

export const stepRefs = (text: string): IIngRef[] =>
    parseStep(text).filter((s): s is IIngRef => typeof s !== 'string' && s.kind === 'ing');

export interface IRefContext {
    ingredients: RecipeIngredientDto[];
    /** All step texts (for "reszta"). */
    steps: string[];
    units: IUnit[];
    /** Chosen servings / recipe servings. */
    ratio: number;
    subs: SubsMap;
}

export interface IRefView {
    /** Quantity with unit, e.g. "350 g", "½ + ½ łyżki". */
    qty: string;
    /** Substitute shown instead of the main ingredient: "wieprzowina mielona" / "wołowina + wieprzowina". */
    lead: string | null;
    /** Group, only when the same ingredient is in more than one row. */
    group: string | null;
    optional: boolean;
    title: string;
}

const lowerFirst = (s: string) => s.charAt(0).toLowerCase() + s.slice(1);
const firstWord = (s: string) => lowerFirst(s.split(/\s+/)[0] ?? s);

/** Sum of the fractions of a row used across all steps ("reszta" and whole refs excluded). */
export function usedFraction(key: string, steps: string[]): number {
    let sum = 0;
    for (const t of steps)
        for (const r of stepRefs(t))
            if (r.key === key && r.part && r.part !== 'rest') sum += PART_VALUE[r.part];
    return sum;
}

export function partValue(ref: IIngRef, steps: string[]): number {
    if (!ref.part) return 1;
    if (ref.part === 'rest') return Math.max(0, 1 - usedFraction(ref.key, steps));
    return PART_VALUE[ref.part];
}

const amountText = (a: { qty: string; unit: string }) => (a.unit ? `${a.qty} ${a.unit}` : a.qty);

/**
 * What the ref shows: null = nothing (a row without a quantity, "do smaku"; or a row that's gone —
 * the API doesn't publish such a recipe, the form marks it as "brak składnika").
 */
export function resolveRef(ref: IIngRef, ctx: IRefContext): IRefView | null {
    const ing = ctx.ingredients.find((i) => i.key === ref.key);
    if (!ing) return null;
    const unit = ctx.units.find((u) => u.id === ing.unitId);
    const part = partValue(ref, ctx.steps);
    const whole = scaleQuantity(ing.quantity, ing.unit, unit, ctx.ratio);
    // nothing left for "reszta" (the fractions already add up to the whole): the formatter would round 0 up
    // to ¼, so show nothing; the form warns about it
    if (!whole.qty || part <= 1e-9) return null;
    const amount = scaleQuantity(ing.quantity, ing.unit, unit, ctx.ratio * part);

    const sub = chosenSub(ing, ctx.subs);
    let qty = amountText(amount);
    let lead: string | null = null;
    if (sub?.mix) {
        const h = scaleQuantity(ing.quantity, ing.unit, unit, ctx.ratio * part * 0.5);
        qty = h.unit ? `${h.qty} + ${h.qty} ${h.unit}` : `${h.qty} + ${h.qty}`;
        lead = `${firstWord(ing.ingredientName)} + ${firstWord(sub.name)}`;
    } else if (sub) lead = lowerFirst(sub.name);

    const dup =
        ctx.ingredients.filter((i) => i.ingredientId === ing.ingredientId).length > 1 &&
        !!ing.section;
    const partText = ref.part
        ? ` · ${ref.part === 'rest' ? 'reszta' : partLabel(ref.part)} z ${amountText(whole)}`
        : '';
    return {
        qty,
        lead,
        group: dup ? ing.section!.toLowerCase() : null,
        optional: ing.optional,
        title: `${ing.ingredientName}${ing.section ? ` · ${ing.section}` : ''}${partText}`,
    };
}

/** The ref as plain text, as the reader sees it (for the form's length counter). */
export function refText(v: IRefView | null): string {
    if (!v) return '';
    return (
        ' ' +
        (v.lead ? `⇄ ${v.lead} · ` : '') +
        v.qty +
        (v.group ? ` · ${v.group}` : '') +
        (v.optional ? ' (opcjonalnie)' : '')
    );
}

/**
 * What a step shows: text and resolved refs. A ref that shows nothing ("sól do smaku") takes the author's
 * space before it along, so "Dopraw solą {{ing:x}}." reads "Dopraw solą.", not "Dopraw solą .".
 */
export function viewSegments(text: string, ctx: IRefContext): (string | IRefView)[] {
    const out: (string | IRefView)[] = [];
    /** Index of the text before a ref that showed nothing. */
    let before = -1;
    const trimBefore = () => {
        const prev = out[before];
        if (typeof prev === 'string') out[before] = prev.replace(/\s+$/, '');
    };
    for (const seg of parseStep(text)) {
        if (typeof seg === 'string') {
            if (before >= 0 && /^[,.;:!?)]/.test(seg)) trimBefore();
            before = -1;
            out.push(seg);
            continue;
        }
        const v = seg.kind === 'ing' ? resolveRef(seg, ctx) : null;
        if (v) {
            before = -1;
            out.push(v);
        } else if (out.length && typeof out[out.length - 1] === 'string') before = out.length - 1;
    }
    if (before >= 0) trimBefore();
    return out;
}

/** Step length as the reader sees it: refs replaced by their text, whitespace collapsed like in HTML. */
export function visibleLength(text: string, ctx: IRefContext): number {
    const s = viewSegments(text, ctx)
        .map((seg) => (typeof seg === 'string' ? seg : refText(seg)))
        .join('')
        .replace(/\s+/g, ' ')
        .trim();
    return [...s].length;
}
