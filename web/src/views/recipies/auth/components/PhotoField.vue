<script lang="ts" setup>
import { computed, ref } from 'vue';

import CropModal from './CropModal.vue';

import UFieldError from '@/components/form/UFieldError.vue';
import { API } from '@/data/api/API';
import { uploadPhoto } from '@/data/api/photos/fetch';

// Zdjęcie 1:1: walidacja pliku → kadrowanie → wysyłka z paskiem postępu. v-model = adres w API.
const photo = defineModel<string | null>({ required: true });

const MAX_MB = 10;
const TYPES = /^image\/(jpeg|png|webp)$/;

const fileInput = ref<HTMLInputElement>();
const cropSrc = ref<string | null>(null);
const error = ref('');
const progress = ref<number | null>(null);
const preview = computed(() => API.resolveUrl(photo.value));

function pick(e: Event) {
    const input = e.target as HTMLInputElement;
    const f = input.files?.[0];
    input.value = '';
    if (!f) return;
    if (!TYPES.test(f.type)) {
        error.value = 'Obsługiwane formaty: JPG, PNG, WEBP.';
        return;
    }
    if (f.size > MAX_MB * 1024 * 1024) {
        error.value = `Plik ma ${(f.size / 1048576).toFixed(1)} MB. Limit to ${MAX_MB} MB.`;
        return;
    }
    error.value = '';
    cropSrc.value = URL.createObjectURL(f);
}

function closeCrop() {
    if (cropSrc.value) URL.revokeObjectURL(cropSrc.value);
    cropSrc.value = null;
}

async function onCropped(blob: Blob) {
    closeCrop();
    progress.value = 0;
    try {
        const res = await uploadPhoto(blob, (p) => (progress.value = p));
        photo.value = res.url;
    } catch {
        error.value = 'Nie udało się wysłać zdjęcia. Spróbuj ponownie.';
    } finally {
        progress.value = null;
    }
}
</script>

<template>
    <div class="photo-field">
        <div
            class="photo-field__box"
            :class="{ 'has-photo': preview, 'is-invalid': error }"
        >
            <img
                v-if="preview"
                :src="preview"
                alt="Zdjęcie przepisu"
                class="photo-field__img"
            />
            <button
                v-else
                type="button"
                class="photo-field__empty"
                @click="fileInput?.click()"
            >
                <svg
                    width="28"
                    height="28"
                    viewBox="0 0 28 28"
                    aria-hidden="true"
                >
                    <rect
                        x="3"
                        y="5"
                        width="22"
                        height="18"
                        rx="2"
                        stroke="currentColor"
                        stroke-width="1.4"
                    />
                    <circle
                        cx="9"
                        cy="11"
                        r="2"
                        stroke="currentColor"
                        stroke-width="1.4"
                    />
                    <path
                        d="M3 19l6-5 6 4 4-3 6 5"
                        stroke="currentColor"
                        stroke-width="1.4"
                        stroke-linejoin="round"
                    />
                </svg>
                <span class="photo-field__cta">Dodaj zdjęcie</span>
                <span class="photo-field__formats">JPG · PNG · WEBP · do {{ MAX_MB }} MB</span>
            </button>
            <div
                v-if="preview && progress == null"
                class="photo-field__actions"
            >
                <button
                    type="button"
                    class="photo-field__overlay-btn"
                    @click="fileInput?.click()"
                >
                    Zmień
                </button>
                <button
                    type="button"
                    class="photo-field__overlay-btn"
                    @click="photo = null"
                >
                    Usuń
                </button>
            </div>
            <div
                v-if="progress != null"
                class="photo-field__progress"
            >
                <span>Wysyłanie… {{ progress }}%</span>
                <div
                    class="photo-field__bar"
                    role="progressbar"
                    aria-label="Postęp wysyłania zdjęcia"
                    :aria-valuenow="progress"
                    aria-valuemin="0"
                    aria-valuemax="100"
                >
                    <div
                        class="photo-field__bar-fill"
                        :style="{ width: `${progress}%` }"
                    />
                </div>
            </div>
        </div>
        <input
            ref="fileInput"
            type="file"
            accept="image/jpeg,image/png,image/webp"
            class="sr-only"
            aria-label="Wybierz zdjęcie"
            tabindex="-1"
            @change="pick"
        />
        <UFieldError v-if="error">{{ error }}</UFieldError>
        <span
            v-else
            class="field-hint"
        >
            Kadrujesz do kwadratu po wybraniu.
        </span>
        <CropModal
            v-if="cropSrc"
            :src="cropSrc"
            @cancel="closeCrop"
            @done="onCropped"
        />
    </div>
</template>

<style scoped>
.photo-field {
    display: flex;
    flex-direction: column;
    gap: 8px;
}
.photo-field__box {
    aspect-ratio: 1 / 1;
    width: 100%;
    position: relative;
    overflow: hidden;
    border-radius: var(--r-lg);
    border: 1.5px dashed var(--rule-strong);
    background: repeating-linear-gradient(45deg, var(--bg) 0 8px, var(--bg-alt) 8px 9px);
}
.photo-field__box.has-photo {
    border: 1px solid var(--rule);
}
.photo-field__box.is-invalid {
    border-color: var(--danger);
}
.photo-field__img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}
.photo-field__empty {
    position: absolute;
    inset: 0;
    background: transparent;
    border: none;
    cursor: pointer;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 10px;
    color: var(--ink-muted);
}
.photo-field__cta {
    font-size: 14px;
    font-weight: 600;
    color: var(--ink);
}
.photo-field__formats {
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--ink-muted);
}
.photo-field__actions {
    position: absolute;
    left: 8px;
    right: 8px;
    bottom: 8px;
    display: flex;
    gap: 6px;
}
.photo-field__overlay-btn {
    flex: 1;
    background: rgba(10, 10, 15, 0.72);
    color: #fff;
    border: none;
    border-radius: var(--r-md);
    padding: 7px 0;
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
    backdrop-filter: blur(6px);
}
.photo-field__progress {
    position: absolute;
    inset: 0;
    background: rgba(10, 10, 15, 0.45);
    display: flex;
    flex-direction: column;
    justify-content: flex-end;
    padding: 12px;
    gap: 6px;
    font-family: var(--font-mono);
    font-size: 12px;
    color: #fff;
}
.photo-field__bar {
    height: 4px;
    background: rgba(255, 255, 255, 0.3);
    border-radius: var(--r-pill);
    overflow: hidden;
}
.photo-field__bar-fill {
    height: 100%;
    background: #fff;
    transition: width 0.09s linear;
}
</style>
