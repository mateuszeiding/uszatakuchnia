<template>
    <div class="skeleton-card">
        <div class="sk skeleton-card__img" />
        <div class="skeleton-card__body">
            <div class="sk skeleton-card__title" />
            <div class="sk skeleton-card__sub" />
            <div class="skeleton-card__meta">
                <div class="sk skeleton-card__pill" />
                <div class="sk skeleton-card__pill" />
            </div>
        </div>
    </div>
</template>

<style scoped>
.skeleton-card {
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: 8px;
    overflow: hidden;
    display: flex;
    flex-direction: column;
}
.skeleton-card__img {
    width: 100%;
    aspect-ratio: 4 / 5;
    border-radius: 0;
}
.skeleton-card__body {
    padding: 12px 14px 14px;
    display: flex;
    flex-direction: column;
    gap: 8px;
}
.skeleton-card__title {
    height: 14px;
    width: 72%;
}
.skeleton-card__sub {
    height: 11px;
    width: 48%;
}
.skeleton-card__meta {
    display: flex;
    gap: 8px;
    margin-top: 4px;
}
.skeleton-card__pill {
    height: 10px;
    width: 44px;
    border-radius: 999px;
}
</style>
