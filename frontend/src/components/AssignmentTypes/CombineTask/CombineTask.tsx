import React, { useState, useEffect } from "react";
import styles from "./CombineTask.module.scss";

const CombineTask = ({ content, onContentChange, submission, isReview }) => {
  const { prompt, instructions, example, subjectA, subjectB } = content || {};

  const [object1, setObject1] = useState(subjectA || "");
  const [object2, setObject2] = useState(subjectB || "");
  const [combinedName, setCombinedName] = useState("");
  const [description, setDescription] = useState("");

  useEffect(() => {
    if (isReview && submission) {
      setObject1(submission.objects?.object1 || "");
      setObject2(submission.objects?.object2 || "");
      setCombinedName(submission.combinedName || "");
      setDescription(submission.description || "");
    }
  }, [submission, isReview]);

  // Update parent component whenever content changes
  useEffect(() => {
    if (typeof onContentChange === "function" && !isReview) {
      const hasContent =
        object1.trim() ||
        object2.trim() ||
        combinedName.trim() ||
        description.trim();

      if (hasContent) {
        const submissionContent = {
          type: "Два в одном",
          objects: {
            object1: object1.trim(),
            object2: object2.trim(),
          },
          combinedName: combinedName.trim(),
          description: description.trim(),
          combination: `${object1.trim()} + ${object2.trim()} = ${combinedName.trim()}`,
        };
        onContentChange(submissionContent);
      } else {
        onContentChange(null);
      }
    }
  }, [
    object1,
    object2,
    combinedName,
    description,
    onContentChange,
    isReview,
  ]);

  const isFormComplete = () => {
    return (
      object1.trim() &&
      object2.trim() &&
      combinedName.trim() &&
      description.trim().length >= 10
    );
  };

  return (
    <div className={`${styles.taskTypeContainer} ${styles.combineTask}`}>
      <p className={styles.instructions}>
        {prompt ||
          "Возьмите два предмета и подумайте, как их можно объединить в один полезный или забавный предмет."}
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

      <div className={styles.combineForm}>
        {/* Objects Selection */}
        <div className={`${styles.section} ${styles.objectsSection}`}>
          <h3>Выберите два предмета</h3>
          <p className={styles.sectionHint}>
            Выберите любые два предмета, которые хотите объединить
          </p>

          <div className={styles.objectsGrid}>
            <div className="form-group">
              <label htmlFor="object1">Первый предмет:</label>
              <input
                id="object1"
                type="text"
                className="form-group__input"
                placeholder="Например: зонт"
                value={object1}
                onChange={(e) => setObject1(e.target.value)}
                disabled={isReview}
              />
            </div>

            <div className={styles.plusDivider}>
              <span className={styles.plusIcon}>+</span>
            </div>

            <div className="form-group">
              <label htmlFor="object2">Второй предмет:</label>
              <input
                id="object2"
                type="text"
                className="form-group__input"
                placeholder="Например: фонарик"
                value={object2}
                onChange={(e) => setObject2(e.target.value)}
                disabled={isReview}
              />
            </div>
          </div>
        </div>

        {/* Combination Result */}
        <div className={`${styles.section} ${styles.resultSection}`}>
          <h3>Результат объединения</h3>
          <p className={styles.sectionHint}>
            Придумайте название и опишите получившийся предмет
          </p>

          <div className="form-group">
            <label htmlFor="combinedName">Название нового предмета:</label>
            <input
              id="combinedName"
              type="text"
              className={`form-group__input ${styles.combinedInput}`}
              placeholder="Например: Светящийся зонт"
              value={combinedName}
              onChange={(e) => setCombinedName(e.target.value)}
              disabled={isReview}
            />
          </div>

          {/* Visual Formula */}
          {(object1 || object2 || combinedName) && (
            <div className={styles.combinationFormula}>
              <div className={styles.formulaItem}>
                <span className={styles.formulaLabel}>Предмет 1</span>
                <span className={styles.formulaValue}>{object1 || "?"}</span>
              </div>
              <span className={styles.formulaOperator}>+</span>
              <div className={styles.formulaItem}>
                <span className={styles.formulaLabel}>Предмет 2</span>
                <span className={styles.formulaValue}>{object2 || "?"}</span>
              </div>
              <span className={styles.formulaOperator}>=</span>
              <div className={`${styles.formulaItem} ${styles.result}`}>
                <span className={styles.formulaLabel}>Результат</span>
                <span className={styles.formulaValue}>
                  {combinedName || "?"}
                </span>
              </div>
            </div>
          )}
        </div>

        {/* Description Section */}
        <div className={`${styles.section} ${styles.descriptionSection}`}>
          <h3>Описание</h3>
          <p className={styles.sectionHint}>
            Опишите, как работает новый предмет и для чего его можно
            использовать
          </p>

          <div className="form-group">
            <label htmlFor="description">Подробное описание:</label>
            <textarea
              id="description"
              className="form-group__textarea"
              rows={5}
              placeholder="Например: Светящийся зонт для безопасных прогулок в темноте. Встроенный фонарик освещает путь впереди, а также делает вас заметнее для водителей."
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              disabled={isReview}
            />
            <div className={styles.characterCount}>
              {description.length} символов (минимум 10)
            </div>
          </div>
        </div>

        {/* Validation Hints */}
        {!isFormComplete() &&
          (object1 || object2 || combinedName || description) &&
          !isReview && (
            <div className={styles.validationHints}>
              <p className={styles.hintTitle}>
                Для завершения задания необходимо:
              </p>
              <ul>
                {!object1.trim() && <li>Ввести первый предмет</li>}
                {!object2.trim() && <li>Ввести второй предмет</li>}
                {!combinedName.trim() && (
                  <li>Придумать название нового предмета</li>
                )}
                {description.trim().length < 10 && (
                  <li>Написать описание (минимум 10 символов)</li>
                )}
              </ul>
            </div>
          )}

        {/* Completion Indicator */}
        {isFormComplete() && !isReview && (
          <div className={styles.completionIndicator}>
            ✓ Отлично! Ваша креативная комбинация готова к отправке.
          </div>
        )}
      </div>
    </div>
  );
};

export default CombineTask;
