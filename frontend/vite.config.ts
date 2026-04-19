import { defineConfig } from "vite";
import path from "path";
export default defineConfig({
    root: "src",
    build: {
        rollupOptions: {
            input: "./src/index.html",
        },
    },
    server: {
        host: true,
        port: 3000,
    },
    resolve: {
        alias: {
            "@": path.resolve(__dirname, "./src"),
        },
    },
    define: {
        "process.env.NODE_ENV": JSON.stringify(
            process.env.NODE_ENV || "development",
        ),
        "process.env.BUN_PUBLIC_SSO_API_URL": JSON.stringify(
            process.env.BUN_PUBLIC_SSO_API_URL || "http://localhost:8080",
        ),
        "process.env.BUN_PUBLIC_APP_ID": JSON.stringify(
            process.env.BUN_PUBLIC_APP_ID || "1",
        ),
    },
    css: {
        modules: {
            localsConvention: "camelCase",
        },
    },
});
