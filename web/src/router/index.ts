import { authGuard } from '@auth0/auth0-vue';
import RecipesView from '@view/recipies/list/RecipesView.vue';
import { createRouter, createWebHistory } from 'vue-router';

import { isOwner } from '@/auth/createAuth';
import AdminIngredientEdit from '@/views/admin/AdminIngredientEdit.vue';
import AdminIngredients from '@/views/admin/AdminIngredients.vue';
import AdminLayout from '@/views/admin/AdminLayout.vue';
import AdminRecipes from '@/views/admin/AdminRecipes.vue';
import AdminTaxonomy from '@/views/admin/AdminTaxonomy.vue';
import ChallengesView from '@/views/challenges/ChallengesView.vue';
import UpsertRecipe from '@/views/recipies/auth/UpsertRecipe.vue';
import RecipeDetails from '@/views/recipies/details/RecipeDetails.vue';
import ShoppingListView from '@/views/recipies/shopping/ShoppingListView.vue';

// Po zalogowaniu (authGuard) wpuszcza tylko konto właściciela; inne wracają na listę.
const ownerGuard = () => (isOwner.value ? true : '/recipes');
const adminGuards = [authGuard, ownerGuard];

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    scrollBehavior: (to, from, saved) => {
        if (saved) return saved;
        // zmiana samych filtrów listy (query) nie przewija strony
        if (to.path === from.path) return false;
        return { top: 0 };
    },
    routes: [
        { path: '/', redirect: '/recipes' },
        { path: '/recipes', component: RecipesView },
        { path: '/recipes/new', component: UpsertRecipe, beforeEnter: adminGuards },
        { path: '/recipes/:id', name: 'recipe-details', component: RecipeDetails, props: true },
        {
            path: '/recipes/:id/zakupy',
            name: 'recipe-shopping',
            component: ShoppingListView,
            props: true,
        },
        {
            path: '/recipes/:id/edit',
            name: 'recipe-edit',
            component: UpsertRecipe,
            props: true,
            beforeEnter: adminGuards,
        },
        { path: '/challenges', component: ChallengesView },
        {
            path: '/admin',
            component: AdminLayout,
            beforeEnter: adminGuards,
            children: [
                { path: '', redirect: '/admin/recipes' },
                { path: 'recipes', component: AdminRecipes },
                { path: 'ingredients', component: AdminIngredients },
                { path: 'ingredients/new', component: AdminIngredientEdit },
                { path: 'ingredients/:id(\\d+)', component: AdminIngredientEdit, props: true },
                { path: 'taxonomy', component: AdminTaxonomy },
                // stare adresy
                { path: 'settings', redirect: '/admin/taxonomy' },
            ],
        },
        { path: '/:pathMatch(.*)*', redirect: '/recipes' },
    ],
});

export default router;
