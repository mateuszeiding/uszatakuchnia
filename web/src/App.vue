<script setup lang="ts">
import { onMounted } from 'vue';
import { RouterView } from 'vue-router';

import './scss/_!import.scss';
import UTopBar from './components/UTopBar.vue';

onMounted(() => {
    window.hideLoader?.();
});
</script>

<template>
    <UTopBar />
    <main>
        <RouterView />
    </main>
</template>
