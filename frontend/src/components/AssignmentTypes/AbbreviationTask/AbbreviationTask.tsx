import React, { useState, useEffect } from "react";
import styles from "./AbbreviationTask.module.scss";
import AbbreviationItem from "./AbbreviationItem";
import TaskDescription from "@/components/TaskDescription";

const AbbreviationTask = ({ onContentChange, submission, isReview }) => {
  const getInitialState = () => {
    if (isReview && submission && submission.abbreviations) {
      return submission.abbreviations;
    }
    return [
      { word: "", expansion: "" },
      { word: "", expansion: "" },
      { word: "", expansion: "" },
      { word: "", expansion: "" },
    ];
  };

  const [abbreviations, setAbbreviations] = useState(getInitialState);

  // Update parent component whenever content changes
  useEffect(() => {
    if (typeof onContentChange === "function" && !isReview) {
      // Create JSON content with student's answer
      const submissionContent = {
        type: "Аббревиатуры",
        abbreviations: abbreviations
          .filter((item) => item.word.trim() && item.expansion.trim())
          .map((item) => ({
            word: item.word.trim(),
            expansion: item.expansion.trim(),
          })),
        totalItems: abbreviations.filter(
          (item) => item.word.trim() && item.expansion.trim(),
        ).length,
      };

      // Only send if there's at least some content
      if (submissionContent.totalItems > 0) {
        onContentChange(submissionContent);
      } else {
        onContentChange(null);
      }
    }
  }, [abbreviations, onContentChange, isReview]);

  const updateAbbreviation = (index, field, value) => {
    const newAbbreviations = [...abbreviations];
    newAbbreviations[index] = {
      ...newAbbreviations[index],
      [field]: value,
    };
    setAbbreviations(newAbbreviations);
  };

  const isValidWord = (word) => {
    return word.length >= 3 && word.length <= 5 && /^[А-ЯЁа-яё]+$/.test(word);
  };

  const isValidExpansion = (expansion, word) => {
    if (!expansion.trim() || !word.trim()) return true;

    const expansionWords = expansion.trim().split(/\s+/);
    const wordLetters = word.trim().toUpperCase().split("");

    if (expansionWords.length !== wordLetters.length) return false;

    return expansionWords.every(
      (expWord, index) =>
        expWord.charAt(0).toUpperCase() === wordLetters[index],
    );
  };

  const abbreviationInstructions =
    "Напишите несколько коротких слов по 3–5 букв. Затем к каждому слову придумайте расшифровку, как будто это аббревиатура.";
  const abbreviationExample =
    "РОТ – Рыбный Омлет с Тмином\nГОРА – Грустный Омбудсмен, Работающий Автомехаником";

  return (
    <div className={`${styles.taskTypeContainer} ${styles.abbreviationTask}`}>
      <TaskDescription
        instructions={abbreviationInstructions}
        example={abbreviationExample}
      />

      {abbreviations.map((item, index) => (
        <AbbreviationItem
          key={index}
          index={index}
          item={item}
          updateAbbreviation={updateAbbreviation}
          isValidWord={isValidWord}
          isValidExpansion={isValidExpansion}
          isReview={isReview}
        />
      ))}
    </div>
  );
};

export default AbbreviationTask;
