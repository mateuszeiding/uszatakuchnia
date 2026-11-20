<script setup lang="ts">
import { computed, onMounted } from 'vue';
import { RouterView, useRoute } from 'vue-router';

import './scss/_!import.scss';
import UConfirmHost from './components/UConfirmHost.vue';
import UFooter from './components/UFooter.vue';
import UToastHost from './components/UToastHost.vue';
import UTopBar from './components/UTopBar.vue';

// stopka na stronach publicznych; panel ma własny układ
const route = useRoute();
const isPublic = computed(() => !route.path.startsWith('/admin'));

onMounted(() => {
    window.hideLoader?.();
});
</script>

<template>
    <UTopBar />
    <main>
        <RouterView />
    </main>
    <UFooter v-if="isPublic" />
    <UConfirmHost />
    <UToastHost />
</template>
