<script lang="ts" setup>
import { useQueryClient } from '@tanstack/vue-query';
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import { onBeforeRouteLeave, onBeforeRouteUpdate, useRouter } from 'vue-router';

import {
    formFromChallenge,
    type IChallengeForm,
    mapServerErrors,
    snapshot,
    toPayload,
    validateForm,
} from './challengeForm';
import ChallengeBasics from './components/ChallengeBasics.vue';
import ChallengeIngredientsEditor from './components/ChallengeIngredientsEditor.vue';

import BusyLabel from '@/components/BusyLabel.vue';
import CautionsEditor from '@/components/CautionsEditor.vue';
import UFieldError from '@/components/form/UFieldError.vue';
import USpinner from '@/components/USpinner.vue';
import { confirmAction } from '@/composables/useConfirm';
import { showToast } from '@/composables/useToast';
import {
    createChallenge,
    deleteChallenge,
    fetchChallengeWithClient,
    invalidateChallenges,
    updateChallenge,
} from '@/data/api/challenges/fetch';
import { fetchDictionariesWithClient, invalidateDictionaries } from '@/data/api/taxonomy/fetch';
import type { IChallenge, IChallengeFormat } from '@/data/dtos/challenge/ChallengeDto';

const props = defineProps<{ id?: string }>();
const LIST = '/admin/challenges';

const qc = useQueryClient();
const router = useRouter();

const challenge = ref<IChallenge>();
const formats = ref<IChallengeFormat[]>([]);
const loading = ref(true);
const notFound = ref(false);
const form = ref<IChallengeForm>(formFromChallenge());
const saved = ref(snapshot(form.value));
const errors = ref<Record<string, string>>({});
const serverError = ref('');
const saving = ref(false);
let leaving = false;

const challengeId = computed(() => (props.id ? Number(props.id) : null));

async function load() {
    loading.value = true;
    notFound.value = false;
    challenge.value = undefined;
    try {
        const [dicts, c] = await Promise.all([
            fetchDictionariesWithClient(qc),
            challengeId.value != null ? fetchChallengeWithClient(qc, challengeId.value) : undefined,
        ]);
        formats.value = dicts.challengeFormats;
        challenge.value = c;
    } catch {
        notFound.value = challengeId.value != null;
    }
    form.value = formFromChallenge(challenge.value);
    saved.value = snapshot(form.value);
    errors.value = {};
    serverError.value = '';
    loading.value = false;
}
watch(challengeId, load, { immediate: true });

const dirty = computed(() => !loading.value && snapshot(form.value) !== saved.value);

async function afterChange() {
    // format usage counts (Classification) and the public list
    await Promise.all([invalidateChallenges(qc), invalidateDictionaries(qc)]);
}

async function save() {
    errors.value = validateForm(form.value);
    serverError.value = '';
    if (Object.keys(errors.value).length) {
        serverError.value = 'Popraw zaznaczone pola.';
        return;
    }
    saving.value = true;
    try {
        const payload = toPayload(form.value);
        if (challengeId.value != null) await updateChallenge(challengeId.value, payload);
        else await createChallenge(payload);
        await afterChange();
        showToast(challengeId.value != null ? 'Wyzwanie zapisane' : 'Wyzwanie dodane');
        leaving = true;
        router.push(LIST);
    } catch (e) {
        const res = (
            e as {
                response?: {
                    status?: number;
                    data?: { fields?: { field: string; message: string }[] };
                };
            }
        ).response;
        if (res?.status === 422 && res.data?.fields) {
            errors.value = mapServerErrors(res.data.fields, form.value);
            serverError.value = 'Popraw zaznaczone pola.';
        } else {
            serverError.value = 'Nie udało się zapisać wyzwania.';
        }
    } finally {
        saving.value = false;
    }
}

async function remove() {
    const c = challenge.value;
    if (!c) return;
    const ok = await confirmAction({
        title: `Usunąć „${c.title}”?`,
        body: 'Wyzwanie zniknie ze strony. Tej operacji nie da się cofnąć.',
        confirmLabel: 'Usuń wyzwanie',
        danger: true,
    });
    if (!ok) return;
    try {
        await deleteChallenge(c.id);
        await afterChange();
        showToast('Wyzwanie usunięte');
        leaving = true;
        router.push(LIST);
    } catch {
        showToast('Nie udało się usunąć wyzwania');
    }
}

const onBeforeUnload = (e: BeforeUnloadEvent) => {
    if (!dirty.value) return;
    e.preventDefault();
    e.returnValue = '';
};
onMounted(() => window.addEventListener('beforeunload', onBeforeUnload));
onBeforeUnmount(() => window.removeEventListener('beforeunload', onBeforeUnload));

const guard = () => {
    if (leaving || !dirty.value) return true;
    return confirmAction({
        title: 'Wyjść bez zapisywania?',
        body: 'Zmiany w wyzwaniu przepadną.',
        confirmLabel: 'Wyjdź',
        danger: true,
    });
};
onBeforeRouteLeave(guard);
onBeforeRouteUpdate(guard);
</script>

