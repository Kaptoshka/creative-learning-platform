import React, { useState, useEffect } from "react";
import Button from "@/components/Button";
import styles from "./StoryTask.module.scss";

const StoryTask = ({ content, onContentChange, submission, isReview }) => {
  const { prompt, instructions, example, wordCount } = content || {};

  const [text, setText] = useState("");
  const [selectedWordCount, setSelectedWordCount] = useState(wordCount || 100);
  const wordCountOptions = [50, 75, 100];

  useEffect(() => {
    if (isReview && submission) {
      setText(submission.story || "");
      setSelectedWordCount(submission.wordCount || 100);
    }
  }, [submission, isReview]);

  const countWords = (str) => {
    return str.trim().split(/\s+/).filter(Boolean).length;
  };

  const currentWordCount = countWords(text);
  const isExactMatch = currentWordCount === selectedWordCount;
  const difference = currentWordCount - selectedWordCount;

  useEffect(() => {
    if (typeof onContentChange === "function" && !isReview) {
      if (text.trim()) {
        const submissionContent = {
          type: "Рассказ из 100 слов",
          wordCount: selectedWordCount,
          actualWordCount: currentWordCount,
          story: text.trim(),
          isComplete: isExactMatch,
        };
        onContentChange(submissionContent);
      } else {
        onContentChange(null);
      }
    }
  }, [
    text,
    selectedWordCount,
    currentWordCount,
    isExactMatch,
    onContentChange,
    isReview,
  ]);

  const getProgressPercentage = () => {
    if (currentWordCount > selectedWordCount) return 100;
    return Math.min((currentWordCount / selectedWordCount) * 100, 100);
  };

  const getProgressColor = () => {
    if (isExactMatch) return styles.success;
    if (currentWordCount > selectedWordCount) return styles.over;
    if (currentWordCount >= selectedWordCount * 0.8) return styles.close;
    return styles.default;
  };

  return (
    <div className={`${styles.taskTypeContainer} ${styles.storyTask}`}>
      <p className={styles.instructions}>
        {prompt ||
          `Составьте рассказ из точно ${selectedWordCount} слов. Ни одного слова больше, ни одного слова меньше.`}
      </p>

      {example && (
        <div className={styles.example}>
          <strong>Примеры тем:</strong>
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

      <div className={styles.storyForm}>
        {!isReview && (
          <div className={`${styles.section} ${styles.wordCountSection}`}>
            <h3>Выберите уровень сложности</h3>
            <p className={styles.sectionHint}>
              Выберите количество слов для вашего рассказа
            </p>

            <div className={styles.wordCountOptions}>
              {wordCountOptions.map((count) => (
                <Button
                  key={count}
                  type="button"
                  variant="word-count-button"
                  className={selectedWordCount === count ? styles.active : ""}
                  onClick={() => setSelectedWordCount(count)}
                  disabled={isReview}
                >
                  <span className={styles.countNumber}>{count}</span>
                  <span className={styles.countLabel}>
                    {count === 50
                      ? "Сложно"
                      : count === 75
                        ? "Средне"
                        : "Стандарт"}
                  </span>
                </Button>
              ))}
            </div>
          </div>
        )}

        <div className={`${styles.section} ${styles.storySection}`}>
          <h3>Напишите ваш рассказ</h3>
          {!isReview && (
            <p className={styles.sectionHint}>
              Рассказ должен иметь начало, развитие и концовку
            </p>
          )}

          <div className={styles.formGroup}>
            <label htmlFor="story">Ваш рассказ:</label>
            <textarea
              id="story"
              className={`form-group__textarea ${styles.storyTextarea} ${
                isExactMatch
                  ? styles.valid
                  : text.trim() && currentWordCount > 0
                    ? styles.inProgress
                    : ""
              }`}
              rows={12}
              placeholder="Начните писать ваш рассказ здесь...&#10;&#10;Помните: у рассказа должно быть начало, развитие и концовка."
              value={text}
              onChange={(e) => setText(e.target.value)}
              disabled={isReview}
            />
          </div>

          <div className={`${styles.wordCounter} ${getProgressColor()}`}>
            <div className={styles.counterHeader}>
              <span className={styles.counterLabel}>Количество слов:</span>
              <span className={styles.counterValue}>
                {currentWordCount} / {selectedWordCount}
                {isExactMatch && <span className={styles.checkIcon}> ✓</span>}
              </span>
            </div>

            <div className={styles.progressBarWrapper}>
              <div className={styles.progressBar}>
                <div
                  className={`${styles.progressFill} ${getProgressColor()}`}
                  style={{ width: `${getProgressPercentage()}%` }}
                />
              </div>
            </div>

            <div className={styles.counterStatus}>
              {isExactMatch && (
                <span className={styles.statusSuccess}>
                  ✓ Отлично! Ровно {selectedWordCount} слов
                </span>
              )}
              {!isExactMatch && currentWordCount > 0 && (
                <span
                  className={`${styles.statusInfo} ${
                    difference > 0 ? styles.over : styles.under
                  }`}
                >
                  {difference > 0 ? (
                    <>
                      На {difference}{" "}
                      {difference === 1
                        ? "слово"
                        : difference < 5
                          ? "слова"
                          : "слов"}{" "}
                      больше
                    </>
                  ) : difference < 0 ? (
                    <>
                      Осталось {Math.abs(difference)}{" "}
                      {Math.abs(difference) === 1
                        ? "слово"
                        : Math.abs(difference) < 5
                          ? "слова"
                          : "слов"}
                    </>
                  ) : null}
                </span>
              )}
              {currentWordCount === 0 && !isReview && (
                <span className={styles.statusEmpty}>
                  Начните писать ваш рассказ
                </span>
              )}
            </div>
          </div>
        </div>

        {!isReview && (
          <div className={`${styles.section} ${styles.tipsSection}`}>
            <h3>Советы по написанию</h3>
            <div className={styles.tipsGrid}>
              <div className={styles.tipCard}>
                <div className={styles.tipIcon}>📖</div>
                <div className={styles.tipContent}>
                  <strong>Структура</strong>
                  <p>Начало → Развитие → Концовка</p>
                </div>
              </div>
              <div className={styles.tipCard}>
                <div className={styles.tipIcon}>✂️</div>
                <div className={styles.tipContent}>
                  <strong>Редактирование</strong>
                  <p>Пишите сначала больше, потом сокращайте</p>
                </div>
              </div>
              <div className={styles.tipCard}>
                <div className={styles.tipIcon}>🎯</div>
                <div className={styles.tipContent}>
                  <strong>Фокус</strong>
                  <p>Одна главная идея или событие</p>
                </div>
              </div>
            </div>
          </div>
        )}

        {!isExactMatch && text.trim() && !isReview && (
          <div className={styles.validationHints}>
            <p className={styles.hintTitle}>
              {difference > 0
                ? "⚠️ Нужно сократить текст"
                : "📝 Продолжайте писать"}
            </p>
            <ul>
              {difference > 0 ? (
                <li>
                  Удалите {difference}{" "}
                  {difference === 1
                    ? "слово"
                    : difference < 5
                      ? "слова"
                      : "слов"}{" "}
                  или перефразируйте текст короче
                </li>
              ) : (
                <li>
                  Добавьте еще {Math.abs(difference)}{" "}
                  {Math.abs(difference) === 1
                    ? "слово"
                    : Math.abs(difference) < 5
                      ? "слова"
                      : "слов"}{" "}
                  к вашему рассказу
                </li>
              )}
            </ul>
          </div>
        )}

        {isExactMatch && !isReview && (
          <div className={styles.completionIndicator}>
            ✓ Превосходно! Ваш рассказ содержит ровно {selectedWordCount} слов и
            готов к отправке.
          </div>
        )}
      </div>
    </div>
  );
};

export default StoryTask;
