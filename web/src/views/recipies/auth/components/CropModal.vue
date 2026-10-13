<script lang="ts" setup>
import { computed, onMounted, reactive, ref } from 'vue';

import UModal from '@/components/UModal.vue';

// Kadrowanie do kwadratu: przeciąganie / strzałki, suwak powiększenia. Wynik: 1200×1200 JPEG.
const props = defineProps<{ src: string }>();
const emit = defineEmits<{ cancel: []; done: [blob: Blob] }>();

const V = 360; // rozmiar podglądu w px
const OUT = 1200;

const img = ref<HTMLImageElement>();
const zoom = ref(1);
const pos = reactive({ x: 0, y: 0 });
let drag: { sx: number; sy: number; px: number; py: number } | null = null;

const base = computed(() => (img.value ? V / Math.min(img.value.width, img.value.height) : 1));
const scale = computed(() => base.value * zoom.value);

function clamp(x: number, y: number, sc = scale.value) {
    if (!img.value) return { x, y };
    return {
        x: Math.min(0, Math.max(V - img.value.width * sc, x)),
        y: Math.min(0, Math.max(V - img.value.height * sc, y)),
    };
}

function setPos(x: number, y: number, sc?: number) {
    Object.assign(pos, clamp(x, y, sc));
}

onMounted(() => {
    const i = new Image();
    i.onload = () => {
        img.value = i;
        const s = V / Math.min(i.width, i.height);
        setPos((V - i.width * s) / 2, (V - i.height * s) / 2, s);
    };
    i.src = props.src;
});

function onZoom(e: Event) {
    const z = Number((e.target as HTMLInputElement).value);
    const ns = base.value * z;
    const k = ns / scale.value;
    // powiększenie względem środka kadru
    setPos(V / 2 - (V / 2 - pos.x) * k, V / 2 - (V / 2 - pos.y) * k, ns);
    zoom.value = z;
}

function onDown(e: PointerEvent) {
    (e.currentTarget as HTMLElement).setPointerCapture(e.pointerId);
    drag = { sx: e.clientX, sy: e.clientY, px: pos.x, py: pos.y };
}
function onMove(e: PointerEvent) {
    if (drag) setPos(drag.px + e.clientX - drag.sx, drag.py + e.clientY - drag.sy);
}
function onKey(e: KeyboardEvent) {
    const d: Record<string, [number, number]> = {
        ArrowLeft: [10, 0],
        ArrowRight: [-10, 0],
        ArrowUp: [0, 10],
        ArrowDown: [0, -10],
    };
    const k = d[e.key];
    if (!k) return;
    e.preventDefault();
    setPos(pos.x + k[0], pos.y + k[1]);
}

function confirm() {
    if (!img.value) return;
    const c = document.createElement('canvas');
    c.width = c.height = OUT;
    const s = scale.value;
    c.getContext('2d')!.drawImage(img.value, -pos.x / s, -pos.y / s, V / s, V / s, 0, 0, OUT, OUT);
    c.toBlob((b) => b && emit('done', b), 'image/jpeg', 0.88);
}

const imgStyle = computed(() =>
    img.value
        ? {
              left: `${pos.x}px`,
              top: `${pos.y}px`,
              width: `${img.value.width * scale.value}px`,
              height: `${img.value.height * scale.value}px`,
          }
        : {}
);
</script>

<template>
    <UModal
        title="Wykadruj zdjęcie"
        :width="V + 56"
        @close="emit('cancel')"
    >
        <div
            class="crop"
            tabindex="0"
            aria-label="Obszar kadru. Przeciągnij lub użyj strzałek."
            @keydown="onKey"
            @pointerdown="onDown"
            @pointermove="onMove"
            @pointerup="drag = null"
        >
            <img
                v-if="img"
                :src="src"
                alt=""
                draggable="false"
                class="crop__img"
                :style="imgStyle"
            />
            <div class="crop__grid" />
        </div>
        <label class="crop__zoom">
            Powiększenie
            <input
                type="range"
                min="1"
                max="3"
                step="0.01"
                :value="zoom"
                @input="onZoom"
            />
        </label>
        <template #footer>
            <button
                type="button"
                class="btn btn--ghost"
                @click="emit('cancel')"
            >
                Anuluj
            </button>
            <button
                type="button"
                class="btn btn--accent"
                :disabled="!img"
                @click="confirm"
            >
                Użyj zdjęcia
            </button>
        </template>
    </UModal>
</template>

<style scoped>
.crop {
    width: 360px;
    height: 360px;
    position: relative;
    overflow: hidden;
    border-radius: var(--r-lg);
    background: var(--bg);
    cursor: grab;
    touch-action: none;
    user-select: none;
}
.crop__img {
    position: absolute;
    max-width: none;
    pointer-events: none;
}
.crop__grid {
    position: absolute;
    inset: 0;
    pointer-events: none;
    background-image:
        linear-gradient(rgba(255, 255, 255, 0.35) 1px, transparent 1px),
        linear-gradient(90deg, rgba(255, 255, 255, 0.35) 1px, transparent 1px);
    background-size: 120px 120px;
    background-position: -1px -1px;
}
.crop__zoom {
    display: flex;
    align-items: center;
    gap: 12px;
    font-size: 13px;
    color: var(--ink-muted);
}
.crop__zoom input {
    flex: 1;
    accent-color: var(--accent);
}
</style>
