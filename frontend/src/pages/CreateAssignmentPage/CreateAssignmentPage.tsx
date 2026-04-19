import React from "react";
import { useNavigate } from "react-router-dom";
import DatePicker, { registerLocale } from "react-datepicker";
import { ru } from "date-fns/locale/ru";

import { useCreateAssignment } from "@/hooks/useCreateAssignment";
import Button from "@/components/Button";
import Loading from "@/components/ui/Loading";
import ErrorMessage from "@/components/ui/ErrorMessage";

import styles from "./CreateAssignmentPage.module.scss";

registerLocale("ru", ru);

const CreateAssignmentPage = () => {
  const navigate = useNavigate();
  const {
    title,
    description,
    deadline,
    selectedTemplate,
    selectedStudent,
    studentQuery,
    searchResults,
    showStudentSearch,
    error,
    success,
    isLoading,
    templates,
    setTitle,
    setDescription,
    setDeadline,
    setSelectedTemplate,
    setStudentQuery,
    setSelectedStudent,
    setShowStudentSearch,
    searchStudents,
    createAssignment,
  } = useCreateAssignment();

  const handleSearchChange = (value: string) => {
    setStudentQuery(value);
    if (value.length >= 2) {
      searchStudents(value);
      setShowStudentSearch(true);
    } else {
      setShowStudentSearch(false);
    }
  };

  const selectStudent = (student: { id: number; first_name: string; last_name: string; email: string }) => {
    setSelectedStudent(student);
    setShowStudentSearch(false);
    setStudentQuery(`${student.first_name} ${student.last_name}`);
  };

  if (isLoading) {
    return <Loading text="Создание..." fullPage />;
  }

  return (
    <div className={styles.createTaskPage}>
      <div className={styles.createTaskPageContainer}>
        <h1>Создать новое задание</h1>

        {error && <ErrorMessage error={error} />}
        {success && <div className={styles.successMessage}>{success}</div>}

        <form
          className={styles.createTaskPageForm}
          onSubmit={(e) => {
            e.preventDefault();
            createAssignment();
          }}
        >
          {/* Student Selection */}
          <div className={styles.taskFormSection}>
            <h3>Назначить студенту (опционально)</h3>
            <div className={styles.formGroup}>
              <input
                type="text"
                placeholder="Введите имя или email студента..."
                value={studentQuery}
                onChange={(e) => handleSearchChange(e.target.value)}
                className={styles.formInput}
              />
              {showStudentSearch && searchResults.length > 0 && (
                <div className={styles.searchResults}>
                  {searchResults.slice(0, 5).map((student) => (
                    <div
                      key={student.id}
                      className={styles.searchResultItem}
                      onClick={() => selectStudent(student)}
                    >
                      {student.first_name} {student.last_name} - {student.email}
                    </div>
                  ))}
                </div>
              )}
              {selectedStudent && (
                <button
                  type="button"
                  className={styles.clearButton}
                  onClick={() => {
                    setSelectedStudent(null);
                    setStudentQuery("");
                  }}
                >
                  ×
                </button>
              )}
            </div>
          </div>

          {/* Basic Info */}
          <div className={styles.taskFormSection}>
            <h3>Основная информация</h3>
            <div className={styles.formGroup}>
              <input
                type="text"
                placeholder="Например: Креативные аббревиатуры"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                className={styles.formInput}
              />
              <textarea
                placeholder="Опишите цель задания и что студент должен делать..."
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                className={styles.formTextarea}
              />
            </div>
          </div>

          {/* Template Selection */}
          <div className={`${styles.taskFormSection} ${styles.taskFormSectionTemplates}`}>
            <h3>Выберите шаблон задания</h3>
            <div className={styles.templatesGrid}>
              {templates.map((template) => (
                <div
                  key={template.id}
                  className={`${styles.templateCard} ${
                    selectedTemplate?.id === template.id ? styles.templateCardSelected : ""
                  }`}
                  onClick={() => setSelectedTemplate(template)}
                >
                  <h4>{template.name}</h4>
                  <p>{template.description}</p>
                </div>
              ))}
            </div>
          </div>

          {/* Deadline */}
          <div className={styles.taskFormSection}>
            <h3>Срок сдачи (опционально)</h3>
            <DatePicker
              selected={deadline}
              onChange={(date) => setDeadline(date as Date)}
              showTimeSelect
              dateFormat="Pp"
              locale="ru"
              minDate={new Date()}
              placeholderText="Выберите дату и время"
              className={styles.formInput}
            />
          </div>

          <div className={styles.formActions}>
            <Button type="submit" variant="primary" disabled={isLoading}>
              {isLoading ? "Создание..." : "Создать задание"}
            </Button>
            <Button variant="outline" onClick={() => navigate("/tasks")}>
              Отмена
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default CreateAssignmentPage;