import React, { useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { useTask } from "@/hooks/useTask";
import Button from "@/components/Button";
import Loading from "@/components/ui/Loading";
import ErrorMessage from "@/components/ui/ErrorMessage";
import PageContainer from "@/components/ui/PageContainer";

const TaskDetailPage = () => {
    const { id } = useParams();
    const navigate = useNavigate();

    const [submissionContent, setSubmissionContent] = useState("");
    const { task, loading, error, submit } = useTask({ taskId: id });

    const handleSubmit = async (e) => {
        e.preventDefault();

        const success = await submit(submissionContent);
        if (success) {
            navigate("/tasks", { viewTransition: true });
        }
    };

    if (loading) {
        return <Loading text="Загрузка задания..." fullPage />;
    }

    if (error || !task) {
        return (
            <PageContainer>
                <ErrorMessage error={error || "Задание не найдено."} />
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
                <h3>Your Submission</h3>
                <form onSubmit={handleSubmit} noValidate>
                    {error && <ErrorMessage error={error} />}
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
                            rows="10"
                            placeholder="Start writing your creative response here..."
                            required
                        />
                    </div>
                    <Button type="submit" variant="primary">
                        Submit Task
                    </Button>
                </form>
            </div>
        </PageContainer>
    );
};

export default TaskDetailPage;
