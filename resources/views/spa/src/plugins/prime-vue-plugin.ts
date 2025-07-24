import type { App } from 'vue';

import ToastService from 'primevue/toastservice'
import PrimeVue from 'primevue/config'
import 'primeicons/primeicons.css'
import noirThemes from "./noir-themes.ts";
import 'primeflex/primeflex.css'



export default {
    install(app: App) {
        app.use(PrimeVue, {
            ripple: true,
            inputStyle: 'outlined',
            theme: { 
                preset: noirThemes,
                options: {
                    darkModeSelector: '.dark-theme' // Add this
                }
            }
        })
        app.use(ToastService)
    }
}