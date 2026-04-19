import React, { useState, useEffect } from "react";
import { Link, useParams, useNavigate, useLocation } from "react-router-dom";
import styles from "./ReviewPage.module.scss";

import AbbreviationTask from "@/components/AssignmentTypes/AbbreviationTask";
import AlliterationTask from "@/components/AssignmentTypes/AlliterationTask";
import CombineTask from "@/components/AssignmentTypes/CombineTask";
import StoryTask from "@/components/AssignmentTypes/StoryTask";
import UnexpectedConnectionsTask from "@/components/AssignmentTypes/UnexpectedConnectionsTask";
import UseCaseTask from "@/components/AssignmentTypes/UseCaseTask";

import apiClient from "@/services/apiClient";

const ReviewPage = () => {
  const { submissionId } = useParams();
  const navigate = useNavigate();
  const location = useLocation();

  const [loading, setLoading] = useState(true);
  const [submission, setSubmission] = useState(null);
  const [task, setTask] = useState(null);
  const [selectedCategories, setSelectedCategories] = useState([]);
  const [strengths, setStrengths] = useState("");
  const [improvements, setImprovements] = useState("");
  const [nextSteps, setNextSteps] = useState("");
  const [encouragement, setEncouragement] = useState("");
  const [showSuccess, setShowSuccess] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const advisoryCategories = [
    { id: "creativity", label: "Creative Approach", icon: "🎨" },
    { id: "originality", label: "Originality", icon: "💡" },
    { id: "technique", label: "Technical Execution", icon: "✨" },
    { id: "composition", label: "Composition/Structure", icon: "🎯" },
    { id: "expression", label: "Artistic Expression", icon: "🌟" },
    { id: "improvement", label: "Growth Potential", icon: "🚀" },
  ];

  const renderTaskByType = () => {
    if (!task || !submission || !submission.content) return null;

    const componentProps = {
      content: task.content,
      submission: submission.content,
      isReview: true,
    };

    const taskType = submission.taskType || (task.content && task.content.type);

    switch (taskType) {
      case "Аббревиатуры":
        return <AbbreviationTask {...componentProps} />;

      case "Неожиданные связи":
        return <UnexpectedConnectionsTask {...componentProps} />;

      case "Нестандартные применения":
        return <UseCaseTask {...componentProps} />;

      case "Рассказ из 100 слов":
        return <StoryTask {...componentProps} />;

      case "Два в одном":
        return <CombineTask {...componentProps} />;

      case "Одна буква":
        return <AlliterationTask {...componentProps} />;

      default:
        return (
          <div className={styles.errorMessage}>
            Неизвестный тип задания: {taskType}
          </div>
        );
    }
  };

  useEffect(() => {
    if (location.state && location.state.submission) {
      const { submission: passedSubmission } = location.state;
      setSubmission(passedSubmission);
      setTask({
        title: passedSubmission.taskTitle,
        type: passedSubmission.taskType,
        difficulty: "N/A",
        category: "N/A",
        estimatedTime: "N/A",
        content: {
          type: passedSubmission.taskType,
          content: passedSubmission.preview,
        },
      });
      setLoading(false);
    } else if (submissionId && submissionId !== "undefined") {
      setLoading(true);
      apiClient
        .get(`/submissions/${submissionId}`)
        .then((res) => {
          setSubmission(res.data.submission);
          setTask(res.data.task);
        })
        .catch((err) => {
          console.error("Fetch error:", err);
          setSubmission(null);
          setTask(null);
        })
        .finally(() => {
          setLoading(false);
        });
    } else {
      setLoading(false);
    }
  }, [submissionId, location.state]);

  const handleCategoryToggle = (categoryId) => {
    setSelectedCategories((prev) =>
      prev.includes(categoryId)
        ? prev.filter((id) => id !== categoryId)
        : [...prev, categoryId],
    );
  };

  const handleProvideFeedback = async (e) => {
    e.preventDefault();

    if (!strengths.trim() && !improvements.trim() && !nextSteps.trim()) {
      alert("Please provide at least one section of feedback");
      return;
    }

    setSubmitting(true);

    const feedbackData = {
      categories: selectedCategories,
      strengths: strengths,
      improvements: improvements,
      nextSteps: nextSteps,
      encouragement: encouragement,
    };

    try {
      await apiClient.post(
        `/submissions/${submissionId}/feedback`,
        feedbackData,
      );
      setShowSuccess(true);
      setTimeout(() => {
        navigate("/submissions");
      }, 2500);
    } catch (error) {
      console.error("Failed to send feedback:", error);
      alert("Failed to send feedback. Please try again.");
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <div className={styles.feedbackPage}>
        <div className={styles.feedbackPageContainer}>
          <div className={styles.feedbackLoading}>
            <div className={styles.feedbackLoadingSpinner}></div>
            <span className={styles.feedbackLoadingText}>
              Загрузка ответа...
            </span>
          </div>
        </div>
      </div>
    );
  }

  if (!submission || !task) {
    return (
      <div className={styles.feedbackPage}>
        <div className={styles.feedbackPageContainer}>
          <div className={`${styles.alert} ${styles.alertWarning}`}>
            <svg
              className={styles.alertIcon}
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <div className={styles.alertContent}>
              <div className={styles.alertTitle}>Submission not found</div>
              <div className={styles.alertMessage}>
                Ответ, который вы ищете, не существует или был удален.
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.feedbackPage}>
      <div className={styles.feedbackPageContainer}>
        {/* Header */}
        <header className={styles.feedbackPageHeader}>
          <div className={styles.feedbackPageHeaderTop}>
            <Link to="/submissions" className={styles.feedbackPageHeaderBack}>
              <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M15 19l-7-7 7-7"
                />
              </svg>
              Назад к ответам
            </Link>
          </div>

          <h1 className={styles.feedbackPageHeaderTitle}>
            Предоставьте обратную связь
          </h1>
          <p className={styles.feedbackPageHeaderSubtitle}>
            Предоставьте конструктивную обратную связь, чтобы помочь студенту
            улучшить и развить свои творческие навыки
          </p>

          <div className={styles.feedbackPageHeaderMeta}>
            <div className={styles.feedbackPageHeaderItem}>
              <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
                />
              </svg>
              <span>{task.title}</span>
            </div>
            <div className={styles.feedbackPageHeaderItem}>
              <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              <span>
                Отправлено {new Date(submission.submittedAt).toLocaleString()}
              </span>
            </div>
            <div
              className={`${styles.statusBadge} ${styles.statusBadgeSubmitted}`}
            >
              <span className={styles.statusBadgeDot}></span>
              Ожидает отзыва
            </div>
          </div>
        </header>

        {/* Success Alert */}
        {showSuccess && (
          <div className={`${styles.alert} ${styles.alertSuccess}`}>
            <svg
              className={styles.alertIcon}
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <div className={styles.alertContent}>
              <div className={styles.alertTitle}>
                Обратная связь отправлена успешно! 🎉
              </div>
              <div className={styles.alertMessage}>
                Студент получит вашу обратную связь и советы. Перенаправление...
              </div>
            </div>
          </div>
        )}

        {/* Info Alert */}
        <div className={`${styles.alert} ${styles.alertInfo}`}>
          <svg
            className={styles.alertIcon}
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <div className={styles.alertContent}>
            <div className={styles.alertTitle}>
              Сконцентрируйтесь на творческом росте
            </div>
            <div className={styles.alertMessage}>
              Оцените творческую работу студента и предоставьте конструктивную
              обратную связь, которая стимулирует экспериментацию и развитие
              творческих способностей.
            </div>
          </div>
        </div>

        {/* Main Layout */}
        <div className={styles.feedbackLayout}>
          {/* Main Content */}
          <div className={styles.feedbackLayoutMain}>
            {/* Submission Details */}
            <div className={styles.submissionCard}>
              <div className={styles.submissionCardHeader}>
                <h2 className={styles.submissionCardTitle}>
                  <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                    />
                  </svg>
                  Ответ на задание
                </h2>

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
                {/* Student's Written Content */}
                {submission.preview && (
                  <div className={styles.submissionCardSection}>
                    <h3 className={styles.submissionCardSectionTitle}>
                      <svg
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M4 6h16M4 12h16m-7 6h7"
                        />
                      </svg>
                      Описание студента
                    </h3>
                    <div className={styles.submissionCardContent}>
                      {submission.preview}
                    </div>
                  </div>
                )}

                <div className={styles.submissionCardSection}>
                  <h3 className={styles.submissionCardSectionTitle}>
                    <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                      />
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                      />
                    </svg>
                    Предоставленная работа
                  </h3>
                  {renderTaskByType()}
                </div>
              </div>
            </div>

            <form
              onSubmit={handleProvideFeedback}
              className={styles.feedbackForm}
            >
              <div className={styles.feedbackFormHeader}>
                <h2 className={styles.feedbackFormTitle}>
                  <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
                    />
                  </svg>
                  Оцените работу
                </h2>
                <p className={styles.feedbackFormDescription}>
                  Оцените работу студента и предоставьте конструктивную обратную
                  связь, чтобы помочь ему развить свои творческие навыки.
                  Основывайтесь на сильных сторонах, конструктивных
                  рекомендациях и конкретных шагах, которые могут быть выполнены
                  в будущем.
                </p>
              </div>

              <div className={styles.feedbackFormGroup}>
                <label className={styles.feedbackFormLabel}>
                  Выберите категорию, которую вы оцениваете
                </label>
                <div className={styles.advisoryCategories}>
                  {advisoryCategories.map((category) => (
                    <div
                      key={category.id}
                      className={`${styles.advisoryCategoriesItem} ${
                        selectedCategories.includes(category.id)
                          ? styles.advisoryCategoriesItemActive
                          : ""
                      }`}
                      onClick={() => handleCategoryToggle(category.id)}
                    >
                      <div className={styles.advisoryCategoriesCheckbox}></div>
                      <span className={styles.advisoryCategoriesLabel}>
                        {category.icon} {category.label}
                      </span>
                    </div>
                  ))}
                </div>
              </div>

              <div className={styles.feedbackFormGroup}>
                <label
                  className={`${styles.feedbackFormLabel} ${styles.feedbackFormLabelRequired}`}
                >
                  ✅ Что сделано хорошо
                </label>
                <textarea
                  className={styles.feedbackFormTextarea}
                  placeholder="Выделите сильные стороны работы. Какие техники, подходы или арт-решения были особенно эффективны?"
                  value={strengths}
                  onChange={(e) => setStrengths(e.target.value)}
                  style={{ minHeight: "120px" }}
                />
              </div>

              {/* Areas for Growth */}
              <div className={styles.feedbackFormGroup}>
                <label className={styles.feedbackFormLabel}>
                  🎯 Куда можно расти
                </label>
                <textarea
                  className={styles.feedbackFormTextarea}
                  placeholder="Конструктивные предложения для развития творческих навыков. Основывайтесь на конкретных, действенных комментариях..."
                  value={improvements}
                  onChange={(e) => setImprovements(e.target.value)}
                  style={{ minHeight: "120px" }}
                />
              </div>

              {/* Next Steps */}
              <div className={styles.feedbackFormGroup}>
                <label className={styles.feedbackFormLabel}>
                  🚀 Рекомендованные следующие шаги
                </label>
                <textarea
                  className={styles.feedbackFormTextarea}
                  placeholder="Предлагайте конкретные действия, упражниения или ресурсы для дальнейшего творческого роста..."
                  value={nextSteps}
                  onChange={(e) => setNextSteps(e.target.value)}
                  style={{ minHeight: "120px" }}
                />
              </div>

              {/* Encouragement */}
              <div className={styles.feedbackFormGroup}>
                <label className={styles.feedbackFormLabel}>
                  💪 Слова поддержки
                </label>
                <textarea
                  className={styles.feedbackFormTextarea}
                  placeholder="Закончите мотивационными словами, чтобы вдохновить на дальнейшее развитие..."
                  value={encouragement}
                  onChange={(e) => setEncouragement(e.target.value)}
                  style={{ minHeight: "80px" }}
                />
              </div>

              {/* Actions */}
              <div className={styles.feedbackFormActions}>
                <button
                  type="submit"
                  className={styles.feedbackFormSubmit}
                  disabled={
                    submitting ||
                    (!strengths.trim() &&
                      !improvements.trim() &&
                      !nextSteps.trim())
                  }
                >
                  {submitting ? (
                    <>
                      <div
                        style={{
                          width: "1rem",
                          height: "1rem",
                          border: "2px solid white",
                          borderTop: "2px solid transparent",
                          borderRadius: "50%",
                          animation: "spin 1s linear infinite",
                        }}
                      ></div>
                      Sending Feedback...
                    </>
                  ) : (
                    <>
                      <svg
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"
                        />
                      </svg>
                      Отправить
                    </>
                  )}
                </button>
                <button
                  type="button"
                  className={styles.feedbackFormCancel}
                  onClick={() => navigate("/submissions")}
                  disabled={submitting}
                >
                  Отмена
                </button>
              </div>
            </form>
          </div>

          {/* Sidebar */}
          <div className={styles.feedbackLayoutSidebar}>
            {/* Task Information */}
            <div className={styles.infoSidebarCard}>
              <h3 className={styles.infoSidebarTitle}>
                <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
                Информация о задании
              </h3>
              <div>
                <div className={styles.infoSidebarItem}>
                  <span className={styles.infoSidebarLabel}>Тип</span>
                  <span className={styles.infoSidebarValue}>{task.type}</span>
                </div>
                <div className={styles.infoSidebarItem}>
                  <span className={styles.infoSidebarLabel}>
                    Затраченное время
                  </span>
                  <span className={styles.infoSidebarValue}>
                    {task.estimatedTime}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ReviewPage;
