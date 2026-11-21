<script lang="ts" setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref, useId, watch } from 'vue';

import { chipLabel, fold, type IRefRow } from '../stepRefsForm';

import IngRefMenu from './IngRefMenu.vue';
import IngRefPicker from './IngRefPicker.vue';
import NewTermModal from './NewTermModal.vue';
import TermMarkMenu from './TermMarkMenu.vue';
import TermPicker from './TermPicker.vue';

import BookIcon from '@/components/BookIcon.vue';
import { type ITerm, parseStep, type RefPart, refToken, termToken } from '@/shared/utils/stepRefs';

// Step (or hint) text field with ingredient refs and glossary terms as atomic elements. A contenteditable of our
// own (the project has no editor library): the model stays the stored text ({{ing:<key>|part}}, {{term:<id>|word}}),
// the DOM is text nodes + one non-editable span per ref. The DOM is rebuilt only when the text changes from
// outside, so typing keeps the caret. A marked word is atomic too: to fix a typo in it, remove the mark first.
const model = defineModel<string>({ required: true });
const props = withDefaults(
    defineProps<{
        rows: IRefRow[];
        terms: ITerm[];
        label: string;
        placeholder: string;
        describedby?: string;
        warn?: boolean;
        broken?: boolean;
        /** Hints: no "@" (quantities live in the step text only), terms only. */
        ingRefs?: boolean;
        /** Smaller field (hints). */
        compact?: boolean;
        /** id of the editable field (focused from outside, e.g. after "+ wskazówka"). */
        fieldId?: string;
    }>(),
    {
        describedby: undefined,
        warn: false,
        broken: false,
        ingRefs: true,
        compact: false,
        fieldId: undefined,
    }
);
const emit = defineEmits<{ 'term-created': [term: ITerm] }>();

const field = ref<HTMLElement>();
const wrap = ref<HTMLElement>();
const listId = useId();
let emitted: string | null = null;

const byKey = computed(() => new Map(props.rows.map((r) => [r.key, r])));

const WARN_SVG =
    '<svg width="12" height="12" viewBox="0 0 16 16" fill="none" aria-hidden="true"><path d="M8 2.2L14.4 13.3H1.6L8 2.2z" stroke="currentColor" stroke-width="1.4" stroke-linejoin="round"/><path d="M8 6.4v3.3M8 11.6v.1" stroke="currentColor" stroke-width="1.6" stroke-linecap="round"/></svg>';

const isChip = (n: Node | null | undefined): n is HTMLElement =>
    n instanceof HTMLElement && n.classList.contains('ref-chip');
const chipPart = (el: HTMLElement) => (el.dataset.part as RefPart | undefined) ?? null;

function paint(el: HTMLElement) {
    const row = byKey.value.get(el.dataset.key ?? '');
    el.classList.toggle('is-broken', !row);
    if (row) el.textContent = chipLabel(row, chipPart(el));
    else el.innerHTML = `${WARN_SVG} brak składnika`;
    el.setAttribute(
        'aria-label',
        `${row ? chipLabel(row, chipPart(el)).slice(1) : 'brak składnika'}: zmień albo usuń`
    );
}

function chip(key: string, part: RefPart | null): HTMLElement {
    const el = document.createElement('span');
    el.className = 'ref-chip';
    el.contentEditable = 'false';
    // a keyboard target: Enter / Space opens the same menu as a click, Backspace / Delete removes it
    el.tabIndex = 0;
    el.setAttribute('role', 'button');
    el.dataset.key = key;
    if (part) el.dataset.part = part;
    paint(el);
    return el;
}

// ── glossary marks ──
const termById = computed(() => new Map(props.terms.map((t) => [t.id, t])));
const isMark = (n: Node | null | undefined): n is HTMLElement =>
    n instanceof HTMLElement && n.classList.contains('term-mark');
const markWord = (el: HTMLElement) => el.dataset.word ?? '';

