<script lang="ts">
import { defineComponent, h, type PropType, type VNodeChild } from 'vue';

import SwapIcon from '@/components/SwapIcon.vue';
import { type IRefContext, viewSegments } from '@/shared/utils/stepRefs';

// Step text with ingredient quantities. The author writes the word ("dodaj mąkę"), the ref adds the quantity
// after it — a plain inline label, deliberately without hover or click (it must not compete with checking off
// the step). A render function, so no template whitespace ends up inside the label. The label wraps only at
// its " · " separators (a long "1 sztuka · do posmarowania" in 26 px is wider than a phone), never inside a part.
export default defineComponent({
    name: 'StepText',
    props: {
        text: { type: String, required: true },
        ctx: { type: Object as PropType<IRefContext>, required: true },
    },
    setup(props) {
        return () =>
            viewSegments(props.text, props.ctx).flatMap((v): VNodeChild[] => {
                if (typeof v === 'string') return [v];
                const part = (c: VNodeChild[], cls = '') =>
                    h('span', { class: `ing-ref__part ${cls}`.trim() }, c);
                const label = h('span', { class: 'ing-ref', title: v.title }, [
                    ...(v.lead
                        ? [part([h(SwapIcon, { class: 'ing-ref__swap' }), ` ${v.lead}`]), ' · ']
                        : []),
                    part([v.qty]),
                    ...(v.group ? [' ', part([`· ${v.group}`], 'ing-ref__group')] : []),
                ]);
                const opt = v.optional
                    ? [h('em', { class: 'ing-ref__opt' }, ' (opcjonalnie)')]
                    : [];
                return [' ', label, ...opt];
            });
    },
});
</script>

<style scoped>
.ing-ref {
    font-family: var(--font-mono);
    font-size: 0.88em;
    font-weight: 600;
    line-height: 1;
    letter-spacing: -0.2px;
    color: var(--ref-ink);
    background: var(--ref-bg);
    border-radius: 4px;
    padding: 0.1em 0.34em;
    box-decoration-break: clone;
    -webkit-box-decoration-break: clone;
}
.ing-ref__part {
    white-space: nowrap;
}
/* the base reset makes svg a block */
.ing-ref__swap {
    display: inline;
    width: 0.9em;
    height: 0.9em;
    vertical-align: -0.1em;
}
.ing-ref__group {
    font-weight: 500;
    opacity: 0.85;
}
.ing-ref__opt {
    color: var(--ink-muted);
    font-size: 0.9em;
}
</style>
