import { useState, useCallback } from "react";
import { useFetch } from "./useFetch";
import apiClient from "@/services/apiClient";

interface Task {
  ID: number;
  id: number;
  title: string;
  description: string;
  content: Record<string, unknown>;
  deadline_time: string;
  status: string;
}

interface UseTaskOptions {
  taskId?: string;
}

interface UseTaskReturn {
  task: Task | null;
  loading: boolean;
  error: string | null;
  submit: (content: string) => Promise<boolean>;
  refetch: () => Promise<void>;
}

export const useTask = (options: UseTaskOptions = {}): UseTaskReturn => {
  const { taskId } = options;

  const url = taskId ? `/assignments/${taskId}` : null;
  const { data, loading, error, refetch } = useFetch<Task>(url, { immediate: !!taskId });

  const [isSubmitting, setIsSubmitting] = useState(false);

  const submit = useCallback(
    async (content: string): Promise<boolean> => {
      if (!taskId) return false;

      setIsSubmitting(true);

      try {
        const payload = {
          assignment_id: Number(taskId),
          content,
        };

        const response = await apiClient.post("/submissions", payload);
        return response.status === 201;
      } catch (err) {
        console.error("Submit error:", err);
        return false;
      } finally {
        setIsSubmitting(false);
      }
    },
    [taskId]
  );

  return {
    task: data as Task | null,
    loading,
    error,
    submit,
    refetch,
  };
};
