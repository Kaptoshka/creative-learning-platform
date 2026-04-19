import React, { useContext, useState, useEffect, useRef } from "react";
import { Link } from "react-router-dom";
import { AuthContext } from "@/context/AuthContext";
import TasksPage from "@/pages/TasksPage";
import Button from "@/components/Button";
import styles from "./DashboardPage.module.scss";

import {
  Plus,
  List,
  Users,
  User,
  BarChart2,
  CheckCircle,
  Clock,
  AlertTriangle,
  HelpCircle,
  ExternalLink,
} from "lucide-react";

const DashboardPage = () => {
  const { user } = useContext(AuthContext);
  const [currentTime, setCurrentTime] = useState(new Date());
  const [loading, setLoading] = useState(true);

  const tasksRef = useRef(null);
  const progressRef = useRef(null);

  const scrollToSection = (elementRef) => {
    if (elementRef.current) {
      elementRef.current.scrollIntoView({
        behavior: "smooth",
        block: "start",
      });
    }
  };

  useEffect(() => {
    const timer = setInterval(() => setCurrentTime(new Date()), 60000);
    const loadingTimer = setTimeout(() => setLoading(false), 1500);
    return () => {
      clearInterval(timer);
      clearTimeout(loadingTimer);
    };
  }, []);

  const getGreeting = () => {
    const hour = currentTime.getHours();
    if (hour < 12) return "Доброе утро";
    if (hour < 18) return "Добрый день";
    return "Добрый вечер";
  };

  const getFormattedDate = () => {
    return currentTime.toLocaleDateString("ru-RU", {
      weekday: "long",
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  };

  if (!user || loading) {
    return (
      <div className={styles.dashboard}>
        <div className={styles.dashboardContainer}>
          <div className="loading">
            <div className="loading__spinner"></div>
            <span className="loading__text">
              {loading ? "Загрузка личного кабинета..." : "Загрузка..."}
            </span>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.dashboard}>
      <div className={styles.dashboardContainer}>
        <header className={styles.dashboardHeader}>
          <div className={styles.dashboardWelcome}>
            <h1 className={styles.dashboardWelcomeTitle}>
              {getGreeting()}, {user.name || user.email}! 👋
            </h1>
            <p className={styles.dashboardWelcomeSubtitle}>
              Добро пожаловать в твою {user.role} информационную панель
            </p>
            <div className={styles.dashboardWelcomeStats}>
              <div className={styles.dashboardWelcomeStat}>
                <span className={styles.dashboardWelcomeStatText}>
                  📅 {getFormattedDate()}
                </span>
              </div>
              <div className={styles.dashboardWelcomeStat}>
                <span className={styles.dashboardWelcomeStatText}>
                  👤 Роль:{" "}
                  {user.role.charAt(0).toUpperCase() + user.role.slice(1)}
                </span>
              </div>
              <div className={styles.dashboardWelcomeStat}>
                <span className={styles.dashboardWelcomeStatText}>
                  ⏰{" "}
                  {currentTime.toLocaleTimeString("en-US", {
                    hour: "2-digit",
                    minute: "2-digit",
                    hourCycle: "h23",
                  })}
                </span>
              </div>
            </div>
          </div>
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
              <Button to="/tasks" variant="outline" icon={List} viewTransition>
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

        <div className="quick-stats">
          {user.role === "student" && (
            <>
              <div className="quick-stats__card">
                <div className="quick-stats__header">
                  <div className="quick-stats__icon">
                    <List size={24} />
                  </div>
                </div>
                <div className="quick-stats__value">8</div>
                <div className="quick-stats__label">Доступных заданий</div>
              </div>
            </>
          )}
        </div>

        <main className={styles.dashboardContent}>
          {/* <div className={styles.dashboardSection}>
            <div className={styles.dashboardSectionHeader}>
              <h2 className={styles.dashboardSectionTitle}>
                <Clock className={styles.dashboardSectionTitleIcon} />
                Недавняя Активность
              </h2>
              <Link
                to="/activity"
                className={styles.dashboardSectionAction}
                viewTransition
              >
                Посмотреть все
              </Link>
            </div>
          </div>*/}

          <div className={styles.dashboardSection} ref={progressRef}>
            <div className={styles.dashboardSectionHeader}>
              <h2 className={styles.dashboardSectionTitle}>
                <BarChart2 className={styles.dashboardSectionTitleIcon} />
                {user.role === "teacher" ? "Успехи учеников" : "Твой прогресс"}
              </h2>
              <Link
                to="/analytics"
                className={styles.dashboardSectionAction}
                viewTransition
              >
                Посмотреть детали
              </Link>
            </div>

            <div style={{ marginBottom: "1rem" }}>
              <div
                style={{
                  display: "flex",
                  justifyContent: "space-between",
                  marginBottom: "0.5rem",
                }}
              >
                <span>Общий прогресс</span>
                <span style={{ fontWeight: 600 }}>65%</span>
              </div>
              <div className="progress">
                <div className="progress__bar" style={{ width: "65%" }}></div>
              </div>
            </div>
          </div>

          <div className={styles.dashboardSection} ref={tasksRef}>
            {user.role === "teacher" ? <TasksPage /> : <TasksPage />}
          </div>
        </main>

        <div className={styles.dashboardSection} style={{ marginTop: "2rem" }}>
          <div className={styles.dashboardSectionHeader}>
            <h2 className={styles.dashboardSectionTitle}>
              <svg
                className={styles.dashboardSectionTitleIcon}
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              Quick Help & Tips
            </h2>
            <Link to="/help" className={styles.dashboardSectionAction}>
              Get help
            </Link>
          </div>

          <div
            style={{
              display: "grid",
              gap: "1rem",
              gridTemplateColumns: "repeat(auto-fit, minmax(250px, 1fr))",
            }}
          >
            <div
              style={{
                background: "var(--background-color)",
                padding: "1rem",
                borderRadius: "var(--radius)",
                border: "1px solid var(--border-color)",
              }}
            >
              <h4
                style={{ marginBottom: "0.5rem", color: "var(--text-primary)" }}
              >
                {user.role === "teacher"
                  ? "📝 Create Engaging Tasks"
                  : "🚀 Getting Started"}
              </h4>
              <p
                style={{
                  fontSize: "0.875rem",
                  color: "var(--text-secondary)",
                  marginBottom: "0.5rem",
                }}
              >
                {user.role === "teacher"
                  ? "Use interactive elements and real-world examples to make your tasks more engaging for students."
                  : "New to the platform? Start with the basics and work your way up through our structured learning path."}
              </p>
              <Link
                to={
                  user.role === "teacher"
                    ? "/guide/creating-tasks"
                    : "/guide/getting-started"
                }
                style={{
                  fontSize: "0.75rem",
                  color: "var(--primary-color)",
                  textDecoration: "none",
                }}
              >
                Learn more →
              </Link>
            </div>

            <div
              style={{
                background: "var(--background-color)",
                padding: "1rem",
                borderRadius: "var(--radius)",
                border: "1px solid var(--border-color)",
              }}
            >
              <h4
                style={{ marginBottom: "0.5rem", color: "var(--text-primary)" }}
              >
                {user.role === "teacher"
                  ? "📊 Track Progress"
                  : "💡 Study Tips"}
              </h4>
              <p
                style={{
                  fontSize: "0.875rem",
                  color: "var(--text-secondary)",
                  marginBottom: "0.5rem",
                }}
              >
                {user.role === "teacher"
                  ? "Monitor your students' progress with detailed analytics and provide timely feedback."
                  : "Break down complex problems into smaller steps and practice regularly for best results."}
              </p>
              <Link
                to={
                  user.role === "teacher"
                    ? "/guide/analytics"
                    : "/guide/study-tips"
                }
                style={{
                  fontSize: "0.75rem",
                  color: "var(--primary-color)",
                  textDecoration: "none",
                }}
              >
                Learn more →
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default DashboardPage;
