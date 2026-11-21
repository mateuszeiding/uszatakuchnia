<script lang="ts" setup>
import { computed, nextTick, onMounted, ref, useId, watch } from 'vue';

import { chipLabel, fold, type IRefRow } from '../stepRefsForm';

import IngRefMenu from './IngRefMenu.vue';
import IngRefPicker from './IngRefPicker.vue';

import { parseStep, type RefPart, refToken } from '@/shared/utils/stepRefs';

// Step text field with ingredient refs as atomic labels. A contenteditable of our own (the project has no
// editor library): the model stays the stored text ({{ing:<key>|part}}), the DOM is text nodes + one
// non-editable span per ref. The DOM is rebuilt only when the text changes from outside, so typing keeps the caret.
const model = defineModel<string>({ required: true });
const props = defineProps<{
    rows: IRefRow[];
    label: string;
    placeholder: string;
    describedby?: string;
    warn?: boolean;
    broken?: boolean;
}>();

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

function render(text: string) {
    if (!field.value) return;
    field.value.replaceChildren(
        ...parseStep(text).map((s) =>
            typeof s === 'string'
                ? document.createTextNode(s)
                : s.kind === 'ing'
                  ? chip(s.key, s.part)
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
onMounted(() => render(model.value));

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
    const found = findAt();
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
function siblingChip(dir: 'before' | 'after'): HTMLElement | null {
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
    return isChip(node) ? node : null;
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
        const el = siblingChip(e.key === 'Backspace' ? 'before' : 'after');
        if (el) {
            e.preventDefault();
            el.remove();
            sync();
        }
    }
}

function onInput() {
    sync();
    updateAt();
}

function onPaste(e: ClipboardEvent) {
    e.preventDefault();
    const text = (e.clipboardData?.getData('text/plain') ?? '').replace(/\r?\n/g, ' ');
    document.execCommand('insertText', false, text);
    // pasted refs ({{ing:…}} copied from another step) become labels; the caret goes back after the paste
    if (/\{\{ing:/.test(text)) {
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
}
function onSelection() {
    if (at.value) updateAt();
}
</script>

<template>
    <div
        ref="wrap"
        class="step-editor"
    >
        <div
            ref="field"
            class="textarea step-editor__field"
            :class="{ 'is-warn': warn, 'is-broken': broken }"
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
            aria-hidden="true"
        >
            {{ placeholder }}
        </span>
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
.step-editor__pop {
    position: absolute;
    z-index: 20;
    outline: none;
}
.step-editor__pop > :deep(*) {
    position: static;
}
</style>
