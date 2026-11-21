<template>
    <!-- The same layout as RecipeCard: the 1:1 frame with a line (the real .badge class keeps the badge height). -->
    <div
        class="skc"
        aria-hidden="true"
    >
        <div class="skc__photo"><div class="sk skc__img" /></div>
        <div class="skc__body">
            <div class="skc__flags">
                <span class="badge skc__badge">
                    <span class="sk skc__fill" />
                    &nbsp;
                </span>
            </div>
            <div class="skc__title">
                <span class="sk skc__line" />
                <span class="sk skc__line skc__line--short" />
            </div>
            <div class="skc__meta">
                <span class="sk skc__line skc__line--meta" />
            </div>
        </div>
    </div>
</template>

<style scoped>
.skc {
    display: flex;
    flex-direction: column;
    height: 100%;
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    overflow: hidden;
}
/* like .r-card__photo: a 1:1 box with a line, the fill inside */
.skc__photo {
    box-sizing: content-box;
    aspect-ratio: 1 / 1;
    border-bottom: 1px solid var(--rule);
}
.skc__img {
    width: 100%;
    height: 100%;
    border-radius: 0;
}
.skc__body {
    padding: 16px;
    display: flex;
    flex-direction: column;
    gap: 10px;
    flex: 1;
}
/* like .r-card__flags (12 px) */
.skc__flags {
    display: flex;
    font-size: 12px;
}
/* shape and height from .badge, the skeleton fill inside, no color dot */
.skc__badge {
    position: relative;
    width: 92px;
    color: transparent;
    background: transparent;
    border-color: transparent;
}
.skc__badge::before {
    visibility: hidden;
}
.skc__fill {
    position: absolute;
    inset: 0;
    border-radius: var(--r-pill);
}
/* like .r-card__title: 18 px × 1.25, two lines reserved */
.skc__title {
    flex: 1;
    min-height: calc(2 * 18px * 1.25);
    display: flex;
    flex-direction: column;
    justify-content: space-around;
}
.skc__line {
    display: block;
    height: 14px;
    width: 86%;
}
.skc__line--short {
    width: 54%;
}
/* like .r-card__meta: a 20 px line */
.skc__meta {
    height: 20px;
    display: flex;
    align-items: center;
}
.skc__line--meta {
    height: 11px;
    width: 62%;
}
</style>
