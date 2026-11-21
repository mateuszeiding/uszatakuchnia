export class RecipeStepDto {
    stepNo: number = 0;
    section: string | null = null;
    text: string = '';
    /** Step to do ahead (marinade, soaking). */
    prep: boolean = false;
    prepLead: string | null = null;
    /** The cook may skip it: progress and "needs prep" count only the other steps. */
    optional: boolean = false;
    hints: string[] = [];

    constructor(obj: Partial<RecipeStepDto>) {
        Object.assign(this, obj);
        this.hints = obj.hints ?? [];
    }
}