function paintMark(el: HTMLElement) {
    const t = termById.value.get(Number(el.dataset.id));
    el.classList.toggle('is-missing', !t);
    el.textContent = markWord(el);
    el.title = t ? t.name : 'Pojęcia nie ma w słowniczku';
    el.setAttribute(
        'aria-label',
        `${markWord(el)}: ${t ? `pojęcie „${t.name}”` : 'pojęcia nie ma w słowniczku'} — zmień albo usuń`
    );
}

function mark(id: number, word: string): HTMLElement {
    const el = document.createElement('span');
    el.className = 'term-mark';
    el.contentEditable = 'false';
    el.tabIndex = 0;
    el.setAttribute('role', 'button');
    el.dataset.id = String(id);
    el.dataset.word = word;
    paintMark(el);
    return el;
}

function render(text: string) {
    if (!field.value) return;
    cancelTermUi();
    field.value.replaceChildren(
        ...parseStep(text).map((s) =>
            typeof s === 'string'
                ? document.createTextNode(s)
                : s.kind === 'ing'
                  ? chip(s.key, s.part)
                  : s.kind === 'term'
                    ? mark(s.id, s.word)
                    : document.createTextNode(s.raw)
        )
    );
}

function serialize(node: Node): string {
    let out = '';
    node.childNodes.forEach((n, i) => {
        if (n.nodeType === Node.TEXT_NODE)
            out += (n.nodeValue ?? '').replace(/\u00a0/g, ' ').replace(/\u200b/g, '');
        else if (isChip(n)) out += refToken(n.dataset.key ?? '', chipPart(n));
        else if (isMark(n)) out += termToken(Number(n.dataset.id), markWord(n));
        else if (n instanceof HTMLBRElement) out += i === node.childNodes.length - 1 ? '' : '\n';
        else out += serialize(n);
    });
    return out;
}

function sync() {
    if (!field.value) return;
    emitted = serialize(field.value);
    model.value = emitted;
}

watch(model, (v) => {
    if (v !== emitted) render(v);
});
watch(byKey, () => field.value?.querySelectorAll<HTMLElement>('.ref-chip').forEach(paint));
watch(termById, () => field.value?.querySelectorAll<HTMLElement>('.term-mark').forEach(paintMark));
onMounted(() => {
    render(model.value);
    document.addEventListener('selectionchange', onDocSelection);
});
onBeforeUnmount(() => document.removeEventListener('selectionchange', onDocSelection));

// ── "@" picker ──
interface IAt {
    node: Text;
    start: number;
    end: number;
    q: string;
}
const at = ref<IAt | null>(null);
/** "Zmień składnik" on an existing ref (no typed query). */
const changing = ref<HTMLElement | null>(null);
const active = ref(0);
const part = ref<RefPart | null>(null);
const pos = ref({ top: 0, left: 0 });
/** "@" closed with Esc: don't reopen until the caret leaves it. */
let dismissed: { node: Text; start: number } | null = null;

const pickerOpen = computed(() => !!at.value || !!changing.value);
const matches = computed(() => {
    const q = fold(at.value?.q ?? '').trim();
    return q ? props.rows.filter((r) => fold(r.name).includes(q)) : props.rows;
});

function place(rect: DOMRect) {
    const box = wrap.value!.getBoundingClientRect();
    pos.value = {
        top: rect.bottom - box.top + 6,
        left: Math.max(0, Math.min(rect.left - box.left, box.width - 380)),
    };
}

