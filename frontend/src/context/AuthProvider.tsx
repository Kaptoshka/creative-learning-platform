import React, { useState, useEffect, useCallback } from "react";
import { jwtDecode } from "jwt-decode";
import { AuthContext } from "./AuthContext";
import { authApi } from "@/shared/api/auth";

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
    children,
}) => {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);

    const decodeAndSetUser = (token: string) => {
        try {
            const decoded: any = jwtDecode(token);
            setUser({
                id: decoded.uid,
                email: decoded.email,
                role: decoded.role,
            });
        } catch (err) {
            console.error("Failed to decode token:", err);
            setUser(null);
        }
    };
    useEffect(() => {
        const initAuth = () => {
            const token = localStorage.getItem("access_token");
            if (token) {
                decodeAndSetUser(token);
            }
            setLoading(false);
        };

        initAuth();
    }, []);

    const login = (accessToken: string, refreshToken: string) => {
        localStorage.setItem("access_token", accessToken);
        localStorage.setItem("refresh_token", refreshToken);

        decodeAndSetUser(accessToken);
    };

    const logout = useCallback(async () => {
        try {
            await authApi.logout();
        } catch (err) {
            console.error("Logout API error:", err);
        } finally {
            localStorage.removeItem("access_token");
            localStorage.removeItem("refresh_token");
            setUser(null);
        }
    }, []);

    const isAuthenticated = () => {
        return !!user && !!localStorage.getItem("access_token");
    };

    return (
        <AuthContext.Provider
            value={{
                user,
                loading,
                login,
                logout,
                isAuthenticated,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
};
