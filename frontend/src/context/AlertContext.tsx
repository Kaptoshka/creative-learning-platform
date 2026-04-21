import {
    createContext,
    useCallback,
    useContext,
    useRef,
    useState,
} from "react";
import { createPortal } from "react-dom";
import { Alert } from "@/components/Alert";
import type { AlertVariant } from "@/components/Alert";
import styles from "./AlertContext.module.scss";

interface AlertItem {
    id: number;
    variant: AlertVariant;
    title?: string;
    message: string;
    duration?: number;
}

interface AlertContextValue {
    show: (opts: Omit<AlertItem, "id">) => void;
    error: (message: string, title?: string) => void;
    success: (message: string, title?: string) => void;
    warning: (message: string, title?: string) => void;
    info: (message: string, title?: string) => void;
    hint: (message: string, title?: string) => void;
}

const AlertContext = createContext<AlertContextValue | null>(null);

export const AlertProvider = ({ children }: { children: React.ReactNode }) => {
    const [alerts, setAlerts] = useState<AlertItem[]>([]);
    const counter = useRef(0);

    const remove = useCallback((id: number) => {
        setAlerts((prev) => prev.filter((a) => a.id !== id));
    }, []);

    const show = useCallback((opts: Omit<AlertItem, "id">) => {
        const id = ++counter.current;
        setAlerts((prev) => [...prev, { ...opts, id }]);
    }, []);

    // shortcuts
    const error = (msg: string, title?: string) =>
        show({ variant: "error", message: msg, title });
    const success = (msg: string, title?: string) =>
        show({ variant: "success", message: msg, title });
    const warning = (msg: string, title?: string) =>
        show({ variant: "warning", message: msg, title });
    const info = (msg: string, title?: string) =>
        show({ variant: "info", message: msg, title });
    const hint = (msg: string, title?: string) =>
        show({ variant: "hint", message: msg, title });

    return (
        <AlertContext.Provider
            value={{ show, error, success, warning, info, hint }}
        >
            {children}
            {createPortal(
                <div className={styles["alert-stack"]}>
                    {alerts.map((a) => (
                        <Alert key={a.id} {...a} onClose={() => remove(a.id)} />
                    ))}
                </div>,
                document.body,
            )}
        </AlertContext.Provider>
    );
};

export function useAlert() {
    const ctx = useContext(AlertContext);
    if (!ctx) throw new Error("useAlert must be used inside <AlertProvider>");
    return ctx;
}
