import { useState, useEffect, useCallback } from "react";
import { useLocation, useParams } from "react-router-dom";
import apiClient from "@/services/apiClient";

interface ReviewForm {
  selectedCategories: string[];
  strengths: string;
  improvements: string;
  nextSteps: string;
  encouragement: string;
}

interface Task {
  id: number;
  title: string;
  type: string;
  content: Record<string, unknown>;
}

interface Submission {
  id: number;
  taskTitle: string;
  taskType: string;
  content: Record<string, unknown>;
  studentName: string;
  studentEmail: string;
}

interface UseReviewReturn {
  submission: Submission | null;
  task: Task | null;
  review: ReviewForm;
  loading: boolean;
  error: string | null;
  submitting: boolean;
  showSuccess: boolean;
  setReviewField: (field: keyof ReviewForm, value: unknown) => void;
  submitReview: () => Promise<void>;
}

export const useReview = (): UseReviewReturn => {
  const { submissionId } = useParams();
  const location = useLocation();

  const [submission, setSubmission] = useState<Submission | null>(null);
  const [task, setTask] = useState<Task | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [showSuccess, setShowSuccess] = useState(false);

  const [review, setReview] = useState<ReviewForm>({
    selectedCategories: [],
    strengths: "",
    improvements: "",
    nextSteps: "",
    encouragement: "",
  });

  useEffect(() => {
    const fetchReviewData = async () => {
      if (location.state && (location.state as { submission?: Submission }).submission) {
        const passedSubmission = (location.state as { submission: Submission }).submission;
        setSubmission(passedSubmission);
        setTask({
          id: 0,
          title: passedSubmission.taskTitle,
          type: passedSubmission.taskType,
          content: { type: passedSubmission.taskType },
        });
        setLoading(false);
      } else if (submissionId && submissionId !== "undefined") {
        try {
          setLoading(true);
          const res = await apiClient.get(`/submissions/${submissionId}`);
          setSubmission(res.data.submission);
          setTask(res.data.task);
        } catch (err) {
          setError("Failed to load submission");
          console.error(err);
        } finally {
          setLoading(false);
        }
      } else {
        setLoading(false);
      }
    };

    fetchReviewData();
  }, [submissionId, location.state]);

  const setReviewField = useCallback(
    (field: keyof ReviewForm, value: unknown) => {
      setReview((prev) => ({ ...prev, [field]: value }));
    },
    []
  );

  const submitReview = useCallback(async () => {
    if (!submission) return;

    setSubmitting(true);
    setError(null);

    try {
      const payload = {
        submission_id: submission.id,
        categories: review.selectedCategories,
        strengths: review.strengths,
        improvements: review.improvements,
        next_steps: review.nextSteps,
        encouragement: review.encouragement,
      };

      await apiClient.post(`/reviews`, payload);
      setShowSuccess(true);

      setTimeout(() => {
        setSubmitting(false);
      }, 2000);
    } catch (err) {
      setError("Failed to submit review");
      console.error(err);
      setSubmitting(false);
    }
  }, [submission, review]);

  return {
    submission,
    task,
    review,
    loading,
    error,
    submitting,
    showSuccess,
    setReviewField,
    submitReview,
  };
};
