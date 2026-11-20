<script lang="ts" setup>
import { useAuth0 } from '@auth0/auth0-vue';

// Stopka stron publicznych (20.11): po lewej „uszatakuchnia.”, po prawej „Zaloguj” — tylko na desktopie
// (admin pracuje na desktopie). Zalogowany widzi „Wyloguj”. Bez tekstów o koncie i „czyje to”.
const { isAuthenticated, isLoading, loginWithRedirect, logout } = useAuth0();

function onAuth() {
    if (isAuthenticated.value) logout({ logoutParams: { returnTo: globalThis.location.origin } });
    else loginWithRedirect();
}
</script>

<template>
    <footer class="footer">
        <div class="footer__inner">
            <span>
                uszatakuchnia
                <span class="footer__dot">.</span>
            </span>
            <button
                type="button"
                class="footer__auth"
                :disabled="isLoading"
                @click="onAuth"
            >
                {{ isAuthenticated ? 'Wyloguj' : 'Zaloguj' }}
            </button>
        </div>
    </footer>
</template>

<style scoped>
/* prototyp pz-main.jsx (Footer) */
.footer {
    margin-top: 72px;
    border-top: 1px solid var(--rule);
}
.footer__inner {
    max-width: 1280px;
    margin: 0 auto;
    padding: 28px 56px 36px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
    font-size: 13px;
    color: var(--ink-muted);
}
.footer__dot {
    color: var(--accent);
}
.footer__auth {
    padding: 0;
    border: none;
    background: none;
    font: inherit;
    color: var(--ink-muted);
    text-decoration: underline;
    text-underline-offset: 3px;
    cursor: pointer;
}
.footer__auth:hover {
    color: var(--ink);
}
.footer__auth:disabled {
    cursor: default;
}

@media (max-width: 900px) {
    .footer {
        margin-top: 48px;
    }
    .footer__inner {
        padding: 24px 20px 32px;
    }
    .footer__auth {
        display: none;
    }
}
</style>
