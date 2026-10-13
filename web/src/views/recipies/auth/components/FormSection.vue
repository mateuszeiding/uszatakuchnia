<script lang="ts" setup>
import UFieldError from '@/components/form/UFieldError.vue';

defineProps<{
    id?: string;
    n: string;
    title: string;
    hint?: string;
    optional?: boolean;
    error?: string;
}>();
</script>

<template>
    <section
        :id="id"
        class="form-section"
        :class="{ 'is-invalid': error }"
    >
        <div class="form-section__head">
            <span class="form-section__n">{{ n }}</span>
            <h2 class="form-section__title">{{ title }}</h2>
            <span
                v-if="optional"
                class="form-section__opt"
            >
                opcjonalne
            </span>
        </div>
        <p
            v-if="hint"
            class="form-section__hint"
        >
            {{ hint }}
        </p>
        <UFieldError
            v-if="error"
            class="form-section__error"
        >
            {{ error }}
        </UFieldError>
        <slot />
    </section>
</template>

<style scoped>
.form-section {
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    padding: 28px;
}
.form-section.is-invalid {
    border-color: var(--danger);
}
.form-section__head {
    display: flex;
    align-items: baseline;
    gap: 12px;
    margin-bottom: 20px;
}
.form-section__hint + .form-section__error,
.form-section__head:has(+ .form-section__hint) {
    margin-bottom: 6px;
}
.form-section__n {
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--accent);
    font-weight: 600;
}
.form-section__title {
    font-size: 20px;
    font-weight: 600;
    letter-spacing: -0.5px;
}
.form-section__opt {
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--ink-muted);
    letter-spacing: 0.5px;
    text-transform: uppercase;
}
.form-section__hint {
    font-size: 13px;
    color: var(--ink-muted);
    line-height: 1.45;
    margin-bottom: 20px;
}
.form-section__error {
    margin: -8px 0 14px;
}
</style>
