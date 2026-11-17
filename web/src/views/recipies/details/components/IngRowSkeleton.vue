<script lang="ts" setup>
// Szkielet wiersza składnika — siatka, paddingi i wysokość linii jak w IngRow.vue
// (18 px pole | 92 px ilość | nazwa; 8 px / telefon 11 px; nazwa 15 px × 1,4).
defineProps<{ wide?: boolean }>();
</script>

<template>
    <li
        class="sk-ing"
        aria-hidden="true"
    >
        <span class="sk-ing__box" />
        <span class="sk-ing__cell"><span class="sk sk-ing__qty" /></span>
        <span class="sk-ing__cell">
            <span
                class="sk sk-ing__name"
                :class="{ 'sk-ing__name--wide': wide }"
            />
        </span>
    </li>
</template>

<style scoped>
.sk-ing {
    display: grid;
    grid-template-columns: 18px 92px minmax(0, 1fr);
    gap: 10px;
    padding: 8px 0;
    border-bottom: 1px solid var(--rule);
}
.sk-ing__box {
    width: 16px;
    height: 16px;
    margin-top: 2.5px;
    border-radius: var(--r-sm);
    border: 1.5px solid var(--rule-strong);
}
/* wysokość linii jak w IngRow (nazwa 15 px × 1,4 wyrównana do linii bazowej z ilością mono = 23 px) */
.sk-ing__cell {
    height: 23px;
    display: flex;
    align-items: center;
}
.sk-ing__qty {
    width: 52px;
    height: 12px;
}
.sk-ing__name {
    width: 60%;
    height: 13px;
}
.sk-ing__name--wide {
    width: 82%;
}

@media (max-width: 900px) {
    .sk-ing {
        padding: 11px 0;
    }
}
</style>
