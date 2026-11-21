import { pluralPL } from '@/shared/utils/quantity';

/** "N z M" and the bar: mandatory steps only (an optional one can be checked, but counts for nothing). */
export function stepProgress(steps: { stepNo: number; optional: boolean }[], done: number[]) {
    const mandatory = steps.filter((s) => !s.optional);
    return {
        done: mandatory.filter((s) => done.includes(s.stepNo)).length,
        total: mandatory.length,
    };
}

/**
 * Cooking mode resumes right after the last checked step, so a skipped optional step isn't offered again;
 * past the end, at the first unchecked mandatory step. -1 = nothing checked.
 */
export function resumeIndex(
    steps: { stepNo: number; optional: boolean }[],
    done: number[]
): number {
    let last = -1;
    steps.forEach((s, i) => {
        if (done.includes(s.stepNo)) last = i;
    });
    if (last < 0) return -1;
    if (last + 1 < steps.length) return last + 1;
    const open = steps.findIndex((s) => !s.optional && !done.includes(s.stepNo));
    return open >= 0 ? open : steps.length - 1;
}

/** Toast after "Odznacz wszystko": "Odznaczono 3 kroki i 5 składników."; when one count is 0, only the other. */
export function uncheckedText(steps: number, ings: number): string {
    const parts = [
        steps ? `${steps} ${pluralPL(steps, 'krok', 'kroki', 'kroków')}` : '',
        ings ? `${ings} ${pluralPL(ings, 'składnik', 'składniki', 'składników')}` : '',
    ].filter(Boolean);
    return `Odznaczono ${parts.join(' i ')}.`;
}
