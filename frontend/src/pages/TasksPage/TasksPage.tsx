import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { X, Clock, User, BookOpen, CheckCircle } from "lucide-react";
import TaskCard from "@/components/TaskCard";
import Modal from "@/components/Modal";
import { useAuth } from "@/hooks/useAuth";
import { useFetch } from "@/hooks/useFetch";
import apiClient from "@/services/apiClient";
import Loading from "@/components/ui/Loading";
import ErrorMessage from "@/components/ui/ErrorMessage";

import styles from "./TasksPage.module.scss";
import dashboardStyles from "@/pages/DashboardPage/DashboardPage.module.scss";

const TasksPage = () => {
  const [tasks, setTasks] = useState([]);
  const [selectedTask, setSelectedTask] = useState(null);
  const navigate = useNavigate();
  const { user } = useAuth();

  const url = user
    ? user.role === "student"
      ? `/assignments/student/${user.id}`
      : `/assignments/teacher/${user.id}`
    : null;

  const { data, loading, error, refetch } = useFetch(url);

  useEffect(() => {
    if (data) {
      setTasks(data || []);
    }
  }, [data]);

  const modalStats = selectedTask
    ? [
        {
          label: "Срок сдачи задания",
          value: new Date(selectedTask.deadline_time).toLocaleDateString(),
          icon: Clock,
          color: "blue",
        },
        {
          label: "Пример",
          value: selectedTask.content.example,
          icon: User,
          color: "purple",
        },
        {
          label: "Тип задания",
          value: selectedTask.content.type,
          icon: BookOpen,
          color: "green",
        },
      ]
    : [];

  const actionButtons = selectedTask
    ? [
        user.role === "teacher" && {
          text: "Редактировать задание",
          variant: "secondary",
          onClick: () => {
            console.log("edit");
            // add edit endpoint
            navigate(`/tasks/${selectedTask.ID}`, { viewTransition: true });
          },
        },
        user.role === "student" &&
          selectedTask.status === "available" && {
            text: "Начать задание",
            variant: "primary",
            onClick: () => {
              console.log("started id:", selectedTask.ID);
              navigate(`/tasks/${selectedTask.ID}`, { viewTransition: true });
            },
          },
        selectedTask.status === "in_progress" && {
          text: "Продолжить",
          variant: "primary",
          onClick: () => {
            console.log("continued");
            navigate(`/tasks/${selectedTask.ID}`, { viewTransition: true });
          },
        },
        selectedTask.status === "completed" && {
          text: "Просмотреть результат",
          variant: "secondary",
          onClick: () => {
            console.log("result");
            navigate(`/tasks/${selectedTask.ID}`, { viewTransition: true });
          },
        },
      ].filter(Boolean)
    : [];

  if (loading) {
    return <Loading text="Загрузка заданий..." />;
  }

  if (error) {
    return <ErrorMessage error={error} />;
  }

  return (
    <div className={styles.tasksPage}>
      <div className={styles.tasksPageContainer}>
        <div className={dashboardStyles.dashboardSectionHeader}>
          <h2 className={dashboardStyles.dashboardSectionTitle}>
            <CheckCircle />
            Мои задания
          </h2>
        </div>

        <div className={styles.tasksGrid}>
          {tasks.length > 0 ? (
            tasks.map((task) => (
              <TaskCard
                key={task.id || task.ID || Math.random()}
                task={task}
                onSelect={setSelectedTask}
              />
            ))
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

export default TasksPage;
