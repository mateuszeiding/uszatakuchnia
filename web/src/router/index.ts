import { authGuard } from '@auth0/auth0-vue';
import IngredientsView from '@view/ingredients/IngredientsView.vue';
import RecipesView from '@view/recipies/RecipesView.vue';
import { createRouter, createWebHistory } from 'vue-router';

import RootView from '@/views/RootView.vue';
import UpsertType from '@/views/type/auth/UpsertType.vue';

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            component: RootView,
        },
        {
            path: '/recipes',
            component: RecipesView,
        },
        {
            path: '/type/new',
            component: UpsertType,
            beforeEnter: authGuard,
        },
        {
            path: '/ingredients',
            component: IngredientsView,
        },
        {
            path: '/:pathMatch(.*)*',
            redirect: '/recipes',
        },
    ],
});

export default router;
