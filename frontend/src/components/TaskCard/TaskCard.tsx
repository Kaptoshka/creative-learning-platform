import React from "react";
import { Clock } from "lucide-react";
import styles from "./TaskCard.module.scss";

const statusText = {
  available: "Доступно",
  in_progress: "В процессе",
  submitted: "Завершено",
};

const TaskCard = ({ task, onSelect }) => {
  const taskStatus = task.status || "available";

  return (
    <div className={styles.taskCard} onClick={() => onSelect(task)}>
      <div className={styles.taskCardHeader}>
        <h3 className={styles.taskCardTitle}>{task.content.title}</h3>
        <span
          className={`${styles.taskCardStatus} ${styles[`task_card__status--${taskStatus}`]}`}
        >
          {statusText[taskStatus]}
        </span>
      </div>
      <p className={styles.taskCardDescription}>{task.content.description}</p>
      <p className={styles.taskCardPrompt}>{task.content.prompt}</p>
      <div className={styles.taskCardFooter}>
        <div className={styles.taskCardTime}>
          <Clock className={styles.taskCardTimeIcon} />
          {new Date(task.deadline_time).toLocaleDateString()}
        </div>
        <span className={styles.taskCardDifficulty}>
          Учитель ID: {task.teacher_id}
        </span>
      </div>
    </div>
  );
};

export default TaskCard;
