import { useEffect, useRef, useState } from "react";
import styles from "./Alert.module.scss";

export type AlertVariant = "error" | "success" | "warning" | "info" | "hint";

export interface AlertProps {
    variant: AlertVariant;
    title?: string;
    message: string;
    duration?: number; // ms, 0 = manual close only
    onClose?: () => void;
    closable?: boolean;
}

const ICONS: Record<AlertVariant, JSX.Element> = {
    error: (
        <svg viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
            <path
                fillRule="evenodd"
                d="M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16zm-.75-4.75a.75.75 0 0 0 1.5 0v-4.5a.75.75 0 0 0-1.5 0v4.5zm.75-7.25a1 1 0 1 1 0 2 1 1 0 0 1 0-2z"
                clipRule="evenodd"
            />
        </svg>
    ),
    success: (
        <svg viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
            <path
                fillRule="evenodd"
                d="M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16zm3.857-9.809a.75.75 0 0 0-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 1 0-1.06 1.061l2.5 2.5a.75.75 0 0 0 1.137-.089l4-5.5z"
                clipRule="evenodd"
            />
        </svg>
    ),
    warning: (
        <svg viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
            <path
                fillRule="evenodd"
                d="M8.485 2.495c.673-1.167 2.357-1.167 3.03 0l6.28 10.875c.673 1.167-.17 2.625-1.516 2.625H3.72c-1.347 0-2.189-1.458-1.515-2.625L8.485 2.495zM10 5a.75.75 0 0 1 .75.75v3.5a.75.75 0 0 1-1.5 0v-3.5A.75.75 0 0 1 10 5zm0 9a1 1 0 1 0 0-2 1 1 0 0 0 0 2z"
                clipRule="evenodd"
            />
        </svg>
    ),
    info: (
        <svg viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
            <path
                fillRule="evenodd"
                d="M18 10a8 8 0 1 1-16 0 8 8 0 0 1 16 0zm-7-4a1 1 0 1 1-2 0 1 1 0 0 1 2 0zM9 9a.75.75 0 0 0 0 1.5h.253a.25.25 0 0 1 .244.304l-.459 2.066A1.75 1.75 0 0 0 10.747 15H11a.75.75 0 0 0 0-1.5h-.253a.25.25 0 0 1-.244-.304l.459-2.066A1.75 1.75 0 0 0 9.253 9H9z"
                clipRule="evenodd"
            />
        </svg>
    ),
    hint: (
        <svg viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
            <path d="M10 2a6 6 0 0 0-3.5 10.84V14a1 1 0 0 0 1 1h5a1 1 0 0 0 1-1v-1.16A6 6 0 0 0 10 2zm-1 14h2v1a1 1 0 0 1-2 0v-1z" />
        </svg>
    ),
};

const CLOSE_ICON = (
    <svg viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
        <path d="M6.28 5.22a.75.75 0 0 0-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 1 0 1.06 1.06L10 11.06l3.72 3.72a.75.75 0 1 0 1.06-1.06L11.06 10l3.72-3.72a.75.75 0 0 0-1.06-1.06L10 8.94 6.28 5.22z" />
    </svg>
);

export const Alert = ({
    variant,
    title,
    message,
    duration = 5000,
    onClose,
    closable = true,
}: AlertProps) => {
    const [visible, setVisible] = useState(true);
    const [exiting, setExiting] = useState(false);
    const [progress, setProgress] = useState(100);
    const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
    const startRef = useRef<number>(Date.now());
    const remainingRef = useRef<number>(duration);

    const handleClose = () => {
        setExiting(true);
        setTimeout(() => {
            setVisible(false);
            onClose?.();
        }, 300);
    };

    const startTimer = () => {
        if (!duration) return;
        startRef.current = Date.now();
        timerRef.current = setTimeout(handleClose, remainingRef.current);
    };

    const pauseTimer = () => {
        if (!duration) return;
        if (timerRef.current) clearTimeout(timerRef.current);
        remainingRef.current -= Date.now() - startRef.current;
    };

    useEffect(() => {
        if (!duration) return;

        startTimer();

        let raf: number;
        const tick = () => {
            const elapsed = Date.now() - startRef.current;
            const pct = Math.max(0, 100 - (elapsed / duration) * 100);
            setProgress(pct);
            if (pct > 0) raf = requestAnimationFrame(tick);
        };
        raf = requestAnimationFrame(tick);

        return () => {
            if (timerRef.current) clearTimeout(timerRef.current);
            cancelAnimationFrame(raf);
        };
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    if (!visible) return null;

    const variantClass =
        styles[`alert--${variant}` as keyof typeof styles] ??
        `alert--${variant}`;

    return (
        <div
            className={[
                styles.alert,
                variantClass,
                exiting ? styles["alert--exiting"] : "",
            ]
                .filter(Boolean)
                .join(" ")}
            role={variant === "error" ? "alert" : "status"}
            aria-live={variant === "error" ? "assertive" : "polite"}
            onMouseEnter={pauseTimer}
            onMouseLeave={startTimer}
        >
            <span className={styles["alert__icon"]}>{ICONS[variant]}</span>

            <div className={styles["alert__body"]}>
                {title && <p className={styles["alert__title"]}>{title}</p>}
                <p className={styles["alert__message"]}>{message}</p>
            </div>

            {closable && (
                <button
                    className={styles["alert__close"]}
                    onClick={handleClose}
                    aria-label="Закрыть уведомление"
                    type="button"
                >
                    {CLOSE_ICON}
                </button>
            )}

            {!!duration && (
                <div className={styles["alert__progress-track"]}>
                    <div
                        className={styles["alert__progress-bar"]}
                        style={{ width: `${progress}%` }}
                    />
                </div>
            )}
        </div>
    );
};
