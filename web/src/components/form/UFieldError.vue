<template>
    <span
        class="field-error"
        role="alert"
    >
        <svg
            width="12"
            height="12"
            viewBox="0 0 12 12"
            aria-hidden="true"
        >
            <circle
                cx="6"
                cy="6"
                r="5"
                stroke="currentColor"
                stroke-width="1.3"
            />
            <path
                d="M6 3.5v3M6 8.3v.2"
                stroke="currentColor"
                stroke-width="1.4"
                stroke-linecap="round"
            />
        </svg>
        <slot />
    </span>
</template>
