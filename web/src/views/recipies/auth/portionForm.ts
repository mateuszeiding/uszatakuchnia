// "Czym jest porcja" in the recipe form. Fields of every mode stay in the form and go to the API together,
// so switching modes back and forth loses nothing. Validation mirrors services/portion.go (publish only here;
// the API also checks drafts).

import type { IUpsertRecipeRequest, RecipeDto } from '@/data/dtos/recipe/RecipeDto';
import type { IUnit } from '@/data/dtos/taxonomy/DictionariesDto';
import { MAX_CM, MIN_CM, type PortionMode } from '@/shared/utils/portion';

let seq = 0;

export interface ISizeRow {
    id: number;
    d: string;
    g: string;
    base: boolean;
}

export interface IPortionForm {
    mode: PortionMode;
    pieceUnitId: number | null;
    sizes: ISizeRow[];
    shape: 'round' | 'rect';
    d: string;
    a: string;
    b: string;
    pieces: string;
}

export const blankSize = (base = false): ISizeRow => ({ id: ++seq, d: '', g: '', base });

export const newPortion = (): IPortionForm => ({
    mode: 'servings',
    pieceUnitId: null,
    sizes: [],
    shape: 'round',
    d: '',
    a: '',
    b: '',
    pieces: '',
});

const str = (n: number | null | undefined) => (n ? String(n) : '');

export function portionFromRecipe(r: RecipeDto): IPortionForm {
    const p = r.portion;
    if (!p) return newPortion();
    return {
        mode: p.mode,
        pieceUnitId: p.pieceUnit?.id ?? null,
        sizes: p.sizes.map((s) => ({
            id: ++seq,
            d: String(s.diameterCm),
            g: String(s.doughG),
            base: s.isBase,
        })),
        shape: p.form?.shape ?? 'round',
        d: str(p.form?.dCm),
        a: str(p.form?.aCm),
        b: str(p.form?.bCm),
        pieces: str(p.form?.pieces),
    };
}

const int = (s: string) => {
    const n = Number(s.trim());
    return s.trim() && Number.isInteger(n) ? n : null;
};

export function portionToRequest(p: IPortionForm): NonNullable<IUpsertRecipeRequest['portion']> {
    const d = int(p.d);
    const a = int(p.a);
    const b = int(p.b);
    const pieces = int(p.pieces);
    const hasForm = p.shape === 'round' ? d != null : a != null || b != null;
    return {
        mode: p.mode,
        pieceUnitId: p.pieceUnitId,
        sizes: p.sizes
            .filter((s) => s.d.trim() || s.g.trim())
            .map((s) => ({ diameterCm: int(s.d) ?? 0, doughG: int(s.g) ?? 0, isBase: s.base })),
        form:
            hasForm || pieces != null
                ? {
                      shape: p.shape,
                      dCm: p.shape === 'round' ? d : null,
                      aCm: p.shape === 'rect' ? a : null,
                      bCm: p.shape === 'rect' ? b : null,
                      pieces,
                  }
                : null,
    };
}

const cmOk = (s: string) => {
    const n = int(s);
    return n != null && n >= MIN_CM && n <= MAX_CM;
};

/** Units the piece can be named with: inflected, with a quantity (placek / placki / placków). */
export const pieceUnits = (units: IUnit[]) =>
    units.filter((u) => u.forms?.length === 3 && !u.noQty);

/**
 * What blocks publishing (one message for the "portion" field); null = fine.
 * servingsOk: the recipe has a servings / pieces count.
 */
export function portionError(p: IPortionForm, servingsOk: boolean): string | null {
    if (p.mode === 'sized') {
        if (!servingsOk) return 'Podaj, ile sztuk jest w przepisie.';
        if (!p.pieceUnitId) return 'Wybierz nazwę sztuki.';
        const rows = p.sizes.filter((s) => s.d.trim() || s.g.trim());
        if (!rows.length) return 'Dodaj co najmniej jeden rozmiar.';
        if (rows.some((s) => !cmOk(s.d)))
            return `Średnica rozmiaru: pełne centymetry od ${MIN_CM} do ${MAX_CM}.`;
        if (rows.some((s) => !(int(s.g)! > 0))) return 'Podaj masę ciasta na sztukę.';
        const ds = rows.map((s) => int(s.d));
        if (new Set(ds).size !== ds.length) return 'Ten sam rozmiar jest na liście dwa razy.';
        if (rows.filter((s) => s.base).length !== 1) return 'Zaznacz jeden rozmiar „w przepisie”.';
    }
    if (p.mode === 'form') {
        if (p.shape === 'round' && !cmOk(p.d))
            return `Podaj średnicę formy: pełne centymetry od ${MIN_CM} do ${MAX_CM}.`;
        if (p.shape === 'rect' && !(cmOk(p.a) && cmOk(p.b)))
            return `Podaj długość i szerokość formy: pełne centymetry od ${MIN_CM} do ${MAX_CM}.`;
        if (p.pieces.trim() && !(int(p.pieces)! >= 1 && int(p.pieces)! <= 100))
            return 'Kawałków z formy: od 1 do 100.';
    }
    return null;
}
