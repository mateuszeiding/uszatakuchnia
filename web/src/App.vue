<script setup lang="ts">
import { RouterView } from 'vue-router';

import './scss/_!import.scss';
import UTopBar from './components/UTopBar.vue';
</script>

<template>
    <UTopBar />
    <main>
        <RouterView />
    </main>
</template>
