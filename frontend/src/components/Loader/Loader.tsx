import styles from "./Loader.module.scss";

type LoaderSize = "sm" | "md" | "lg";
type LoaderVariant = "spinner" | "dots" | "bar";

interface LoaderProps {
    size?: LoaderSize;
    variant?: LoaderVariant;
    label?: string;
    fullPage?: boolean;
    overlay?: boolean;
    inline?: boolean;
}

const Spinner = ({ sizeClass }: { sizeClass: string }) => (
    <div className={`${styles["loader__spinner"]} ${sizeClass}`} />
);

const Dots = () => (
    <div className={styles["loader__dots"]}>
        <span />
        <span />
        <span />
    </div>
);

const Bar = () => (
    <div className={styles["loader__bar-track"]}>
        <div className={styles["loader__bar-fill"]} />
    </div>
);

const Loader = ({
    size = "md",
    variant = "spinner",
    label,
    fullPage = false,
    overlay = false,
    inline = false,
}: LoaderProps) => {
    const sizeClass = styles[`loader__spinner--${size}`] ?? "";

    const inner = (
        <div className={styles["loader__inner"]}>
            {variant === "spinner" && <Spinner sizeClass={sizeClass} />}
            {variant === "dots" && <Dots />}
            {variant === "bar" && <Bar />}
            {label && <p className={styles["loader__label"]}>{label}</p>}
        </div>
    );

    if (fullPage) {
        return (
            <div
                className={styles["loader--full-page"]}
                role="status"
                aria-live="polite"
            >
                {inner}
                <span className={styles["loader__sr"]}>Загрузка…</span>
            </div>
        );
    }

    if (overlay) {
        return (
            <div
                className={styles["loader--overlay"]}
                role="status"
                aria-live="polite"
            >
                {inner}
                <span className={styles["loader__sr"]}>Загрузка…</span>
            </div>
        );
    }

    if (inline || (!fullPage && !overlay && !label)) {
        return (
            <span
                className={styles["loader--inline"]}
                role="status"
                aria-live="polite"
            >
                {variant === "spinner" && <Spinner sizeClass={sizeClass} />}
                {variant === "dots" && <Dots />}
                {variant === "bar" && <Bar />}
                <span className={styles["loader__sr"]}>Загрузка…</span>
            </span>
        );
    }

    return (
        <div className={styles["loader"]} role="status" aria-live="polite">
            {inner}
            <span className={styles["loader__sr"]}>Загрузка…</span>
        </div>
    );
};

export default Loader;
