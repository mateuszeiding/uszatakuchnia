import { authGuard } from '@auth0/auth0-vue';
import IngredientsView from '@view/ingredients/IngredientsView.vue';
import RecipesView from '@view/recipies/list/RecipesView.vue';
import { createRouter, createWebHistory } from 'vue-router';

import UpsertRecipe from '@/views/recipies/auth/UpsertRecipe.vue';
import RecipeDetails from '@/views/recipies/details/RecipeDetails.vue';
import UpsertType from '@/views/type/auth/UpsertType.vue';

const PlaceholderView = {
    template:
        '<div class="container" style="padding-top:80px"><p class="t-h2 muted">Wkrótce</p></div>',
};

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        { path: '/', redirect: '/recipes' },
        { path: '/recipes', component: RecipesView },
        { path: '/recipes/new', component: UpsertRecipe, beforeEnter: authGuard },
        { path: '/recipes/:id', name: 'recipe-details', component: RecipeDetails, props: true },
        {
            path: '/recipes/:id/edit',
            name: 'recipe-edit',
            component: UpsertRecipe,
            props: true,
            beforeEnter: authGuard,
        },
        { path: '/challenges', component: PlaceholderView },
        { path: '/wiki', component: PlaceholderView },
        { path: '/ingredients', component: IngredientsView },
        { path: '/type/new', component: UpsertType, beforeEnter: authGuard },
        { path: '/:pathMatch(.*)*', redirect: '/recipes' },
    ],
});

export default router;
