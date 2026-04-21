import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { Clock, User } from "lucide-react";
import Modal from "@/components/Modal";
import Loader from "@/components/Loader";
import { useAlert } from "@/context/AlertContext";
import { useAuth } from "@/hooks/useAuth";
import { useFetch } from "@/hooks/useFetch";

const TeacherDashboardPage = () => {
    const [tasks, setTasks] = useState([]);
    const [selectedTask, setSelectedTask] = useState(null);
    const navigate = useNavigate();
    const alert = useAlert();
    const { user } = useAuth();

    const url = user ? `/assignments/student/${user.id}` : null;
    const { data, loading, error } = useFetch(url);

    useEffect(() => {
        if (data?.assignments) {
            setTasks(data.assignments);
        }
    }, [data]);

    useEffect(() => {
        if (error) {
            alert.error(
                error?.message ?? "Не удалось загрузить задания.",
                "Ошибка загрузки",
            );
        }
    }, [error]);

    const modalStats = selectedTask
        ? [
              {
                  label: "Время выполнения",
                  value: `${selectedTask.timeLimit} мин`,
                  icon: Clock,
                  color: "blue",
              },
              {
                  label: "Сложность",
                  value: selectedTask.difficulty,
                  icon: User,
                  color: "purple",
              },
          ]
        : [];

    const actionButtons = selectedTask
        ? [
              selectedTask.status === "available" && {
                  text: "Начать задание",
                  className: "button button--primary button--flex",
                  onClick: () =>
                      navigate(`/tasks/${selectedTask.id}`, {
                          viewTransition: true,
                      }),
              },
              selectedTask.status === "in_progress" && {
                  text: "Продолжить",
                  className: "button button--success button--flex",
                  onClick: () =>
                      navigate(`/tasks/${selectedTask.id}`, {
                          viewTransition: true,
                      }),
              },
              selectedTask.status === "completed" && {
                  text: "Просмотреть результат",
                  className: "button button--secondary button--flex",
                  onClick: () =>
                      navigate(`/tasks/${selectedTask.id}`, {
                          viewTransition: true,
                      }),
              },
          ].filter(Boolean)
        : [];

    if (loading) {
        return <Loader fullPage variant="bar" label="Загрузка…" />;
    }

    return (
        <div className="tasks-page">
            <div className="tasks-page__container">
                <div className="tasks-page__header">
                    <h1 className="tasks-page__title">Мои задания</h1>
                    <p className="tasks-page__subtitle">
                        Выберите задание для выполнения
                    </p>
                </div>

                <div className="tasks-grid">
                    {tasks.length > 0
                        ? tasks.map((task) => console.log(task)) // TODO: заменить на карточку задания
                        : !error && <p>У вас пока нет заданий.</p>}
                </div>

                {selectedTask && (
                    <Modal
                        task={selectedTask}
                        onClose={() => setSelectedTask(null)}
                        onAction={{ stats: modalStats, actionButtons }}
                    />
                )}
            </div>
        </div>
    );
};

export default TeacherDashboardPage;
