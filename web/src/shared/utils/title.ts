// Browser tab title: one format everywhere, "<page> · Uszata Kuchnia"; without a page, just "Uszata Kuchnia".
// The server sends the same format to link-preview bots (handlers/share.go).
export const SITE_NAME = 'Uszata Kuchnia';

export const pageTitle = (page?: string | null) => (page ? `${page} · ${SITE_NAME}` : SITE_NAME);

export function setTitle(page?: string | null) {
    document.title = pageTitle(page);
}
