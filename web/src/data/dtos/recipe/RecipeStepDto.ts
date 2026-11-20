export class RecipeStepDto {
    stepNo: number = 0;
    section: string | null = null;
    text: string = '';
    /** Step to do ahead (marinade, soaking). */
    prep: boolean = false;
    /** How far ahead, e.g. "12 h", "dzień wcześniej". */
    prepLead: string | null = null;
    /** Hints below the step text (0–n). */
    hints: string[] = [];

    constructor(obj: Partial<RecipeStepDto>) {
        Object.assign(this, obj);
        this.hints = obj.hints ?? [];
    }
}
