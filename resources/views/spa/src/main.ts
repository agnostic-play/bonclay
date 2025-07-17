import {createApp} from 'vue'
import './style.css';
import App from './App.vue';
import PrimeVuePlugin from "./plugins/prime-vue-plugin.ts";
import router from "@/router";


const app = createApp(App);
app.use(PrimeVuePlugin)
app.use(router)
app.mount('#app')



