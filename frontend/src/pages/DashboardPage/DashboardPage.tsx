import React, { useRef } from "react";
import { Link } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";
import { useCurrentTime } from "@/hooks/useCurrentTime";
import TasksPage from "@/pages/TasksPage";
import Button from "@/components/Button";
import Loading from "@/components/ui/Loading";
import Greeting from "@/components/ui/Greeting";
import QuickStats from "@/components/ui/QuickStats";
import ProgressBar from "@/components/ui/ProgressBar";
import styles from "./DashboardPage.module.scss";

import { Plus, List, Users, User, BarChart2, CheckCircle } from "lucide-react";

const DashboardPage = () => {
    const { user, loading: authLoading } = useAuth();
    const { greeting, formattedDate, formattedTime } = useCurrentTime();

    const tasksRef = useRef(null);
    const progressRef = useRef(null);

    const scrollToSection = (elementRef: React.RefObject<HTMLDivElement>) => {
        if (elementRef.current) {
            elementRef.current.scrollIntoView({
                behavior: "smooth",
                block: "start",
            });
        }
    };

    if (!user || authLoading) {
        return (
            <div className={styles.dashboard}>
                <div className={styles.dashboardContainer}>
                    <Loading text="Загрузка личного кабинета..." fullPage />
                </div>
            </div>
        );
    }

    return (
        <div className={styles.dashboard}>
            <div className={styles.dashboardContainer}>
                <header className={styles.dashboardHeader}>
                    <Greeting
                        name={user.name || user.email}
                        role={user.role}
                        greeting={greeting}
                        formattedDate={formattedDate}
                        formattedTime={formattedTime}
                    />
                </header>

                <div className={styles.dashboardActions}>
                    {user.role === "teacher" ? (
                        <>
                            <Button
                                to="/tasks/create"
                                variant="gradient"
                                icon={Plus}
                                viewTransition
                            >
                                Создать новую задачу
                            </Button>
                            <Button
                                to="/tasks"
                                variant="outline"
                                icon={List}
                                viewTransition
                            >
                                Управление задачами
                            </Button>
                            <Button
                                to="/submissions"
                                variant="secondary"
                                icon={Users}
                                viewTransition
                            >
                                Проверка ответов
                            </Button>
                        </>
                    ) : (
                        <>
                            <Button
                                onClick={() => scrollToSection(tasksRef)}
                                variant="gradient"
                                icon={List}
                            >
                                Посмотреть Задания
                            </Button>
                            <Button
                                to="/profile"
                                variant="outline"
                                icon={User}
                                viewTransition
                            >
                                Мой Профиль
                            </Button>
                            <Button
                                onClick={() => scrollToSection(progressRef)}
                                variant="success"
                                icon={BarChart2}
                            >
                                Мой Прогресс
                            </Button>
                        </>
                    )}
                </div>

                {user.role === "student" && (
                    <QuickStats
                        items={[
                            {
                                label: "Доступных заданий",
                                value: 8,
                                icon: List,
                            },
                        ]}
                    />
                )}

                <main className={styles.dashboardContent}>
                    <div className={styles.dashboardSection} ref={progressRef}>
                        <div className={styles.dashboardSectionHeader}>
                            <h2 className={styles.dashboardSectionTitle}>
                                <BarChart2
                                    className={styles.dashboardSectionTitleIcon}
                                />
                                {user.role === "teacher"
                                    ? "Успехи учеников"
                                    : "Твой прогресс"}
                            </h2>
                            <Link
                                to="/analytics"
                                className={styles.dashboardSectionAction}
                                viewTransition
                            >
                                Посмотреть детали
                            </Link>
                        </div>
                        <ProgressBar label="Общий прогресс" value={65} />
                    </div>

                    <div className={styles.dashboardSection} ref={tasksRef}>
                        <TasksPage />
                    </div>
                </main>

                <div
                    className={styles.dashboardSection}
                    style={{ marginTop: "2rem" }}
                >
                    <div className={styles.dashboardSectionHeader}>
                        <h2 className={styles.dashboardSectionTitle}>
                            <CheckCircle
                                className={styles.dashboardSectionTitleIcon}
                            />
                            Quick Help & Tips
                        </h2>
                        <Link
                            to="/help"
                            className={styles.dashboardSectionAction}
                        >
                            Get help
                        </Link>
                    </div>
                    <p style={{ color: "var(--text-secondary)" }}>
                        {user.role === "teacher"
                            ? "Создавайте интересные задания для учеников"
                            : "Начните с простых заданий и постепенно усложняйте"}
                    </p>
                </div>
            </div>
        </div>
    );
};

export default DashboardPage;
