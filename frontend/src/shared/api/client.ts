import axios, { type AxiosRequestConfig } from "axios";
import { config } from "@/config";
import type { RefreshResponse } from "./types";

export const api = axios.create({
    baseURL: config.APIURL ?? "http://localhost:8000",
    headers: { "Content-Type": "application/json" },
});

// --- Request interceptor: adds access token ---
api.interceptors.request.use((config) => {
    const token = localStorage.getItem("access_token");
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// --- Response interceptor: auto-refresh on 401 ---

type QueueEntry = {
    resolve: (token: string) => void;
    reject: (err: unknown) => void;
};

let isRefreshing = false;
let queue: QueueEntry[] = [];

const flushQueue = (token: string) => {
    queue.forEach(({ resolve }) => resolve(token));
    queue = [];
};

const rejectQueue = (err: unknown) => {
    queue.forEach(({ reject }) => reject(err));
    queue = [];
};

const clearTokens = () => {
    localStorage.removeItem("access_token");
    localStorage.removeItem("refresh_token");
};

api.interceptors.response.use(
    (response) => response,
    async (error) => {
        const original: AxiosRequestConfig & { _retry?: boolean } =
            error.config;

        if (error.response?.status !== 401 || original._retry) {
            return Promise.reject(error);
        }

        if (isRefreshing) {
            return new Promise<string>((resolve, reject) => {
                queue.push({ resolve, reject });
            }).then((token) => {
                original.headers = {
                    ...original.headers,
                    Authorization: `Bearer ${token}`,
                };
                return api(original);
            });
        }

        original._retry = true;
        isRefreshing = true;

        try {
            const refreshToken = localStorage.getItem("refresh_token");
            if (!refreshToken) throw new Error("No refresh token");

            const { data } = await axios.post<RefreshResponse>(
                "/api/v1/auth/refresh",
                { RefreshToken: refreshToken },
            );

            localStorage.setItem("access_token", data.access_token);
            localStorage.setItem("refresh_token", data.refresh_token);

            flushQueue(data.access_token);

            original.headers = {
                ...original.headers,
                Authorization: `Bearer ${data.access_token}`,
            };

            return api(original);
        } catch (err) {
            rejectQueue(err);
            clearTokens();
            window.location.href = "/login";
            return Promise.reject(err);
        } finally {
            isRefreshing = false;
        }
    },
);
