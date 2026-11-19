<script lang="ts" setup>
import { computed, ref } from 'vue';
import { RouterLink, useRoute } from 'vue-router';

import AuthButton from './auth/AuthButton.vue';

import { isOwner } from '@/auth/createAuth';

const route = useRoute();
// Link aktywny dla całej sekcji: /challenges i /challenges/:id (RouterLink zaznacza tylko trasy zagnieżdżone).
const LINKS = [
    { to: '/recipes', label: 'Przepisy' },
    { to: '/challenges', label: 'Wyzwania' },
];
const section = computed(
    () => LINKS.find((l) => route.path === l.to || route.path.startsWith(`${l.to}/`))?.to
);

const dark = ref(localStorage.getItem('theme') === 'dark');

function toggleTheme() {
    dark.value = !dark.value;
    const theme = dark.value ? 'dark' : 'light';
    localStorage.setItem('theme', theme);
    document.documentElement.setAttribute('data-theme', theme);
}

if (dark.value) {
    document.documentElement.setAttribute('data-theme', 'dark');
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
                    aria-label="Tryb ciemny"
                    :aria-pressed="dark"
                    @click="toggleTheme"
                >
                    <span
                        class="topbar__theme-label"
                        aria-hidden="true"
                    >
                        {{ dark ? 'dark' : 'light' }}
                    </span>
                    <span
                        class="topbar__theme-dot"
                        aria-hidden="true"
                    >
                        {{ dark ? '☾' : '☀' }}
                    </span>
                </button>

                <!-- admin pracuje na desktopie: na telefonie bez „Panel” i „Zaloguj” -->
                <template v-if="isOwner">
                    <RouterLink
                        to="/admin"
                        class="topbar__admin-badge topbar__desktop"
                    >
                        Panel
                    </RouterLink>
                </template>

                <div class="topbar__desktop">
                    <AuthButton />
                </div>
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
    padding: 14px 56px;
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

.topbar__theme {
    display: flex;
    align-items: center;
    gap: 8px;
    background: var(--bg-alt);
    color: var(--ink);
    border: 1px solid var(--rule);
    border-radius: var(--r-pill);
    padding: 5px 4px 5px 12px;
    font-size: 12px;
    font-family: var(--font-mono);
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

/* Telefon: logo, „Przepisy”, „Wyzwania”, motyw (sama ikona). Cele 44 px; wysokość paska jak dotąd
   (57 px z linią) — przyklejone elementy strony przepisu i listy zakupów liczą od niej. */
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
    .topbar__theme-label,
    .topbar__desktop {
        display: none;
    }
    .topbar__theme {
        width: 44px;
        height: 44px;
        padding: 0;
        justify-content: center;
        background: transparent;
        border: none;
    }
}
</style>
