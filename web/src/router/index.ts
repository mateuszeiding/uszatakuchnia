import { authGuard } from '@auth0/auth0-vue';
import RecipesView from '@view/recipies/list/RecipesView.vue';
import { createRouter, createWebHistory } from 'vue-router';

import { isOwner } from '@/auth/createAuth';
import { setTitle } from '@/shared/utils/title';
import AdminChallengeEdit from '@/views/admin/AdminChallengeEdit.vue';
import AdminChallenges from '@/views/admin/AdminChallenges.vue';
import AdminIngredientEdit from '@/views/admin/AdminIngredientEdit.vue';
import AdminIngredients from '@/views/admin/AdminIngredients.vue';
import AdminLayout from '@/views/admin/AdminLayout.vue';
import AdminRecipes from '@/views/admin/AdminRecipes.vue';
import AdminTaxonomy from '@/views/admin/AdminTaxonomy.vue';
import ChallengesView from '@/views/challenges/ChallengesView.vue';
import ChallengeView from '@/views/challenges/ChallengeView.vue';
import NotFoundView from '@/views/errors/NotFoundView.vue';
import UpsertRecipe from '@/views/recipies/auth/UpsertRecipe.vue';
import RecipeDetails from '@/views/recipies/details/RecipeDetails.vue';
import ShoppingListView from '@/views/recipies/shopping/ShoppingListView.vue';

// After login (authGuard) lets only the owner account in; others go back to the list.
const ownerGuard = () => (isOwner.value ? true : '/recipes');
const adminGuards = [authGuard, ownerGuard];

/** Waits (up to ~3 s) for the condition: the view's content loads after navigation (skeleton → data).
 *  setTimeout, not requestAnimationFrame: rAF doesn't tick in a background tab. */
function waitFor(ok: () => boolean, ms = 3000): Promise<void> {
    const end = performance.now() + ms;
    return new Promise((resolve) => {
        const tick = () => (ok() || performance.now() > end ? resolve() : setTimeout(tick, 30));
        tick();
    });
}

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    scrollBehavior: async (to, from, saved) => {
        // back / forward: the same position once the page grows after data loads; instant, like the native "back"
        // (no CSS smooth scrolling)
        if (saved) {
            await waitFor(
                () => document.documentElement.scrollHeight >= saved.top + window.innerHeight
            );
            return { ...saved, behavior: 'instant' };
        }
        // a change of list filters only (query) doesn't scroll the page
        if (to.path === from.path) return false;
        // anchor (e.g. #skladniki from the shopping list): scroll once the section has loaded (the skeleton has aria-busy);
        // scrollIntoView so that scroll-margin-top under the sticky bar applies
        if (to.hash) {
            window.scrollTo({ top: 0, behavior: 'instant' });
            const ready = () => document.querySelector(`${to.hash}:not([aria-busy="true"])`);
            await waitFor(() => !!ready());
            await new Promise((r) => setTimeout(r, 0)); // layout after the content is inserted
            ready()?.scrollIntoView({ block: 'start', behavior: 'instant' });
            return false;
        }
        return { top: 0 };
    },
    routes: [
        { path: '/', redirect: '/recipes' },
        { path: '/recipes', component: RecipesView, meta: { title: 'Przepisy' } },
        { path: '/recipes/new', component: UpsertRecipe, beforeEnter: adminGuards },
        { path: '/recipes/:id', name: 'recipe-details', component: RecipeDetails, props: true },
        {
            path: '/recipes/:id/zakupy',
            name: 'recipe-shopping',
            component: ShoppingListView,
            props: true,
            meta: { title: 'Lista zakupów' },
        },
        {
            path: '/recipes/:id/edit',
            name: 'recipe-edit',
            component: UpsertRecipe,
            props: true,
            beforeEnter: adminGuards,
        },
        { path: '/challenges', component: ChallengesView, meta: { title: 'Wyzwania' } },
        {
            path: '/challenges/:id',
            component: ChallengeView,
            props: true,
            meta: { title: 'Wyzwania' },
        },
        {
            path: '/admin',
            component: AdminLayout,
            meta: { title: 'Panel' },
            beforeEnter: adminGuards,
            children: [
                { path: '', redirect: '/admin/recipes' },
                { path: 'recipes', component: AdminRecipes },
                { path: 'ingredients', component: AdminIngredients },
                { path: 'ingredients/new', component: AdminIngredientEdit },
                { path: 'ingredients/:id(\\d+)', component: AdminIngredientEdit, props: true },
                { path: 'challenges', component: AdminChallenges },
                { path: 'challenges/new', component: AdminChallengeEdit },
                { path: 'challenges/:id(\\d+)', component: AdminChallengeEdit, props: true },
                { path: 'taxonomy', component: AdminTaxonomy },
                // old URLs
                { path: 'settings', redirect: '/admin/taxonomy' },
            ],
        },
        // unknown URL: "Nie ma takiej strony" without a redirect, so the URL in the address bar stays
        { path: '/:pathMatch(.*)*', name: 'not-found', component: NotFoundView },
    ],
});

// Tab title: from the route meta (child meta overrides the parent); data views set it after loading.
// A query-only change (list filters, servings on the shopping list) leaves the title alone.
router.afterEach((to, from) => {
    if (to.path === from.path && from.matched.length) return;
    setTitle(to.meta.title as string | undefined);
});

export default router;