function findAt(): IAt | null {
    const sel = getSelection();
    if (!sel?.rangeCount || !sel.isCollapsed) return null;
    const { startContainer: n, startOffset: off } = sel.getRangeAt(0);
    if (n.nodeType !== Node.TEXT_NODE || !field.value?.contains(n)) return null;
    const m = (n.nodeValue ?? '').slice(0, off).match(/(?:^|[\s(„"])@([^@\n]{0,30})$/);
    if (!m) return null;
    return { node: n as Text, start: off - m[1]!.length - 1, end: off, q: m[1]! };
}

function updateAt() {
    const found = props.ingRefs ? findAt() : null;
    const same = (x: { node: Text; start: number } | null) =>
        !!found && !!x && x.node === found.node && x.start === found.start;
    if (!same(dismissed)) dismissed = null;
    if (!found || dismissed) {
        at.value = null;
        return;
    }
    if (!at.value || at.value.start !== found.start || at.value.node !== found.node) {
        active.value = 0;
        part.value = null;
    }
    at.value = found;
    const r = document.createRange();
    r.setStart(found.node, found.start);
    r.setEnd(found.node, found.start + 1);
    place(r.getBoundingClientRect());
}

watch(
    matches,
    () => (active.value = Math.min(active.value, Math.max(0, matches.value.length - 1)))
);

/** Caret right after a ref. With nothing after it, a zero-width space holds the caret (browsers can't put it
 *  after a non-editable element at the end); serialize() drops it, so the text gets no stray space. */
function caretAfter(el: Node) {
    let next = el.nextSibling;
    let at = 0;
    if (!next || (next.nodeType === Node.TEXT_NODE && !next.nodeValue)) {
        next = document.createTextNode('​');
        el.parentNode!.insertBefore(next, el.nextSibling);
        at = 1;
    }
    const r = document.createRange();
    if (next.nodeType === Node.TEXT_NODE) r.setStart(next, at);
    else r.setStartAfter(el);
    r.collapse(true);
    const sel = getSelection()!;
    sel.removeAllRanges();
    sel.addRange(r);
}

function pick(key: string) {
    if (changing.value) {
        const el = changing.value;
        el.dataset.key = key;
        if (part.value) el.dataset.part = part.value;
        else delete el.dataset.part;
        paint(el);
        closePicker();
        field.value?.focus();
        caretAfter(el);
    } else if (at.value) {
        const { node, start, end } = at.value;
        const r = document.createRange();
        r.setStart(node, start);
        r.setEnd(node, end);
        r.deleteContents();
        const el = chip(key, part.value);
        r.insertNode(el);
        closePicker();
        caretAfter(el);
    }
    sync();
}

function closePicker() {
    at.value = null;
    changing.value = null;
}

/** Keys while the picker is open (the caret stays in the field, or the picker holds focus when changing a ref). */
function pickerKey(e: KeyboardEvent) {
    const n = matches.value.length;
    if (e.key === 'ArrowDown' || e.key === 'ArrowUp') {
        e.preventDefault();
        if (n) active.value = (active.value + (e.key === 'ArrowDown' ? 1 : n - 1)) % n;
    } else if (e.key === 'ArrowLeft' || e.key === 'ArrowRight') {
        e.preventDefault();
        const parts: (RefPart | null)[] = [null, '1/2', '1/3', '1/4', '2/3', '3/4', 'rest'];
        const i = parts.indexOf(part.value);
        part.value = parts[(i + (e.key === 'ArrowRight' ? 1 : parts.length - 1)) % parts.length]!;
    } else if (e.key === 'Enter' || e.key === 'Tab') {
        const row = matches.value[active.value];
        if (!row) return;
        e.preventDefault();
        pick(row.key);
    } else if (e.key === 'Escape') {
        e.preventDefault();
        e.stopPropagation();
        if (at.value) dismissed = { node: at.value.node, start: at.value.start };
        const el = changing.value;
        closePicker();
        if (el) {
            field.value?.focus();
            caretAfter(el);
        }
    }
}

// ── editing an existing ref ──
const menu = ref<HTMLElement | null>(null);
const menuRow = computed(() =>
    menu.value ? (byKey.value.get(menu.value.dataset.key ?? '') ?? null) : null
);

function openMenu(el: HTMLElement) {
    closePicker();
    menu.value = el;
    place(el.getBoundingClientRect());
}
function onClick(e: MouseEvent) {
    const el = (e.target as HTMLElement).closest('.ref-chip');
    if (isChip(el)) openMenu(el);
    const m = (e.target as HTMLElement).closest('.term-mark');
    if (isMark(m)) openTermMenu(m);
}
function closeMenu(refocus = true) {
    const el = menu.value;
    menu.value = null;
    if (refocus && el?.isConnected) {
        field.value?.focus();
        caretAfter(el);
    }
}
async function menuChange() {
    const el = menu.value!;
    menu.value = null;
    changing.value = el;
    part.value = chipPart(el);
    active.value = Math.max(
        0,
        props.rows.findIndex((r) => r.key === el.dataset.key)
    );
    await nextTick();
    wrap.value?.querySelector<HTMLElement>('.step-editor__pop')?.focus();
}
function menuPart(p: RefPart | null) {
    const el = menu.value!;
    if (p) el.dataset.part = p;
    else delete el.dataset.part;
    paint(el);
    sync();
    closeMenu();
}
function menuRemove() {
    const el = menu.value!;
    menu.value = null;
    removeChip(el);
}

// ── keyboard / clipboard ──
/** The atomic element (ref or marked word) right before / after the caret. */
function siblingAtom(dir: 'before' | 'after'): HTMLElement | null {
    const sel = getSelection();
    if (!sel?.rangeCount || !sel.isCollapsed) return null;
    const { startContainer: n, startOffset: off } = sel.getRangeAt(0);
    let node: Node | null = null;
    if (n.nodeType === Node.TEXT_NODE) {
        // the zero-width caret holder counts as nothing
        const v = n.nodeValue ?? '';
        const empty = (s: string) => !s.replace(/\u200b/g, '');
        if (dir === 'before' && empty(v.slice(0, off))) node = n.previousSibling;
        if (dir === 'after' && empty(v.slice(off))) node = n.nextSibling;
    } else node = n.childNodes[dir === 'before' ? off - 1 : off] ?? null;
    return isChip(node) || isMark(node) ? node : null;
}

function removeChip(el: HTMLElement) {
    const prev = el.previousSibling;
    el.remove();
    sync();
    field.value?.focus();
    if (prev && field.value?.contains(prev)) {
        const r = document.createRange();
        if (prev.nodeType === Node.TEXT_NODE) r.setStart(prev, prev.nodeValue?.length ?? 0);
        else r.setStartAfter(prev);
        getSelection()?.removeAllRanges();
        getSelection()?.addRange(r);
    }
}

function onKeydown(e: KeyboardEvent) {
    const target = e.target as HTMLElement;
    if ((e.ctrlKey || e.metaKey) && e.shiftKey && e.key.toLowerCase() === 'l') {
        e.preventDefault();
        const s = readSel();
        sel.value = s;
        if (s?.ok) openTermPick();
        else if (s) placeBar(s.range);
        return;
    }
    if (isMark(target)) {
        if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault();
            openTermMenu(target);
        } else if (e.key === 'Backspace' || e.key === 'Delete') {
            e.preventDefault();
            unwrapMark(target, e.key === 'Backspace');
        }
        return;
    }
    if (isChip(target)) {
        if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault();
            openMenu(target);
        } else if (e.key === 'Backspace' || e.key === 'Delete') {
            e.preventDefault();
            removeChip(target);
        }
        return;
    }
    if (pickerOpen.value && !changing.value) {
        pickerKey(e);
        if (e.defaultPrevented) return;
    }
    if (e.key === 'Enter') {
        e.preventDefault(); // one paragraph per step
        return;
    }
    if (e.key === 'Backspace' || e.key === 'Delete') {
        const el = siblingAtom(e.key === 'Backspace' ? 'before' : 'after');
        if (el?.classList.contains('term-mark')) {
            // the first Backspace only removes the mark: the author's word stays
            e.preventDefault();
            unwrapMark(el, e.key === 'Backspace');
        } else if (el) {
            e.preventDefault();
            el.remove();
            sync();
        }
    }
}

function onInput() {
    // an edit can move or remove the word a pending "Pojęcie…" would mark (a saved Range follows the DOM)
    // or the mark an open menu points at
    cancelTermUi();
    sync();
    updateAt();
}

function onPaste(e: ClipboardEvent) {
    e.preventDefault();
    const text = (e.clipboardData?.getData('text/plain') ?? '').replace(/\r?\n/g, ' ');
    document.execCommand('insertText', false, text);
    // pasted refs ({{ing:…}} copied from another step) become labels; the caret goes back after the paste
    if (/\{\{(ing|term):/.test(text)) {
        const at = caretOffset();
        sync();
        // not model.value: it returns the parent's prop until the parent re-renders
        render(emitted ?? '');
        if (at != null) setCaretOffset(at);
    }
}

/** Caret position as an offset in the stored text (a ref counts as its {{…}} token). */
function caretOffset(): number | null {
    const sel = getSelection();
    if (!field.value || !sel?.rangeCount) return null;
    const r = document.createRange();
    r.setStart(field.value, 0);
    r.setEnd(sel.getRangeAt(0).endContainer, sel.getRangeAt(0).endOffset);
    const box = document.createElement('div');
    box.append(r.cloneContents());
    return serialize(box).length;
}

function setCaretOffset(offset: number) {
    let left = offset;
    const r = document.createRange();
    for (const n of [...(field.value?.childNodes ?? [])]) {
        const text = n.nodeType === Node.TEXT_NODE;
        const len = text
            ? (n.nodeValue ?? '').length
            : isChip(n)
              ? refToken(n.dataset.key ?? '', chipPart(n)).length
              : isMark(n)
                ? termToken(Number(n.dataset.id), markWord(n)).length
                : 0;
        if (text && left <= len) {
            r.setStart(n, left);
            break;
        }
        left -= len;
        r.setStartAfter(n);
    }
    r.collapse(true);
    getSelection()?.removeAllRanges();
    getSelection()?.addRange(r);
}

function onCopy(e: ClipboardEvent, cut = false) {
    const sel = getSelection();
    if (!sel?.rangeCount || sel.isCollapsed) return;
    const frag = document.createElement('div');
    frag.append(sel.getRangeAt(0).cloneContents());
    e.preventDefault();
    e.clipboardData?.setData('text/plain', serialize(frag));
    if (cut) {
        sel.getRangeAt(0).deleteContents();
        sync();
    }
}

function onBlur(e: FocusEvent) {
    if (wrap.value?.contains(e.relatedTarget as Node)) return;
    at.value = null;
    sel.value = null;
}
function onSelection(e?: Event) {
    if (at.value) updateAt();
    // Tab / Shift move focus to the "Pojęcie…" bar: keep it
    if (e instanceof KeyboardEvent && (e.key === 'Tab' || e.key === 'Shift')) return;
    if (!termPick.value && !termMenu.value && !newTerm.value) {
        sel.value = readSel();
        if (sel.value) placeBar(sel.value.range);
    }
}

// ── marking a selected word as a term ──
interface ISel {
    range: Range;
    word: string;
    ok: boolean;
    why: string;
}
/** The "Pojęcie…" bar over a selection (or why the selection can't be marked). */
const sel = ref<ISel | null>(null);
const barPos = ref({ top: 0, left: 0 });

function readSel(): ISel | null {
    const s = getSelection();
    if (!s?.rangeCount || s.isCollapsed || !field.value) return null;
    const r = s.getRangeAt(0).cloneRange();
    if (!field.value.contains(r.commonAncestorContainer)) return null;
    const n = r.startContainer;
    const one = n === r.endContainer && n.nodeType === Node.TEXT_NODE;
    if (one) {
        // a double-click selects the word with the space after it
        const v = n.nodeValue ?? '';
        let a = r.startOffset;
        let b = r.endOffset;
        while (a < b && /[\s\u200b]/.test(v[a]!)) a++;
        while (b > a && /[\s\u200b]/.test(v[b - 1]!)) b--;
        r.setStart(n, a);
        r.setEnd(n, b);
    }
    const word = r.toString();
    if (one && !word) return null;
    const why = !one
        ? 'Zaznacz samo słowo — bez ilości i bez innego pojęcia.'
        : /[{}|]/.test(word)
          ? 'Zaznaczenie nie może zawierać znaków { } |.'
          : [...word].length > 60
            ? 'Zaznacz najwyżej 60 znaków.'
            : '';
    return { range: r, word, ok: !why, why };
}

// the selection left the field (a click elsewhere, typing in another field): the bar goes away
function onDocSelection() {
    if (!sel.value) return;
    const s = getSelection();
    if (
        !s?.rangeCount ||
        s.isCollapsed ||
        !field.value?.contains(s.getRangeAt(0).commonAncestorContainer)
    )
        sel.value = null;
}

function placeBar(r: Range) {
    const rect = r.getBoundingClientRect();
    const box = wrap.value!.getBoundingClientRect();
    barPos.value = {
        top: rect.top - box.top - 40,
        left: Math.max(0, Math.min(rect.left - box.left, box.width - 220)),
    };
}

interface ITermPick {
    word: string;
    /** New mark: the selection to wrap. */
    range?: Range;
    /** "Zmień pojęcie" on an existing mark. */
    el?: HTMLElement;
    current: number | null;
}
const termPick = ref<ITermPick | null>(null);
const newTerm = ref<ITermPick | null>(null);
const termMenu = ref<HTMLElement | null>(null);
const termMenuTerm = computed(() =>
    termMenu.value ? (termById.value.get(Number(termMenu.value.dataset.id)) ?? null) : null
);

function openTermPick() {
    const s = sel.value;
    if (!s?.ok) return;
    closePicker();
    menu.value = null;
    termPick.value = { word: s.word, range: s.range, current: null };
    place(s.range.getBoundingClientRect());
    sel.value = null;
}

function caretAtEnd(node: Text, end: boolean) {
    const r = document.createRange();
    r.setStart(node, end ? (node.nodeValue?.length ?? 0) : 0);
    r.collapse(true);
    getSelection()?.removeAllRanges();
    getSelection()?.addRange(r);
}

function applyTerm(p: ITermPick, id: number) {
    let el = p.el;
    if (el) {
        el.dataset.id = String(id);
        paintMark(el);
    } else if (p.range) {
        p.range.deleteContents();
        el = mark(id, p.word);
        p.range.insertNode(el);
    }
    field.value?.focus();
    if (el) caretAfter(el);
    sync();
}

function pickTerm(id: number) {
    const p = termPick.value!;
    termPick.value = null;
    applyTerm(p, id);
}

function closeTermPick() {
    const p = termPick.value;
    termPick.value = null;
    field.value?.focus();
    if (p?.el?.isConnected) caretAfter(p.el);
    else if (p?.range) {
        getSelection()?.removeAllRanges();
        getSelection()?.addRange(p.range);
    }
}

function createTerm(name: string) {
    newTerm.value = { ...termPick.value!, word: termPick.value!.word };
    newName.value = name;
    termPick.value = null;
}
const newName = ref('');

function termCreated(t: ITerm) {
    const p = newTerm.value!;
    newTerm.value = null;
    emit('term-created', t);
    applyTerm(p, t.id);
}
function newTermClosed() {
    const p = newTerm.value;
    newTerm.value = null;
    field.value?.focus();
    if (p?.range) {
        getSelection()?.removeAllRanges();
        getSelection()?.addRange(p.range);
    }
}

function openTermMenu(el: HTMLElement) {
    closePicker();
    menu.value = null;
    termPick.value = null;
    sel.value = null;
    termMenu.value = el;
    place(el.getBoundingClientRect());
}
function closeTermMenu(refocus = true) {
    const el = termMenu.value;
    termMenu.value = null;
    if (refocus && el?.isConnected) {
        field.value?.focus();
        caretAfter(el);
    }
}
function termMenuChange() {
    const el = termMenu.value!;
    termMenu.value = null;
    termPick.value = { word: markWord(el), el, current: Number(el.dataset.id) };
    place(el.getBoundingClientRect());
}

/** The text was replaced from outside: a pending pick or an open menu would act on nodes that are gone. */
function cancelTermUi() {
    termPick.value = null;
    termMenu.value = null;
    sel.value = null;
}

/** Removes the mark, the word stays as plain text; the caret goes to its end (Backspace) or start (Delete). */
function unwrapMark(el: HTMLElement, caretEnd = true) {
    const t = document.createTextNode(markWord(el));
    el.replaceWith(t);
    termMenu.value = null;
    field.value?.focus();
    caretAtEnd(t, caretEnd);
    sync();
}
</script>

<template>
    <div
        ref="wrap"
        class="step-editor"
    >
        <div
            :id="fieldId"
            ref="field"
            class="textarea step-editor__field"
            :class="{ 'is-warn': warn, 'is-broken': broken, 'is-compact': compact }"
            contenteditable="true"
            role="textbox"
            aria-multiline="true"
            :aria-label="label"
            :aria-describedby="describedby"
            :aria-controls="pickerOpen ? listId : undefined"
            :aria-activedescendant="
                pickerOpen && matches.length ? `${listId}-o${active}` : undefined
            "
            spellcheck="true"
            @input="onInput"
            @keydown="onKeydown"
            @keyup="onSelection"
            @click="onClick"
            @mouseup="onSelection"
            @paste="onPaste"
            @copy="onCopy"
            @cut="onCopy($event, true)"
            @drop.prevent
            @blur="onBlur"
        />
        <span
            v-if="!model.trim()"
            class="step-editor__ph"
            :class="{ 'is-compact': compact }"
            aria-hidden="true"
        >
            {{ placeholder }}
        </span>
        <div
            v-if="sel && !termPick && !termMenu && !newTerm && !pickerOpen && !menu"
            class="term-bar"
            :style="{ top: `${barPos.top}px`, left: `${barPos.left}px` }"
        >
            <button
                v-if="sel.ok"
                type="button"
                class="term-bar__btn"
                title="Oznacz pojęcie (Ctrl+Shift+L)"
                aria-keyshortcuts="Control+Shift+L"
                @mousedown.prevent
                @click="openTermPick"
            >
                <BookIcon :size="13" />
                Pojęcie…
            </button>
            <span
                v-else
                class="term-bar__why"
            >
                {{ sel.why }}
            </span>
        </div>
        <div
            v-if="termPick"
            class="step-editor__pop"
            :style="{ top: `${pos.top}px`, left: `${pos.left}px` }"
            @focusout="!wrap?.contains($event.relatedTarget as Node) && (termPick = null)"
        >
            <TermPicker
                :terms="terms"
                :word="termPick.word"
                :current="termPick.current"
                @pick="pickTerm"
                @create="createTerm"
                @close="closeTermPick"
            />
        </div>
        <div
            v-if="termMenu"
            class="step-editor__pop"
            :style="{ top: `${pos.top}px`, left: `${pos.left}px` }"
            @focusout="!wrap?.contains($event.relatedTarget as Node) && closeTermMenu(false)"
        >
            <TermMarkMenu
                :term="termMenuTerm"
                :word="termMenu.dataset.word ?? ''"
                @change="termMenuChange"
                @remove="unwrapMark(termMenu!)"
                @close="closeTermMenu()"
            />
        </div>
        <NewTermModal
            v-if="newTerm"
            :word="newName"
            @created="termCreated"
            @close="newTermClosed"
        />
        <div
            v-if="pickerOpen"
            class="step-editor__pop"
            :style="{ top: `${pos.top}px`, left: `${pos.left}px` }"
            tabindex="-1"
            @keydown="changing && pickerKey($event)"
            @focusout="changing && !wrap?.contains($event.relatedTarget as Node) && closePicker()"
        >
            <IngRefPicker
                v-model:active="active"
                v-model:part="part"
                :list-id="listId"
                :rows="matches"
                :query="at?.q ?? ''"
                @pick="pick"
            />
        </div>
        <div
            v-if="menu"
            class="step-editor__pop"
            :style="{ top: `${pos.top}px`, left: `${pos.left}px` }"
            @focusout="!wrap?.contains($event.relatedTarget as Node) && closeMenu(false)"
        >
            <IngRefMenu
                :row="menuRow"
                :part="(menu.dataset.part as RefPart | undefined) ?? null"
                @change="menuChange"
                @part="menuPart"
                @remove="menuRemove"
                @close="closeMenu()"
            />
        </div>
    </div>
</template>

<style scoped>
.step-editor {
    position: relative;
}
.step-editor__field {
    min-height: 64px;
    font-size: 15px;
    line-height: 1.7;
    white-space: pre-wrap;
    overflow-wrap: anywhere;
    cursor: text;
    resize: none;
}
/* over the field, at its padding (a ::before inside contenteditable would push the caret) */
.step-editor__ph {
    position: absolute;
    top: 11px;
    left: 13px;
    font-size: 15px;
    line-height: 1.7;
    color: var(--ink-faint);
    pointer-events: none;
}
.step-editor__field.is-warn {
    border-color: var(--prep-line);
}
.step-editor__field.is-broken {
    border-color: var(--danger);
}
.step-editor__field :deep(.ref-chip) {
    display: inline;
    white-space: nowrap;
    font-size: 13px;
    font-weight: 600;
    padding: 2px 6px;
    border-radius: 4px;
    background: var(--ref-bg);
    color: var(--ref-ink);
    border: 1px solid transparent;
    cursor: pointer;
    user-select: all;
    box-decoration-break: clone;
    -webkit-box-decoration-break: clone;
}
.step-editor__field :deep(.ref-chip.is-broken) {
    background: var(--danger-soft);
    color: var(--danger);
    border-color: var(--danger);
}
.step-editor__field :deep(.ref-chip svg) {
    display: inline; /* the base reset makes svg a block */
    vertical-align: -0.1em;
}
.step-editor__field.is-compact,
.step-editor__ph.is-compact {
    min-height: 40px;
    font-size: 13px;
    line-height: 1.6;
}
.step-editor__ph.is-compact {
    top: 10px;
}
.step-editor__field :deep(.term-mark) {
    display: inline;
    cursor: pointer;
    text-decoration-line: underline;
    text-decoration-style: dotted;
    text-decoration-color: var(--accent);
    text-decoration-thickness: 1.5px;
    text-underline-offset: 0.22em;
    border-radius: 3px;
    user-select: all;
}
.step-editor__field :deep(.term-mark.is-missing) {
    text-decoration-color: var(--danger);
    color: var(--danger);
}
.term-bar {
    position: absolute;
    z-index: 21;
    display: flex;
    padding: 3px;
    border-radius: var(--r-lg);
    background: var(--ink);
    box-shadow: 0 6px 18px rgba(10, 10, 15, 0.25);
}
.term-bar__btn {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 5px 10px;
    border: none;
    border-radius: var(--r-md);
    background: color-mix(in srgb, var(--bg) 14%, transparent);
    color: var(--bg);
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
}
.term-bar__why {
    padding: 5px 10px;
    font-size: 12.5px;
    color: var(--bg);
}
.step-editor__pop {
    position: absolute;
    z-index: 20;
    outline: none;
}
.step-editor__pop > :deep(*) {
    position: static;
}
</style>
