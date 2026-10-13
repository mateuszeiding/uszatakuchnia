<script lang="ts" setup>
import UModal from './UModal.vue';

import { pendingConfirm, settleConfirm } from '@/composables/useConfirm';
</script>

<template>
    <UModal
        v-if="pendingConfirm"
        :title="pendingConfirm.title"
        @close="settleConfirm(false)"
    >
        <p
            v-if="pendingConfirm.body"
            class="modal__text"
        >
            {{ pendingConfirm.body }}
        </p>
        <template #footer>
            <button
                type="button"
                class="btn btn--ghost"
                @click="settleConfirm(false)"
            >
                Anuluj
            </button>
            <button
                type="button"
                class="btn"
                :class="pendingConfirm.danger ? 'btn--danger-fill' : 'btn--accent'"
                autofocus
                @click="settleConfirm(true)"
            >
                {{ pendingConfirm.confirmLabel ?? 'Potwierdź' }}
            </button>
        </template>
    </UModal>
</template>
