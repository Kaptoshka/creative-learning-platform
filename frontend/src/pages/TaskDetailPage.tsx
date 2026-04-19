import React, { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import apiClient from "@/services/apiClient";
import Button from "@/components/Button";

const TaskDetailPage = () => {
  const { id } = useParams();
  const navigate = useNavigate();

  const [task, setTask] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const [submissionContent, setSubmissionContent] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    const fetchTask = async () => {
      if (!id) return;
      try {
        setLoading(true);
        const response = await apiClient.get(`/assignments/${id}`);
        setTask(response.data);
      } catch (err) {
        setError(
          "Failed to load the task. It may not exist or you may not have permission.",
        );
        console.error("Task fetch failed:", err);
      } finally {
        setLoading(false);
      }
    };
    fetchTask();
  }, [id]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSubmitting(true);
    setError("");

    try {
      const payload = {
        assignment_id: Number(id),
        content: submissionContent,
      };

      const response = await apiClient.post("/submissions", payload);

      if (response.status === 201) {
        navigate("/tasks", { viewTransition: true });
      }
    } catch (err) {
      if (err.response && err.response.data && err.response.data.error) {
        setError(err.response.data.error);
      } else {
        setError("An unexpected error occurred while submitting.");
      }
      console.error("Submission failed:", err);
    } finally {
      setIsSubmitting(false);
    }
  };

  if (loading) {
    return <div className="page-container">Loading task...</div>;
  }

  if (error || !task) {
    return (
      <div className="page-container error-message">
        {error || "Task not found."}
      </div>
    );
  }

  return (
    <div className="page-container task-detail-page">
      <header className="task-detail-header">
        <h1>{task.title}</h1>
        <p>{task.description}</p>
      </header>

      <div className="submission-form-container">
        <h3>Your Submission</h3>
        <form onSubmit={handleSubmit} noValidate>
          {error && <p className="error-message">{error}</p>}
          <div className="form-group">
            <label htmlFor="submissionContent">Enter your work below:</label>
            <textarea
              id="submissionContent"
              value={submissionContent}
              onChange={(e) => setSubmissionContent(e.target.value)}
              rows="10"
              placeholder="Start writing your creative response here..."
              required
              disabled={isSubmitting}
            />
          </div>
          <Button type="submit" variant="primary" disabled={isSubmitting}>
            {isSubmitting ? "Submitting..." : "Submit Task"}
          </Button>
        </form>
      </div>
    </div>
  );
};

export default TaskDetailPage;
