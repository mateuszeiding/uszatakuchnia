import type { IDictionaries } from './DictionariesDto';

/**
 * Słowniki na czas ładowania formularzy (szkielet w docelowym układzie): 14 alergenów UE jest stałych
 * (rozporządzenie 1169/2011), więc chipy alergenów mają od razu docelową liczbę i szerokości.
 * Reszta pusta — pola wyboru mają tę samą wysokość bez opcji.
 */
export const LOADING_DICTIONARIES: IDictionaries = {
    axes: [],
    units: [],
    ingredientTypes: [],
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
