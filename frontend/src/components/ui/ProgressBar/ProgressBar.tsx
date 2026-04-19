import React from "react";
import styles from "./ProgressBar.module.scss";

interface ProgressBarProps {
  label?: string;
  value: number;
  max?: number;
}

const ProgressBar = ({ label, value, max = 100 }: ProgressBarProps) => {
  const percentage = Math.min((value / max) * 100, 100);

  return (
    <div className={styles.progressBar}>
      {label && (
        <div className={styles.progressBarLabel}>
          <span>{label}</span>
          <span className={styles.progressBarValue}>{percentage}%</span>
        </div>
      )}
      <div className={styles.progressBarTrack}>
        <div
          className={styles.progressBarFill}
          style={{ width: `${percentage}%` }}
        />
      </div>
    </div>
  );
};

export default ProgressBar;
