import React, { useState, useEffect } from "react";
import styles from "./UnexpectedConnectionsTask.module.scss";

const UnexpectedConnectionsTask = ({
  content,
  onContentChange,
  submission,
  isReview,
}) => {
  const { prompt, instructions, example } = content || {};

  const [word1, setWord1] = useState("");
  const [word2, setWord2] = useState("");
  const [differences, setDifferences] = useState(["", "", ""]);
  const [similarities, setSimilarities] = useState(["", "", ""]);
  const [sentences, setSentences] = useState("");

  useEffect(() => {
    if (isReview && submission) {
      setWord1(submission.words?.word1 || "");
      setWord2(submission.words?.word2 || "");
      setSentences(submission.sentences || "");

      // Pad arrays to ensure 3 items for the inputs
      const paddedDifferences = [...(submission.differences || [])];
      while (paddedDifferences.length < 3) paddedDifferences.push("");
      setDifferences(paddedDifferences);

      const paddedSimilarities = [...(submission.similarities || [])];
      while (paddedSimilarities.length < 3) paddedSimilarities.push("");
      setSimilarities(paddedSimilarities);
    }
  }, [submission, isReview]);

  useEffect(() => {
    if (typeof onContentChange === "function" && !isReview) {
      const hasContent =
        word1.trim() ||
        word2.trim() ||
        differences.some((d) => d.trim()) ||
        similarities.some((s) => s.trim()) ||
        sentences.trim();

      if (hasContent) {
        const submissionContent = {
          type: "Неожиданные связи",
          words: {
            word1: word1.trim(),
            word2: word2.trim(),
          },
          differences: differences.filter((d) => d.trim()).map((d) => d.trim()),
          similarities: similarities
            .filter((s) => s.trim())
            .map((s) => s.trim()),
          sentences: sentences.trim(),
          totalDifferences: differences.filter((d) => d.trim()).length,
          totalSimilarities: similarities.filter((s) => s.trim()).length,
        };
        onContentChange(submissionContent);
      } else {
        onContentChange(null);
      }
    }
  }, [
    word1,
    word2,
    differences,
    similarities,
    sentences,
    onContentChange,
    isReview,
  ]);

  const updateDifference = (index, value) => {
    const newDifferences = [...differences];
    newDifferences[index] = value;
    setDifferences(newDifferences);
  };

  const updateSimilarity = (index, value) => {
    const newSimilarities = [...similarities];
    newSimilarities[index] = value;
    setSimilarities(newSimilarities);
  };

  const isFormComplete = () => {
    return (
      word1.trim() &&
      word2.trim() &&
      differences.filter((d) => d.trim()).length >= 3 &&
      similarities.filter((s) => s.trim()).length >= 3 &&
      sentences.trim().length > 0
    );
  };

  return (
    <div
      className={`${styles.taskTypeContainer} ${styles.unexpectedConnectionsTask}`}
    >
      <p className={styles.example}>
        {prompt ||
          "Возьмите два случайных предмета и найдите между ними неожиданные связи."}
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

      <div className={styles.connectionsForm}>
        <div className={`${styles.section} ${styles.wordsSection}`}>
          <h3>Два слова</h3>
          <p className={styles.sectionHint}>
            Выберите два случайных слова из словаря или книги
          </p>

          <div className={styles.wordsGrid}>
            <div className="form-group">
              <label htmlFor="word1">Первое слово:</label>
              <input
                id="word1"
                type="text"
                className="form-group__input"
                placeholder="Например: автомат"
                value={word1}
                onChange={(e) => setWord1(e.target.value)}
                disabled={isReview}
              />
            </div>

            <div className="form-group">
              <label htmlFor="word2">Второе слово:</label>
              <input
                id="word2"
                type="text"
                className="form-group__input"
                placeholder="Например: морковь"
                value={word2}
                onChange={(e) => setWord2(e.target.value)}
                disabled={isReview}
              />
            </div>
          </div>
        </div>

        {/* Differences Section */}
        <div className={`${styles.section} ${styles.differencesSection}`}>
          <h3>Отличия (3)</h3>
          <p className={styles.sectionHint}>
            Найдите и опишите 3 отличия между этими предметами
          </p>

          {differences.map((diff, index) => (
            <div key={index} className="form-group">
              <label htmlFor={`difference-${index}`}>
                Отличие {index + 1}:
              </label>
              <input
                id={`difference-${index}`}
                type="text"
                className="form-group__input"
                placeholder={`Например: ${index === 0 ? "размеры" : index === 1 ? "способ применения" : "морковь можно вырастить, а автомат нельзя"}`}
                value={diff}
                onChange={(e) => updateDifference(index, e.target.value)}
                disabled={isReview}
              />
            </div>
          ))}
        </div>

        {/* Similarities Section */}
        <div className={`${styles.section} ${styles.similaritiesSection}`}>
          <h3>Сходства (3)</h3>
          <p className={styles.sectionHint}>
            Найдите и опишите 3 сходства между этими предметами
          </p>

          {similarities.map((sim, index) => (
            <div key={index} className="form-group">
              <label htmlFor={`similarity-${index}`}>
                Сходство {index + 1}:
              </label>
              <input
                id={`similarity-${index}`}
                type="text"
                className="form-group__input"
                placeholder={`Например: ${index === 0 ? "можно держать в руках" : index === 1 ? "ствол напоминает морковь" : "у АК-47 рукоятка оранжевого цвета"}`}
                value={sim}
                onChange={(e) => updateSimilarity(index, e.target.value)}
                disabled={isReview}
              />
            </div>
          ))}
        </div>

        {/* Sentences Section */}
        <div className={`${styles.section} ${styles.sentencesSection}`}>
          <h3>Предложения с обоими словами</h3>
          <p className={styles.sectionHint}>
            Составьте несколько предложений, используя оба слова
          </p>

          <div className="form-group">
            <label htmlFor="sentences">Ваши предложения:</label>
            <textarea
              id="sentences"
              className="form-group__textarea"
              rows={5}
              placeholder="Например: Солдат начал стрелять из автомата по банке, но вдруг понял, что стреляет в зайца, бегающего среди грядок моркови."
              value={sentences}
              onChange={(e) => setSentences(e.target.value)}
              disabled={isReview}
            />
            <div className={styles.characterCount}>
              {sentences.length} символов
            </div>
          </div>
        </div>

        {/* Completion Indicator */}
        {!isFormComplete() &&
          (word1 ||
            word2 ||
            differences.some((d) => d) ||
            similarities.some((s) => s) ||
            sentences) &&
          !isReview && (
            <div className={styles.validationHints}>
              <p className={styles.hintTitle}>
                Для завершения задания необходимо:
              </p>
              <ul>
                {!word1.trim() && <li>Ввести первое слово</li>}
                {!word2.trim() && <li>Ввести второе слово</li>}
                {differences.filter((d) => d.trim()).length < 3 && (
                  <li>
                    Заполнить все 3 отличия (
                    {differences.filter((d) => d.trim()).length}/3)
                  </li>
                )}
                {similarities.filter((s) => s.trim()).length < 3 && (
                  <li>
                    Заполнить все 3 сходства (
                    {similarities.filter((s) => s.trim()).length}/3)
                  </li>
                )}
                {!sentences.trim() && (
                  <li>Написать предложения с обоими словами</li>
                )}
              </ul>
            </div>
          )}

        {isFormComplete() && !isReview && (
          <div className={styles.completionIndicator}>
            ✓ Задание заполнено и готово к отправке
          </div>
        )}
      </div>
    </div>
  );
};

export default UnexpectedConnectionsTask;
