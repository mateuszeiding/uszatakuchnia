<script lang="ts" setup>
import { computed } from 'vue';

// A removed row that a step refers to: the row stays faded in place with "Cofnij", so the ref in the step
// can be saved from turning into "brak składnika". Changing the group, order or quantity needs no notice.
const props = defineProps<{
    name: string;
    group: string | null;
    qty: string;
    unit: string;
    /** Steps (editor numbering) that refer to the row. */
    steps: number[];
}>();
const emit = defineEmits<{ undo: [] }>();

const where = computed(() => {
    const s = props.steps;
    if (s.length === 1) return `w kroku ${s[0]}`;
    return `w krokach ${s.slice(0, -1).join(', ')} i ${s[s.length - 1]}`;
});
</script>

<template>
    <div class="removed">
        <div
            class="removed__ghost"
            aria-hidden="true"
        >
            <span class="removed__cell removed__cell--mono">{{ qty || '—' }}</span>
            <span class="removed__cell">{{ unit || 'bez jednostki' }}</span>
            <span class="removed__cell">{{ name }}</span>
        </div>
        <div
            class="removed__alert"
            role="alert"
        >
            <svg
                class="removed__warn"
                width="14"
                height="14"
                viewBox="0 0 16 16"
                fill="none"
                aria-hidden="true"
            >
                <path
                    d="M8 2.2L14.4 13.3H1.6L8 2.2z"
                    stroke="currentColor"
                    stroke-width="1.4"
                    stroke-linejoin="round"
                />
                <path
                    d="M8 6.4v3.3M8 11.6v.1"
                    stroke="currentColor"
                    stroke-width="1.6"
                    stroke-linecap="round"
                />
            </svg>
            <span class="removed__text">
                Usunięto „{{ name }}{{ group ? ` · ${group}` : '' }}”. Ilość tego składnika jest
                {{ where }} — tam zostanie znacznik „brak składnika”.
            </span>
            <button
                type="button"
                class="removed__undo"
                @click="emit('undo')"
            >
                Cofnij
            </button>
        </div>
    </div>
</template>

<style scoped>
.removed {
    display: flex;
    flex-direction: column;
    gap: 10px;
}
.removed__ghost {
    display: grid;
    grid-template-columns: 84px 124px minmax(0, 1fr);
    gap: 8px;
    margin-left: 28px;
    opacity: 0.5;
}
.removed__cell {
    padding: 8px 10px;
    border: 1px solid var(--rule);
    border-radius: var(--r-md);
    background: var(--bg);
    font-size: 14px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}
.removed__cell--mono {
    font-family: var(--font-mono);
}
.removed__alert {
    display: flex;
    gap: 10px;
    align-items: flex-start;
    padding: 10px 12px;
    border-radius: var(--r-lg);
    background: var(--prep-soft);
    border: 1px solid var(--prep-line);
    font-size: 13.5px;
    line-height: 1.5;
}
.removed__warn {
    flex-shrink: 0;
    margin-top: 2px;
    color: var(--prep);
}
.removed__text {
    flex: 1;
}
.removed__undo {
    background: none;
    border: none;
    padding: 0;
    font: inherit;
    font-weight: 600;
    color: var(--accent);
    cursor: pointer;
    white-space: nowrap;
}
</style>
