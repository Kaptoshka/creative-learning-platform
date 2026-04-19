import React, { useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { useTaskPage } from "@/hooks/useTaskPage";
import offlineQueue from "@/services/offlineQueue";

import AbbreviationTask from "@/components/AssignmentTypes/AbbreviationTask";
import AlliterationTask from "@/components/AssignmentTypes/AlliterationTask";
import CombineTask from "@/components/AssignmentTypes/CombineTask";
import StoryTask from "@/components/AssignmentTypes/StoryTask";
import UnexpectedConnectionsTask from "@/components/AssignmentTypes/UnexpectedConnectionsTask";
import UseCaseTask from "@/components/AssignmentTypes/UseCaseTask";

import Button from "@/components/Button";
import Loading from "@/components/ui/Loading";
import ErrorMessage from "@/components/ui/ErrorMessage";

import styles from "./TaskPage.module.scss";

const TaskPage = () => {
  const navigate = useNavigate();
  const {
    task,
    loading,
    error,
    submissionContent,
    isSubmitting,
    submitSuccess,
    successMessage,
    setSubmissionContent,
    handleSubmit,
    isFormValid,
  } = useTaskPage();

  const handleContentChange = useCallback(
    (content: unknown) => {
      setSubmissionContent((prev) => {
        if (JSON.stringify(prev) === JSON.stringify(content)) return prev;
        return content;
      });
    },
    [setSubmissionContent]
  );

  const renderTaskByType = () => {
    if (!task) return null;

    const taskProps = {
      content: task.content,
      onContentChange: handleContentChange,
    };

    switch (task.content?.type) {
      case "Аббревиатуры":
        return <AbbreviationTask {...taskProps} />;
      case "Аллитерация":
        return <AlliterationTask {...taskProps} />;
      case "Два в одном":
        return <CombineTask {...taskProps} />;
      case "Рассказ из 100 слов":
        return <StoryTask {...taskProps} />;
      case "Неожиданные связи":
        return <UnexpectedConnectionsTask {...taskProps} />;
      case "Нестандартные применения":
        return <UseCaseTask {...taskProps} />;
      default:
        return <div>Unknown task type</div>;
    }
  };

  if (loading) {
    return (
      <div className={`${styles.pageContainer} loading`}>
        <Loading text="Загрузка задания..." fullPage />
      </div>
    );
  }

  if (error && !task) {
    return (
      <div className={`${styles.pageContainer}`}>
        <ErrorMessage error={error} />
        <Button variant="outline" onClick={() => navigate("/tasks")}>
          Вернуться к заданиям
        </Button>
      </div>
    );
  }

  if (!task) {
    return (
      <div className={`${styles.pageContainer}`}>
        <ErrorMessage error="Задание не найдено" />
        <Button variant="outline" onClick={() => navigate("/tasks")}>
          Вернуться к заданиям
        </Button>
      </div>
    );
  }

  return (
    <div className={styles.pageContainer}>
      <header className={styles.pageHeader}>
        <h1 className={styles.pageTitle}>{task.title}</h1>
        <p className={styles.pageDescription}>{task.description}</p>
      </header>

      <form
        className={styles.form}
        onSubmit={(e) => {
          e.preventDefault();
          if (isFormValid()) {
            handleSubmit(e);
          }
        }}
      >
        {renderTaskByType()}

        {submitSuccess && (
          <div className={styles.successMessage}>{successMessage}</div>
        )}

        {error && !submitSuccess && <ErrorMessage error={error} />}

        <div className={styles.formActions}>
          <Button
            type="submit"
            variant="primary"
            disabled={!isFormValid() || isSubmitting}
          >
            {isSubmitting ? "Отправка..." : "Отправить задание"}
          </Button>
          <Button
            variant="outline"
            onClick={() => navigate("/tasks")}
          >
            Назад к заданиям
          </Button>
        </div>
      </form>
    </div>
  );
};

export default TaskPage;