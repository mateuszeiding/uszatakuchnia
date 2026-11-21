<script setup lang="ts">
// Displayed quantity ("1,5", "1/3", "1,25", "0,25–0,5"). Thirds are plain digits (no ⅓ glyphs),
// so the parent's font works; the quantity only must not wrap between "1" and "1/3".
defineProps<{ text?: string }>();
</script>

<template>
    <span class="qty-text">{{ text ?? '' }}</span>
</template>

<style scoped>
.qty-text {
    white-space: nowrap;
}
</style>
