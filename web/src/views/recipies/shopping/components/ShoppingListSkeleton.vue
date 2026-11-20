<script lang="ts" setup>
import ShoppingHeader from './ShoppingHeader.vue';
import ShoppingRowSkeleton from './ShoppingRowSkeleton.vue';

import USpinner from '@/components/USpinner.vue';

// Rows of ShoppingRow height; the .shop container as in ShoppingList.vue.
defineProps<{ recipeId: number }>();
const AISLES = [3, 2, 2];
</script>

<template>
    <ShoppingHeader
        :recipe-id="recipeId"
        recipe-name=""
        :servings="0"
        :done="0"
        :total="0"
        loading
    />
    <div
        class="shop"
        aria-busy="true"
    >
        <p class="shop__hint">
            Odhacz to, co masz w domu albo już w koszyku. Ilości zsumowane ze wszystkich grup
            przepisu.
        </p>
        <div class="shop__cols">
            <section
                v-for="(n, i) in AISLES"
                :key="i"
                class="aisle"
                aria-hidden="true"
            >
                <h2 class="aisle__title"><span class="sk shop__sk-title">&nbsp;</span></h2>
                <ul class="aisle__list">
                    <ShoppingRowSkeleton
                        v-for="r in n"
                        :key="r"
                    />
                </ul>
            </section>
            <div
                class="shop__loading"
                role="status"
            >
                <USpinner :size="28" />
                <span>ładowanie listy zakupów…</span>
            </div>
        </div>
    </div>
</template>

<!-- .shop / .aisle as in ShoppingList.vue: keep them in sync -->
<style lang="scss" scoped>
.shop {
    max-width: 640px;
    margin: 0 auto;
    padding: 8px 0 40px;
    display: flex;
    flex-direction: column;
    gap: 24px;
}
.shop__hint {
    margin-top: 14px;
    font-size: 13px;
    line-height: 1.5;
    color: var(--ink-muted);
}
.shop__cols {
    position: relative;
    display: flex;
    flex-direction: column;
    gap: 24px;
}
.aisle__title {
    display: flex;
    padding-bottom: 6px;
    border-bottom: 1px solid var(--rule-strong);
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 600;
    letter-spacing: 0.5px;
    text-transform: uppercase;
}
/* line height from real text (nbsp) */
.shop__sk-title {
    display: block;
    width: 110px;
    color: transparent;
}
.aisle__list {
    list-style: none;
    padding: 0;
    margin: 0;
}
.shop__loading {
    position: absolute;
    inset: 0;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 12px;
    font-size: 14px;
    color: var(--ink-muted);
}

@media (max-width: 900px) {
    .shop {
        padding: 4px 16px 28px;
        gap: 20px;
    }
    .shop__hint {
        margin-top: 10px;
    }
    .shop__cols {
        gap: 20px;
    }
}
</style>
