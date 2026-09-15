import { defineConfig, loadEnv } from "vite";
import path from "path";

export default defineConfig(({ mode }) => {
    const env = loadEnv(mode, process.cwd());

    return {
        root: "src",
        build: {
            outDir: "../dist",
            emptyOutDir: true,
            rollupOptions: {
                input: "./index.html",
            },
        },
        server: {
            host: true,
            port: 3000,
        },
        proxy: {
            "/api/v1": {
                target: "http://gateway-service:8000/api/v1",
                changeOrigin: true,
            },
        },
        resolve: {
            alias: {
                "@": path.resolve(__dirname, "./src"),
            },
        },
        define: {
            "import.meta.env.VITE_API_URL": JSON.stringify(
                env.VITE_API_URL || "/",
            ),
            "import.meta.env.VITE_APP_ID": JSON.stringify(
                env.VITE_APP_ID || "1",
            ),
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
            preprocessorOptions: {
                scss: {
                    api: "modern-compiler",
                },
            },
        },
    };
});
