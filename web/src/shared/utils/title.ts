// Tytuł karty przeglądarki: jeden format wszędzie — „<strona> · Uszata Kuchnia”, bez strony „Uszata Kuchnia”.
// Ten sam format oddaje serwer botom podglądu linków (handlers/share.go).
export const SITE_NAME = 'Uszata Kuchnia';

export const pageTitle = (page?: string | null) => (page ? `${page} · ${SITE_NAME}` : SITE_NAME);

export function setTitle(page?: string | null) {
    document.title = pageTitle(page);
}
