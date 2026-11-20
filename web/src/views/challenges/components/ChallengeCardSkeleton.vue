<template>
    <!-- The same layout as ChallengeCard: badge (the real .badge-attr class, same height),
         title and description 2 lines each, footer with a border; the "done" circle in the corner. -->
    <div
        class="skc"
        aria-hidden="true"
    >
        <span class="skc__format">
            <span class="badge-attr skc__badge">
                <span class="sk skc__fill" />
                &nbsp;
            </span>
        </span>
        <span class="skc__title">
            <span class="sk skc__line" />
            <span class="sk skc__line skc__line--short" />
        </span>
        <span class="skc__desc">
            <span class="sk skc__line skc__line--thin" />
            <span class="sk skc__line skc__line--thin skc__line--mid" />
        </span>
        <span class="skc__foot">
            <span class="sk skc__line skc__line--meta" />
            &nbsp;
        </span>
        <span class="skc__dot" />
    </div>
</template>

<style scoped>
.skc {
    position: relative;
    display: flex;
    flex-direction: column;
    gap: 10px;
    height: 100%;
    padding: 16px 18px 14px;
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-xl);
}
.skc__format {
    display: flex;
}
.skc__badge {
    position: relative;
    width: 128px;
    border-color: transparent;
}
.skc__fill {
    position: absolute;
    inset: 0;
    border-radius: var(--r-pill);
}
/* like .ch-card__title: 18 px × 1.25, two lines */
.skc__title {
    font-size: 18px;
    min-height: calc(2 * 18px * 1.25);
    display: flex;
    flex-direction: column;
    justify-content: space-around;
}
/* like .ch-card__desc: 14 px × 1.5, two lines */
.skc__desc {
    flex: 1;
    min-height: calc(2 * 14px * 1.5);
    display: flex;
    flex-direction: column;
    justify-content: space-around;
}
.skc__line {
    display: block;
    height: 14px;
    width: 82%;
}
.skc__line--short {
    width: 48%;
}
.skc__line--thin {
    height: 10px;
    width: 94%;
}
.skc__line--mid {
    width: 70%;
}
/* like .ch-card__foot: border, 10 px gap, a 13 px text line (&nbsp; keeps the line height) */
.skc__foot {
    font-size: 13px;
    padding-top: 10px;
    border-top: 1px solid var(--rule);
    display: flex;
    align-items: center;
}
.skc__line--meta {
    height: 11px;
    width: 52%;
}
.skc__dot {
    position: absolute;
    top: 16px;
    right: 16px;
    width: 24px;
    height: 24px;
    border-radius: var(--r-pill);
    border: 1.5px solid var(--rule);
}
</style>
