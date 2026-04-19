import React from "react";
import styles from "./FeatureItem.module.scss";

// eslint-disable-next-line no-unused-vars
const FeatureItem = ({ icon: Icon, colorClass, title, description }) => (
  <div className={styles.featuresItem}>
    <div className={`${styles.featuresIcon} ${colorClass}`}>
      <Icon className={styles.featuresIconSvg} />
    </div>
    <h3 className={styles.featuresTitle}>{title}</h3>
    <p className={styles.featuresDescription}>{description}</p>
  </div>
);

export default FeatureItem;
