<script lang="ts" setup>
import { useAuth0 } from '@auth0/auth0-vue';
import type { RecipeBaseDto } from '@/data/dtos/recipe/RecipeDto';

defineProps<{ recipe: RecipeBaseDto }>();

const { isAuthenticated } = useAuth0();

function fmtTime(min: number | null) {
    if (!min) return null;
    if (min < 60) return `${min} MIN`;
    const h = Math.floor(min / 60);
    const m = min % 60;
    return m ? `${h}H ${m}MIN` : `${h}H`;
}

const diffLabel: Record<number, string> = { 1: 'ŁATWE', 2: 'ŚREDNIE', 3: 'TRUDNE' };
</script>

<template>
    <RouterLink class="card recipe-card" :to="{ name: 'recipe-details', params: { id: recipe.id } }">
        <div class="photo" style="--hue: 220; position: relative;">
            <img
                v-if="recipe.photoUrl"
                :src="recipe.photoUrl"
                :alt="recipe.name"
                style="position:absolute;inset:0;width:100%;height:100%;object-fit:cover;"
            />
        </div>

        <div class="body">
            <div style="margin-bottom: 10px;">
                <span v-if="recipe.category" class="badge" :class="`cat-${recipe.category}`">
                    {{ recipe.category }}<template v-if="recipe.region"> · {{ recipe.region }}</template>
                </span>
            </div>
            <div class="title balance">{{ recipe.name }}</div>
            <div class="meta">
                <template v-if="fmtTime(recipe.timeMinutes)">
                    <span>{{ fmtTime(recipe.timeMinutes) }}</span>
                    <span class="sep">·</span>
                </template>
                <span v-if="recipe.difficulty">{{ diffLabel[recipe.difficulty] ?? '' }}</span>
            </div>
        </div>

        <div v-if="isAuthenticated" class="recipe-card__actions">
            <RouterLink
                :to="{ name: 'recipe-edit', params: { id: recipe.id } }"
                class="btn btn--secondary btn--sm"
                @click.prevent.stop
            >Edytuj</RouterLink>
        </div>
    </RouterLink>
</template>

<style scoped>
.recipe-card { position: relative; }

.recipe-card__actions {
    display: none;
    position: absolute;
    bottom: 12px;
    right: 12px;
    gap: 8px;
}

.recipe-card:hover .recipe-card__actions {
    display: flex;
}
</style>
