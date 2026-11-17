<script lang="ts" setup>
// Szkielet panelu filtrów przy pierwszym wejściu (słowniki jeszcze nie w cache):
// 4 sekcje — pasek nagłówka 10 px + 2 rzędy chipów 24 px, żeby panel nie rósł skokowo.
const ROWS = [3, 2];
</script>

<template>
    <div
        class="fsk"
        aria-hidden="true"
    >
        <div
            v-for="s in 4"
            :key="s"
            class="fsk__group"
        >
            <div class="sk fsk__head" />
            <div
                v-for="(n, ri) in ROWS"
                :key="ri"
                class="fsk__row"
            >
                <div
                    v-for="c in n"
                    :key="c"
                    class="sk fsk__chip"
                />
            </div>
        </div>
    </div>
</template>

<style scoped>
.fsk__group {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-bottom: 22px;
}
.fsk__head {
    height: 10px;
    width: 88px;
    margin-bottom: 4px;
}
.fsk__group:nth-child(even) .fsk__head {
    width: 64px;
}
.fsk__row {
    display: flex;
    gap: 6px;
}
.fsk__chip {
    height: 24px;
    width: 64px;
    border-radius: var(--r-pill);
}
.fsk__chip:nth-child(2) {
    width: 52px;
}
.fsk__chip:nth-child(3) {
    width: 78px;
}
.fsk__row:nth-child(odd) .fsk__chip:first-child {
    width: 86px;
}
</style>
