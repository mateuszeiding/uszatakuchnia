import IngredientsView from '@view/ingredients/IngredientsView.vue'
import RecipesView from '@view/recipies/RecipesView.vue'
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/recipes',
      component: RecipesView,
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
})

export default router
