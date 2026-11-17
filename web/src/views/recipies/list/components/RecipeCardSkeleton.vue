<template>
    <!-- Ten sam układ co RecipeCard: zdjęcie 1:1 z linią, body 16 px / odstęp 10 px,
         badge typu dania (prawdziwa klasa .badge — ta sama wysokość), tytuł na 2 linie, linia metadanych. -->
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
/* jak .r-card__photo: kontener 1:1 z linią, wypełnienie w środku */
.skc__photo {
    box-sizing: content-box;
    aspect-ratio: 1 / 1;
    position: relative;
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
/* jak .r-card__flags (12 px) */
.skc__flags {
    display: flex;
    font-size: 12px;
}
/* kształt i wysokość z .badge, wypełnienie szkieletu w środku, bez kropki koloru */
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
/* jak .r-card__title: 18 px × 1,25, dwie linie zarezerwowane */
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
/* jak .r-card__meta: linia 20 px */
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
