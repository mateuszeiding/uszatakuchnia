export class RecipeStepDto {
    stepNo: number = 0;
    section: string | null = null;
    text: string = '';
    /** Krok do zrobienia wcześniej (marynata, namaczanie). */
    prep: boolean = false;
    /** Ile wcześniej, np. "12 h", "dzień wcześniej". */
    prepLead: string | null = null;

    constructor(obj: Partial<RecipeStepDto>) {
        Object.assign(this, obj);
    }
}
