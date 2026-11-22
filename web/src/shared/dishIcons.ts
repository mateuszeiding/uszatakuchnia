// Dish-type icons for the recipe card placeholder frame (Claude Design, "ikony-typow/*.svg", 21.11.2026).
// 64×64 line icons: stroke 2.4, round caps and joins, currentColor. Geometry copied 1:1 from the SVG files.

export type DishIconShape = { d: string } | { cx: number; cy: number; r: number };

export const DISH_ICONS = {
    sniadanie: [
        {
            d: 'M14 34c0-11 9-18 19-17 6-6 18-3 20 6 6 3 6 13 1 17-1 8-10 11-17 8-7 4-17 1-20-6-3-1-3-6-3-8z',
        },
        { cx: 32, cy: 32, r: 7 },
    ],
    zupa: [
        { d: 'M8 30h48c0 13-10 22-24 22S8 43 8 30z' },
        { d: 'M24 56h16' },
        { d: 'M24 22c-3-4 3-6 0-10' },
        { d: 'M32 22c-3-4 3-6 0-10' },
        { d: 'M40 22c-3-4 3-6 0-10' },
    ],
    glowne: [
        { cx: 32, cy: 34, r: 16 },
        { cx: 32, cy: 34, r: 10 },
        { d: 'M8 12v12a3 3 0 006 0V12' },
        { d: 'M11 24v30' },
        { d: 'M56 12c-4 3-5 9-5 15h5' },
        { d: 'M56 12v42' },
    ],
    salatka: [
        { d: 'M8 34h48c0 12-10 20-24 20S8 46 8 34z' },
        { d: 'M22 34c-2-9 3-16 11-18' },
        { d: 'M34 34c1-11 9-17 17-16-1 8-6 14-14 16' },
        { d: 'M27 26c4 2 6 5 7 8' },
    ],
    przekaska: [
        { d: 'M10 54L54 10' },
        { cx: 22, cy: 42, r: 6 },
        { d: 'M27 32l5-5 5 5-5 5z' },
        { cx: 42, cy: 22, r: 6 },
    ],
    deser: [
        { d: 'M18 34h28l-4 20H22z' },
        { d: 'M16 34c0-9 7-15 16-15s16 6 16 15' },
        { d: 'M32 19v-5' },
        { cx: 32, cy: 11, r: 3 },
    ],
    wypiek: [
        { d: 'M10 46c0-14 10-24 22-24s22 10 22 24v6H10z' },
        { d: 'M22 32l4 6' },
        { d: 'M30 29l4 6' },
        { d: 'M38 30l4 6' },
    ],
    sos: [
        { d: 'M18 14h28v8H18z' },
        {
            d: 'M20 22h24v3c4 3 6 7 6 12v15a4 4 0 01-4 4H18a4 4 0 01-4-4V37c0-5 2-9 6-12z',
        },
        { d: 'M20 42h24' },
    ],
    napoj: [{ d: 'M18 16h28l-4 40H22z' }, { d: 'M20 30h24' }, { d: 'M36 24l6-16h8' }],
} satisfies Record<string, DishIconShape[]>;

export type DishIconKey = keyof typeof DISH_ICONS;

export const isDishIconKey = (k: string | null | undefined): k is DishIconKey =>
    !!k && Object.prototype.hasOwnProperty.call(DISH_ICONS, k);

// The icon is stored on the taxonomy value (taxonomy_values.icon, set by migration 0013 for today's dish types).
// An unknown or missing key = no icon: the frame shows only the type's color.
export const dishIconOf = (v: { icon?: string | null } | null): DishIconKey | null =>
    isDishIconKey(v?.icon) ? v.icon : null;
