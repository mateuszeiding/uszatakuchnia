import type { BeforeSendEvent } from '@vercel/analytics/vue';

// The admin panel is owner-only traffic and would skew the stats.
export function skipAdmin(event: BeforeSendEvent): BeforeSendEvent | null {
    return new URL(event.url).pathname.startsWith('/admin') ? null : event;
}
