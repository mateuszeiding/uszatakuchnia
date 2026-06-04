import { VueQueryPlugin } from '@tanstack/vue-query';
import { createApp } from 'vue';

import App from './App.vue';
import { createAuth } from './auth/createAuth';
import router from './router';

const app = createApp(App);

app.use(createAuth()).use(router).use(VueQueryPlugin);

app.mount('#app');
