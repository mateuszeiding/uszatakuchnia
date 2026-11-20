<script lang="ts" setup>
// Challenge page loading: the same layout as ChallengeView (format chip, title, metadata, button,
// paragraphs; the "Składniki" card on the right). On a phone the card goes below the description, the button full width.
defineProps<{ narrow: boolean }>();
</script>

<template>
    <div
        class="chs"
        aria-hidden="true"
    >
        <div class="chs__main">
            <span class="badge-attr chs__badge">
                <span class="sk chs__fill" />
                &nbsp;
            </span>
            <div class="chs__title">
                <span class="sk chs__line" />
                <span class="sk chs__line chs__line--short" />
            </div>
            <div class="chs__meta">
                <span class="sk chs__meta-line" />
                <span class="sk chs__meta-line" />
            </div>
            <span
                v-if="!narrow"
                class="sk chs__btn"
            />
            <div class="chs__desc">
                <span
                    v-for="n in 4"
                    :key="n"
                    class="sk chs__text"
                />
            </div>
            <div
                v-if="narrow"
                class="chs__card"
            >
                <span class="sk chs__card-title" />
                <span
                    v-for="n in 5"
                    :key="n"
                    class="sk chs__row"
                />
            </div>
            <span
                v-if="narrow"
                class="sk chs__btn"
            />
        </div>
        <div
            v-if="!narrow"
            class="chs__card"
        >
            <span class="sk chs__card-title" />
            <span
                v-for="n in 5"
                :key="n"
                class="sk chs__row"
            />
        </div>
    </div>
</template>

<style scoped>
.chs {
    display: grid;
    grid-template-columns: minmax(0, 1fr) 340px;
    gap: 56px;
    align-items: start;
}
.chs__main {
    display: flex;
    flex-direction: column;
    gap: 22px;
    min-width: 0;
}
.chs__badge {
    position: relative;
    align-self: flex-start;
    width: 140px;
    border-color: transparent;
}
.chs__fill {
    position: absolute;
    inset: 0;
    border-radius: var(--r-pill);
}
.chs__title {
    display: flex;
    flex-direction: column;
    gap: 14px;
    padding: 6px 0;
}
.chs__line {
    display: block;
    height: 36px;
    width: 80%;
}
.chs__line--short {
    width: 45%;
}
.chs__meta {
    display: flex;
    gap: 36px;
    padding-top: 16px;
    border-top: 1px solid var(--rule);
    height: 64px;
}
.chs__meta-line {
    display: block;
    width: 140px;
    height: 36px;
}
.chs__btn {
    display: block;
    width: 200px;
    height: 40px;
    border-radius: var(--r-lg);
}
.chs__desc {
    display: flex;
    flex-direction: column;
    gap: 13px;
    max-width: 680px;
}
.chs__text {
    display: block;
    height: 16px;
}
.chs__text:last-child {
    width: 60%;
}
.chs__card {
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-xl);
    padding: 20px 24px;
    display: flex;
    flex-direction: column;
    gap: 18px;
}
.chs__card-title {
    display: block;
    width: 110px;
    height: 22px;
}
.chs__row {
    display: block;
    height: 14px;
    width: 70%;
}

@media (max-width: 900px) {
    .chs {
        grid-template-columns: minmax(0, 1fr);
        gap: 24px;
    }
    .chs__main {
        gap: 18px;
    }
    .chs__line {
        height: 26px;
    }
    .chs__meta {
        gap: 20px;
    }
    .chs__btn {
        width: 100%;
        height: 48px;
    }
}
</style>
