import React, { useState, useEffect } from "react";
import { Link, useNavigate } from "react-router-dom";
import styles from "./SubmissionsPage.module.scss";
import apiClient from "@/services/apiClient";
import { useAuth } from "@/hooks/useAuth";
import Button from "@/components/Button/Button.jsx";

const SubmissionsPage = () => {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [loading, setLoading] = useState(true);
  const [submissions, setSubmissions] = useState([]);
  const [filter, setFilter] = useState("all");
  const [searchQuery, setSearchQuery] = useState("");
  const [sortBy, setSortBy] = useState("date");

  useEffect(() => {
    const fetchSubmissions = async () => {
      if (!user) {
        setLoading(false);
        return;
      }
      setLoading(true);
      try {
        const submissionsResponse = await apiClient.get(
          `/submissions/teacher/${user.id}`,
        );
        const submissionsData = submissionsResponse.data || [];

        const detailedSubmissions = await Promise.all(
          submissionsData.map(async (sub) => {
            try {
              const assignmentPromise = apiClient.get(
                `/assignments/${sub.AssignmentID}`,
              );

              const studentPromise = apiClient.get(`/users/${sub.StudentID}`);

              const [assignmentRes, studentRes] = await Promise.all([
                assignmentPromise,
                studentPromise,
              ]);

              const assignment = assignmentRes.data;
              const student = Array.isArray(studentRes.data)
                ? studentRes.data[0]
                : studentRes.data;

              const submittedAt = new Date(sub.SubmittedAt);
              const daysWaiting = Math.floor(
                (new Date() - submittedAt) / (1000 * 60 * 60 * 24),
              );
              let priority = "normal";
              if (daysWaiting >= 3) {
                priority = "urgent";
              } else if (daysWaiting === 0) {
                priority = "new";
              }

              const assignmentContent = assignment.content || {};

              return {
                id: sub.ID,
                taskTitle: assignmentContent.title || "Untitled Assignment",
                taskType: assignmentContent.type || "unknown",
                studentName: `${student?.first_name} ${student?.last_name}`,
                studentEmail: student?.email || "Unknown Email",
                studentAvatar: student?.first_name
                  ? student.first_name[0]
                  : "?",
                submittedAt: submittedAt,
                daysWaiting: daysWaiting,
                priority: priority,
                hasContent: !!sub.Content,
                content: sub.Content, // Pass full submission content
                preview:
                  sub.Content?.text ||
                  sub.Content?.story ||
                  sub.Content?.sentence ||
                  "",
              };
            } catch (error) {
              console.error(
                `Failed to fetch details for submission ${sub.ID}:`,
                error,
              );
              return null;
            }
          }),
        );

        setSubmissions(detailedSubmissions.filter(Boolean));
      } catch (error) {
        console.error("Failed to fetch submissions:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchSubmissions();
  }, [user]);

  const handleReview = (submission) => {
    console.log("submission:", submission);
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
      filtered.sort((a, b) => b.submittedAt - a.submittedAt);
    } else if (sortBy === "student") {
      filtered.sort((a, b) => a.studentName.localeCompare(b.studentName));
    } else if (sortBy === "task") {
      filtered.sort((a, b) => a.taskTitle.localeCompare(b.taskTitle));
    }

    return filtered;
  };

  const getPriorityColor = (priority) => {
    switch (priority) {
      case "urgent":
        return styles.submissionCardUrgent;
      case "new":
        return styles.submissionCardNew;
      default:
        return "";
    }
  };

  const filteredSubmissions = getFilteredSubmissions();

  if (loading) {
    return (
      <div className={styles.submissionsPage}>
        <div className={styles.submissionsPageContainer}>
          <div className={styles.loading}>
            <div className={styles.loadingSpinner}></div>
            <span className={styles.loadingText}>Загрузка ответов...</span>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.submissionsPage}>
      <div className={styles.submissionsPageContainer}>
        {/* Header */}
        <header className={styles.submissionsPageHeader}>
          <div>
            <h1 className={styles.submissionsPageHeaderTitle}>
              Ответы на задания
            </h1>
            <p className={styles.submissionsPageHeaderSubtitle}>
              Оценивайте и дайте обратную связь студентам, чтобы они могли
              дальше развивать свои креативные навыки
            </p>
          </div>
          <div className={styles.submissionsPageHeaderStats}>
            <div className={styles.statItem}>
              <div className={styles.statItemValue}>{submissions.length}</div>
              <div className={styles.statItemLabel}>Всего ответов</div>
            </div>
            <div className={styles.statItem}>
              <div
                className={`${styles.statItemValue} ${styles.statItemValueUrgent}`}
              >
                {submissions.filter((s) => s.priority === "urgent").length}
              </div>
              <div className={styles.statItemLabel}>Давно отправлены</div>
            </div>
            <div className={styles.statItem}>
              <div
                className={`${styles.statItemValue} ${styles.statItemValueNew}`}
              >
                {submissions.filter((s) => s.priority === "new").length}
              </div>
              <div className={styles.statItemLabel}>Отправлены сегодня</div>
            </div>
          </div>
        </header>

        {/* Filters and Search */}
        <div className={styles.controls}>
          <div className={styles.controlsFilters}>
            <Button
              variant={filter === "all" ? "active" : "default"}
              onClick={() => setFilter("all")}
              className={styles.controlsFilterBtn}
            >
              <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M4 6h16M4 12h16M4 18h16"
                />
              </svg>
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
              <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              Давно отправлены
              <span
                className={`${styles.controlsBadge} ${styles.controlsBadgeUrgent}`}
              >
                {submissions.filter((s) => s.priority === "urgent").length}
              </span>
            </Button>
            <Button
              variant={filter === "recent" ? "active" : "default"}
              onClick={() => setFilter("recent")}
              className={`${styles.controlsFilterBtn} ${
                filter === "recent" ? styles.controlsFilterBtnActive : ""
              }`}
            >
              <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              Недавно отправлены
              <span className={styles.controlsBadge}>
                {submissions.filter((s) => s.daysWaiting <= 1).length}
              </span>
            </Button>
          </div>

          <div className={styles.controlsRight}>
            <div className={styles.controlsSearch}>
              <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                />
              </svg>
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

        {/* Submissions List */}
        {filteredSubmissions.length === 0 ? (
          <div className={styles.emptyState}>
            <div className={styles.emptyStateIcon}>
              <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                />
              </svg>
            </div>
            <h3 className={styles.emptyStateTitle}>
              Ответы на задания не найдены
            </h3>
            <p className={styles.emptyStateMessage}>
              {searchQuery
                ? "Try adjusting your search query"
                : "All submissions have been reviewed!"}
            </p>
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
                  <div
                    className={`${styles.priorityBadge} ${styles.priorityBadgeUrgent}`}
                  >
                    <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                      />
                    </svg>
                    Давно отправленны
                  </div>
                )}
                {submission.priority === "new" && (
                  <div
                    className={`${styles.priorityBadge} ${styles.priorityBadgeNew}`}
                  >
                    <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z"
                      />
                    </svg>
                    New
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
                    <span className={styles.submissionCardTaskIcon}>
                      {getTaskIcon(submission.taskType)}
                    </span>
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
                      <svg
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                      </svg>
                      <span>
                        {submission.daysWaiting === 0
                          ? "Сегодня"
                          : `${submission.daysWaiting} дней ${submission.daysWaiting > 1 ? " " : ""} назад`}
                      </span>
                    </div>
                    <div className={styles.submissionCardMetaItem}>
                      <svg
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                        />
                      </svg>
                      <span>{submission.submittedAt.toLocaleDateString()}</span>
                    </div>
                  </div>

                  <Button
                    to={`/submissions/${submission.id}`}
                    variant="primary"
                  >
                    <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                      />
                    </svg>
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
