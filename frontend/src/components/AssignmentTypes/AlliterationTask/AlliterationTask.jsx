import React, { useState, useEffect } from "react";
import Button from "@/components/Button";
import styles from "./AlliterationTask.module.scss";

const AlliterationTask = ({ content, onContentChange, submission, isReview }) => {
  const { prompt, instructions, example, letter } = content || {};
  const [sentence, setSentence] = useState("");
  const [selectedLetter, setSelectedLetter] = useState(letter || "");

  const russianAlphabet = [
    "А",
    "Б",
    "В",
    "Г",
    "Д",
    "Е",
    "Ё",
    "Ж",
    "З",
    "И",
    "Й",
    "К",
    "Л",
    "М",
    "Н",
    "О",
    "П",
    "Р",
    "С",
    "Т",
    "У",
    "Ф",
    "Х",
    "Ц",
    "Ч",
    "Ш",
    "Щ",
    "Э",
    "Ю",
    "Я",
  ];

  useEffect(() => {
    if (isReview && submission) {
      setSentence(submission.sentence || "");
      setSelectedLetter(submission.letter || "");
    }
  }, [submission, isReview]);

  const countWords = (str) => {
    return str.trim().split(/\s+/).filter(Boolean).length;
  };

  const currentWordCount = countWords(sentence);
  const minWords = 5;

  const isValid =
    sentence.trim() &&
    sentence
      .trim()
      .split(/\s+/)
      .every(
        (w) =>
          w.length > 0 && w[0].toLowerCase() === selectedLetter.toLowerCase(),
      );

  useEffect(() => {
    if (typeof onContentChange === "function" && !isReview) {
      if (sentence.trim() && selectedLetter) {
        const submissionContent = {
          type: "Одна буква",
          letter: selectedLetter,
          sentence: sentence.trim(),
          wordCount: currentWordCount,
          isValid: isValid,
        };
        onContentChange(submissionContent);
      } else {
        onContentChange(null);
      }
    }
  }, [
    sentence,
    selectedLetter,
    currentWordCount,
    isValid,
    onContentChange,
    isReview,
  ]);

  const isFormComplete = () => {
    return selectedLetter && isValid && currentWordCount >= minWords;
  };

  const getInvalidWords = () => {
    if (!sentence.trim() || !selectedLetter) return [];
    return sentence
      .trim()
      .split(/\s+/)
      .filter(
        (w) =>
          w.length > 0 && w[0].toLowerCase() !== selectedLetter.toLowerCase(),
      );
  };

  const invalidWords = getInvalidWords();

  return (
    <div className={`${styles.taskTypeContainer} ${styles.alliterationTask}`}>
      <p className={styles.instructions}>
        {prompt ||
          "Составьте предложение, где все слова будут начинаться на одну букву."}
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

      <div className={styles.alliterationForm}>
        {/* Letter Selection */}
        {!letter && !isReview && (
          <div className={`${styles.section} ${styles.letterSection}`}>
            <h3>Выберите букву</h3>
            <p className={styles.sectionHint}>
              Выберите букву, на которую будут начинаться все слова в
              предложении
            </p>

            <div className={styles.letterGrid}>
              {russianAlphabet.map((char) => (
                <Button
                  key={char}
                  type="button"
                  variant="letter"
                  className={selectedLetter === char ? styles.active : ""}
                  onClick={() => setSelectedLetter(char)}
                  disabled={isReview}
                >
                  {char}
                </Button>
              ))}
            </div>
          </div>
        )}

        {/* Display selected or preset letter */}
        {(selectedLetter || letter) && (
          <div className={styles.selectedLetterDisplay}>
            <div className={styles.letterBadge}>
              <span className={styles.badgeLabel}>Выбранная буква:</span>
              <span className={styles.badgeLetter}>
                {selectedLetter || letter}
              </span>
            </div>
            {!letter && !isReview && (
              <Button
                type="button"
                variant="outline"
                size="small"
                onClick={() => {
                  setSelectedLetter("");
                  setSentence("");
                }}
                disabled={isReview}
              >
                Изменить букву
              </Button>
            )}
          </div>
        )}

        {/* Sentence Writing */}
        {(selectedLetter || letter) && (
          <div className={`${styles.section} ${styles.sentenceSection}`}>
            <h3>Напишите предложение</h3>
            <p className={styles.sectionHint}>
              Все слова должны начинаться на букву «{selectedLetter || letter}».
              Минимум {minWords} слов.
            </p>

            <div className={styles.formGroup}>
              <label htmlFor="sentence">Ваше предложение:</label>
              <textarea
                id="sentence"
                className={`form-group__textarea ${
                  isValid && currentWordCount >= minWords
                    ? styles.valid
                    : sentence.trim() && invalidWords.length > 0
                      ? styles.invalid
                      : ""
                }`}
                rows={6}
                placeholder={`Например: Разительного роста растение росло рядом с рощей...`}
                value={sentence}
                onChange={(e) => setSentence(e.target.value)}
                disabled={isReview}
              />

              {/* Word Counter */}
              <div className={styles.wordInfo}>
                <span
                  className={`${styles.wordCount} ${
                    currentWordCount >= minWords ? styles.sufficient : ""
                  }`}
                >
                  Слов: {currentWordCount}{" "}
                  {currentWordCount >= minWords ? "✓" : `/ ${minWords}`}
                </span>
                {isValid && currentWordCount >= minWords && (
                  <span className={styles.validationCheck}>
                    ✓ Все слова на букву «{selectedLetter || letter}»
                  </span>
                )}
              </div>
            </div>

            {/* Validation Messages */}
            {sentence.trim() && invalidWords.length > 0 && !isReview && (
              <div className={styles.errorMessage}>
                ⚠️ Некоторые слова не начинаются на букву «
                {selectedLetter || letter}»:
                <span className={styles.invalidWords}>
                  {invalidWords.join(", ")}
                </span>
              </div>
            )}

            {sentence.trim() &&
              currentWordCount < minWords &&
              invalidWords.length === 0 &&
              !isReview && (
                <div className={styles.warningMessage}>
                  Добавьте еще {minWords - currentWordCount}{" "}
                  {minWords - currentWordCount === 1
                    ? "слово"
                    : minWords - currentWordCount < 5
                      ? "слова"
                      : "слов"}
                </div>
              )}
          </div>
        )}

        {/* Tips Section */}
        {!isReview && (
          <div className={`${styles.section} ${styles.tipsSection}`}>
            <h3>Советы</h3>
            <div className={styles.tipsList}>
              <div className={styles.tipItem}>
                <span className={styles.tipIcon}>💡</span>
                <span>
                  Смысл предложения не так важен, главное — все слова на одну
                  букву
                </span>
              </div>
              <div className={styles.tipItem}>
                <span className={styles.tipIcon}>📖</span>
                <span>
                  Используйте прилагательные и наречия для увеличения количества
                  слов
                </span>
              </div>
              <div className={styles.tipItem}>
                <span className={styles.tipIcon}>✨</span>
                <span>
                  Попробуйте создать забавное или необычное предложение
                </span>
              </div>
            </div>
          </div>
        )}

        {/* Completion Indicator */}
        {isFormComplete() && !isReview && (
          <div className={styles.completionIndicator}>
            ✓ Отлично! Ваше предложение готово к отправке.
          </div>
        )}
      </div>
    </div>
  );
};

export default AlliterationTask;
