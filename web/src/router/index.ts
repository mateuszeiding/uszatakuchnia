import { authGuard } from '@auth0/auth0-vue';
import RecipesView from '@view/recipies/list/RecipesView.vue';
import { createRouter, createWebHistory } from 'vue-router';

import AdminDashboard from '@/views/admin/AdminDashboard.vue';
import AdminIngredients from '@/views/admin/AdminIngredients.vue';
import AdminRecipes from '@/views/admin/AdminRecipes.vue';
import AdminSettings from '@/views/admin/AdminSettings.vue';
import UpsertRecipe from '@/views/recipies/auth/UpsertRecipe.vue';
import RecipeDetails from '@/views/recipies/details/RecipeDetails.vue';

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
        { path: '/admin', component: AdminDashboard, beforeEnter: authGuard },
        { path: '/admin/recipes', component: AdminRecipes, beforeEnter: authGuard },
        { path: '/admin/ingredients', component: AdminIngredients, beforeEnter: authGuard },
        { path: '/admin/settings', component: AdminSettings, beforeEnter: authGuard },
        { path: '/challenges', component: PlaceholderView },
        { path: '/wiki', component: PlaceholderView },
        { path: '/:pathMatch(.*)*', redirect: '/recipes' },
    ],
});

export default router;
