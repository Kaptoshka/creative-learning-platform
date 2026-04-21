import React, { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { useTask } from "@/hooks/useTask";
import { useAlert } from "@/context/AlertContext";
import Loader from "@/components/Loader";
import Button from "@/components/Button";
import PageContainer from "@/components/ui/PageContainer";

const TaskDetailPage = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const alert = useAlert();

    const [submissionContent, setSubmissionContent] = useState("");
    const { task, loading, error, submit } = useTask({ taskId: id });

    useEffect(() => {
        if (error) {
            alert.error(
                error?.message ?? "Не удалось загрузить задание.",
                "Ошибка",
            );
        }
    }, [error]);

    const handleSubmit = async (e) => {
        e.preventDefault();

        if (!submissionContent.trim()) {
            alert.warning("Заполните поле перед отправкой.");
            return;
        }

        const success = await submit(submissionContent);
        if (success) {
            alert.success("Задание успешно отправлено!");
            navigate("/tasks", { viewTransition: true });
        } else {
            alert.error(
                "Не удалось отправить задание. Попробуйте снова.",
                "Ошибка отправки",
            );
        }
    };

    if (loading) {
        return <Loader fullPage variant="bar" label="Загрузка задания…" />;
    }

    if (error || !task) {
        return (
            <PageContainer>
                <Button variant="outline" onClick={() => navigate("/tasks")}>
                    Вернуться к заданиям
                </Button>
            </PageContainer>
        );
    }

    return (
        <PageContainer className="task-detail-page">
            <header className="task-detail-header">
                <h1>{task.title}</h1>
                <p>{task.description}</p>
            </header>

            <div className="submission-form-container">
                <h3>Твои задания</h3>
                <form onSubmit={handleSubmit} noValidate>
                    <div className="form-group">
                        <label htmlFor="submissionContent">
                            Enter your work below:
                        </label>
                        <textarea
                            id="submissionContent"
                            value={submissionContent}
                            onChange={(e) =>
                                setSubmissionContent(e.target.value)
                            }
                            rows={10}
                            placeholder="Start writing your creative response here..."
                            required
                        />
                    </div>
                    <Button type="submit" variant="primary">
                        Отправь работу
                    </Button>
                </form>
            </div>
        </PageContainer>
    );
};

export default TaskDetailPage;
