import React from "react";
import { Link } from "react-router-dom";
import styles from "./Button.module.scss";

const Button = ({
  children,
  onClick,
  to,
  viewTransition,
  variant = "default",
  size = "medium",
  isActive = false,
  icon: Icon,
  fullWidth = false,
  className = "",
  isLoading = false,
  loadingText = "Загрузка...",
  ...props
}) => {
  const classes = [
    styles.button,
    styles[`button--${variant}`],
    styles[`button--${size}`],
    isActive ? styles.active : "",
    fullWidth ? styles["button--full-width"] : "",
    Icon ? styles["button--icon"] : "",
    isLoading ? styles["button--loading"] : "",
    className,
  ]
    .filter(Boolean)
    .join(" ");

  const content = (
    <>
      {isLoading ? (
        <>
          <span className={styles.spinner}></span>
          {loadingText}
        </>
      ) : (
        <>
          {Icon && <Icon className={styles.icon} />}
          {children}
        </>
      )}
    </>
  );

  if (to) {
    return (
      <Link
        to={to}
        className={classes}
        viewTransition={viewTransition}
        {...props}
      >
        {content}
      </Link>
    );
  }

  return (
    <button onClick={onClick} className={classes} {...props}>
      {content}
    </button>
  );
};

export default Button;
