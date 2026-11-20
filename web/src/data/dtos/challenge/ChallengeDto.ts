/** Challenge ingredient item: plain text; required = "Musisz użyć" (must use). */
export interface IChallengeIngredient {
    name: string;
    required: boolean;
}

/** Challenge (the list and the details share one shape). */
export interface IChallenge {
    id: number;
    title: string;
    formatId: number;
    format: string;
    /** Paragraphs separated by a blank line. */
    description: string;
    /** null = no limit (informational only, no timer). */
    timeLimitMin: number | null;
    /** 1 warm-up … 5 masterclass. */
    difficulty: number;
    source: string | null;
    /** "Don't mess this up" notes in the author's order. */
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
    /** No empty items; missing field = unchanged. */
    cautions?: string[];
    ingredients: IChallengeIngredient[];
}

/** Format dictionary item; usedIn = number of challenges using it. */
export interface IChallengeFormat {
    id: number;
    name: string;
    usedIn: number;
}
