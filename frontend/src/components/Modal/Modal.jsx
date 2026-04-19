import React from "react";
import ReactDOM from "react-dom";
import { X } from "lucide-react";
import Button from "@/components/Button";
import styles from "./Modal.module.scss";

const Modal = ({ task, onClose, onAction }) => {
  const modalRoot = document.getElementById("modal-root");

  if (!modalRoot) {
    console.error(
      "Элемент #modal-root не найден в DOM. Модальное окно не будет отображаться корректно.",
    );
    return null;
  }

  return ReactDOM.createPortal(
    <div className={styles.modal}>
      <div className={styles.modalContent}>
        <div className={styles.modalHeader}>
          <h2 className={styles.modalTitle}>{task.content.title}</h2>
          <Button onClick={onClose} className={styles.modalClose} size="small">
            <X className={styles.modalCloseIcon} />
          </Button>
        </div>

        <div className={styles.modalBody}>
          <p className={styles.modalDescription}>{task.content.description}</p>
          <div className={styles.modalStats}>
            {onAction.stats.map((stat, index) => (
              <div key={stat.label || index} className={styles.statCard}>
                <div className={styles.statCardIcon}>
                  <stat.icon className={styles.statCardIconSvg} />
                  <span className={styles.statCardLabel}>{stat.label}</span>
                </div>
                <p
                  className={`${styles.statCardValue} ${styles[`stat_card_value--${stat.color}`]}`}
                >
                  {stat.value}
                </p>
              </div>
            ))}
          </div>
        </div>

        <div className={styles.modalFooter}>
          {onAction.actionButtons.map((btn, index) => (
            <Button
              key={btn.text || index}
              onClick={btn.onClick}
              variant={btn.variant || "default"}
              className={btn.className}
            >
              {btn.text}
            </Button>
          ))}

          <Button onClick={onClose} variant="outline">
            Закрыть
          </Button>
        </div>
      </div>
    </div>,
    modalRoot,
  );
};

export default Modal;
