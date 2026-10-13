import { onBeforeUnmount, onMounted, ref } from 'vue';

// Blokada wygaszania ekranu (Wake Lock API) na czas życia komponentu.
// Przeglądarka zwalnia blokadę po schowaniu karty — odnawiamy ją po powrocie.
export function useWakeLock() {
    const active = ref(false);
    let lock: WakeLockSentinel | null = null;
    let alive = true;

    async function acquire() {
        try {
            if (!('wakeLock' in navigator)) return;
            lock = await navigator.wakeLock.request('screen');
            if (!alive) {
                await lock.release();
                return;
            }
            active.value = true;
            lock.addEventListener('release', () => (active.value = false));
        } catch {
            active.value = false;
        }
    }

    const onVisibility = () => {
        if (document.visibilityState === 'visible') acquire();
    };

    onMounted(() => {
        acquire();
        document.addEventListener('visibilitychange', onVisibility);
    });

    onBeforeUnmount(() => {
        alive = false;
        document.removeEventListener('visibilitychange', onVisibility);
        lock?.release().catch(() => undefined);
    });

    return active;
}
