import React, { useEffect, useState, useCallback } from "react";
import { useParams, useNavigate } from "react-router-dom";
import apiClient from "@/services/apiClient";
import offlineQueue from "@/services/offlineQueue";

import AbbreviationTask from "@/components/AssignmentTypes/AbbreviationTask";
import AlliterationTask from "@/components/AssignmentTypes/AlliterationTask";
import CombineTask from "@/components/AssignmentTypes/CombineTask";
import StoryTask from "@/components/AssignmentTypes/StoryTask";
import UnexpectedConnectionsTask from "@/components/AssignmentTypes/UnexpectedConnectionsTask";
import UseCaseTask from "@/components/AssignmentTypes/UseCaseTask";

import Button from "@/components/Button";

import styles from "./TaskPage.module.scss";
import tasksPageStyles from "@/pages/TasksPage/TasksPage.module.scss";

const TaskPage = () => {
  const { id } = useParams();
  const navigate = useNavigate();

  const [task, setTask] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [submissionContent, setSubmissionContent] = useState(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitSuccess, setSubmitSuccess] = useState(false);
  const [successMessage, setSuccessMessage] = useState("");

  useEffect(() => {
    setLoading(true);
    apiClient
      .get(`/assignments/${id}`)
      .then((res) => {
        setTask(res.data);
      })
      .catch((err) => {
        setError(
          "Failed to load task. It might not exist or you may not have permission.",
        );
        console.error(err);
      })
      .finally(() => {
        setLoading(false);
      });
  }, [id]);

  const handleContentChange = useCallback((content) => {
    setSubmissionContent((prev) => {
      if (JSON.stringify(prev) === JSON.stringify(content)) return prev;
      return content;
    });
  }, []);

  const handleSubmit = async (e) => {
    if (e) e.preventDefault();

    if (!submissionContent) {
      setError("Пожалуйста, заполните задание перед отправкой.");
      return;
    }

    setIsSubmitting(true);
    setError("");
    setSubmitSuccess(false);

    try {
      const contentToSubmit =
        typeof submissionContent === "string"
          ? submissionContent
          : JSON.stringify(submissionContent);

      const payload = {
        assignment_id: Number(id),
        content: contentToSubmit,
      };

      const response = await offlineQueue.send("/submissions", payload);

      if (response.status === 201 || response.status === "offline") {
        const message =
          response.status === "offline"
            ? "Интернета нет. Задание сохранено и будет отправлено позже. Перенаправление..."
            : "✓ Задание успешно отправлено! Перенаправление...";
        setSuccessMessage(message);
        setSubmitSuccess(true);

        setTimeout(() => {
          setIsSubmitting(false);
          navigate("/dashboard", { viewTransition: true });
        }, 2000);
      }
    } catch (err) {
      if (err.response?.data?.error) {
        setError(err.response.data.error);
      } else {
        setError("An unexpected error occurred during submission.");
      }
      console.error("Submission failed:", err);
      setIsSubmitting(false);
    }
  };

  const isFormValid = () => {
    if (!submissionContent) return false;

    if (typeof submissionContent === "object") {
      return JSON.stringify(submissionContent).length > 10;
    }

    return submissionContent.trim().length > 0;
  };

  const renderTaskByType = () => {
    const taskProps = {
      content: task.content,
      onContentChange: handleContentChange,
    };

    switch (task.content.type) {
      case "Аббревиатуры":
        return <AbbreviationTask {...taskProps} />;

      case "Неожиданные связи":
        return <UnexpectedConnectionsTask {...taskProps} />;

      case "Нестандартные применения":
        return <UseCaseTask {...taskProps} />;

      case "Рассказ из 100 слов":
        return <StoryTask {...taskProps} />;

      case "Два в одном":
        return <CombineTask {...taskProps} />;

      case "Одна буква":
        return <AlliterationTask {...taskProps} />;

      default:
        return (
          <div className={styles.errorMessage}>
            Неизвестный тип задания: {task.content.type}
          </div>
        );
    }
  };

  if (loading) {
    return (
      <div className={`${styles.pageContainer} loading`}>
        <div className="loading__spinner"></div>
        <p className="loading__text">Загрузка задания...</p>
      </div>
    );
  }

  if (error && !task) {
    return (
      <div className={`${styles.pageContainer} "error-container"`}>
        <div className="error-message">
          <h2>Ошибка</h2>
          <p>{error}</p>
          <Button
            variant="outline"
            onClick={() => navigate("/tasks", { viewTransition: true })}
          >
            Вернуться к заданиям
          </Button>
        </div>
      </div>
    );
  }

  if (!task) {
    return (
      <div className={`${styles.pageContainer} "error-container"`}>
        <div className="error-message">
          <h2>Задание не найдено</h2>
          <p>Запрошенное задание не существует или недоступно.</p>
          <Button
            variant="outline"
            onClick={() => navigate("/tasks", { viewTransition: true })}
          >
            Вернуться к заданиям
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className={`${styles.pageContainer} ${tasksPageStyles.tasksPage}`}>
      <header className={styles.taskPageHeader}>
        <div className={styles.headerTop}>
          <h1 className={styles.taskPageTitle}>
            {task.content.title || "Задание"}
          </h1>

          <Button
            variant="outline"
            size="small"
            onClick={() => navigate("/dashboard", { viewTransition: true })}
          >
            ✕ Закрыть
          </Button>
        </div>

        {task.content.description && (
          <p className={styles.taskPageDescription}>
            {task.content.description}
          </p>
        )}

        <div className={styles.taskPageMeta}>
          <span className={styles.taskStatus} data-status={task.status}>
            {task.status === "completed"
              ? "Выполнено"
              : task.status === "in_progress"
                ? "В процессе"
                : "Доступно"}
          </span>
        </div>
      </header>

      <main className={styles.taskPageContent}>{renderTaskByType()}</main>

      <footer className={styles.taskPageFooter}>
        <div className={styles.taskPageSubmission}>
          <form onSubmit={handleSubmit} className={styles.submissionForm}>
            <div className={styles.formActions}>
              <Button
                type="submit"
                variant="primary"
                fullWidth={true}
                disabled={
                  isSubmitting || task.status === "submitted" || !isFormValid()
                }
                isLoading={isSubmitting}
                loadingText="Отправка..."
              >
                Отправить задание
              </Button>
            </div>

            {task.status === "submitted" && (
              <p className={styles.completionNotice}>
                ✓ Это задание уже выполнено.
              </p>
            )}

            {!isFormValid() && submissionContent !== null && (
              <p className={styles.validationNotice}>
                Заполните все обязательные поля задания
              </p>
            )}
          </form>
        </div>
      </footer>

      {error && <div className={"status-message error-message"}>{error}</div>}

      {submitSuccess && (
        <div className={"status-message success-message"}>{successMessage}</div>
      )}
    </div>
  );
};

export default TaskPage;
