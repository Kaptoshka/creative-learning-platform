import React from "react";
import styles from "./ReviewForm.module.scss";

interface ReviewFormProps {
  review: {
    selectedCategories: string[];
    strengths: string;
    improvements: string;
    nextSteps: string;
    encouragement: string;
  };
  onChange: (field: string, value: unknown) => void;
  disabled?: boolean;
}

const categories = [
  { id: "creativity", label: "Creative Approach", icon: "🎨" },
  { id: "originality", label: "Originality", icon: "💡" },
  { id: "technique", label: "Technical Execution", icon: "✨" },
  { id: "composition", label: "Composition/Structure", icon: "🎯" },
  { id: "expression", label: "Artistic Expression", icon: "🌟" },
  { id: "improvement", label: "Growth Potential", icon: "🚀" },
];

const ReviewForm = ({ review, onChange, disabled }: ReviewFormProps) => {
  const toggleCategory = (categoryId: string) => {
    const current = review.selectedCategories;
    const updated = current.includes(categoryId)
      ? current.filter((id) => id !== categoryId)
      : [...current, categoryId];
    onChange("selectedCategories", updated);
  };

  return (
    <div className={styles.reviewForm}>
      <div className={styles.categorySection}>
        <h3 className={styles.sectionTitle}>
          Что было сделано хорошо?
        </h3>
        <div className={styles.categoryGrid}>
          {categories.map((cat) => (
            <button
              key={cat.id}
              type="button"
              className={`${styles.categoryChip} ${
                review.selectedCategories.includes(cat.id) ? styles.categoryChipActive : ""
              }`}
              onClick={() => toggleCategory(cat.id)}
              disabled={disabled}
            >
              <span className={styles.categoryIcon}>{cat.icon}</span>
              <span>{cat.label}</span>
            </button>
          ))}
        </div>
      </div>

      <div className={styles.textSection}>
        <div className={styles.textField}>
          <label>Сильные стороны работы</label>
          <textarea
            value={review.strengths}
            onChange={(e) => onChange("strengths", e.target.value)}
            placeholder="Опишите, что было сделано хорошо..."
            disabled={disabled}
          />
        </div>

        <div className={styles.textField}>
          <label>Что можно улучшить</label>
          <textarea
            value={review.improvements}
            onChange={(e) => onChange("improvements", e.target.value)}
            placeholder="Предложите области для роста..."
            disabled={disabled}
          />
        </div>

        <div className={styles.textField}>
          <label>Следующие шаги</label>
          <textarea
            value={review.nextSteps}
            onChange={(e) => onChange("nextSteps", e.target.value)}
            placeholder="Рекомендации для дальнейшего развития..."
            disabled={disabled}
          />
        </div>

        <div className={styles.textField}>
          <label>Поощрение</label>
          <textarea
            value={review.encouragement}
            onChange={(e) => onChange("encouragement", e.target.value)}
            placeholder="Вдохновляющее сообщение студенту..."
            disabled={disabled}
          />
        </div>
      </div>
    </div>
  );
};

export default ReviewForm;
