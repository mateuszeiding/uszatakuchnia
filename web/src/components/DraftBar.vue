<script lang="ts" setup>
import UIcon from './UIcon.vue';

import type { DraftBarKind } from '@/composables/useDraftBar';

// The "Szkic" bar above the nav: only the owner sees it, on a draft's page and on its shopping list.
// Not sticky: it scrolls away, and below it the page looks as if it were published.
defineProps<{ recipeId: number; kind: DraftBarKind }>();
</script>

<template>
    <div
        class="draft-bar"
        role="region"
        aria-label="Podgląd szkicu"
    >
        <div class="draft-bar__inner">
            <span class="draft-bar__label">
                <UIcon
                    name="lock"
                    :size="14"
                />
                Szkic
            </span>
            <p
                v-if="kind === 'list'"
                class="draft-bar__text"
            >
                <span class="draft-bar__long">
                    Ta lista działa tylko u ciebie — link nie otworzy się u innych, dopóki nie
                    opublikujesz przepisu.
                </span>
                <span class="draft-bar__short">Ta lista działa tylko u ciebie.</span>
            </p>
            <p
                v-else
                class="draft-bar__text"
            >
                <span class="draft-bar__long">
                    Widzisz go tylko ty. Poniżej strona wygląda tak, jak zobaczą ją czytelnicy po
                    publikacji.
                </span>
                <span class="draft-bar__short">Widzisz go tylko ty.</span>
            </p>
            <RouterLink
                :to="{ name: 'recipe-edit', params: { id: recipeId } }"
                class="draft-bar__edit"
            >
                <UIcon
                    name="edit"
                    :size="14"
                />
                Edytuj
            </RouterLink>
        </div>
    </div>
</template>

<style lang="scss" scoped>
@use '@/scss/mixins/media' as m;

.draft-bar {
    background: var(--bar-bg);
    border-bottom: 1px solid var(--prep-line);
}

.draft-bar__inner {
    max-width: 1280px;
    min-height: 48px;
    margin: 0 auto;
    padding: 8px 56px;
    display: flex;
    align-items: center;
    gap: 12px;
}

.draft-bar__label {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    font-family: var(--font-mono);
    font-size: 13px;
    font-weight: 700;
    letter-spacing: 0.4px;
    text-transform: uppercase;
    color: var(--prep);
}

.draft-bar__text {
    flex: 1;
    margin: 0;
    font-size: 14px;
    line-height: 1.35;
    color: var(--ink);
}

.draft-bar__short {
    display: none;
}

.draft-bar__edit {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    min-height: 40px;
    padding: 0 14px;
    border: 1px solid var(--rule-strong);
    border-radius: var(--r-md);
    background: var(--bg-alt);
    color: var(--ink);
    font-size: 14px;
    font-weight: 600;
    text-decoration: none;
    white-space: nowrap;
    transition: border-color 0.15s;
}

.draft-bar__edit:hover {
    border-color: var(--ink-muted);
}

.draft-bar__edit:focus-visible {
    outline: 2px solid var(--accent);
    outline-offset: 2px;
}

@media (max-width: 900px) {
    .draft-bar__inner {
        min-height: 52px;
        padding: 6px 8px 6px 16px;
    }
    .draft-bar__long {
        display: none;
    }
    .draft-bar__short {
        display: inline;
    }
}

@include m.print {
    .draft-bar {
        display: none;
    }
}
</style>
