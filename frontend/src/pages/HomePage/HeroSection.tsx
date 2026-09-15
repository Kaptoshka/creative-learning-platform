import React, { useContext } from "react";
import { useNavigate } from "react-router-dom";
import Button from "@/components/Button";
import { AuthContext } from "@/context/AuthContext";
import styles from "./HomePage.module.scss";

const WORDS = ["идеи", "решения", "образы", "смыслы"];

const HeroSection = () => {
    const navigate = useNavigate();
    const { user } = useContext(AuthContext);

    const handleStartLearning = () => {
        navigate(user ? "/dashboard" : "/auth", { viewTransition: true });
    };

    return (
        <section className={styles.hero}>
            <div className={styles.heroNoise} aria-hidden="true" />

            <div className={styles.heroGrid} aria-hidden="true">
                {Array.from({ length: 24 }).map((_, i) => (
                    <div key={i} className={styles.heroGridCell} />
                ))}
            </div>

            <div className={styles.heroContent}>
                <div className={styles.heroBadge}>
                    <span className={styles.heroBadgeDot} />
                    Платформа творческого обучения
                </div>

                <h1 className={styles.heroTitle}>
                    <span className={styles.heroTitleLine}>Рождайте</span>
                    <span
                        className={styles.heroTitleMarquee}
                        aria-hidden="true"
                    >
                        {WORDS.map((w) => (
                            <span key={w} className={styles.heroTitleWord}>
                                {w}
                            </span>
                        ))}
                    </span>
                    <span className={styles.heroTitleLine}>
                        каждый <em className={styles.heroTitleEm}>день</em>
                    </span>
                </h1>

                <p className={styles.heroSubtitle}>
                    Задания, которые будят воображение. Прогресс, который виден.
                    Обратная связь, которая развивает.
                </p>

                <div className={styles.heroActions}>
                    <Button variant="gradient" onClick={handleStartLearning}>
                        Начать бесплатно
                    </Button>
                    <Button variant="secondary">Посмотреть задания</Button>
                </div>

                <div className={styles.heroStats}>
                    <div className={styles.heroStat}>
                        <span className={styles.heroStatValue}>6</span>
                        <span className={styles.heroStatLabel}>
                            типов заданий
                        </span>
                    </div>
                    <div
                        className={styles.heroStatDivider}
                        aria-hidden="true"
                    />
                    <div className={styles.heroStat}>
                        <span className={styles.heroStatValue}>∞</span>
                        <span className={styles.heroStatLabel}>
                            идей внутри
                        </span>
                    </div>
                    <div
                        className={styles.heroStatDivider}
                        aria-hidden="true"
                    />
                    <div className={styles.heroStat}>
                        <span className={styles.heroStatValue}>100%</span>
                        <span className={styles.heroStatLabel}>
                            творческой свободы
                        </span>
                    </div>
                </div>
            </div>

            <div className={styles.heroVisual} aria-hidden="true">
                <div className={styles.heroCard}>
                    <div className={styles.heroCardHeader}>
                        <span
                            className={styles.heroCardDot}
                            style={{ background: "#ef4444" }}
                        />
                        <span
                            className={styles.heroCardDot}
                            style={{ background: "#f59e0b" }}
                        />
                        <span
                            className={styles.heroCardDot}
                            style={{ background: "#22c55e" }}
                        />
                    </div>
                    <div className={styles.heroCardTask}>
                        <p className={styles.heroCardLabel}>Задание дня</p>
                        <p className={styles.heroCardTitle}>Аббревиатуры</p>
                        <p className={styles.heroCardDesc}>
                            Придумайте расшифровку для слова{" "}
                            <strong>МЫСЛЬ</strong>
                        </p>
                    </div>
                    <div className={styles.heroCardInput}>
                        <span className={styles.heroCardCursor} />
                    </div>
                    <div className={styles.heroCardFooter}>
                        <span className={styles.heroCardTime}>⏱ 5 мин</span>
                        <span className={styles.heroCardLevel}>Начальный</span>
                    </div>
                </div>

                <div className={styles.heroFloatA}>✦</div>
                <div className={styles.heroFloatB}>◈</div>
                <div className={styles.heroFloatC}>⬡</div>
            </div>
        </section>
    );
};

export default HeroSection;
