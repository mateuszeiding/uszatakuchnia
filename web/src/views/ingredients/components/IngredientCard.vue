<script lang="ts" setup>
import UBadge from '@/components/UBadge.vue';
import type { IngredientDto } from '@/data/dtos/ingredients/IngredientDto';

defineProps<{
    ingredient: IngredientDto;
}>();
</script>
<template>
    <div class="card ingredient-card">
        <div class="image-wrapper">
            <img
                alt="img"
                :src="ingredient.image.urls.small"
            />
            <a
                class="photo-credit"
                :href="ingredient.image.author.profileURL"
            >
                {{ ingredient.image.author.name ?? ingredient.image.author.username }}
            </a>
        </div>
        <div class="body">
            <div class="name">{{ ingredient.name }}</div>
            <u-badge
                :val="ingredient.type.code"
                :name="ingredient.type.name"
            />
        </div>
    </div>
</template>
<style lang="scss" scoped>
.ingredient-card {
    display: flex;
    flex-direction: row;
    overflow: hidden;
    transition: transform .15s, box-shadow .15s;

    &:hover {
        transform: translateY(-2px);
        box-shadow: var(--shadow-card);
    }
}

.image-wrapper {
    width: 120px;
    flex-shrink: 0;
    position: relative;
    overflow: hidden;

    img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        object-position: center;
        display: block;
    }
}

.photo-credit {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    padding: 2px 6px;
    font-size: 9px;
    font-family: var(--font-mono);
    color: var(--ink-muted);
    background: color-mix(in srgb, var(--bg) 80%, transparent);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    text-decoration: none;

    &:hover {
        color: var(--ink);
    }
}

.body {
    padding: 14px 16px;
    display: flex;
    flex-direction: column;
    gap: 8px;
    justify-content: center;
}

.name {
    font-size: 15px;
    font-weight: 600;
    letter-spacing: -0.2px;
    color: var(--ink);
}
</style>
