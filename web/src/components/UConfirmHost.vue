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
        <ul
            v-if="pendingConfirm.items?.length"
            class="confirm__list"
        >
            <li
                v-for="(it, i) in pendingConfirm.items"
                :key="i"
            >
                {{ it }}
            </li>
        </ul>
        <p
            v-if="pendingConfirm.note"
            class="confirm__note"
        >
            {{ pendingConfirm.note }}
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

<style scoped>
.confirm__list {
    margin: 12px 0 0;
    padding-left: 18px;
    font-size: 14px;
    line-height: 1.6;
    color: var(--ink);
}
.confirm__note {
    margin-top: 12px;
    font-size: 13px;
    color: var(--ink-muted);
}
</style>
