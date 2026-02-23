<script setup lang="ts">
import { RouterView } from 'vue-router';

import UTopBar from './components/UTopBar.vue';
</script>

<template>
    <UTopBar />
    <main class="container pt-7 pb-6">
        <RouterView />
    </main>
</template>
