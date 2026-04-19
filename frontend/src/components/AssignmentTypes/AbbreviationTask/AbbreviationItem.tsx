import React from "react";
import styles from "./AbbreviationTask.module.scss"; // Предполагается, что стили будут общими или переданы

const AbbreviationItem = ({
  index,
  item,
  updateAbbreviation,
  isValidWord,
  isValidExpansion,
  isReview,
}) => {
  return (
    <div key={index} className={styles.abbreviationItem}>
      <div className={styles.formGroup}>
        <label htmlFor={`word-${index}`}>{index + 1}. Слово (3-5 букв):</label>
        <input
          id={`word-${index}`}
          type="text"
          className={styles.formGroupInput}
          placeholder="Например: РОТ"
          value={item.word}
          onChange={(e) =>
            updateAbbreviation(index, "word", e.target.value.toUpperCase())
          }
          maxLength={5}
          disabled={isReview}
        />
        {item.word && !isValidWord(item.word) && !isReview && (
          <p className={styles.errorMessage}>
            Слово должно содержать 3-5 русских букв
          </p>
        )}
      </div>

      <div className={styles.formGroup}>
        <label htmlFor={`expansion-${index}`}>Расшифровка:</label>
        <input
          id={`expansion-${index}`}
          type="text"
          className={styles.formGroupInput}
          placeholder="Например: Рыбный Омлет с Тмином"
          value={item.expansion}
          onChange={(e) =>
            updateAbbreviation(index, "expansion", e.target.value)
          }
          disabled={isReview}
        />
        {item.word &&
          item.expansion &&
          !isValidExpansion(item.expansion, item.word) &&
          !isReview && (
            <p className={styles.errorMessage}>
              Первые буквы слов в расшифровке должны соответствовать буквам
              аббревиатуры
            </p>
          )}
      </div>
    </div>
  );
};

export default AbbreviationItem;