<template>
    <nav
        class="crumbs"
        aria-label="Ścieżka"
    >
        <RouterLink :to="LIST">Wyzwania</RouterLink>
        <span aria-hidden="true">/</span>
        <span class="crumbs__current">
            <span
                v-if="loading && challengeId != null"
                class="sk crumbs__sk"
            />
            {{ challengeId != null ? challenge?.title : 'Nowe wyzwanie' }}
        </span>
    </nav>
    <h1 class="admin-title">{{ challengeId != null ? 'Edytuj wyzwanie' : 'Nowe wyzwanie' }}</h1>

    <p
        v-if="!loading && notFound"
        class="edit__missing"
    >
        Nie ma takiego wyzwania.
        <RouterLink :to="LIST">Wróć do listy</RouterLink>
    </p>
    <form
        v-else
        class="edit"
        :inert="loading"
        :aria-busy="loading"
        novalidate
        @submit.prevent="save"
    >
        <div
            v-if="loading"
            class="edit__loading"
            role="status"
            aria-label="Wczytywanie wyzwania"
        >
            <USpinner :size="22" />
        </div>
        <section class="edit__section">
            <div class="edit__head">
                <span class="edit__no">01</span>
                <h2 class="edit__title">Podstawowe</h2>
            </div>
            <ChallengeBasics
                v-model="form"
                :formats="formats"
                :errors="errors"
            />
        </section>

        <section class="edit__section">
            <div class="edit__head">
                <span class="edit__no">02</span>
                <h2 class="edit__title">Zadanie</h2>
            </div>
            <div class="field">
                <label
                    for="ch-desc"
                    class="field-label"
                >
                    Opis
                    <span
                        class="req"
                        aria-hidden="true"
                    >
                        *
                    </span>
                </label>
                <textarea
                    id="ch-desc"
                    v-model="form.description"
                    v-autosize
                    class="textarea textarea--auto"
                    :class="{ 'is-invalid': errors.description }"
                    rows="4"
                    required
                    aria-describedby="ch-desc-hint"
                />
                <UFieldError v-if="errors.description">{{ errors.description }}</UFieldError>
                <span
                    id="ch-desc-hint"
                    class="field-hint"
                >
                    Akapity oddzielaj pustą linią.
                </span>
            </div>
        </section>

        <section class="edit__section">
            <div class="edit__head">
                <span class="edit__no">03</span>
                <h2 class="edit__title">Uwagi</h2>
                <span class="edit__opt">opcjonalne</span>
            </div>
            <CautionsEditor
                v-model="form.cautions"
                label="wyzwaniu"
                :errors="errors"
            />
        </section>

        <section class="edit__section">
            <div class="edit__head">
                <span class="edit__no">04</span>
                <h2 class="edit__title">Składniki</h2>
                <span class="edit__opt">opcjonalne</span>
            </div>
            <p class="edit__hint">
                Zwykły tekst („ryba białego mięsa”, „coś kwaśnego”). Puste = bez listy (np. format
                „na temat”). „Obowiązkowy” = trafia do „Musisz użyć”.
            </p>
            <ChallengeIngredientsEditor
                v-model="form.ingredients"
                :errors="errors"
            />
        </section>

        <div class="edit__bar">
            <button
                v-if="challengeId != null"
                type="button"
                class="btn btn--danger"
                @click="remove"
            >
                Usuń wyzwanie
            </button>
            <span
                v-if="serverError"
                class="edit__bar-error"
                role="alert"
            >
                {{ serverError }}
            </span>
            <span class="edit__bar-right">
                <RouterLink
                    :to="LIST"
                    class="btn btn--ghost"
                >
                    Anuluj
                </RouterLink>
                <button
                    type="submit"
                    class="btn btn--accent"
                    :class="{ 'btn--saving': saving }"
                    :disabled="saving"
                >
                    <BusyLabel
                        :busy="saving"
                        label="Zapisz wyzwanie"
                        busy-label="Zapisywanie…"
                        spinner-color="#fff"
                    />
                </button>
            </span>
        </div>
    </form>
</template>

<style scoped>
.crumbs {
    display: flex;
    gap: 8px;
    align-items: center;
    font-size: 13px;
    color: var(--ink-muted);
}
.crumbs a {
    color: var(--ink-muted);
}
.crumbs__current {
    color: var(--ink);
    font-weight: 600;
}
.crumbs__sk {
    display: inline-block;
    width: 120px;
    height: 11px;
    vertical-align: middle;
}
.admin-title {
    font-size: 28px;
    font-weight: 600;
    letter-spacing: -0.8px;
    margin: 10px 0 22px;
}
.edit {
    position: relative;
    display: flex;
    flex-direction: column;
    gap: 16px;
    max-width: 820px;
}
.edit__loading {
    position: absolute;
    left: 0;
    right: 0;
    top: 120px;
    z-index: 1;
    display: flex;
    justify-content: center;
    pointer-events: none;
}
.edit__section {
    background: var(--bg-alt);
    border: 1px solid var(--rule);
    border-radius: var(--r-lg);
    padding: 24px;
}
.edit__head {
    display: flex;
    align-items: baseline;
    gap: 12px;
    margin-bottom: 16px;
}
.edit__head:has(+ .edit__hint) {
    margin-bottom: 6px;
}
.edit__no {
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 600;
    color: var(--accent);
}
.edit__title {
    margin: 0;
    font-size: 20px;
    font-weight: 600;
    letter-spacing: -0.5px;
}
.edit__opt {
    font-family: var(--font-mono);
    font-size: 12px;
    font-weight: 500;
    letter-spacing: 0.5px;
    text-transform: uppercase;
    color: var(--ink-muted);
}
.edit__hint {
    margin: 0 0 16px;
    font-size: 13px;
    line-height: 1.45;
    color: var(--ink-muted);
}
.edit__missing {
    color: var(--ink-muted);
}
.edit__bar {
    position: sticky;
    bottom: 0;
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 12px 0;
    background: var(--bg);
    border-top: 1px solid var(--rule);
}
.edit__bar-error {
    font-size: 13px;
    font-weight: 500;
    color: var(--danger);
}
.edit__bar-right {
    margin-left: auto;
    display: flex;
    gap: 10px;
}
</style>
