<script lang="ts" setup>
import { nextTick, onMounted, ref } from 'vue';

import { useNoindex } from '@/composables/useNoindex';
import { setTitle } from '@/shared/utils/title';

// Error page: eyebrow, h1, description, buttons, an optional section below.
// Tab title = the heading, noindex, focus on h1 after load.
const props = defineProps<{
    kicker: string;
    title: string;
    /** Eyebrow in --prep (connection error), otherwise --ink-muted. */
    warn?: boolean;
}>();

const h1 = ref<HTMLElement>();
useNoindex();
onMounted(async () => {
    setTitle(props.title);
    await nextTick();
    h1.value?.focus({ preventScroll: true });
});
</script>

<template>
    <div class="container err">
        <div class="err__box">
            <div
                class="err__kicker"
                :class="{ 'err__kicker--warn': warn }"
            >
                {{ kicker }}
            </div>
            <h1
                ref="h1"
                class="err__title"
                tabindex="-1"
            >
                {{ title }}
            </h1>
            <p class="err__text"><slot /></p>
            <div class="err__actions"><slot name="actions" /></div>
            <slot name="more" />
        </div>
    </div>
</template>

<style scoped>
.err {
    padding-top: 88px;
}
.err__box {
    max-width: 560px;
}
.err__kicker {
    margin-bottom: 14px;
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 600;
    letter-spacing: 0.5px;
    text-transform: uppercase;
    color: var(--ink-muted);
}
.err__kicker--warn {
    color: var(--prep);
}
.err__title {
    margin: 0;
    font-size: 56px;
    line-height: 1.02;
    letter-spacing: -2px;
    font-weight: 500;
}
.err__title:focus {
    outline: none;
}
.err__text {
    margin: 16px 0 0;
    font-size: 18px;
    line-height: 1.55;
    color: var(--ink-dim);
}
.err__text :deep(code) {
    font-family: var(--font-mono);
    font-size: 0.9em;
    overflow-wrap: anywhere;
}
.err__actions {
    display: flex;
    gap: 10px;
    margin-top: 28px;
}
.err__actions :deep(.btn) {
    min-height: 44px;
    padding: 0 18px;
    font-size: 15px;
    gap: 8px;
}

@media (max-width: 900px) {
    .err {
        padding-top: 40px;
    }
    .err__title {
        font-size: 36px;
        letter-spacing: -1.2px;
    }
    .err__text {
        font-size: 16px;
    }
    .err__actions {
        flex-direction: column;
    }
}
</style>
