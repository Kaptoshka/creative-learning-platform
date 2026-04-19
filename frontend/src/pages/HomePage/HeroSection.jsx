import React, { useContext } from "react";
import { useNavigate } from "react-router-dom";
import Button from "@/components/Button";
import { AuthContext } from "@/context/AuthContext";

import styles from "./HomePage.module.scss";

const HeroSection = () => {
  const navigate = useNavigate();
  const { user } = useContext(AuthContext);

  const handleStartLearning = () => {
    if (user) {
      navigate("/dashboard", { viewTransition: true });
    } else {
      navigate("/auth", { viewTransition: true });
    }
  };

  return (
    <div className={styles.homepageHero}>
      <h1 className={styles.homepageTitle}>
        Развивайте свою{" "}
        <span className={styles.homepageTitleAccent}>креативность</span>
      </h1>
      <p className={styles.homepageSubtitle}>
        Платформа для творческого обучения с интересными заданиями и
        отслеживанием прогресса
      </p>
      <div className={styles.homepageButtons}>
        <Button variant="primary" onClick={handleStartLearning}>
          Начать обучение
        </Button>
        <Button variant="secondary">Узнать больше</Button>
      </div>
    </div>
  );
};

export default HeroSection;
