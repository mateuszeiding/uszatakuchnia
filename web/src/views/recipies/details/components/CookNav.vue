<script lang="ts" setup>
// "Dalej" takes 3/4 of the width and 72 px of height, so it can be hit with a knuckle. No swipe and no
// "tap anywhere": a stray touch would check off a step, and scrolling a long step would look like the gesture.
defineProps<{
    /** The "Zanim zaczniesz" screen: one "Zaczynam" button that checks nothing. */
    start: boolean;
    canBack: boolean;
    last: boolean;
}>();
const emit = defineEmits<{ back: []; next: [] }>();
</script>

<template>
    <div
        class="cook-nav"
        :class="{ 'cook-nav--start': start }"
    >
        <button
            v-if="!start"
            type="button"
            class="cook-nav__btn cook-nav__back"
            aria-label="Wstecz"
            :disabled="!canBack"
            @click="emit('back')"
        >
            ←
        </button>
        <button
            type="button"
            class="cook-nav__btn cook-nav__next"
            @click="emit('next')"
        >
            <template v-if="start">Zaczynam →</template>
            <template v-else-if="last">
                Gotowe
                <svg
                    width="20"
                    height="20"
                    viewBox="0 0 22 22"
                    fill="none"
                    aria-hidden="true"
                >
                    <path
                        d="M5 11.5L9.5 16L17 7"
                        stroke="currentColor"
                        stroke-width="2.4"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    />
                </svg>
            </template>
            <template v-else>Dalej →</template>
        </button>
    </div>
</template>

<style scoped>
.cook-nav {
    grid-area: nav;
    flex-shrink: 0;
    display: flex;
    gap: 10px;
    padding: 12px 12px calc(12px + env(safe-area-inset-bottom));
    border-top: 1px solid var(--rule);
    background: var(--bg);
}
.cook-nav__btn {
    min-height: 72px;
    border-radius: 14px;
    font-size: 20px;
    font-weight: 600;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 10px;
}
.cook-nav__back {
    flex: 1;
    background: var(--bg-alt);
    color: var(--ink);
    border: 1px solid var(--rule-strong);
}
.cook-nav__back:disabled {
    opacity: 0.4;
    cursor: default;
}
.cook-nav__next {
    flex: 3;
    background: var(--accent-fill);
    color: var(--on-accent);
    border: none;
}

/* desktop: the same buttons, centred under the content */
@media (min-width: 901px) {
    .cook-nav {
        width: 100%;
        max-width: 800px;
        margin: 0 auto;
        padding: 16px 20px calc(16px + env(safe-area-inset-bottom));
    }
}

/* landscape phone: a 150 px column on the right, "Dalej" on top (3/4 of the height) */
@media (orientation: landscape) and (max-height: 500px) {
    .cook-nav {
        flex-direction: column-reverse;
        gap: 8px;
        padding: 10px calc(10px + env(safe-area-inset-right)) 10px 10px;
        border-top: none;
        border-left: 1px solid var(--rule);
    }
    .cook-nav__btn {
        min-height: 0;
    }
    .cook-nav__back {
        font-size: 18px;
    }
}
</style>
