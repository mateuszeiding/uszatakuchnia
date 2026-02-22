import { fileURLToPath, URL } from 'node:url';

import vue from '@vitejs/plugin-vue';
import { defineConfig } from 'vite';
import vueDevTools from 'vite-plugin-vue-devtools';

// https://vite.dev/config/
export default defineConfig({
    plugins: [vue(), vueDevTools()],
    base: '/',
    server: {
        proxy: {
            '/api': {
                target: 'http://localhost:3000',
                changeOrigin: true,
            },
        },
    },
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url)),
            '@shared': fileURLToPath(new URL('./src/shared', import.meta.url)),
            '@cmp': fileURLToPath(new URL('./src/components', import.meta.url)),
            '@view': fileURLToPath(new URL('./src/views', import.meta.url)),
            '@api': fileURLToPath(new URL('./src/data/api', import.meta.url)),
            '@dto': fileURLToPath(new URL('./src/data/dto', import.meta.url)),
            '@contract': fileURLToPath(new URL('./src/data/contracts', import.meta.url)),
        },
    },
});
