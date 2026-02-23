<script lang="ts" setup>
import UBadge from './UBadge.vue';

import type { IngredientDto } from '@/data/dtos/ingredients/IngredientDto';

defineProps<{
    ingredient: IngredientDto;
}>();
</script>
<template>
    <div class="card p-relative border radius-4 pe-5 w-100 h-100 bg-gray-ghost-white">
        <div class="image-wrapper p-relative">
            <img
                alt="img"
                class="radius-top-s-4 radius-bot-s-4"
                :src="ingredient.image.urls.small"
            />
            <a
                class="ellipsis max-w-75 fs-xxs p-absolute left-0 bottom-0 bg-purple-magnolia py-1 px-3 radius-top-e-2 tx-violet-french-violet"
                :href="ingredient.image.author.profileURL"
            >
                {{ ingredient.image.author.name ?? ingredient.image.author.username }}
            </a>
        </div>
        <div class="pt-2">
            <div class="fs-lg pb-2">
                {{ ingredient.name }}
            </div>
            <u-badge
                class="d-inline-flex"
                :val="ingredient.type.code"
                :name="ingredient.type.name"
            />
        </div>
    </div>
</template>
<style lang="scss" scoped>
.card {
    display: flex;
    flex-direction: row;
    overflow: hidden;
    gap: var(--space-3);

    &::before {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        z-index: -1;
        scale: 100%;
        border-radius: 1.2rem;
        box-shadow: 0 4px 12px 0 rgba(31, 38, 135, 0.37);

        animation: growing 0.2s 1 normal;
    }
}
@keyframes growing {
    from {
        scale: 95%;
    }
    to {
        scale: 100%;
    }
}

.image-wrapper {
    width: 120px;
    overflow: hidden;
    position: relative;

    img {
        width: 100%;
        height: 100px;
        object-fit: cover;
        object-position: center center;
        display: block;
    }
}
</style>
