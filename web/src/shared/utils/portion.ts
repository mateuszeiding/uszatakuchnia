// Portion size: one place that turns the reader's choice (servings, pieces in a size, their own form) into the
// scale used by the ingredient list, quantities in steps, the shopping list and the macros.
// The link preview mirrors the URL format and names in the API: handlers/share_portion.go.

import type { IUnit } from '@/data/dtos/taxonomy/DictionariesDto';
import { pluralPL } from '@/shared/utils/quantity';

export type PortionMode = 'servings' | 'sized' | 'form';

export interface IPieceSize {
    diameterCm: number;
    doughG: number;
    isBase: boolean;
}

export interface IFormDto {
    shape: 'round' | 'rect';
    dCm: number | null;
    aCm: number | null;
    bCm: number | null;
    /** Pieces cut from the form (per-piece macros); null = not given. */
    pieces: number | null;
}

export interface IPortion {
    mode: PortionMode;
    pieceUnit: IUnit | null;
    sizes: IPieceSize[];
    form: IFormDto | null;
}

/** The reader's form (or the recipe's). */
export type FormSize = { shape: 'round'; d: number } | { shape: 'rect'; a: number; b: number };
/** The reader's pieces: how many and which of the author's diameters. */
export type PieceChoice = { n: number; d: number };
/** Stored in uk:progress:<id>.size; missing = as in the recipe. */
export type SizeChoice = PieceChoice | FormSize;

/** What the portion code needs from a recipe (detail DTO and list card alike). */
export interface IPortionRecipe {
    servings: number;
    portion?: IPortion | null;
}

export const MIN_CM = 10;
export const MAX_CM = 60;
export const MAX_PIECES = 20;
export const ROUND_SIZES = [18, 20, 22, 24, 26, 28];
export const RECT_SIZES: [number, number][] = [
    [20, 20],
    [20, 30],
    [24, 36],
    [25, 35],
    [30, 40],
];

const SERVINGS: IPortion = { mode: 'servings', pieceUnit: null, sizes: [], form: null };
export const portionOf = (r: IPortionRecipe): IPortion => r.portion ?? SERVINGS;

const cmOk = (n: unknown): n is number =>
    typeof n === 'number' && Number.isInteger(n) && n >= MIN_CM && n <= MAX_CM;

export const baseSize = (p: IPortion): IPieceSize | undefined => p.sizes.find((s) => s.isBase);

/** The recipe's form, or null when it isn't complete. */
export function recipeForm(p: IPortion): FormSize | null {
    const f = p.form;
    if (!f) return null;
    if (f.shape === 'round') return f.dCm ? { shape: 'round', d: f.dCm } : null;
    return f.aCm && f.bCm ? { shape: 'rect', a: f.aCm, b: f.bCm } : null;
}

/**
 * The mode the page works in. A draft in "sized" without a base size (or "form" without a form) behaves like
 * servings, so the preview never divides by nothing.
 */
export function effectiveMode(r: IPortionRecipe): PortionMode {
    const p = portionOf(r);
    if (p.mode === 'sized' && baseSize(p) && r.servings > 0) return 'sized';
    if (p.mode === 'form' && recipeForm(p)) return 'form';
    return 'servings';
}

export const isForm = (s: SizeChoice | null | undefined): s is FormSize =>
    !!s && 'shape' in s && (s.shape === 'round' || s.shape === 'rect');
export const isPieces = (s: SizeChoice | null | undefined): s is PieceChoice =>
    !!s && 'n' in s && 'd' in s;

/** The recipe's own size (null in servings mode). */
export function defaultSize(r: IPortionRecipe): SizeChoice | null {
    const p = portionOf(r);
    switch (effectiveMode(r)) {
        case 'sized':
            return { n: r.servings, d: baseSize(p)!.diameterCm };
        case 'form':
            return recipeForm(p);
        default:
            return null;
    }
}

/**
 * The size to use for a stored or linked choice. A diameter the author no longer offers falls back to the
 * recipe's (the number of pieces stays); a form out of 10–60 cm or of the wrong kind falls back to the recipe's.
 */
export function resolveSize(r: IPortionRecipe, stored: unknown): SizeChoice | null {
    const def = defaultSize(r);
    if (!def) return null;
    const s = stored as Partial<PieceChoice & { shape: string; a: number; b: number }> | null;
    if (isPieces(def)) {
        const n =
            s && Number.isInteger(s.n) && s.n! >= 1 && s.n! <= MAX_PIECES ? (s.n as number) : def.n;
        const d =
            s && portionOf(r).sizes.some((x) => x.diameterCm === s.d) ? (s.d as number) : def.d;
        return { n, d };
    }
    if (s?.shape === 'round' && cmOk(s.d)) return { shape: 'round', d: s.d };
    if (s?.shape === 'rect' && cmOk(s.a) && cmOk(s.b)) return { shape: 'rect', a: s.a, b: s.b };
    return def;
}

export const formArea = (f: FormSize) =>
    f.shape === 'round' ? Math.PI * (f.d / 2) ** 2 : f.a * f.b;

const doughOf = (p: IPortion, d: number) => p.sizes.find((s) => s.diameterCm === d)?.doughG ?? 0;

/**
 * Scale of every quantity against the recipe as written.
 * servings: chosen / recipe servings (as before); sized: n × dough(d) / (servings × dough(base));
 * form: area of the reader's form / area of the recipe's.
 */
