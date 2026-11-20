import type { AxisCode, IDictionaries, ITaxonomyAxis } from './DictionariesDto';

// Axes shown while the recipe form loads: value counts and lengths match the real dictionaries,
// so the chips wrap into the same number of rows. Chip labels are invisible in the skeleton.
const AXES: [AxisCode, string, boolean, string[]][] = [
    [
        'dish',
        'Typ dania',
        true,
        [
            'śniadanie',
            'zupa',
            'danie główne',
            'sałatka',
            'przekąska',
            'deser',
            'pieczywo i wypieki',
            'sos / dodatek',
            'napój',
        ],
    ],
    [
        'base',
        'Baza',
        false,
        [
            'mięso',
            'drób',
            'ryby',
            'owoce morza',
            'jajka',
            'nabiał',
            'warzywa',
            'strączki',
            'grzyby',
            'makaron',
            'ryż i kasze',
        ],
    ],
    [
        'cuisine',
        'Kuchnia',
        false,
        [
            'polska',
            'włoska',
            'francuska',
            'hiszpańska',
            'grecka',
            'bliskowschodnia',
            'indyjska',
            'azjatycka',
            'meksykańska',
            'amerykańska',
            'skandynawska',
        ],
    ],
    [
        'diet',
        'Diety',
        false,
        [
            'wegetariańskie',
            'wegańskie',
            'bez glutenu',
            'bez laktozy',
            'bez cukru',
            'wysokobiałkowe',
            'niskokaloryczne',
            'keto',
            'low fodmap',
        ],
    ],
    [
        'practical',
        'Praktyczne',
        false,
        [
            '1 patelnia',
            '1 garnek',
            'do 5 składników',
            'można mrozić',
            'na imprezę',
            'dla dzieci',
            'meal prep',
        ],
    ],
];
const loadingAxes: ITaxonomyAxis[] = AXES.map(([code, label, required, values], ai) => ({
    id: -1 - ai,
    code,
    label,
    multi: code !== 'cuisine',
    required,
    hasColor: code === 'dish',
    filterMode: code === 'diet' ? 'and' : 'or',
    // the dish type has a color (the dot changes the chip's width); the dot is invisible in the skeleton
    values: values.map((v, i) => ({
        id: -100 * (ai + 1) - i,
        label: v,
        ...(code === 'dish' ? { color: 'brown' } : {}),
    })),
}));

/**
 * Dictionaries while the forms load (skeleton in the final layout): the 14 EU allergens are fixed
 * (Regulation 1169/2011), so the allergen chips get their final count and widths right away.
 * The rest is empty: select fields have the same height without options.
 */
export const LOADING_DICTIONARIES: IDictionaries = {
    axes: loadingAxes,
    units: [],
    ingredientTypes: [],
    equipment: [],
    challengeFormats: [],
    allergens: [
        'gluten',
        'mleko',
        'jaja',
        'ryby',
        'skorupiaki',
        'mięczaki',
        'orzeszki ziemne',
        'orzechy',
        'soja',
        'seler',
        'gorczyca',
        'sezam',
        'siarczyny',
        'łubin',
    ].map((label, i) => ({ code: `loading-${i}`, label })),
};
