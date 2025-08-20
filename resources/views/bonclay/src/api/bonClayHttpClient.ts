import {BaseHttpClient} from "@/lib/httpClient/http.ts";
import axios, {type AxiosRequestConfig} from "axios";

export class BonClayHttpClient extends BaseHttpClient{
    constructor() {
        // this.baseURL = baseURL;
        //
        // const defaultConfig: AxiosRequestConfig = {
        //     baseURL,
        //     timeout: 10_000,
        //     headers: { 'Content-Type': 'application/json' },
        //     ...config,
        // };
        //
        // this.instance = axios.create(defaultConfig);
        // this.setupInterceptors();
        super(Axios{})
    }

}