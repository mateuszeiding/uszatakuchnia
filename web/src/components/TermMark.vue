<script lang="ts" setup>
import { computed, inject, ref } from 'vue';

import { COOK_TERMS, openTerm, showTerm } from '@/composables/useTermExplain';
import type { ITerm } from '@/shared/utils/stepRefs';

// A glossary term in step or hint text: the author's word with a dotted underline. A button, so it works from the
// keyboard; it opens on click only — on hover, explanations would pop up while reading, and touch has no hover.
const props = defineProps<{ term: ITerm; word: string }>();

const cook = inject(COOK_TERMS, false);
const el = ref<HTMLButtonElement>();
const open = computed(() => !!el.value && openTerm.value?.anchor === el.value);

// the step row is a click target in places; the term must not check it off
const onClick = (e: MouseEvent) => {
    e.stopPropagation();
    showTerm(props.term, el.value!, cook);
};
</script>

<template>
    <button
        ref="el"
        type="button"
        class="term"
        :class="{ 'term--cook': cook, 'is-open': open }"
        aria-haspopup="dialog"
        :aria-expanded="open"
        :aria-label="`${word} — wyjaśnienie`"
        @click="onClick"
        v-text="word"
    />
</template>

<style scoped>
.term {
    display: inline;
    padding: 0;
    margin: 0;
    border: 0;
    border-radius: 3px;
    background: none;
    font: inherit;
    color: inherit;
    letter-spacing: inherit;
    text-align: inherit;
    cursor: pointer;
    text-decoration-line: underline;
    text-decoration-style: dotted;
    text-decoration-color: var(--accent);
    text-decoration-thickness: 1.5px;
    text-underline-offset: 0.22em;
    box-decoration-break: clone;
    -webkit-box-decoration-break: clone;
}
.term--cook {
    text-decoration-thickness: 2px;
}
.term.is-open {
    background: var(--accent-soft);
}
</style>
