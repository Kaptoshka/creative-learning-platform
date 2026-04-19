import React, { useState, useEffect, useContext } from "react";
import { useNavigate } from "react-router-dom";
import { Clock, User } from "lucide-react";
// import TaskCard from "../components/TaskCard";
import Modal from "@/components/Modal";
import { AuthContext } from "@/context/AuthContext";
import apiClient from "@/services/apiClient";

const TeacherDashboardPage = () => {
  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [selectedTask, setSelectedTask] = useState(null);
  const navigate = useNavigate();
  const { user } = useContext(AuthContext);

  useEffect(() => {
    if (!user) return;

    const fetchTasks = async () => {
      try {
        setLoading(true);
        const response = await apiClient.get(`/assignments/student/${user.id}`);
        setTasks(response.data.assignments || []);
      } catch (err) {
        setError("Failed to load tasks.");
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchTasks();
  }, [user]);

  // Подготовка данных для Modal
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

  // Кнопки действий для Modal
  const actionButtons = selectedTask
    ? [
        selectedTask.status === "available" && {
          text: "Начать задание",
          className: "button button--primary button--flex",
          onClick: () => {
            console.log("started");
            navigate(`/tasks/${selectedTask.id}`, { viewTransition: true });
          },
        },
        selectedTask.status === "in_progress" && {
          text: "Продолжить",
          className: "button button--success button--flex",
          onClick: () => {
            console.log("continued");
            navigate(`/tasks/${selectedTask.id}`, { viewTransition: true });
          },
        },
        selectedTask.status === "completed" && {
          text: "Просмотреть результат",
          className: "button button--secondary button--flex",
          onClick: () => {
            console.log("result");
            navigate(`/tasks/${selectedTask.id}`, { viewTransition: true });
          },
        },
      ].filter(Boolean)
    : [];

  if (loading) {
    return <div className="tasks-page">Загрузка...</div>;
  }

  if (error) {
    return <div className="tasks-page error-message">{error}</div>;
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
          {tasks.length > 0 ? (
            tasks.map((task) => console.log(task))
          ) : (
            <p>У вас пока нет заданий.</p>
          )}
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
