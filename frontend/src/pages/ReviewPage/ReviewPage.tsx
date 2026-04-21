import AbbreviationTask from "@/components/AssignmentTypes/AbbreviationTask";
import AlliterationTask from "@/components/AssignmentTypes/AlliterationTask";
import CombineTask from "@/components/AssignmentTypes/CombineTask";
import StoryTask from "@/components/AssignmentTypes/StoryTask";
import UnexpectedConnectionsTask from "@/components/AssignmentTypes/UnexpectedConnectionsTask";
import UseCaseTask from "@/components/AssignmentTypes/UseCaseTask";

import Button from "@/components/Button";
import Loader from "@/components/Loader";
import { useAlert } from "@/context/AlertContext";
import ReviewForm from "@/components/ui/ReviewForm";

import styles from "./ReviewPage.module.scss";

const ReviewPage = () => {
    const navigate = useNavigate();
    const alert = useAlert();
    const {
        submission,
        task,
        review,
        loading,
        error,
        submitting,
        showSuccess,
        setReviewField,
        submitReview,
    } = useReview();

    useEffect(() => {
        if (error) {
            alert.error(
                error?.message ?? "Не удалось загрузить ответ.",
                "Ошибка загрузки",
            );
        }
    }, [error]);

    useEffect(() => {
        if (showSuccess) {
            alert.success("Обратная связь успешно отправлена!");
        }
    }, [showSuccess]);

    const renderTaskByType = () => {
        if (!task || !submission?.content) return null;

        const componentProps = {
            content: task.content,
            submission: submission.content,
            isReview: true,
        };

        const taskType = submission.taskType || (task.content?.type as string);

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
                return <div>Неизвестный тип задания: {taskType}</div>;
        }
    };

    if (loading) {
        return <Loader fullPage variant="bar" label="Загрузка ответа…" />;
    }

    if (error && !submission) {
        return (
            <div className={styles.feedbackPage}>
                <Button
                    variant="outline"
                    onClick={() => navigate("/submissions")}
                >
                    Вернуться к ответам
                </Button>
            </div>
        );
    }

    if (!submission || !task) {
        return (
            <div className={styles.feedbackPage}>
                <Button
                    variant="outline"
                    onClick={() => navigate("/submissions")}
                >
                    Вернуться к ответам
                </Button>
            </div>
        );
    }

    return (
        <div className={styles.feedbackPage}>
            <div className={styles.feedbackPageContainer}>
                <header className={styles.feedbackPageHeader}>
                    <h1 className={styles.feedbackPageTitle}>
                        Проверка: {submission.taskTitle}
                    </h1>
                    <p className={styles.feedbackPageSubtitle}>
                        Студент: {submission.studentName}
                    </p>
                </header>

                <div className={styles.submissionCardSection}>
                    <h3>Ответ студента</h3>
                    {renderTaskByType()}
                </div>

                <div className={styles.feedbackFormSection}>
                    <h3>Обратная связь</h3>
                    <ReviewForm
                        review={review}
                        onChange={setReviewField}
                        disabled={submitting}
                    />
                </div>

                <div className={styles.feedbackFormActions}>
                    <Button
                        variant="primary"
                        onClick={submitReview}
                        disabled={submitting}
                    >
                        {submitting ? (
                            <Loader size="sm" variant="dots" inline />
                        ) : (
                            "Отправить обратную связь"
                        )}
                    </Button>
                    <Button
                        variant="outline"
                        onClick={() => navigate("/submissions")}
                    >
                        Назад к ответам
                    </Button>
                </div>
            </div>
        </div>
    );
};

export default ReviewPage;
