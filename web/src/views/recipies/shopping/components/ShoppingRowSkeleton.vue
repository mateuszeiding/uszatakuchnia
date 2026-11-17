<script lang="ts" setup>
// Szkielet pozycji listy zakupów — siatka, paddingi i wysokości linii jak w ShoppingRow.vue
// (nazwa 15 px × 1,35 + linia „do: …” 13 px; telefon: nazwa 17 px, pole 22 px, min. 56 px).
</script>

<template>
    <li
        class="sk-srow"
        aria-hidden="true"
    >
        <div class="sk-srow__grid">
            <span class="sk-srow__box" />
            <span class="sk-srow__text">
                <span class="sk-srow__name"><span class="sk sk-srow__bar" /></span>
                <span class="sk-srow__for"><span class="sk sk-srow__bar sk-srow__bar--for" /></span>
            </span>
            <span class="sk-srow__sum"><span class="sk sk-srow__bar sk-srow__bar--sum" /></span>
        </div>
    </li>
</template>

<style lang="scss" scoped>
.sk-srow {
    border-bottom: 1px solid var(--rule);
}
.sk-srow__grid {
    display: grid;
    grid-template-columns: 18px minmax(0, 1fr) auto;
    gap: 12px;
    align-items: start;
    align-content: center;
    min-height: 44px;
    padding: 8px 0;
}
.sk-srow__box {
    width: 16px;
    height: 16px;
    margin-top: 2px;
    border-radius: var(--r-sm);
    border: 1.5px solid var(--rule-strong);
}
.sk-srow__name,
.sk-srow__for,
.sk-srow__sum {
    display: flex;
    align-items: center;
}
.sk-srow__name {
    height: calc(15px * 1.35);
}
.sk-srow__for {
    margin-top: 1px;
    font-size: 13px;
    height: 1lh;
}
.sk-srow__sum {
    height: 20px;
}
.sk-srow__bar {
    display: block;
    width: 120px;
    height: 13px;
}
.sk-srow__bar--for {
    width: 90px;
    height: 10px;
}
.sk-srow__bar--sum {
    width: 56px;
}

@media (max-width: 900px) {
    .sk-srow__grid {
        grid-template-columns: 24px minmax(0, 1fr) auto;
        min-height: 56px;
        padding: 10px 0;
    }
    .sk-srow__box {
        width: 22px;
        height: 22px;
        margin-top: 0.5px;
        border-radius: var(--r-md);
    }
    .sk-srow__name {
        height: calc(17px * 1.35);
    }
    .sk-srow__sum {
        height: 23px;
    }
}
</style>
