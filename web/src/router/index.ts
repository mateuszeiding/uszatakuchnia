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
import UpsertRecipe from '@/views/recipies/auth/UpsertRecipe.vue';
import RecipeDetails from '@/views/recipies/details/RecipeDetails.vue';
import ShoppingListView from '@/views/recipies/shopping/ShoppingListView.vue';

// Po zalogowaniu (authGuard) wpuszcza tylko konto właściciela; inne wracają na listę.
const ownerGuard = () => (isOwner.value ? true : '/recipes');
const adminGuards = [authGuard, ownerGuard];

/** Czeka (do ~3 s), aż warunek się spełni — treść widoku wczytuje się po nawigacji (szkielet → dane).
 *  setTimeout, nie requestAnimationFrame: ten nie tyka w karcie w tle. */
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
        // wstecz / dalej: to samo miejsce, gdy strona urośnie po wczytaniu danych; od razu, jak systemowe „wstecz”
        // (bez płynnego przewijania z CSS)
        if (saved) {
            await waitFor(
                () => document.documentElement.scrollHeight >= saved.top + window.innerHeight
            );
            return { ...saved, behavior: 'instant' };
        }
        // zmiana samych filtrów listy (query) nie przewija strony
        if (to.path === from.path) return false;
        // kotwica (np. #skladniki z listy zakupów): przewiń, gdy sekcja się wczyta (szkielet ma aria-busy);
        // scrollIntoView, żeby zadziałał scroll-margin-top pod przyklejonym paskiem
        if (to.hash) {
            window.scrollTo({ top: 0, behavior: 'instant' });
            const ready = () => document.querySelector(`${to.hash}:not([aria-busy="true"])`);
            await waitFor(() => !!ready());
            await new Promise((r) => setTimeout(r, 0)); // układ po wstawieniu treści
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
                // stare adresy
                { path: 'settings', redirect: '/admin/taxonomy' },
            ],
        },
        { path: '/:pathMatch(.*)*', redirect: '/recipes' },
    ],
});

// Tytuł karty: z meta trasy (meta dzieci nadpisuje rodzica), widoki z danymi ustawiają go po wczytaniu.
// Zmiana samego query (filtry listy, porcje na liście zakupów) nie rusza tytułu.
router.afterEach((to, from) => {
    if (to.path === from.path && from.matched.length) return;
    setTitle(to.meta.title as string | undefined);
});

export default router;
