import { pluralPL } from '@/shared/utils/quantity';

/** Toast po „Odznacz wszystko”: „Odznaczono 3 kroki i 5 składników.”; gdy jedna liczba to 0 — tylko druga. */
export function uncheckedText(steps: number, ings: number): string {
    const parts = [
        steps ? `${steps} ${pluralPL(steps, 'krok', 'kroki', 'kroków')}` : '',
        ings ? `${ings} ${pluralPL(ings, 'składnik', 'składniki', 'składników')}` : '',
    ].filter(Boolean);
    return `Odznaczono ${parts.join(' i ')}.`;
}
