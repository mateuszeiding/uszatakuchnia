import { describe, expect, it } from 'vitest';

import {
    losParam,
    paragraphs,
    parseLos,
    plural,
    rangeText,
    rollChallenge,
    timeLimitText,
} from '../challenges';

describe('rangeText', () => {
    it('opisuje wybór trudności', () => {
        expect(rangeText([])).toBe('dowolna');
        expect(rangeText([3])).toBe('3');
        expect(rangeText([3, 2])).toBe('2–3');
        expect(rangeText([1, 2, 3, 4, 5])).toBe('1–5');
        expect(rangeText([4, 1])).toBe('1, 4');
        expect(rangeText([1, 2, 4])).toBe('1, 2, 4');
    });
});

describe('parseLos / losParam', () => {
    it('czyta parametr z adresu', () => {
        expect(parseLos(undefined)).toBeNull();
        expect(parseLos('')).toEqual([]);
        expect(parseLos('3,2,2,9,x')).toEqual([2, 3]);
        expect(parseLos(['1,5'])).toEqual([1, 5]);
        expect(losParam([3, 1, 3])).toBe('1,3');
        expect(losParam([])).toBe('');
    });
});

describe('timeLimitText', () => {
    it('formatuje limit', () => {
        expect(timeLimitText(null)).toBe('bez limitu');
        expect(timeLimitText(45)).toBe('45 min');
        expect(timeLimitText(75)).toBe('1 h 15 min');
        expect(timeLimitText(120)).toBe('2 h');
    });
});

describe('rollChallenge', () => {
    const items = [
        { id: 1, difficulty: 1 },
        { id: 2, difficulty: 2 },
        { id: 3, difficulty: 2 },
        { id: 4, difficulty: 3 },
    ];
    const none = new Set<number>();

    it('losuje z wybranych trudności, pomija zrobione', () => {
        expect(rollChallenge(items, [2], new Set([2]), { random: () => 0 })).toEqual({
            kind: 'ok',
            id: 3,
        });
        expect(rollChallenge(items, [], none, { random: () => 0.99 })).toEqual({
            kind: 'ok',
            id: 4,
        });
    });
    it('brak wyzwań o trudności', () => {
        expect(rollChallenge(items, [5], none)).toEqual({ kind: 'none' });
    });
    it('same zrobione → alldone, potem losowanie spośród zrobionych', () => {
        expect(rollChallenge(items, [2], new Set([2, 3]))).toEqual({ kind: 'alldone' });
        expect(
            rollChallenge(items, [2], new Set([2, 3]), { includeDone: true, random: () => 0 })
        ).toEqual({ kind: 'ok', id: 2 });
    });
    it('„Losuj jeszcze raz” nie losuje otwartego', () => {
        expect(rollChallenge(items, [2], none, { exclude: 2, random: () => 0 })).toEqual({
            kind: 'ok',
            id: 3,
        });
        expect(rollChallenge(items, [2], new Set([3]), { exclude: 2 })).toEqual({
            kind: 'onlycurrent',
        });
        expect(rollChallenge(items, [1], none, { exclude: 1 })).toEqual({ kind: 'onlycurrent' });
        expect(rollChallenge(items, [2], new Set([2, 3]), { exclude: 2 })).toEqual({
            kind: 'alldone',
        });
    });
});

describe('paragraphs / plural', () => {
    it('dzieli opis na akapity', () => {
        expect(paragraphs('A\n\nB\nC\n\n\n  \nD')).toEqual(['A', 'B\nC', 'D']);
    });
    it('odmienia', () => {
        const w = (n: number) => plural(n, 'wyzwanie', 'wyzwania', 'wyzwań');
        expect([1, 2, 5, 12, 22, 25].map(w)).toEqual([
            'wyzwanie',
            'wyzwania',
            'wyzwań',
            'wyzwań',
            'wyzwania',
            'wyzwań',
        ]);
    });
});
