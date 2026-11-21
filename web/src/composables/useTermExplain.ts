import { type InjectionKey, shallowRef } from 'vue';

import type { ITerm } from '@/shared/utils/stepRefs';

/** Where the explanation opens: a popover by the word (desktop), a bottom sheet (phone), the big sheet (cooking mode). */
export type ExplainMode = 'popover' | 'sheet' | 'cook';

export interface IOpenTerm {
    term: ITerm;
    anchor: HTMLElement;
    mode: ExplainMode;
}

/** Cooking mode provides true: terms inside it open the big sheet. */
export const COOK_TERMS: InjectionKey<boolean> = Symbol('cook-terms');

// One explanation open at a time, app-wide: one host (TermExplainHost in App.vue) renders it.
export const openTerm = shallowRef<IOpenTerm | null>(null);

export function showTerm(term: ITerm, anchor: HTMLElement, cook: boolean) {
    if (openTerm.value?.anchor === anchor) {
        closeTerm();
        return;
    }
    const narrow = window.matchMedia('(max-width: 900px)').matches;
    openTerm.value = { term, anchor, mode: cook ? 'cook' : narrow ? 'sheet' : 'popover' };
}

/** refocus: Esc, ✕ and "Zamknij" return focus to the word; a click elsewhere leaves it where it went. */
export function closeTerm(refocus = false) {
    const anchor = openTerm.value?.anchor;
    openTerm.value = null;
    if (refocus && anchor?.isConnected) anchor.focus();
}
