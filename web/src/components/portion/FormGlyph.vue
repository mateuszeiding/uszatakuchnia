<script lang="ts" setup>
defineProps<{ shape: 'round' | 'rect' }>();
</script>

<template>
    <span
        class="fglyph"
        :class="`fglyph--${shape}`"
        aria-hidden="true"
    />
</template>

<style scoped>
.fglyph {
    display: inline-block;
    flex-shrink: 0;
    border: 1.6px solid currentColor;
    opacity: 0.8;
}
.fglyph--round {
    width: 14px;
    height: 14px;
    border-radius: var(--r-pill);
}
.fglyph--rect {
    width: 17px;
    height: 12px;
    border-radius: 2px;
}
</style>
