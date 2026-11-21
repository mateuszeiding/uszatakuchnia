<script setup lang="ts">
import { Analytics } from '@vercel/analytics/vue';
import { computed, onMounted } from 'vue';
import { RouterView, useRoute } from 'vue-router';

import './scss/_!import.scss';
import DraftBar from './components/DraftBar.vue';
import UConfirmHost from './components/UConfirmHost.vue';
import UFooter from './components/UFooter.vue';
import UToastHost from './components/UToastHost.vue';
import UTopBar from './components/UTopBar.vue';
import { draftBar } from './composables/useDraftBar';
import { skipAdmin } from './shared/utils/analytics';

// Explicit, so `vite dev` never reports even if process.env.NODE_ENV is not inlined.
const analyticsMode = import.meta.env.DEV ? 'development' : 'production';

// footer on public pages; the admin panel has its own layout
const route = useRoute();
const isPublic = computed(() => !route.path.startsWith('/admin'));

onMounted(() => {
    window.hideLoader?.();
});
</script>

<template>
    <!-- "Szkic" bar above the nav, not sticky (UTopBar is sticky, so it stays on top) -->
    <DraftBar
        v-if="draftBar"
        :recipe-id="draftBar.recipeId"
        :kind="draftBar.kind"
    />
    <UTopBar />
    <main>
        <RouterView />
    </main>
    <UFooter v-if="isPublic" />
    <UConfirmHost />
    <UToastHost />
    <Analytics
        :mode="analyticsMode"
        :before-send="skipAdmin"
    />
</template>
