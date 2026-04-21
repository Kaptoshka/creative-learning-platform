import React, { useRef, useEffect, useState, useCallback } from "react";
import styles from "./ThirtyCirclesTask.module.scss";

interface CircleData {
    id: number;
    drawing: string | null; // base64 SVG path data
    isEmpty: boolean;
}

interface ThirtyCirclesTaskProps {
    content: {
        type: string;
        timeLimit?: number; // seconds, default 180
        title?: string;
        description?: string;
    };
    onContentChange: (content: unknown) => void;
}

const TOTAL_CIRCLES = 30;
const DEFAULT_TIME = 180; // 3 minutes

const ThirtyCirclesTask: React.FC<ThirtyCirclesTaskProps> = ({
    content,
    onContentChange,
}) => {
    const timeLimit = content?.timeLimit ?? DEFAULT_TIME;

    const [circles, setCircles] = useState<CircleData[]>(
        Array.from({ length: TOTAL_CIRCLES }, (_, i) => ({
            id: i,
            drawing: null,
            isEmpty: true,
        })),
    );
    const [activeCircle, setActiveCircle] = useState<number | null>(null);
    const [phase, setPhase] = useState<"ready" | "drawing" | "review">("ready");
    const [timeLeft, setTimeLeft] = useState(timeLimit);
    const [isDrawing, setIsDrawing] = useState(false);
    const [currentPath, setCurrentPath] = useState<string>("");
    const [strokeColor, setStrokeColor] = useState("#0f172a");
    const [strokeWidth, setStrokeWidth] = useState(2.5);
    const [label, setLabel] = useState("");
    const [labels, setLabels] = useState<Record<number, string>>({});

    const canvasRef = useRef<SVGSVGElement>(null);
    const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);
    const pathPointsRef = useRef<{ x: number; y: number }[]>([]);

    const filledCount = circles.filter((c) => !c.isEmpty).length;

    // Timer
    useEffect(() => {
        if (phase === "drawing") {
            timerRef.current = setInterval(() => {
                setTimeLeft((prev) => {
                    if (prev <= 1) {
                        clearInterval(timerRef.current!);
                        setPhase("review");
                        setActiveCircle(null);
                        return 0;
                    }
                    return prev - 1;
                });
            }, 1000);
        }
        return () => {
            if (timerRef.current) clearInterval(timerRef.current);
        };
    }, [phase]);

    // Notify parent
    useEffect(() => {
        onContentChange({
            type: content.type,
            circles,
            labels,
            filledCount,
            phase,
        });
    }, [circles, labels, filledCount, phase]);

    const formatTime = (s: number) => {
        const m = Math.floor(s / 60);
        const sec = s % 60;
        return `${m}:${sec.toString().padStart(2, "0")}`;
    };

    const getCircleCenter = (svg: SVGSVGElement) => {
        const rect = svg.getBoundingClientRect();
        return {
            cx: rect.width / 2,
            cy: rect.height / 2,
            r: rect.width / 2 - 4,
        };
    };

    const getPoint = (
        e: React.MouseEvent | React.TouchEvent,
        svg: SVGSVGElement,
    ) => {
        const rect = svg.getBoundingClientRect();
        if ("touches" in e) {
            return {
                x: e.touches[0].clientX - rect.left,
                y: e.touches[0].clientY - rect.top,
            };
        }
        return {
            x: (e as React.MouseEvent).clientX - rect.left,
            y: (e as React.MouseEvent).clientY - rect.top,
        };
    };

    const startDrawing = (e: React.MouseEvent | React.TouchEvent) => {
        if (phase !== "drawing" || activeCircle === null) return;
        e.preventDefault();
        const svg = canvasRef.current!;
        const pt = getPoint(e, svg);
        pathPointsRef.current = [pt];
        setCurrentPath(`M ${pt.x} ${pt.y}`);
        setIsDrawing(true);
    };

    const continueDrawing = (e: React.MouseEvent | React.TouchEvent) => {
        if (!isDrawing || phase !== "drawing") return;
        e.preventDefault();
        const svg = canvasRef.current!;
        const pt = getPoint(e, svg);
        pathPointsRef.current.push(pt);

        // Smooth bezier
        const pts = pathPointsRef.current;
        if (pts.length < 2) return;
        const last = pts[pts.length - 2];
        const curr = pts[pts.length - 1];
        const cpx = (last.x + curr.x) / 2;
        const cpy = (last.y + curr.y) / 2;
        setCurrentPath((prev) => `${prev} Q ${last.x} ${last.y} ${cpx} ${cpy}`);
    };

    const endDrawing = () => {
        if (!isDrawing) return;
        setIsDrawing(false);
        if (!currentPath || activeCircle === null) return;

        // Save stroke to circle
        setCircles((prev) =>
            prev.map((c) => {
                if (c.id !== activeCircle) return c;
                const existing = c.drawing || "";
                return {
                    ...c,
                    drawing:
                        existing +
                        `|||${strokeColor}|${strokeWidth}|${currentPath}`,
                    isEmpty: false,
                };
            }),
        );
        setCurrentPath("");
        pathPointsRef.current = [];
    };

    const clearCircle = (id: number, e: React.MouseEvent) => {
        e.stopPropagation();
        setCircles((prev) =>
            prev.map((c) =>
                c.id === id ? { ...c, drawing: null, isEmpty: true } : c,
            ),
        );
        setLabels((prev) => {
            const next = { ...prev };
            delete next[id];
            return next;
        });
    };

    const parseStrokes = (data: string | null) => {
        if (!data) return [];
        return data
            .split("|||")
            .filter(Boolean)
            .map((s) => {
                const parts = s.split("|");
                return {
                    color: parts[0],
                    width: parseFloat(parts[1]),
                    path: parts[2],
                };
            });
    };

    const handleCircleClick = (id: number) => {
        if (phase !== "drawing") return;
        setActiveCircle(id === activeCircle ? null : id);
        if (labels[id] !== undefined) setLabel(labels[id]);
        else setLabel("");
    };

    const handleLabelChange = (val: string) => {
        setLabel(val);
        if (activeCircle !== null) {
            setLabels((prev) => ({ ...prev, [activeCircle]: val }));
        }
    };

    const progress = ((timeLimit - timeLeft) / timeLimit) * 100;
    const urgency = timeLeft <= 30;

    return (
        <div className={styles.wrapper}>
            {/* Header strip */}
            <div className={styles.header}>
                <div className={styles.meta}>
                    <span className={styles.filled}>
                        <strong>{filledCount}</strong>
                        <span> / {TOTAL_CIRCLES}</span>
                    </span>
                    <span className={styles.filledLabel}>кругов заполнено</span>
                </div>

                {phase === "drawing" && (
                    <div
                        className={`${styles.timer} ${urgency ? styles.timerUrgent : ""}`}
                    >
                        <svg viewBox="0 0 36 36" className={styles.timerRing}>
                            <circle
                                cx="18"
                                cy="18"
                                r="15.9"
                                fill="none"
                                strokeWidth="2.5"
                                className={styles.timerTrack}
                            />
                            <circle
                                cx="18"
                                cy="18"
                                r="15.9"
                                fill="none"
                                strokeWidth="2.5"
                                className={styles.timerProgress}
                                strokeDasharray={`${100 - progress} 100`}
                                strokeDashoffset="25"
                            />
                        </svg>
                        <span className={styles.timerText}>
                            {formatTime(timeLeft)}
                        </span>
                    </div>
                )}

                {phase === "review" && (
                    <div className={styles.reviewBadge}>✓ Время вышло</div>
                )}
            </div>

            {/* Progress bar */}
            {phase === "drawing" && (
                <div className={styles.progressBar}>
                    <div
                        className={`${styles.progressFill} ${urgency ? styles.progressUrgent : ""}`}
                        style={{ width: `${100 - progress}%` }}
                    />
                </div>
            )}

            {/* Ready screen */}
            {phase === "ready" && (
                <div className={styles.readyScreen}>
                    <div className={styles.readyIcon}>◎</div>
                    <h2 className={styles.readyTitle}>Тридцать кругов</h2>
                    <p className={styles.readyDesc}>
                        Перед вами 30 пустых кругов. За{" "}
                        <strong>3 минуты</strong> нарисуйте в как можно большем
                        количестве кругов узнаваемые объекты: пицца, часы,
                        яблоко — всё, что придёт в голову.
                    </p>
                    <ul className={styles.readyRules}>
                        <li>
                            Нажмите на круг — он станет активным для рисования
                        </li>
                        <li>Рисуйте мышью или пальцем внутри круга</li>
                        <li>Можно подписать каждый объект</li>
                        <li>Можно объединять круги для смелых идей!</li>
                    </ul>
                    <button
                        className={styles.startBtn}
                        onClick={() => setPhase("drawing")}
                    >
                        Начать
                    </button>
                </div>
            )}

            {/* Drawing toolbar */}
            {phase === "drawing" && (
                <div className={styles.toolbar}>
                    <div className={styles.toolGroup}>
                        <span className={styles.toolLabel}>Цвет</span>
                        {[
                            "#0f172a",
                            "#2563eb",
                            "#dc2626",
                            "#16a34a",
                            "#d97706",
                            "#9333ea",
                        ].map((c) => (
                            <button
                                key={c}
                                className={`${styles.colorSwatch} ${strokeColor === c ? styles.colorSwatchActive : ""}`}
                                style={{ background: c }}
                                onClick={() => setStrokeColor(c)}
                            />
                        ))}
                    </div>
                    <div className={styles.toolGroup}>
                        <span className={styles.toolLabel}>Толщина</span>
                        {[1.5, 2.5, 4].map((w) => (
                            <button
                                key={w}
                                className={`${styles.widthBtn} ${strokeWidth === w ? styles.widthBtnActive : ""}`}
                                onClick={() => setStrokeWidth(w)}
                            >
                                <span
                                    style={{
                                        height: w * 3,
                                        background: strokeColor,
                                    }}
                                />
                            </button>
                        ))}
                    </div>
                    {activeCircle !== null && (
                        <div className={styles.toolGroup}>
                            <span className={styles.toolLabel}>Название</span>
                            <input
                                className={styles.labelInput}
                                value={label}
                                onChange={(e) =>
                                    handleLabelChange(e.target.value)
                                }
                                placeholder="что нарисовано?"
                                maxLength={20}
                            />
                        </div>
                    )}
                </div>
            )}

            {/* Grid */}
            {phase !== "ready" && (
                <div className={styles.grid}>
                    {circles.map((circle) => {
                        const isActive = activeCircle === circle.id;
                        const strokes = parseStrokes(circle.drawing);
                        return (
                            <div
                                key={circle.id}
                                className={`${styles.circleWrapper} ${isActive ? styles.circleActive : ""} ${!circle.isEmpty ? styles.circleFilled : ""}`}
                                onClick={() => handleCircleClick(circle.id)}
                            >
                                <svg
                                    ref={isActive ? canvasRef : undefined}
                                    className={styles.circleSvg}
                                    viewBox="0 0 100 100"
                                    onMouseDown={
                                        isActive ? startDrawing : undefined
                                    }
                                    onMouseMove={
                                        isActive ? continueDrawing : undefined
                                    }
                                    onMouseUp={
                                        isActive ? endDrawing : undefined
                                    }
                                    onMouseLeave={
                                        isActive ? endDrawing : undefined
                                    }
                                    onTouchStart={
                                        isActive ? startDrawing : undefined
                                    }
                                    onTouchMove={
                                        isActive ? continueDrawing : undefined
                                    }
                                    onTouchEnd={
                                        isActive ? endDrawing : undefined
                                    }
                                >
                                    {/* Circle border */}
                                    <circle
                                        cx="50"
                                        cy="50"
                                        r="46"
                                        fill={isActive ? "#f8faff" : "white"}
                                        stroke={
                                            isActive
                                                ? "var(--primary-color)"
                                                : "var(--border-color)"
                                        }
                                        strokeWidth={isActive ? "2" : "1.5"}
                                    />

                                    {/* Saved strokes */}
                                    {strokes.map((s, i) => (
                                        <path
                                            key={i}
                                            d={s.path}
                                            fill="none"
                                            stroke={s.color}
                                            strokeWidth={s.width}
                                            strokeLinecap="round"
                                            strokeLinejoin="round"
                                            clipPath={`url(#clip-${circle.id})`}
                                        />
                                    ))}

                                    {/* Current stroke */}
                                    {isActive && currentPath && (
                                        <path
                                            d={currentPath}
                                            fill="none"
                                            stroke={strokeColor}
                                            strokeWidth={strokeWidth}
                                            strokeLinecap="round"
                                            strokeLinejoin="round"
                                            clipPath={`url(#clip-${circle.id})`}
                                        />
                                    )}

                                    {/* Clip path */}
                                    <defs>
                                        <clipPath id={`clip-${circle.id}`}>
                                            <circle cx="50" cy="50" r="45" />
                                        </clipPath>
                                    </defs>

                                    {/* Circle number */}
                                    {circle.isEmpty && !isActive && (
                                        <text
                                            x="50"
                                            y="54"
                                            textAnchor="middle"
                                            fontSize="14"
                                            fill="var(--text-muted)"
                                            fontFamily="inherit"
                                        >
                                            {circle.id + 1}
                                        </text>
                                    )}
                                </svg>

                                {/* Label */}
                                {labels[circle.id] && (
                                    <div className={styles.circleLabel}>
                                        {labels[circle.id]}
                                    </div>
                                )}

                                {/* Clear button */}
                                {!circle.isEmpty && phase === "drawing" && (
                                    <button
                                        className={styles.clearBtn}
                                        onClick={(e) =>
                                            clearCircle(circle.id, e)
                                        }
                                        title="Очистить"
                                    >
                                        ×
                                    </button>
                                )}
                            </div>
                        );
                    })}
                </div>
            )}

            {/* Review panel */}
            {phase === "review" && (
                <div className={styles.reviewPanel}>
                    <div className={styles.reviewStats}>
                        <div className={styles.reviewStat}>
                            <span className={styles.reviewStatValue}>
                                {filledCount}
                            </span>
                            <span className={styles.reviewStatLabel}>
                                заполнено кругов
                            </span>
                        </div>
                        <div className={styles.reviewStat}>
                            <span className={styles.reviewStatValue}>
                                {TOTAL_CIRCLES - filledCount}
                            </span>
                            <span className={styles.reviewStatLabel}>
                                осталось пустых
                            </span>
                        </div>
                        <div className={styles.reviewStat}>
                            <span className={styles.reviewStatValue}>
                                {Math.round(
                                    (filledCount / TOTAL_CIRCLES) * 100,
                                )}
                                %
                            </span>
                            <span className={styles.reviewStatLabel}>
                                заполнено
                            </span>
                        </div>
                    </div>
                    <p className={styles.reviewHint}>
                        Обсудите результаты: есть ли повторяющиеся темы? Кто-то
                        объединил круги? Сколько идей уникальных, а сколько —
                        вариации одного?
                    </p>
                    <button
                        className={styles.resetBtn}
                        onClick={() => {
                            setCircles(
                                Array.from(
                                    { length: TOTAL_CIRCLES },
                                    (_, i) => ({
                                        id: i,
                                        drawing: null,
                                        isEmpty: true,
                                    }),
                                ),
                            );
                            setLabels({});
                            setActiveCircle(null);
                            setTimeLeft(timeLimit);
                            setPhase("ready");
                        }}
                    >
                        Попробовать снова
                    </button>
                </div>
            )}
        </div>
    );
};

export default ThirtyCirclesTask;
