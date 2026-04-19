import React from "react";
import { List } from "lucide-react";
import styles from "./QuickStats.module.scss";

interface QuickStatsProps {
  items: Array<{
    label: string;
    value: string | number;
    icon?: React.ComponentType<{ size?: number }>;
  }>;
}

const QuickStats = ({ items }: QuickStatsProps) => {
  return (
    <div className={styles.quickStats}>
      {items.map((item, index) => (
        <div key={index} className={styles.quickStatsCard}>
          <div className={styles.quickStatsHeader}>
            <div className={styles.quickStatsIcon}>
              {item.icon ? <item.icon size={24} /> : <List size={24} />}
            </div>
          </div>
          <div className={styles.quickStatsValue}>{item.value}</div>
          <div className={styles.quickStatsLabel}>{item.label}</div>
        </div>
      ))}
    </div>
  );
};

export default QuickStats;
