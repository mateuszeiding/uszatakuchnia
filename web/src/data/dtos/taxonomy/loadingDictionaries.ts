import type { AxisCode, IDictionaries, ITaxonomyAxis } from './DictionariesDto';

// Osie na czas ładowania formularza przepisu: liczba i długości wartości jak w słownikach (17.11),
// żeby chipy zawijały się w tyle samo rzędów. Etykiety chipów są w szkielecie niewidoczne.
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
    // typ dania ma kolor (kropka w chipie zmienia jego szerokość) — kropka w szkielecie niewidoczna
    values: values.map((v, i) => ({
        id: -100 * (ai + 1) - i,
        label: v,
        ...(code === 'dish' ? { color: 'brown' } : {}),
    })),
}));

/**
 * Słowniki na czas ładowania formularzy (szkielet w docelowym układzie): 14 alergenów UE jest stałych
 * (rozporządzenie 1169/2011), więc chipy alergenów mają od razu docelową liczbę i szerokości.
 * Reszta pusta — pola wyboru mają tę samą wysokość bez opcji.
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
