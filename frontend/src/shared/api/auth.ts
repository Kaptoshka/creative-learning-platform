import { api } from "./client";
import type {
    LoginRequest,
    LoginResponse,
    LogoutRequest,
    RefreshResponse,
    RegisterRequest,
    RegisterResponse,
} from "./types";

export const authApi = {
    register(data: RegisterRequest) {
        return api
            .post<RegisterResponse>("/auth/register", data)
            .then((res) => res.data);
    },

    login(data: LoginRequest) {
        return api
            .post<LoginResponse>("/auth/login", data)
            .then((res) => {
                localStorage.setItem("access_token", res.data.access_token);
                localStorage.setItem("refresh_token", res.data.refresh_token);
                return res.data;
            });
    },

    refresh(refreshToken: string) {
        return api
            .post<RefreshResponse>("/auth/refresh", {
                RefreshToken: refreshToken,
            } satisfies { RefreshToken: string })
            .then((res) => {
                localStorage.setItem("access_token", res.data.access_token);
                localStorage.setItem("refresh_token", res.data.refresh_token);
                return res.data;
            });
    },

    logout() {
        const refreshToken = localStorage.getItem("refresh_token");
        if (!refreshToken) return Promise.resolve();
        return api
            .post("/auth/logout", {
                refresh_token: refreshToken,
            } satisfies LogoutRequest)
            .then(() => {
                localStorage.removeItem("access_token");
                localStorage.removeItem("refresh_token");
            });
    },

    logoutAll(userId: string) {
        return api
            .post("/auth/logout-all", { UserID: userId })
            .then((res) => res.data);
    },
};
