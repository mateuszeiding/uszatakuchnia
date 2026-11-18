/// <reference types="vite/client" />

// eslint-disable-next-line @typescript-eslint/naming-convention -- nazwa narzucona przez Vite
interface ImportMetaEnv {
    readonly VITE_AUTH0_DOMAIN: string;
    readonly VITE_AUTH0_CLIENT_ID: string;
    readonly VITE_AUTH0_AUDIENCE: string;
    /** `sub` (user_id) konta właściciela w Auth0 — tylko ono widzi UI admina. */
    readonly VITE_AUTH0_OWNER_SUB?: string;
}
