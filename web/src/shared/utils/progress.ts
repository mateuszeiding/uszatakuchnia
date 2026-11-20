import { pluralPL } from '@/shared/utils/quantity';

/** Toast after "Odznacz wszystko": "Odznaczono 3 kroki i 5 składników."; when one count is 0, only the other. */
export function uncheckedText(steps: number, ings: number): string {
    const parts = [
        steps ? `${steps} ${pluralPL(steps, 'krok', 'kroki', 'kroków')}` : '',
        ings ? `${ings} ${pluralPL(ings, 'składnik', 'składniki', 'składników')}` : '',
    ].filter(Boolean);
    return `Odznaczono ${parts.join(' i ')}.`;
}
