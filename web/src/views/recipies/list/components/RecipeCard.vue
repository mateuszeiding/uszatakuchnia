<script lang="ts" setup>
import { useAuth0 } from '@auth0/auth0-vue';

import { deleteRecipe } from '@/data/api/recipe/fetch';
import type { RecipeBaseDto } from '@/data/dtos/recipe/RecipeDto';
import { confirmDialog } from '@/utils/dialog';

const props = defineProps<{ recipe: RecipeBaseDto }>();
const emit = defineEmits<{ deleted: [id: number] }>();

const { isAuthenticated } = useAuth0();

function fmtTime(min: number | null) {
    if (!min) return null;
    if (min < 60) return `${min} MIN`;
    const h = Math.floor(min / 60);
    const m = min % 60;
    return m ? `${h}H ${m}MIN` : `${h}H`;
}

const diffLabel: Record<number, string> = { 1: 'ŁATWE', 2: 'ŚREDNIE', 3: 'TRUDNE' };

async function handleDelete(e: MouseEvent) {
    e.preventDefault();
    e.stopPropagation();

    if (!confirmDialog(`Usuń przepis „${props.recipe.name}"?`)) return;
    await deleteRecipe(props.recipe.id);
    emit('deleted', props.recipe.id);
}
</script>

<template>
    <RouterLink
        class="card recipe-card"
        :to="{ name: 'recipe-details', params: { id: recipe.id } }"
    >
        <div
            class="photo recipe-card__photo"
            style="--hue: 220; position: relative"
        >
            <img
                v-if="recipe.photoUrl"
                :src="recipe.photoUrl"
                :alt="recipe.name"
                style="position: absolute; inset: 0; width: 100%; height: 100%; object-fit: cover"
            />
            <span
                v-if="recipe.status === 'draft'"
                class="draft-badge"
            >
                SZKIC
            </span>
            <div
                v-if="isAuthenticated"
                class="recipe-card__actions"
                @click.prevent.stop
            >
                <RouterLink
                    :to="{ name: 'recipe-edit', params: { id: recipe.id } }"
                    class="card-action-btn"
                    title="Edytuj"
                >
                    <svg
                        width="14"
                        height="14"
                        viewBox="0 0 14 14"
                        fill="none"
                    >
                        <path
                            d="M2 10.5L10 2.5l1.5 1.5L3.5 12H2v-1.5z"
                            stroke="currentColor"
                            stroke-width="1.3"
                            stroke-linejoin="round"
                        />
                    </svg>
                </RouterLink>
                <button
                    class="card-action-btn card-action-btn--danger"
                    title="Usuń"
                    @click="handleDelete"
                >
                    <svg
                        width="14"
                        height="14"
                        viewBox="0 0 14 14"
                        fill="none"
                    >
                        <path
                            d="M3 4h8M5.5 4V2.5h3V4M4 4l0.5 8h5L10 4M6 6v4M8 6v4"
                            stroke="currentColor"
                            stroke-width="1.3"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        />
                    </svg>
                </button>
            </div>
        </div>

        <div class="body">
            <div class="card-top-row">
                <span
                    v-if="recipe.category"
                    class="badge"
                    :class="`cat-${recipe.category}`"
                >
                    {{ recipe.category }}
                    <template v-if="recipe.region">· {{ recipe.region }}</template>
                </span>
                <span class="recipe-id">#{{ String(recipe.id).padStart(3, '0') }}</span>
            </div>
            <div class="title balance">{{ recipe.name }}</div>
            <span
                v-if="recipe.needsPrep"
                class="prep-badge"
            >
                <svg
                    width="9"
                    height="9"
                    viewBox="0 0 10 10"
                    fill="none"
                >
                    <circle
                        cx="5"
                        cy="5"
                        r="4"
                        stroke="currentColor"
                        stroke-width="1.3"
                    />
                    <path
                        d="M5 3v2.5l1.5 1"
                        stroke="currentColor"
                        stroke-width="1.2"
                        stroke-linecap="round"
                    />
                </svg>
                wymaga przygotowania
            </span>
            <div class="meta">
                <template v-if="fmtTime(recipe.timeMinutes)">
                    <span>{{ fmtTime(recipe.timeMinutes) }}</span>
                    <span class="sep">·</span>
                </template>
                <span v-if="recipe.difficulty">{{ diffLabel[recipe.difficulty] ?? '' }}</span>
            </div>
        </div>
    </RouterLink>
</template>

<style scoped>
.recipe-card__photo {
    position: relative;
}
.recipe-card__actions {
    display: none;
    position: absolute;
    top: 10px;
    right: 10px;
    gap: 6px;
    z-index: 2;
}
.recipe-card:hover .recipe-card__actions {
    display: flex;
}
.card-action-btn {
    width: 30px;
    height: 30px;
    border-radius: var(--r-md);
    background: rgba(10, 10, 15, 0.78);
    color: #fff;
    border: none;
    cursor: pointer;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    backdrop-filter: blur(8px);
    text-decoration: none;
    transition: background 0.12s;
}
.card-action-btn--danger {
    background: rgba(197, 55, 40, 0.92);
}
.draft-badge {
    position: absolute;
    top: 10px;
    left: 10px;
    font-family: var(--font-mono);
    font-size: 10px;
    font-weight: 600;
    letter-spacing: 0.5px;
    background: rgba(10, 10, 15, 0.78);
    color: #fff;
    padding: 4px 8px;
    border-radius: var(--r-sm);
    text-transform: uppercase;
}
.card-top-row {
    display: flex;
    justify-content: space-between;
    align-items: baseline;
    margin-bottom: 10px;
    gap: 8px;
}
.recipe-id {
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--ink-faint);
    flex-shrink: 0;
}
.prep-badge {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    margin-top: 8px;
    width: fit-content;
    background: color-mix(in srgb, #d9a441 16%, transparent);
    border: 1px solid color-mix(in srgb, #d9a441 45%, transparent);
    color: #d9a441;
    padding: 2px 9px;
    border-radius: var(--r-pill);
    font-size: 10px;
    font-family: var(--font-mono);
    font-weight: 600;
    letter-spacing: 0.3px;
}
</style>
