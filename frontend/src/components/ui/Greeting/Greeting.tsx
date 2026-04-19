import React from "react";
import styles from "./Greeting.module.scss";

interface GreetingProps {
  name?: string;
  role?: string;
  greeting: string;
  formattedDate: string;
  formattedTime: string;
}

const Greeting = ({ name, role, greeting, formattedDate, formattedTime }: GreetingProps) => {
  return (
    <div className={styles.greeting}>
      <h1 className={styles.greetingTitle}>
        {greeting}, {name || "гость"}! 👋
      </h1>
      <p className={styles.greetingSubtitle}>
        Добро пожаловать в твой {role || "личный"} кабинет
      </p>
      <div className={styles.greetingStats}>
        <span className={styles.greetingStat}>
          📅 {formattedDate}
        </span>
        <span className={styles.greetingStat}>
          👤 Роль: {role ? role.charAt(0).toUpperCase() + role.slice(1) : "—"}
        </span>
        <span className={styles.greetingStat}>
          ⏰ {formattedTime}
        </span>
      </div>
    </div>
  );
};

export default Greeting;
