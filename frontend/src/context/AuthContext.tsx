import React, { createContext } from "react";

interface User {
    id: string;
    email: string;
    role: string;
}

interface AuthContextType {
    user: User | null;
    loading: boolean;
    login: (accessToken: string, refreshToken: string) => void;
    logout: () => Promise<void>;
    isAuthenticated: () => boolean;
}

export const AuthContext = createContext<AuthContextType | null>(null);
