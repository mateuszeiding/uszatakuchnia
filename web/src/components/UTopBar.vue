<script lang="ts" setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue';
import { RouterLink, useRoute } from 'vue-router';

import { isOwner } from '@/auth/createAuth';

const route = useRoute();
// Link active for the whole section: /challenges and /challenges/:id (RouterLink marks only nested routes).
const LINKS = [
    { to: '/recipes', label: 'Przepisy' },
    { to: '/challenges', label: 'Wyzwania' },
];
const section = computed(
    () => LINKS.find((l) => route.path === l.to || route.path.startsWith(`${l.to}/`))?.to
);

// Theme: the saved choice, otherwise the system one (prefers-color-scheme). Button state = the visible theme,
// so the first tap always switches to the opposite. Until there's a choice, the button follows the system.
const readSaved = () => {
    try {
        const t = localStorage.getItem('theme');
        return t === 'dark' || t === 'light' ? t : null;
    } catch {
        return null;
    }
};
const system = window.matchMedia?.('(prefers-color-scheme: dark)');
const saved = readSaved();
const dark = ref(saved ? saved === 'dark' : !!system?.matches);
if (saved) document.documentElement.setAttribute('data-theme', saved);

const onSystem = (e: MediaQueryListEvent) => {
    if (!readSaved()) dark.value = e.matches;
};
onMounted(() => system?.addEventListener?.('change', onSystem));
onBeforeUnmount(() => system?.removeEventListener?.('change', onSystem));

function toggleTheme() {
    dark.value = !dark.value;
    const theme = dark.value ? 'dark' : 'light';
    try {
        localStorage.setItem('theme', theme);
    } catch {
        /* private window: the theme works until the end of the visit */
    }
    document.documentElement.setAttribute('data-theme', theme);
}
</script>

<template>
    <header class="topbar">
        <div class="topbar__inner">
            <div class="topbar__left">
                <RouterLink
                    to="/"
                    class="topbar__logo"
                >
                    uszatakuchnia
                    <span class="topbar__logo-dot">.</span>
                </RouterLink>
                <nav
                    class="topbar__nav"
                    aria-label="Główna"
                >
                    <RouterLink
                        v-for="l in LINKS"
                        :key="l.to"
                        :to="l.to"
                        class="topbar__link"
                        :class="{ 'is-active': section === l.to }"
                        :aria-current="section === l.to ? 'page' : undefined"
                    >
                        {{ l.label }}
                    </RouterLink>
                </nav>
            </div>

            <div class="topbar__right">
                <button
                    type="button"
                    class="topbar__theme"
                    :aria-label="dark ? 'Włącz jasny motyw' : 'Włącz ciemny motyw'"
                    :title="dark ? 'Włącz jasny motyw' : 'Włącz ciemny motyw'"
                    @click="toggleTheme"
                >
                    <span
                        class="topbar__theme-dot"
                        aria-hidden="true"
                    >
                        {{ dark ? '☾' : '☀' }}
                    </span>
                </button>

                <!-- the admin works on desktop: "Panel" on desktop only; "Zaloguj" is in the footer -->
                <template v-if="isOwner">
                    <RouterLink
                        to="/admin"
                        class="topbar__admin-badge topbar__desktop"
                    >
                        Panel
                    </RouterLink>
                </template>
            </div>
        </div>
    </header>
</template>

<style scoped>
.topbar {
    position: sticky;
    top: 0;
    z-index: var(--z-sticky);
    background: color-mix(in srgb, var(--bg) 85%, transparent);
    backdrop-filter: blur(12px);
    border-bottom: 1px solid var(--rule);
}

.topbar__inner {
    max-width: 1280px;
    margin: 0 auto;
    /* 9 + 44 (theme button) + 9: keeps the bar height unchanged */
    padding: 9px 56px;
    display: flex;
    align-items: center;
    justify-content: space-between;
}

.topbar__left {
    display: flex;
    align-items: center;
    gap: 32px;
}

.topbar__logo {
    font-size: 17px;
    font-weight: 600;
    letter-spacing: -0.4px;
    color: var(--ink);
    text-decoration: none;
}

.topbar__logo-dot {
    color: var(--accent);
}

.topbar__nav {
    display: flex;
    gap: 22px;
}

.topbar__link {
    font-size: 14px;
    font-weight: 500;
    color: var(--ink-dim);
    text-decoration: none;
    transition: color var(--dur-fast);
}

.topbar__link:hover,
.topbar__link.is-active {
    color: var(--ink);
}

.topbar__right {
    display: flex;
    align-items: center;
    gap: 12px;
}

/* theme: icon only, a 44 px touch target at every width */
.topbar__theme {
    width: 44px;
    height: 44px;
    padding: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    background: transparent;
    border: none;
    border-radius: var(--r-pill);
    color: var(--ink);
    cursor: pointer;
}

.topbar__theme-dot {
    width: 22px;
    height: 22px;
    border-radius: var(--r-pill);
    background: var(--accent-fill);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    color: var(--on-accent);
    font-size: 12px;
}

.topbar__admin-badge {
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 600;
    color: var(--accent);
    background: var(--accent-soft);
    padding: 3px 10px;
    border-radius: var(--r-pill);
    letter-spacing: 0.3px;
}

/* Phone: logo, "Przepisy", "Wyzwania", theme (icon only). 44 px targets; the bar keeps its height
   (57 px with the border): sticky elements of the recipe page and the shopping list are offset from it. */
@media (max-width: 900px) {
    .topbar__inner {
        padding: 6px 16px;
        gap: 12px;
    }
    .topbar__left {
        gap: 16px;
        min-width: 0;
    }
    .topbar__logo {
        font-size: 16px;
        white-space: nowrap;
    }
    .topbar__nav {
        gap: 16px;
    }
    .topbar__link {
        display: inline-flex;
        align-items: center;
        min-height: 44px;
        white-space: nowrap;
    }
    .topbar__desktop {
        display: none;
    }
}
</style>
