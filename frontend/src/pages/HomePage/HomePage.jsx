import React from "react";
import FeatureItem from "@/components/FeatureItem";
import { BookOpen, Clock, CheckCircle } from "lucide-react";
import HeroSection from "./HeroSection";
import styles from "./HomePage.module.scss";
import featureItemStyles from "@/components/FeatureItem/FeatureItem.module.scss";

const HomePage = () => {
  return (
    <div className={styles.homepage}>
      <div className={styles.homepageContainer}>
        <HeroSection />
        <div className={featureItemStyles.features}>
          <FeatureItem
            icon={BookOpen}
            colorClass={featureItemStyles.featuresIconBlue}
            title="Креативные задания"
            description="Разнообразные задачи для развития творческого мышления и навыков"
          />
          <FeatureItem
            icon={Clock}
            colorClass={featureItemStyles.featuresIconGreen}
            title="Отслеживание времени"
            description="Контроль времени выполнения заданий и анализ продуктивности"
          />
          <FeatureItem
            icon={CheckCircle}
            colorClass={featureItemStyles.featuresIconPurple}
            title="Прогресс"
            description="Отслеживание выполненных заданий и достижений в обучении"
          />
        </div>
      </div>
    </div>
  );
};

export default HomePage;
