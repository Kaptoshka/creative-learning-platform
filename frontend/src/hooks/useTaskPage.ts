import { useState, useEffect, useCallback } from "react";
import { useParams, useNavigate } from "react-router-dom";
import apiClient from "@/services/apiClient";
import offlineQueue from "@/services/offlineQueue";

interface TaskData {
  ID: number;
  id: number;
  title: string;
  description: string;
  content: Record<string, unknown>;
  type: string;
}

interface UseTaskPageReturn {
  task: TaskData | null;
  loading: boolean;
  error: string | null;
  submissionContent: unknown;
  isSubmitting: boolean;
  submitSuccess: boolean;
  successMessage: string;
  setSubmissionContent: (content: unknown) => void;
  handleSubmit: (e?: React.FormEvent) => Promise<void>;
  isFormValid: () => boolean;
}

export const useTaskPage = (): UseTaskPageReturn => {
  const { id } = useParams();
  const navigate = useNavigate();

  const [task, setTask] = useState<TaskData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [submissionContent, setSubmissionContent] = useState<unknown>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitSuccess, setSubmitSuccess] = useState(false);
  const [successMessage, setSuccessMessage] = useState("");

  useEffect(() => {
    setLoading(true);
    apiClient
      .get(`/assignments/${id}`)
      .then((res) => {
        setTask(res.data);
      })
      .catch((err) => {
        setError("Failed to load task");
        console.error(err);
      })
      .finally(() => {
        setLoading(false);
      });
  }, [id]);

  const handleSubmit = useCallback(
    async (e?: React.FormEvent) => {
      if (e) e.preventDefault();

      if (!submissionContent) {
        setError("Пожалуйста, заполните задание перед отправкой.");
        return;
      }

      setIsSubmitting(true);
      setError("");
      setSubmitSuccess(false);

      try {
        const contentToSubmit =
          typeof submissionContent === "string"
            ? submissionContent
            : JSON.stringify(submissionContent);

        const payload = {
          assignment_id: Number(id),
          content: contentToSubmit,
        };

        const response = await offlineQueue.send("/submissions", payload);

        if (response.status === 201 || response.status === "offline") {
          const message =
            response.status === "offline"
              ? "Интернета нет. Задание сохранено и будет отправлено позже."
              : "Задание успешно отправлено!";
          setSuccessMessage(message);
          setSubmitSuccess(true);

          setTimeout(() => {
            setIsSubmitting(false);
            navigate("/dashboard", { viewTransition: true });
          }, 2000);
        }
      } catch (err) {
        if (err.response?.data?.error) {
          setError(err.response.data.error);
        } else {
          setError("An unexpected error occurred during submission.");
        }
        console.error("Submission failed:", err);
        setIsSubmitting(false);
      }
    },
    [submissionContent, id, navigate]
  );

  const isFormValid = useCallback(() => {
    if (!submissionContent) return false;

    if (typeof submissionContent === "object") {
      return JSON.stringify(submissionContent).length > 10;
    }

    return (submissionContent as string).trim().length > 0;
  }, [submissionContent]);

  return {
    task,
    loading,
    error,
    submissionContent,
    isSubmitting,
    submitSuccess,
    successMessage,
    setSubmissionContent,
    handleSubmit,
    isFormValid,
  };
};
