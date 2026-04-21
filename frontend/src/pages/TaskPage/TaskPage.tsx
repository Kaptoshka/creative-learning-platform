import React, { useCallback, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useTaskPage } from "@/hooks/useTaskPage";
import { useAlert } from "@/context/AlertContext";
import Loader from "@/components/Loader";

import AbbreviationTask from "@/components/AssignmentTypes/AbbreviationTask";
import AlliterationTask from "@/components/AssignmentTypes/AlliterationTask";
import CombineTask from "@/components/AssignmentTypes/CombineTask";
import StoryTask from "@/components/AssignmentTypes/StoryTask";
import UnexpectedConnectionsTask from "@/components/AssignmentTypes/UnexpectedConnectionsTask";
import UseCaseTask from "@/components/AssignmentTypes/UseCaseTask";
import ThirtyCirclesTask from "@/components/AssignmentTypes/ThirtyCirclesTask";

import Button from "@/components/Button";

import styles from "./TaskPage.module.scss";

const TaskPage = () => {
    const navigate = useNavigate();
    const alert = useAlert();
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

    useEffect(() => {
        if (error) {
            alert.error(
                error?.message ?? "Не удалось загрузить задание.",
                "Ошибка",
            );
        }
    }, [error]);

    useEffect(() => {
        if (submitSuccess && successMessage) {
            alert.success(successMessage);
        }
    }, [submitSuccess]);

    const handleContentChange = useCallback(
        (content: unknown) => {
            setSubmissionContent((prev) => {
                if (JSON.stringify(prev) === JSON.stringify(content))
                    return prev;
                return content;
            });
        },
        [setSubmissionContent],
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
            case "Тридцать кругов":
                return <ThirtyCirclesTask {...taskProps} />;
            default:
                return <div>Unknown task type</div>;
        }
    };

    if (loading) {
        return (
            <div className={styles.pageContainer}>
                <Loader fullPage variant="bar" label="Загрузка задания…" />
            </div>
        );
    }

    if (error || !task) {
        return (
            <div className={styles.pageContainer}>
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
                    if (isFormValid()) handleSubmit(e);
                }}
            >
                {renderTaskByType()}

                <div className={styles.formActions}>
                    <Button
                        type="submit"
                        variant="primary"
                        disabled={!isFormValid() || isSubmitting}
                    >
                        {isSubmitting ? (
                            <Loader size="sm" variant="dots" />
                        ) : (
                            "Отправить задание"
                        )}
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
