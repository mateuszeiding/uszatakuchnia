/** Pozycja składnika wyzwania — zwykły tekst; required = „Musisz użyć”. */
export interface IChallengeIngredient {
    name: string;
    required: boolean;
}

/** Wyzwanie (lista i szczegóły mają ten sam kształt). */
export interface IChallenge {
    id: number;
    title: string;
    formatId: number;
    format: string;
    /** Akapity rozdzielone pustą linią. */
    description: string;
    /** null = bez limitu (tylko informacja, bez zegara). */
    timeLimitMin: number | null;
    /** 1 rozgrzewka … 5 mistrzowskie. */
    difficulty: number;
    source: string | null;
    /** Uwagi „czego nie zawalić” w kolejności autora. */
    cautions: string[];
    ingredients: IChallengeIngredient[];
    createdAt: string;
    updatedAt: string;
}

export interface IUpsertChallengeRequest {
    title: string;
    formatId: number;
    description: string;
    timeLimitMin: number | null;
    difficulty: number;
    source: string | null;
    /** Bez pustych pozycji; brak pola = bez zmian. */
    cautions?: string[];
    ingredients: IChallengeIngredient[];
}

/** Pozycja słownika formatów; usedIn = w ilu wyzwaniach. */
export interface IChallengeFormat {
    id: number;
    name: string;
    usedIn: number;
}
