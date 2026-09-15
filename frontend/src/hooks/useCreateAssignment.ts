import { useState, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import apiClient from "@/services/apiClient";
import offlineQueue from "@/services/offlineQueue";
import { useAuth } from "./useAuth";
import { creativeTemplates } from "@/data/templates";
import type { CreativeTemplate } from "@/data/templates";

interface Student {
    id: number;
    first_name: string;
    last_name: string;
    email: string;
}

interface UseCreateAssignmentReturn {
    // Form state
    title: string;
    description: string;
    deadline: Date | null;
    selectedTemplate: CreativeTemplate | null;
    selectedStudent: Student | null;

    // Search state
    studentQuery: string;
    searchResults: Student[];
    showStudentSearch: boolean;

    // UI state
    error: string;
    success: string;
    isLoading: boolean;

    // Templates
    templates: CreativeTemplate[];

    // Actions
    setTitle: (value: string) => void;
    setDescription: (value: string) => void;
    setDeadline: (value: Date | null) => void;
    setSelectedTemplate: (template: CreativeTemplate | null) => void;
    setStudentQuery: (value: string) => void;
    setSelectedStudent: (student: Student | null) => void;
    setShowStudentSearch: (show: boolean) => void;
    searchStudents: (query: string) => Promise<void>;
    createAssignment: () => Promise<boolean>;
    validate: () => boolean;
}

export const useCreateAssignment = (): UseCreateAssignmentReturn => {
    const navigate = useNavigate();
    const { user } = useAuth();

    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    const [deadline, setDeadline] = useState<Date | null>(null);
    const [selectedTemplate, setSelectedTemplate] =
        useState<CreativeTemplate | null>(null);
    const [selectedStudent, setSelectedStudent] = useState<Student | null>(
        null,
    );

    const [studentQuery, setStudentQuery] = useState("");
    const [searchResults, setSearchResults] = useState<Student[]>([]);
    const [showStudentSearch, setShowStudentSearch] = useState(false);

    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");
    const [isLoading, setIsLoading] = useState(false);

    const searchStudents = useCallback(async (query: string) => {
        if (!query.trim()) {
            setSearchResults([]);
            return;
        }

        try {
            const response = await apiClient.get(`/users?search=${query}`);
            setSearchResults(response.data || []);
        } catch (err) {
            console.error("Search failed:", err);
        }
    }, []);

    const validate = useCallback(() => {
        if (!title.trim()) {
            setError("Пожалуйста, введите название задания");
            return false;
        }
        if (!description.trim()) {
            setError("Пожалуйста, введите описание задания");
            return false;
        }
        if (!selectedTemplate) {
            setError("Пожалуйста, выберите шаблон");
            return false;
        }
        return true;
    }, [title, description, selectedTemplate]);

    const createAssignment = useCallback(async (): Promise<boolean> => {
        if (!validate()) return false;

        setIsLoading(true);
        setError("");
        setSuccess("");

        try {
            const payload = {
                title,
                description,
                deadline: deadline?.toISOString(),
                template_id: selectedTemplate?.id,
                student_id: selectedStudent?.id,
                content: {
                    type: selectedTemplate?.type,
                    content: selectedTemplate?.content,
                },
            };

            const sendFn = async () => apiClient.post("/assignments", payload);

            let response;
            if (navigator.onLine) {
                response = await sendFn();
            } else {
                const queued = await offlineQueue.queue(
                    "/assignments",
                    payload,
                );
                response = { status: queued ? "offline" : 201 };
            }

            if (response.status === 201 || response.status === "offline") {
                const message =
                    response.status === "offline"
                        ? "Интернета нет. Задание сохранено и будет отправлено автоматически."
                        : "Задание успешно создано!";
                setSuccess(message);
                setTimeout(() => navigate("/tasks"), 2000);
                return true;
            }
            return false;
        } catch (err: unknown) {
            const errMsg = err instanceof Error ? err.message : "Unknown error";
            setError(errMsg);
            return false;
        } finally {
            setIsLoading(false);
        }
    }, [
        title,
        description,
        deadline,
        selectedTemplate,
        selectedStudent,
        navigate,
        validate,
    ]);

    return {
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
        templates: creativeTemplates,
        setTitle,
        setDescription,
        setDeadline,
        setSelectedTemplate,
        setStudentQuery,
        setSelectedStudent,
        setShowStudentSearch,
        searchStudents,
        createAssignment,
        validate,
    };
};
