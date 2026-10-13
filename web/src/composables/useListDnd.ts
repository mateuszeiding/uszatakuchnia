import { type Ref, ref } from 'vue';

// Przeciąganie elementów między grupami (HTML5 DnD, mysz) + zmiana kolejności strzałkami ↑ ↓ na uchwycie.

interface IWithId {
    id: number;
}
interface IGroupLike<T extends IWithId> {
    id: number;
    items: T[];
}

type Pos = 'before' | 'after';
interface IOver {
    gid: number;
    iid: number | null; // null = koniec grupy
    pos: Pos;
}

export function useListDnd<T extends IWithId>(groups: Ref<IGroupLike<T>[]>) {
    const drag = ref<{ gid: number; iid: number } | null>(null);
    const over = ref<IOver | null>(null);

    const reset = () => {
        drag.value = null;
        over.value = null;
    };

    function moveTo(fromIid: number, to: IOver) {
        let moved: T | undefined;
        for (const g of groups.value) {
            const idx = g.items.findIndex((it) => it.id === fromIid);
            if (idx >= 0) [moved] = g.items.splice(idx, 1);
        }
        if (!moved) return;
        const target = groups.value.find((g) => g.id === to.gid);
        if (!target) return;
        const idx =
            to.iid == null
                ? target.items.length
                : target.items.findIndex((it) => it.id === to.iid) + (to.pos === 'after' ? 1 : 0);
        target.items.splice(idx, 0, moved);
    }

    function onDrop(e: DragEvent) {
        e.preventDefault();
        if (drag.value && over.value && over.value.iid !== drag.value.iid)
            moveTo(drag.value.iid, over.value);
        reset();
    }

    /** Atrybuty uchwytu (przycisk z ikoną). */
    const handle = (gid: number, iid: number) => ({
        draggable: true,
        onDragstart: (e: DragEvent) => {
            if (!e.dataTransfer) return;
            e.dataTransfer.effectAllowed = 'move';
            e.dataTransfer.setData('text/plain', '');
            const row = (e.currentTarget as HTMLElement).closest('[data-row]');
            if (row) e.dataTransfer.setDragImage(row, 16, 16);
            drag.value = { gid, iid };
        },
        onDragend: reset,
        onKeydown: (e: KeyboardEvent) => {
            if (e.key !== 'ArrowUp' && e.key !== 'ArrowDown') return;
            e.preventDefault();
            const g = groups.value.find((x) => x.id === gid);
            if (!g) return;
            const i = g.items.findIndex((x) => x.id === iid);
            const j = i + (e.key === 'ArrowUp' ? -1 : 1);
            if (j < 0 || j >= g.items.length) return;
            [g.items[i], g.items[j]] = [g.items[j]!, g.items[i]!];
            const btn = e.currentTarget as HTMLElement;
            requestAnimationFrame(() => btn.focus());
        },
    });

    /** Atrybuty wiersza (cel upuszczenia). */
    const row = (gid: number, iid: number) => ({
        'data-row': true,
        onDragover: (e: DragEvent) => {
            if (!drag.value) return;
            e.preventDefault();
            const r = (e.currentTarget as HTMLElement).getBoundingClientRect();
            const pos: Pos = e.clientY < r.top + r.height / 2 ? 'before' : 'after';
            if (over.value?.iid !== iid || over.value?.pos !== pos) over.value = { gid, iid, pos };
        },
        onDrop,
    });

    /** Atrybuty pola na końcu grupy (upuszczenie do pustej grupy / na koniec). */
    const groupEnd = (gid: number) => ({
        onDragover: (e: DragEvent) => {
            if (!drag.value) return;
            e.preventDefault();
            if (over.value?.gid !== gid || over.value?.iid !== null)
                over.value = { gid, iid: null, pos: 'after' };
        },
        onDrop,
    });

    const lineAt = (gid: number, iid: number): Pos | null =>
        drag.value && over.value?.gid === gid && over.value.iid === iid && iid !== drag.value.iid
            ? over.value.pos
            : null;

    const endLine = (gid: number) =>
        !!drag.value && over.value?.gid === gid && over.value.iid === null;

    const isDragged = (iid: number) => drag.value?.iid === iid;

    return { handle, row, groupEnd, lineAt, endLine, isDragged };
}
