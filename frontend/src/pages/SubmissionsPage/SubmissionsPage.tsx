import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import styles from "./SubmissionsPage.module.scss";
import { useAuth } from "@/hooks/useAuth";
import { useSubmissions } from "@/hooks/useSubmissions";
import Button from "@/components/Button/Button.jsx";
import Loading from "@/components/ui/Loading";

const SubmissionsPage = () => {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [filter, setFilter] = useState("all");
  const [searchQuery, setSearchQuery] = useState("");
  const [sortBy, setSortBy] = useState("date");

  const { submissions, loading, error, refetch } = useSubmissions(user?.id || null);

  const handleReview = (submission: { id: number }) => {
    navigate(`/submissions/${submission.id}`, { state: { submission } });
  };

  const getFilteredSubmissions = () => {
    let filtered = [...submissions];

    if (filter === "urgent") {
      filtered = filtered.filter((sub) => sub.priority === "urgent");
    } else if (filter === "recent") {
      filtered = filtered.filter((sub) => sub.daysWaiting <= 1);
    }

    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(
        (sub) =>
          sub.studentName.toLowerCase().includes(query) ||
          sub.taskTitle.toLowerCase().includes(query) ||
          sub.studentEmail.toLowerCase().includes(query),
      );
    }

    if (sortBy === "date") {
      filtered.sort((a, b) => b.submittedAt.getTime() - a.submittedAt.getTime());
    } else if (sortBy === "student") {
      filtered.sort((a, b) => a.studentName.localeCompare(b.studentName));
    } else if (sortBy === "task") {
      filtered.sort((a, b) => a.taskTitle.localeCompare(b.taskTitle));
    }

    return filtered;
  };

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case "urgent":
        return styles.submissionCardUrgent;
      case "new":
        return styles.submissionCardNew;
      default:
        return "";
    }
  };

  const getTaskIcon = (taskType: string) => {
    return taskType || "📝";
  };

  const filteredSubmissions = getFilteredSubmissions();

  if (loading) {
    return (
      <div className={styles.submissionsPage}>
        <div className={styles.submissionsPageContainer}>
          <Loading text="Загрузка ответов..." fullPage />
        </div>
      </div>
    );
  }

  return (
    <div className={styles.submissionsPage}>
      <div className={styles.submissionsPageContainer}>
        <header className={styles.submissionsPageHeader}>
          <div>
            <h1 className={styles.submissionsPageHeaderTitle}>
              Ответы на задания
            </h1>
            <p className={styles.submissionsPageHeaderSubtitle}>
              Оценивайте и дайте обратную связь студентам
            </p>
          </div>
          <div className={styles.submissionsPageHeaderStats}>
            <div className={styles.statItem}>
              <div className={styles.statItemValue}>{submissions.length}</div>
              <div className={styles.statItemLabel}>Всего ответов</div>
            </div>
            <div className={styles.statItem}>
              <div className={`${styles.statItemValue} ${styles.statItemValueUrgent}`}>
                {submissions.filter((s) => s.priority === "urgent").length}
              </div>
              <div className={styles.statItemLabel}>Давно отправлены</div>
            </div>
            <div className={styles.statItem}>
              <div className={`${styles.statItemValue} ${styles.statItemValueNew}`}>
                {submissions.filter((s) => s.priority === "new").length}
              </div>
              <div className={styles.statItemLabel}>Отправлены сегодня</div>
            </div>
          </div>
        </header>

        <div className={styles.controls}>
          <div className={styles.controlsFilters}>
            <Button
              variant={filter === "all" ? "active" : "default"}
              onClick={() => setFilter("all")}
              className={styles.controlsFilterBtn}
            >
              Все ответы
              <span className={styles.controlsBadge}>{submissions.length}</span>
            </Button>
            <Button
              variant={filter === "urgent" ? "active" : "default"}
              onClick={() => setFilter("urgent")}
              className={`${styles.controlsFilterBtn} ${
                filter === "urgent" ? styles.controlsFilterBtnActive : ""
              }`}
            >
              Давно отправлены
              <span className={`${styles.controlsBadge} ${styles.controlsBadgeUrgent}`}>
                {submissions.filter((s) => s.priority === "urgent").length}
              </span>
            </Button>
            <Button
              variant={filter === "recent" ? "active" : "default"}
              onClick={() => setFilter("recent")}
            >
              Недавно отправлены
              <span className={styles.controlsBadge}>
                {submissions.filter((s) => s.daysWaiting <= 1).length}
              </span>
            </Button>
          </div>

          <div className={styles.controlsRight}>
            <div className={styles.controlsSearch}>
              <input
                type="text"
                placeholder="Поиск по студенту или заданию..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className={styles.controlsSearchInput}
              />
            </div>

            <select
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value)}
              className={styles.controlsSort}
            >
              <option value="date">Сортировать по дате</option>
              <option value="student">Сортировать по студенту</option>
              <option value="task">Сортировать по заданию</option>
            </select>
          </div>
        </div>

        {filteredSubmissions.length === 0 ? (
          <div className={styles.emptyState}>
            <h3 className={styles.emptyStateTitle}>
              Ответы на задания не найдены
            </h3>
          </div>
        ) : (
          <div className={styles.submissionsGrid}>
            {filteredSubmissions.map((submission) => (
              <div
                key={submission.id}
                className={`${styles.submissionCard} ${getPriorityColor(submission.priority)}`}
                onClick={() => handleReview(submission)}
              >
                {submission.priority === "urgent" && (
                  <div className={`${styles.priorityBadge} ${styles.priorityBadgeUrgent}`}>
                    Давно отправлен
                  </div>
                )}
                {submission.priority === "new" && (
                  <div className={`${styles.priorityBadge} ${styles.priorityBadgeNew}`}>
                    Новый
                  </div>
                )}

                <div className={styles.submissionCardHeader}>
                  <div className={styles.submissionCardStudent}>
                    <div className={styles.submissionCardAvatar}>
                      {submission.studentAvatar}
                    </div>
                    <div className={styles.submissionCardStudentInfo}>
                      <div className={styles.submissionCardStudentName}>
                        {submission.studentName}
                      </div>
                      <div className={styles.submissionCardStudentEmail}>
                        {submission.studentEmail}
                      </div>
                    </div>
                  </div>
                </div>

                <div className={styles.submissionCardBody}>
                  <div className={styles.submissionCardTask}>
                    <h3 className={styles.submissionCardTaskTitle}>
                      {submission.taskTitle}
                    </h3>
                  </div>
                  <p className={styles.submissionCardPreview}>
                    {submission.preview}
                  </p>
                </div>

                <div className={styles.submissionCardFooter}>
                  <div className={styles.submissionCardMeta}>
                    <div className={styles.submissionCardMetaItem}>
                      <span>
                        {submission.daysWaiting === 0
                          ? "Сегодня"
                          : `${submission.daysWaiting} дней назад`}
                      </span>
                    </div>
                  </div>

                  <Button to={`/submissions/${submission.id}`} variant="primary">
                    Предоставить обратную связь
                  </Button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default SubmissionsPage;