import { useState, useCallback } from "react";
import styles from "./ArchitectTask.module.scss";

// ─── Constants ────────────────────────────────────────────────
const EXAMPLE_NOUNS = [
  "киви",
  "зима",
  "бронза",
  "маяк",
  "джунгли",
  "пианино",
  "вулкан",
  "зеркало",
  "шёлк",
  "янтарь",
];

const CARD_BORDER_COLORS = [
  "#9333ea",
  "#16a34a",
  "#d97706",
  "#dc2626",
  "#2563eb",
  "#0d9488",
  "#db2777",
  "#9333ea",
  "#16a34a",
  "#d97706",
];

// ─── Component ────────────────────────────────────────────────
export default function ArchitectTask() {
  const [nouns, setNouns] = useState(Array(10).fill(""));
  const [ideas, setIdeas] = useState([]);
  const [userProject, setUserProject] = useState("");
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [submitStatus, setSubmitStatus] =
    (useState < "idle") | "success" | ("error" > "idle");

  const filledCount = nouns.filter((n) => n.trim().length > 0).length;
  const progress = (filledCount / 10) * 100;

  const handleNounChange = useCallback((index, value) => {
    setNouns((prev) => {
      const next = [...prev];
      next[index] = value;
      return next;
    });
  }, []);

  const fillExamples = () => {
    setNouns([...EXAMPLE_NOUNS]);
  };

  const clearAll = () => {
    setNouns(Array(10).fill(""));
    setIdeas([]);
    setSubmitStatus("idle");
  };

  const generateIdeas = async () => {
    const words = nouns.filter((n) => n.trim().length > 0);
    if (words.length < 3) return;

    setLoading(true);
    setIdeas([]);

    const prompt = `Ты — творческий архитектор. Пользователь написал слова как условия заказчика загородного дома. Для каждого слова придумай одну конкретную архитектурную или дизайнерскую идею (интерьер, экстерьер, участок, материалы). Идеи должны быть нестандартными и интересными.

Ответь ТОЛЬКО валидным JSON-массивом без markdown-блоков, пояснений и лишнего текста:
[{"word":"слово","idea":"конкретная идея для дома"}]

Слова: ${words.join(", ")}`;

    try {
      const response = await fetch("https://api.anthropic.com/v1/messages", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          model: "claude-sonnet-4-20250514",
          max_tokens: 1000,
          messages: [{ role: "user", content: prompt }],
        }),
      });

      const data = await response.json();
      const text = data.content
        .map((b) => (b.type === "text" ? b.text : ""))
        .join("");
      const clean = text.replace(/```json|```/g, "").trim();
      const parsed = JSON.parse(clean);
      setIdeas(parsed);
    } catch {
      setIdeas([]);
    } finally {
      setLoading(false);
    }
  };

  const submitProject = async () => {
    if (!userProject.trim()) return;
    setSubmitting(true);
    setSubmitStatus("idle");

    // Here you'd call your backend / save the project.
    // For demo: simulate a network call.
    await new Promise((r) => setTimeout(r, 800));

    // In a real app, handle success/error from API.
    setSubmitStatus("success");
    setSubmitting(false);
  };

  return (
    <div className={styles.taskPage}>
      {/* ─── Header ─────────────────────────────────────── */}
      <div className={styles.header}>
        <div className={styles.header__badge}>
          <span>🏛</span> Метод аналогий
        </div>
        <h1 className={styles.header__title}>Архитектор по методу аналогий</h1>
        <p className={styles.header__description}>
          Вы — архитектор. Напишите 10 любых существительных, которые приходят в
          голову. Представьте их как 10 обязательных условий заказчика
          загородного дома — и спроектируйте дом, воплощая каждое слово
          буквально или метафорически.
        </p>
      </div>

      {/* ─── Step 1: Nouns ───────────────────────────────── */}
      <div className={styles.phase}>
        <div className={styles.stepHeader}>
          <span className={styles.stepBadge}>1</span>
          <span className={styles.stepTitle}>Напишите 10 существительных</span>
        </div>
        <p className={styles.stepHint}>
          Любые слова, которые приходят в голову прямо сейчас. Не думайте —
          просто пишите. Нужно минимум 3 слова для генерации идей.
        </p>

        {/* Progress */}
        <div className={styles.progressBar}>
          <div
            className={styles.progressBar__fill}
            style={{ width: `${progress}%` }}
          />
        </div>

        {/* Grid */}
        <div className={styles.nounsGrid}>
          {nouns.map((val, i) => (
            <div key={i} className={styles.nounSlot}>
              <span className={styles.nounSlot__index}>{i + 1}.</span>
              <input
                type="text"
                className={`${styles.nounSlot__input} ${val.trim() ? styles["nounSlot__input--filled"] : ""}`}
                placeholder="слово..."
                value={val}
                maxLength={20}
                onChange={(e) => handleNounChange(i, e.target.value)}
              />
            </div>
          ))}
        </div>

        {/* Example/clear */}
        <div className={styles.exampleRow} style={{ marginTop: "0.75rem" }}>
          <button className={styles.exampleTag} onClick={fillExamples}>
            Случайные примеры →
          </button>
          {filledCount > 0 && (
            <button className={styles.exampleTag} onClick={clearAll}>
              Очистить всё
            </button>
          )}
        </div>
      </div>

      <hr className={styles.divider} />

      {/* ─── Step 2: Generate ────────────────────────────── */}
      <div className={styles.phase}>
        <div className={styles.stepHeader}>
          <span className={styles.stepBadge}>2</span>
          <span className={styles.stepTitle}>
            Сгенерировать архитектурные идеи
          </span>
        </div>
        <p className={styles.stepHint}>
          Каждое слово — условие заказчика. ИИ предложит, как воплотить его в
          проекте дома. Используйте эти идеи как отправную точку для вашей
          фантазии.
        </p>

        <div className={styles.btnRow}>
          <button
            className={`${styles.btn} ${styles["btn--primary"]}`}
            onClick={generateIdeas}
            disabled={loading || filledCount < 3}
          >
            {loading ? (
              <>
                <span className={styles.spinnerInline} />
                Архитектор думает...
              </>
            ) : (
              "✦ Создать проект дома"
            )}
          </button>
        </div>

        <div className={styles.projectArea}>
          {!loading && ideas.length === 0 && (
            <p className={styles.emptyHint}>
              Идеи появятся здесь после генерации
            </p>
          )}
          {loading && (
            <p className={styles.emptyHint}>
              <span
                className={styles.spinnerInline}
                style={{ marginRight: "0.5rem" }}
              />
              Генерируем идеи...
            </p>
          )}
          {!loading && ideas.length > 0 && (
            <div className={styles.roomList}>
              {ideas.map((item, i) => (
                <div
                  key={i}
                  className={styles.roomCard}
                  style={{
                    borderLeftColor:
                      CARD_BORDER_COLORS[i % CARD_BORDER_COLORS.length],
                    animationDelay: `${i * 0.05}s`,
                  }}
                >
                  <div className={styles.roomCard__tag}>{item.word}</div>
                  <div className={styles.roomCard__idea}>{item.idea}</div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      <hr className={styles.divider} />

      {/* ─── Step 3: User project ────────────────────────── */}
      <div className={styles.phase}>
        <div className={styles.stepHeader}>
          <span className={styles.stepBadge}>3</span>
          <span className={styles.stepTitle}>Ваш проект дома</span>
        </div>
        <p className={styles.stepHint}>
          Используйте сгенерированные идеи как вдохновение. Включите воображение
          и опишите свой уникальный дом своими словами — как он выглядит, какие
          пространства в нём есть, какая атмосфера.
        </p>

        <div className={styles.userProjectArea}>
          <textarea
            value={userProject}
            onChange={(e) => setUserProject(e.target.value)}
            placeholder="Мой дом будет стоять на холме... В центре гостиной — живое дерево... Стены спальни покрыты..."
          />
        </div>

        <div className={styles.btnRow} style={{ marginTop: "1rem" }}>
          <button
            className={`${styles.btn} ${styles["btn--success"]}`}
            onClick={submitProject}
            disabled={submitting || userProject.trim().length === 0}
          >
            {submitting ? (
              <>
                <span className={styles.spinnerInline} />
                Сохранение...
              </>
            ) : (
              "Сохранить проект"
            )}
          </button>
        </div>

        {submitStatus === "success" && (
          <div
            className={`${styles.statusMessage} ${styles["statusMessage--success"]}`}
          >
            Проект сохранён! Отличная работа — ваш дом звучит неповторимо.
          </div>
        )}
        {submitStatus === "error" && (
          <div
            className={`${styles.statusMessage} ${styles["statusMessage--error"]}`}
          >
            Не удалось сохранить. Попробуйте ещё раз.
          </div>
        )}
      </div>

      {/* ─── Footer ─────────────────────────────────────── */}
      <p className={styles.footerNote}>
        Метод аналогий — техника латерального мышления Эдварда де Боно.
        Случайные слова помогают выйти за рамки привычных решений и найти
        неожиданные, но жизнеспособные идеи.
      </p>
    </div>
  );
}
