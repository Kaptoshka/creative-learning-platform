import React, { useState, useEffect } from "react";
import Button from "@/components/Button";
import styles from "./UseCaseTask.module.scss";

const UseCaseTask = ({ content, onContentChange, submission, isReview }) => {
  const { prompt, instructions, example, subject } = content || {};
  const minIdeas = 3;

  const [ideas, setIdeas] = useState(["", "", ""]);
  const [customSubject, setCustomSubject] = useState(subject || "");

  useEffect(() => {
    if (isReview && submission) {
      setCustomSubject(submission.subject || "");
      const submittedIdeas = submission.ideas || [];
      // Pad array to at least minIdeas length for the inputs
      while (submittedIdeas.length < minIdeas) {
        submittedIdeas.push("");
      }
      setIdeas(submittedIdeas);
    }
  }, [submission, isReview]);

  // Update parent component whenever content changes
  useEffect(() => {
    if (typeof onContentChange === "function" && !isReview) {
      const filledIdeas = ideas.filter((idea) => idea.trim());

      if (filledIdeas.length > 0 || customSubject.trim()) {
        const submissionContent = {
          type: "Нестандартные применения",
          subject: customSubject.trim() || "предмет",
          ideas: filledIdeas.map((idea) => idea.trim()),
          totalIdeas: filledIdeas.length,
        };
        onContentChange(submissionContent);
      } else {
        onContentChange(null);
      }
    }
  }, [ideas, customSubject, onContentChange, isReview]);

  const updateIdea = (index, value) => {
    const newIdeas = [...ideas];
    newIdeas[index] = value;
    setIdeas(newIdeas);
  };

  const addIdea = () => {
    setIdeas([...ideas, ""]);
  };

  const removeIdea = (index) => {
    if (ideas.length > minIdeas) {
      const newIdeas = ideas.filter((_, i) => i !== index);
      setIdeas(newIdeas);
    }
  };

  const isFormComplete = () => {
    return (
      customSubject.trim() &&
      ideas.filter((idea) => idea.trim()).length >= minIdeas
    );
  };

  const filledIdeasCount = ideas.filter((idea) => idea.trim()).length;

  return (
    <div className={`${styles.taskTypeContainer} ${styles.useCaseTask}`}>
      <p className={styles.instructions}>
        {prompt ||
          "Придумайте необычные способы применения для обычного предмета."}
      </p>

      {example && (
        <div className={styles.example}>
          <strong>Пример:</strong>
          <p>{example}</p>
        </div>
      )}

      {instructions && instructions.length > 0 && (
        <div className={styles.instructions}>
          <strong>Инструкции:</strong>
          <ol>
            {instructions.map((instruction, index) => (
              <li key={index}>{instruction}</li>
            ))}
          </ol>
        </div>
      )}

      <div className={styles.useCaseForm}>
        {/* Subject Selection */}
        <div className={`${styles.section} ${styles.subjectSection}`}>
          <h3>Выберите предмет</h3>
          <p className={styles.sectionHint}>
            Выберите любой обычный предмет или введите свой вариант
          </p>

          <div className="form-group">
            <label htmlFor="subject">Предмет для исследования:</label>
            <input
              id="subject"
              type="text"
              className={`form-group__input ${styles.subjectInput}`}
              placeholder="Например: карандаш, ложка, скрепка, бутылка..."
              value={customSubject}
              onChange={(e) => setCustomSubject(e.target.value)}
              disabled={isReview}
            />
          </div>
        </div>
        <div className={`${styles.section} ${styles.ideasSection}`}>
          <h3>Ваши идеи (минимум {minIdeas})</h3>
          <p className={styles.sectionHint}>
            Придумайте необычные, креативные способы использования предмета
          </p>

          <div className={styles.ideasList}>
            {ideas.map((idea, index) => (
              <div key={index} className={styles.ideaItem}>
                <div className="form-group">
                  <label htmlFor={`idea-${index}`}>Идея {index + 1}:</label>
                  <div className={styles.inputWithAction}>
                    <input
                      id={`idea-${index}`}
                      type="text"
                      className="form-group__input"
                      placeholder={
                        index === 0
                          ? "Например: закручивать бумагу для поделок"
                          : index === 1
                            ? "Например: размешивать напитки"
                            : "Например: использовать как снаряд для рогатки"
                      }
                      value={idea}
                      onChange={(e) => updateIdea(index, e.target.value)}
                      disabled={isReview}
                    />
                    {ideas.length > minIdeas && !isReview && (
                      <Button
                        type="button"
                        variant="danger"
                        onClick={() => removeIdea(index)}
                        title="Удалить идею"
                      >
                        ✕
                      </Button>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>

          {!isReview && (
            <Button
              type="button"
              variant="outline"
              className={styles.buttonAdd}
              onClick={addIdea}
            >
              + Добавить еще идею
            </Button>
          )}
        </div>

        {/* Progress Indicator */}
        {!isReview && (
          <div className={styles.progressSection}>
            <div className={styles.progressBarContainer}>
              <div className={styles.progressLabel}>
                <span>Прогресс:</span>
                <span className={styles.progressCount}>
                  {filledIdeasCount} / {minIdeas}{" "}
                  {minIdeas === filledIdeasCount ? "✓" : ""}
                </span>
              </div>
              <div className={styles.progressBar}>
                <div
                  className={`${styles.progressFill} ${
                    filledIdeasCount >= minIdeas ? styles.complete : ""
                  }`}
                  style={{
                    width: `${Math.min(
                      (filledIdeasCount / minIdeas) * 100,
                      100,
                    )}%`,
                  }}
                />
              </div>
            </div>
          </div>
        )}

        {/* Validation Hints */}
        {!isFormComplete() &&
          (customSubject || ideas.some((idea) => idea.trim())) &&
          !isReview && (
            <div className={styles.validationHints}>
              <p className={styles.hintTitle}>Для завершения задания:</p>
              <ul>
                {!customSubject.trim() && (
                  <li>Выберите или введите предмет</li>
                )}
                {filledIdeasCount < minIdeas && (
                  <li>
                    Придумайте еще {minIdeas - filledIdeasCount}{" "}
                    {minIdeas - filledIdeasCount === 1
                      ? "идею"
                      : minIdeas - filledIdeasCount < 5
                        ? "идеи"
                        : "идей"}
                  </li>
                )}
              </ul>
            </div>
          )}

        {/* Completion Indicator */}
        {isFormComplete() && !isReview && (
          <div className={styles.completionIndicator}>
            ✓ Отлично! Вы придумали {filledIdeasCount}{" "}
            {filledIdeasCount === 1
              ? "идею"
              : filledIdeasCount < 5
                ? "идеи"
                : "идей"}
            . Задание готово к отправке.
          </div>
        )}
      </div>
    </div>
  );
};

export default UseCaseTask;
