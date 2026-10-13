<script lang="ts" setup>
// Obudowa grupy składników / etapu kroków: nazwa grupy, lista, pole upuszczenia na końcu.
defineProps<{
    label: string;
    index: number;
    groupCount: number;
    placeholder: string;
    noun: string;
    empty: boolean;
    endLine: boolean;
    endProps: Record<string, unknown>;
}>();
const emit = defineEmits<{ 'update:label': [v: string]; remove: [] }>();
</script>

<template>
    <div class="group-shell">
        <div class="group-shell__head">
            <input
                :value="label"
                class="group-shell__label"
                :aria-label="`Nazwa grupy ${index + 1}`"
                :placeholder="groupCount === 1 ? 'Nazwa grupy (opcjonalnie)' : placeholder"
                @input="emit('update:label', ($event.target as HTMLInputElement).value)"
            />
            <button
                v-if="groupCount > 1"
                type="button"
                class="row-icon-btn"
                :aria-label="`Usuń grupę ${label || index + 1}`"
                :title="`Usuń grupę ${label || index + 1}`"
                @click="emit('remove')"
            >
                <svg
                    width="11"
                    height="11"
                    viewBox="0 0 11 11"
                    aria-hidden="true"
                >
                    <path
                        d="M2 2l7 7M9 2L2 9"
                        stroke="currentColor"
                        stroke-width="1.5"
                        stroke-linecap="round"
                    />
                </svg>
            </button>
        </div>
        <div class="group-shell__items">
            <slot />
        </div>
        <div
            v-bind="endProps"
            class="group-shell__end"
            :class="{ 'is-empty': empty }"
        >
            <template v-if="empty">Przeciągnij tu {{ noun }} albo dodaj nowy</template>
            <span
                v-if="endLine"
                class="drop-line"
                :class="empty ? 'drop-line--before' : 'drop-line--after'"
                aria-hidden="true"
            />
        </div>
        <slot name="footer" />
    </div>
</template>

<style scoped>
.group-shell {
    background: var(--bg);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    padding: 16px;
}
.group-shell__head {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 12px;
}
.group-shell__label {
    flex: 1;
    background: transparent;
    border: none;
    font-family: var(--font-mono);
    font-size: 13px;
    font-weight: 600;
    letter-spacing: 0.4px;
    text-transform: uppercase;
    color: var(--ink);
    padding: 4px 0;
}
.group-shell__label::placeholder {
    color: var(--ink-muted);
}
.group-shell__items {
    display: flex;
    flex-direction: column;
    gap: 8px;
}
.group-shell__end {
    position: relative;
}
.group-shell__end.is-empty {
    min-height: 44px;
    margin-top: 4px;
    border: 1px dashed var(--rule-strong);
    border-radius: var(--r-md);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 12px;
    color: var(--ink-muted);
}
</style>
