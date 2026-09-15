import React from "react";
import HeroSection from "./HeroSection";
import FeaturesSection from "./FeaturesSection";
import styles from "./HomePage.module.scss";

const HomePage = () => (
    <div className={styles.homepage}>
        <HeroSection />
        <FeaturesSection />
    </div>
);

export default HomePage;
