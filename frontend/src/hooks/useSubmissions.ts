import { useState, useEffect, useCallback } from "react";
import apiClient from "@/services/apiClient";

interface SubmissionData {
  id: number;
  taskTitle: string;
  taskType: string;
  studentName: string;
  studentEmail: string;
  studentAvatar: string;
  submittedAt: Date;
  daysWaiting: number;
  priority: "new" | "normal" | "urgent";
  hasContent: boolean;
  content: Record<string, unknown>;
  preview: string;
  status: string;
}

interface UseSubmissionsReturn {
  submissions: SubmissionData[];
  loading: boolean;
  error: string | null;
  refetch: () => Promise<void>;
}

export const useSubmissions = (userId: number | null): UseSubmissionsReturn => {
  const [submissions, setSubmissions] = useState<SubmissionData[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchSubmissions = useCallback(async () => {
    if (!userId) {
      setLoading(false);
      return;
    }

    try {
      setLoading(true);
      const response = await apiClient.get(`/submissions/teacher/${userId}`);
      const data = response.data || [];

      const detailed = await Promise.all(
        data.map(async (sub: Record<string, unknown>) => {
          const [assignRes, studentRes] = await Promise.all([
            apiClient.get(`/assignments/${sub.AssignmentID}`),
            apiClient.get(`/users/${sub.StudentID}`),
          ]);

          const assignment = assignRes.data;
          const student = Array.isArray(studentRes.data)
            ? studentRes.data[0]
            : studentRes.data;

          const submittedAt = new Date(sub.SubmittedAt as string);
          const daysWaiting = Math.floor(
            (new Date().getTime() - submittedAt.getTime()) / (1000 * 60 * 60 * 24)
          );

          let priority: "new" | "normal" | "urgent" = "normal";
          if (daysWaiting >= 3) priority = "urgent";
          else if (daysWaiting === 0) priority = "new";

          const content = assignment.content || {};

          return {
            id: sub.ID,
            taskTitle: content.title || "Untitled",
            taskType: content.type || "unknown",
            studentName: `${student?.first_name} ${student?.last_name}`,
            studentEmail: student?.email || "Unknown",
            studentAvatar: student?.first_name?.[0] || "?",
            submittedAt,
            daysWaiting,
            priority,
            hasContent: !!sub.Content,
            content: sub.Content as Record<string, unknown>,
            preview:
              (sub.Content as Record<string, unknown>)?.text ||
              (sub.Content as Record<string, unknown>)?.story ||
              (sub.Content as Record<string, unknown>)?.sentence ||
              "",
            status: sub.Status || "pending",
          };
        })
      );

      setSubmissions(detailed);
      setError(null);
    } catch (err) {
      setError("Failed to load submissions");
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, [userId]);

  useEffect(() => {
    fetchSubmissions();
  }, [fetchSubmissions]);

  return { submissions, loading, error, refetch: fetchSubmissions };
};
