import type { App } from 'vue';

import ToastService from 'primevue/toastservice'
import PrimeVue from 'primevue/config'
import 'primeicons/primeicons.css'
import noirThemes from "./noir-themes.ts";


export default {
    install(app: App) {
        app.use(PrimeVue, {
            ripple: true,
            inputStyle: 'outlined',
            theme: { preset: noirThemes }
        })
        app.use(ToastService)
    }
}