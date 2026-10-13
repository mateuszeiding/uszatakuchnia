<script lang="ts" setup>
import { toastMessage } from '@/composables/useToast';
</script>

<template>
    <div
        class="sr-only"
        role="status"
    >
        {{ toastMessage }}
    </div>
    <div
        v-if="toastMessage"
        class="toast"
        aria-hidden="true"
    >
        {{ toastMessage }}
    </div>
</template>
