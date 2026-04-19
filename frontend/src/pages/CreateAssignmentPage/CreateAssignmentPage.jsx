import React, { useState, useContext } from "react";
import offlineQueue from "@/services/offlineQueue";
import apiClient from "@/services/apiClient";
import { AuthContext } from "@/context/AuthContext";
import Button from "@/components/Button";
import { useNavigate } from "react-router-dom";
import "react-datepicker/dist/react-datepicker.css";
import DatePicker, { registerLocale } from "react-datepicker";
import { ru } from "date-fns/locale/ru";

import styles from "./CreateAssignmentPage.module.scss";

const CreateAssignmentPage = () => {
  registerLocale("ru", ru);

  const { user } = useContext(AuthContext);
  const [studentQuery, setStudentQuery] = useState("");
  const [selectedStudent, setSelectedStudent] = useState(null);
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [deadline, setDeadline] = useState(null);
  const [selectedTemplate, setSelectedTemplate] = useState(null);

  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [showStudentSearch, setShowStudentSearch] = useState(false);
  const [searchResults, setSearchResults] = useState([]);

  const navigate = useNavigate();

  const creativeTemplates = [
    {
      id: 1,
      name: "Аббревиатуры",
      description: "Создание креативных расшифровок для коротких слов",
      instructions: [
        "Выберите 4-6 слов длиной 3-5 букв",
        "Придумайте оригинальные расшифровки",
        "Расшифровки должны быть смешными или неожиданными",
      ],
      example: "РОТ → Рыбный Омлет с Тмином",
      content:
        "Напишите несколько коротких слов по 3–5 букв. Затем к каждому слову придумайте расшифровку, как будто это аббревиатура.",
    },
    {
      id: 2,
      name: "Неожиданные связи",
      description: "Поиск сходств и различий между случайными предметами",
      instructions: [
        "Выберите два случайных предмета",
        "Найдите 3 различия между ними",
        "Найдите 3 неожиданных сходства",
        "Составьте одно творческое предложение, связывающее их",
      ],
      example:
        "Автомат и морковь: оба могут быть оранжевыми, требуют ухода, имеют свою 'ammunition'",
      content:
        "Возьмите два случайных предмета и найдите между ними неожиданные связи.",
    },
    {
      id: 3,
      name: "Нестандартные применения",
      description: "Поиск необычных способов использования обычных предметов",
      instructions: [
        "Выберите один обычный предмет",
        "Придумайте минимум 5 необычных способов его применения",
        "Способы могут быть практичными или фантастическими",
        "Опишите каждый способ в 1-2 предложениях",
      ],
      example:
        "Карандаш: палочка для размешивания краски, закладка для книги, указка на презентации",
      content: "Придумайте необычные способы применения для обычного предмета.",
    },
    {
      id: 4,
      name: "Рассказ из 100 слов",
      description: "Написание истории строго из 100 слов",
      instructions: [
        "Выберите одну из предложенных тем или придумайте свою",
        "Напишите законченную историю",
        "Используйте ровно 100 слов - ни больше, ни меньше",
        "История должна иметь начало, развитие и концовку",
      ],
      example:
        "Темы: путешествие во времени, последний человек на Земле, говорящие животные",
      content:
        "Составьте рассказ из точно 100 слов. Ни одного слова больше, ни одного слова меньше.",
    },
    {
      id: 5,
      name: "Два в одном",
      description: "Объединение двух предметов в один",
      instructions: [
        "Выберите 3 пары различных предметов",
        "Объедините каждую пару в один новый предмет",
        "Опишите получившийся предмет",
        "Объясните его применение и преимущества",
      ],
      example:
        "Зонт + Фонарик = Светящийся зонт для безопасных прогулок в темноте",
      content:
        "Возьмите два предмета и подумайте, как их можно объединить в один полезный или забавный предмет.",
    },
    {
      id: 6,
      name: "Одна буква",
      description: "Составление предложения на одну букву",
      type: "Одна буква",
      instructions: [
        "Выберите букву русского алфавита",
        "Составьте предложение из минимум 5 слов",
        "Все слова должны начинаться на выбранную букву",
        "Смысл предложения не так важен",
      ],
      example: "Разительного роста растение росло рядом с рощей",
      prompt:
        "Составьте предложение, где все слова будут начинаться на одну букву.",
      title: "Одна буква",
    },
  ];

  const handleStudentSearch = async (query) => {
    setStudentQuery(query);
    setError("");

    if (query.length < 2) {
      setShowStudentSearch(false);
      setSearchResults([]);
      return;
    }

    if (!navigator.onLine) {
      return;
    }

    try {
      const response = await apiClient.get(
        `/users/search?role=student&query=${query}`,
      );
      setSearchResults(response.data);
      setShowStudentSearch(true);
    } catch (err) {
      console.error("failed to search students", err);
      if (!navigator.onLine || err.message === "Network Error") {
        setError("Нет интернета. Поиск студентов недоступен.");
      } else {
        setError("Не удалось загрузить список студентов.");
      }
      setSearchResults([]);
    }
  };

  const selectStudent = (student) => {
    setSelectedStudent(student);
    setStudentQuery(
      `${student.last_name} ${student.first_name} (${student.email})`,
    );
    setShowStudentSearch(false);
  };

  const selectTemplate = (template) => {
    setSelectedTemplate(template);
    if (!title) setTitle(template.name);
    if (!description) setDescription(template.description);
  };

  const handleSubmit = async () => {
    setIsLoading(true);
    setError("");
    setSuccess("");

    if (!selectedStudent || !title || !deadline || !selectedTemplate) {
      setError("Пожалуйста, заполните все обязательные поля");
      setIsLoading(false);
      return;
    }

    try {
      console.log("deadline", deadline);
      const payload = {
        teacher_id: user.id,
        student_id: selectedStudent.id,
        deadline: deadline ? deadline.toISOString() : null,
        status: "available",
        content: JSON.stringify({
          title: title.trim(),
          description: description.trim(),
          type: selectedTemplate.name,
          prompt: selectedTemplate.content,
          instructions: selectedTemplate.instructions,
          example: selectedTemplate.example,
        }),
      };

      const response = await offlineQueue.send("/assignments", payload);

      if (response.status === 201 || response.status === "offline") {
        const message =
          response.status === "offline"
            ? "Интернета нет. Задание сохранено и будет отправлено автоматически при появлении сети."
            : "Задание успешно создано!";

        setSuccess(message);

        setTimeout(() => {
          setStudentQuery("");
          setSelectedStudent(null);
          setTitle("");
          setDescription("");
          setDeadline("");
          setSelectedTemplate(null);
          setSuccess("");
        }, 3000);
      }
    } catch (err) {
      setError(
        err.response?.data?.message ||
          "Произошла ошибка при создании задания. Попробуйте еще раз.",
      );
      console.error("Task creation failed:", err);
    } finally {
      setIsLoading(false);
    }
  };

  const getInputClassName = (value, isRequired = false) => {
    let className = "form-group__input";
    if (isRequired && !value) {
      className += " form-group__input--invalid";
    } else if (value) {
      className += " form-group__input--valid";
    }
    return className;
  };

  return (
    <div className={styles.createTaskPage}>
      <div className={styles.createTaskPageContainer}>
        <div className={styles.createTaskPageHeader}>
          <h1 className={styles.createTaskPageTitle}>Создать новое задание</h1>
          <p className={styles.createTaskPageSubtitle}>
            Выберите студента и шаблон для создания творческого задания.
          </p>
        </div>

        <div className={styles.createTaskPageContent}>
          <div className={styles.createTaskPageFormSection}>
            <div className={styles.taskForm}>
              {error && (
                <div className="status-message error-message">{error}</div>
              )}

              {success && (
                <div className="status-message success-message">{success}</div>
              )}

              {/* Student Selection Section */}
              <div className={`${styles.taskFormSection}`}>
                <h3 className={`${styles.taskFormSectionTitle}`}>
                  Выбор студента
                </h3>

                <div className="form-group">
                  <label className="form-group__label" htmlFor="studentSearch">
                    Поиск студента *
                  </label>
                  <div className={styles.studentSearch}>
                    <input
                      id="studentSearch"
                      type="text"
                      className={getInputClassName(selectedStudent, true)}
                      value={studentQuery}
                      onChange={(e) => handleStudentSearch(e.target.value)}
                      disabled={isLoading}
                      placeholder="Введите имя или email студента..."
                    />

                    {showStudentSearch && searchResults.length > 0 && (
                      <div className={styles.studentDropdown}>
                        {searchResults.map((student) => (
                          <div
                            key={student.id}
                            className={styles.studentOption}
                            onClick={() => selectStudent(student)}
                          >
                            <div className={styles.studentOptionName}>
                              {student.last_name} {student.first_name}
                            </div>
                            <div className={styles.studentOptionEmail}>
                              {student.email}
                            </div>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                </div>
              </div>

              {/* Basic Information Section */}
              <div className={styles.taskFormSection}>
                <h3 className={styles.taskFormSectionTitle}>
                  Информация о задании
                </h3>

                <div className="form-group">
                  <label className="form-group__label" htmlFor="title">
                    Название задания *
                  </label>
                  <input
                    id="title"
                    type="text"
                    className={getInputClassName(title, true)}
                    value={title}
                    onChange={(e) => setTitle(e.target.value)}
                    disabled={isLoading}
                    placeholder="Например: Креативные аббревиатуры"
                    maxLength="100"
                  />
                </div>

                <div className="form-group">
                  <label className="form-group__label" htmlFor="description">
                    Краткое описание
                  </label>
                  <textarea
                    id="description"
                    className="form-group__textarea"
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                    rows="3"
                    disabled={isLoading}
                    placeholder="Опишите цель задания и что студент должен делать..."
                    maxLength="300"
                  />
                  <div className="form-group__helper">
                    {description.length}/300 символов
                  </div>
                </div>

                <div className="form-group">
                  <label className="form-group__label" htmlFor="deadline">
                    Крайний срок выполнения *
                  </label>
                  <DatePicker
                    id="deadline"
                    selected={deadline}
                    onChange={(date) => setDeadline(date)}
                    showTimeSelect
                    timeFormat="HH:mm"
                    timeIntervals={15}
                    dateFormat="d MMMM yyyy, HH:mm"
                    timeCaption="Время"
                    locale="ru"
                    minDate={new Date()}
                    placeholderText="Выберите дату и время"
                    className={getInputClassName(deadline, true)}
                    autoComplete="off"
                    filterTime={(date) => {
                      const isToday =
                        new Date().toDateString() === date.toDateString();
                      return isToday ? date > new Date() : true;
                    }}
                  />
                </div>
              </div>

              {/* Action Buttons */}
              <div className={styles.taskFormActions}>
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => navigate("/", { viewTransition: true })}
                  disabled={isLoading}
                >
                  ❌ Отмена
                </Button>
                <Button
                  type="button"
                  variant="primary"
                  disabled={
                    isLoading ||
                    !selectedStudent ||
                    !title ||
                    !deadline ||
                    !selectedTemplate
                  }
                  onClick={handleSubmit}
                  isLoading={isLoading}
                  loadingText="Создание..."
                >
                  🎯 Создать задание
                </Button>
              </div>
            </div>
          </div>

          {/* Templates Section */}
          <div className={styles.createTaskPageTemplatesSection}>
            <div className={styles.templatesPanel}>
              <h3 className={styles.templatesPanelTitle}>
                🎨 Творческие шаблоны
              </h3>
              <p className={styles.templatesPanelSubtitle}>
                Выберите один из готовых шаблонов для создания задания
              </p>

              <div className={styles.templatesList}>
                {creativeTemplates.map((template) => (
                  <div
                    key={template.id}
                    className={`${styles.templateCard} ${selectedTemplate?.id === template.id ? styles.templateCardSelected : ""}`}
                    onClick={() => selectTemplate(template)}
                  >
                    <div className={styles.templateCardHeader}>
                      <div>
                        <h4 className={styles.templateCardTitle}>
                          {template.name}
                        </h4>
                        <p className={styles.templateCardDescription}>
                          {template.description}
                        </p>
                      </div>
                      {selectedTemplate?.id === template.id && (
                        <div className={styles.templateCardSelected}>✅</div>
                      )}
                    </div>

                    <div className={styles.templateCardContent}>
                      <div className={styles.templateCardExample}>
                        <strong>Пример:</strong> {template.example}
                      </div>

                      <div className={styles.templateCardInstructions}>
                        <strong>Инструкции:</strong>
                        <ul>
                          {template.instructions.map((instruction, index) => (
                            <li key={index}>{instruction}</li>
                          ))}
                        </ul>
                      </div>
                    </div>
                  </div>
                ))}
              </div>

              <div className={styles.templatesPanelFooter}>
                <p className={styles.templatesPanelHint}>
                  💡 <strong>Совет:</strong> Выберите шаблон, который лучше
                  всего подходит для развития творческих способностей вашего
                  студента
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CreateAssignmentPage;
