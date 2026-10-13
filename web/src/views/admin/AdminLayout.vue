<script lang="ts" setup>
// Panel admina (desktop): nawigacja modułów po lewej, moduł w <RouterView>.
const MODULES = [
    { to: '/admin/recipes', label: 'Przepisy' },
    { to: '/admin/ingredients', label: 'Składniki' },
    { to: '/admin/taxonomy', label: 'Klasyfikacja' },
];
</script>

<template>
    <div class="admin">
        <nav
            class="admin__nav"
            aria-label="Panel"
        >
            <div class="admin__kicker">Panel</div>
            <RouterLink
                v-for="m in MODULES"
                :key="m.to"
                :to="m.to"
                class="admin__link"
                active-class="is-active"
            >
                {{ m.label }}
            </RouterLink>
            <RouterLink
                to="/recipes"
                class="admin__public"
            >
                ← Strona publiczna
            </RouterLink>
        </nav>
        <div class="admin__main">
            <RouterView />
        </div>
    </div>
</template>

<style scoped>
.admin {
    display: grid;
    grid-template-columns: 220px minmax(0, 1fr);
    min-height: calc(100vh - 60px);
}
.admin__nav {
    position: sticky;
    top: 60px;
    align-self: start;
    height: calc(100vh - 60px);
    background: var(--bg-alt);
    border-right: 1px solid var(--rule);
    padding: 22px 12px;
    display: flex;
    flex-direction: column;
    gap: 2px;
}
.admin__kicker {
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--ink-muted);
    letter-spacing: 0.5px;
    text-transform: uppercase;
    font-weight: 600;
    padding: 0 10px 12px;
}
.admin__link {
    padding: 9px 10px;
    border-radius: var(--r-md);
    color: var(--ink);
    font-size: 14px;
    font-weight: 500;
}
.admin__link:hover {
    background: var(--bg);
    text-decoration: none;
}
.admin__link.is-active {
    background: var(--accent-soft);
    color: var(--accent);
    font-weight: 600;
}
.admin__public {
    margin-top: auto;
    padding: 9px 10px;
    font-size: 13px;
    color: var(--ink-muted);
}
.admin__main {
    padding: 32px 40px 64px;
    max-width: 1040px;
    width: 100%;
}
</style>