export function portionScale(r: IPortionRecipe, servings: number, size: SizeChoice | null): number {
    const p = portionOf(r);
    switch (effectiveMode(r)) {
        case 'sized': {
            if (!isPieces(size)) return 1;
            const base = baseSize(p)!;
            return (size.n * doughOf(p, size.d)) / (r.servings * base.doughG);
        }
        case 'form': {
            const own = recipeForm(p)!;
            return isForm(size) ? formArea(size) / formArea(own) : 1;
        }
        default:
            return r.servings ? servings / r.servings : 1;
    }
}

export function sameSize(a: SizeChoice | null, b: SizeChoice | null): boolean {
    if (!a || !b) return a === b;
    if (isPieces(a) && isPieces(b)) return a.n === b.n && a.d === b.d;
    if (isForm(a) && isForm(b)) {
        if (a.shape === 'round' && b.shape === 'round') return a.d === b.d;
        if (a.shape === 'rect' && b.shape === 'rect')
            return (a.a === b.a && a.b === b.b) || (a.a === b.b && a.b === b.a);
    }
    return false;
}

// ─── Names ───

const PIECE_FORMS: [string, string, string] = ['sztuka', 'sztuki', 'sztuk'];

export const pieceForms = (r: IPortionRecipe): [string, string, string] =>
    portionOf(r).pieceUnit?.forms ?? PIECE_FORMS;

/** "placek" / "placki" / "placków" for n. */
export const pieceNoun = (r: IPortionRecipe, n: number) => {
    const [one, few, many] = pieceForms(r);
    return pluralPL(n, one, few, many);
};

/** "Placki": the counter label (plural). */
export const piecesLabel = (r: IPortionRecipe) => {
    const few = pieceForms(r)[1];
    return few.charAt(0).toUpperCase() + few.slice(1);
};

export const formName = (f: FormSize) =>
    f.shape === 'round' ? `tortownica Ø ${f.d} cm` : `blacha ${f.a} × ${f.b} cm`;
export const formShort = (f: FormSize) =>
    f.shape === 'round' ? `Ø ${f.d} cm` : `${f.a} × ${f.b} cm`;

/** "3 placki 25 cm" / "tortownica Ø 26 cm"; '' in servings mode. */
export function sizeLabel(r: IPortionRecipe, s: SizeChoice | null): string {
    if (isPieces(s)) return `${s.n} ${pieceNoun(r, s.n)} ${s.d} cm`;
    if (isForm(s)) return formName(s);
    return '';
}

/** "×1,17" (two decimals at most, Polish comma). */
export const fmtScale = (s: number) => '×' + String(Math.round(s * 100) / 100).replace('.', ',');

export function scaleWarning(s: number): string | null {
    if (s > 2) return 'Ponad 2× więcej niż w przepisie — rozważ podział na dwie formy.';
    if (s < 0.5) return 'Mniej niż połowa przepisu — ciasto może wyjść za cienkie.';
    return null;
}

// ─── URL (?sztuki=3&rozmiar=25, ?forma=o26 / p30x40) ───

export function sizeToQuery(s: SizeChoice | null): Record<string, string> {
    if (isPieces(s)) return { sztuki: String(s.n), rozmiar: String(s.d) };
    if (isForm(s)) return { forma: s.shape === 'round' ? `o${s.d}` : `p${s.a}x${s.b}` };
    return {};
}

type Query = Record<string, unknown>;
const first = (v: unknown) => String(Array.isArray(v) ? v[0] : (v ?? ''));

/** The size in a link; null = the link has none for this recipe's mode (take the stored one). */
export function parseSizeQuery(r: IPortionRecipe, q: Query): SizeChoice | null {
    const mode = effectiveMode(r);
    if (mode === 'sized') {
        const n = Number(first(q.sztuki));
        const d = Number(first(q.rozmiar));
        if (!first(q.sztuki) && !first(q.rozmiar)) return null;
        const def = defaultSize(r) as PieceChoice;
        return resolveSize(r, { n: first(q.sztuki) ? n : def.n, d: first(q.rozmiar) ? d : def.d });
    }
    if (mode === 'form') {
        const m = first(q.forma).match(/^(?:o(\d{2})|p(\d{2})x(\d{2}))$/);
        if (!m) return null;
        return resolveSize(
            r,
            m[1] ? { shape: 'round', d: +m[1] } : { shape: 'rect', a: +m[2]!, b: +m[3]! }
        );
    }
    return null;
}

// ─── Macros ───

/**
 * Macros per one piece the reader makes, against nutrition.perPortion (one base-size piece, one piece of the
 * recipe's form, or the whole recipe form when it has no pieces): sized → dough(d) / dough(base);
 * form → the area scale (the piece count stays, each piece gets bigger).
 */
export function pieceMacroFactor(r: IPortionRecipe, size: SizeChoice | null): number {
    const p = portionOf(r);
    const mode = effectiveMode(r);
    if (mode === 'sized' && isPieces(size)) return doughOf(p, size.d) / baseSize(p)!.doughG;
    if (mode === 'form') return portionScale(r, 0, size);
    return 1;
}

/** Pieces the reader ends up with: chosen pieces, or the form's pieces (1 = the form without pieces). */
export function pieceCount(r: IPortionRecipe, size: SizeChoice | null): number {
    if (isPieces(size)) return size.n;
    return portionOf(r).form?.pieces ?? 1;
}
