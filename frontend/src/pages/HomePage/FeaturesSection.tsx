import React from "react";
import {
    BookOpen,
    Clock,
    CheckCircle,
    Zap,
    Users,
    BarChart2,
} from "lucide-react";
import styles from "./HomePage.module.scss";

const FEATURES = [
    {
        icon: BookOpen,
        tag: "Контент",
        title: "Креативные задания",
        description:
            "Аббревиатуры, аллитерации, неожиданные связи — каждый формат тренирует отдельную грань мышления.",
        accent: "blue",
    },
    {
        icon: Zap,
        tag: "Скорость",
        title: "Мгновенная обратная связь",
        description:
            "Преподаватель видит работу сразу после отправки и оставляет комментарии точечно.",
        accent: "warning",
    },
    {
        icon: BarChart2,
        tag: "Аналитика",
        title: "Прогресс как на ладони",
        description:
            "Графики выполнения, история заданий и динамика роста — всё в одном дашборде.",
        accent: "success",
    },
    {
        icon: Clock,
        tag: "Тайминг",
        title: "Контроль времени",
        description:
            "Таймер на каждое задание учит работать в условиях ограничений — как в реальной жизни.",
        accent: "purple",
    },
    {
        icon: Users,
        tag: "Команда",
        title: "Учитель и ученик",
        description:
            "Преподаватель создаёт задания, назначает дедлайны и отслеживает прогресс группы.",
        accent: "blue",
    },
    {
        icon: CheckCircle,
        tag: "Результат",
        title: "Фиксация достижений",
        description:
            "Все выполненные работы сохраняются — можно вернуться, перечитать и оценить рост.",
        accent: "success",
    },
];

const accentMap: Record<string, string> = {
    blue: styles.featureAccentBlue,
    warning: styles.featureAccentWarning,
    success: styles.featureAccentSuccess,
    purple: styles.featureAccentPurple,
};

const FeaturesSection = () => (
    <section className={styles.features}>
        <div className={styles.featuresHeader}>
            <p className={styles.featuresEyebrow}>Возможности платформы</p>
            <h2 className={styles.featuresTitle}>
                Всё, что нужно для{" "}
                <span className={styles.featuresTitleAccent}>роста</span>
            </h2>
        </div>

        <div className={styles.featuresGrid}>
            {FEATURES.map(
                ({ icon: Icon, tag, title, description, accent }, i) => (
                    <div
                        key={title}
                        className={styles.featureCard}
                        style={{ animationDelay: `${i * 0.07}s` }}
                    >
                        <div className={styles.featureMeta}>
                            <div
                                className={`${styles.featureIcon} ${accentMap[accent]}`}
                            >
                                <Icon />
                            </div>
                            <span className={styles.featureTag}>{tag}</span>
                        </div>
                        <h3 className={styles.featureTitle}>{title}</h3>
                        <p className={styles.featureDesc}>{description}</p>
                    </div>
                ),
            )}
        </div>
    </section>
);

export default FeaturesSection;
