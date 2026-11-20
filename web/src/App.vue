<script setup lang="ts">
import { computed, onMounted } from 'vue';
import { RouterView, useRoute } from 'vue-router';

import './scss/_!import.scss';
import DraftBar from './components/DraftBar.vue';
import UConfirmHost from './components/UConfirmHost.vue';
import UFooter from './components/UFooter.vue';
import UToastHost from './components/UToastHost.vue';
import UTopBar from './components/UTopBar.vue';
import { draftBar } from './composables/useDraftBar';

// stopka na stronach publicznych; panel ma własny układ
const route = useRoute();
const isPublic = computed(() => !route.path.startsWith('/admin'));

onMounted(() => {
    window.hideLoader?.();
});
</script>

<template>
    <!-- pasek „Szkic” nad nawigacją, nie przyklejony (UTopBar jest sticky, więc zostaje u góry) -->
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
</template>
