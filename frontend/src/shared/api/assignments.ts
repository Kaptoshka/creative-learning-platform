import { api } from "./client";
import type {
    CreateAssignmentRequest,
    GetAssignmentResponse,
    GetStudentSubmissionResponse,
    ListAssignmentSubmissionsResponse,
    ListAssignmentsResponse,
    ListMyAssignmentsResponse,
    PaginationParams,
    ProvideFeedbackRequest,
    SaveDraftRequest,
    SaveDraftResponse,
    StartAssignmentResponse,
    SubmitAssignmentRequest,
    SubmitAssignmentResponse,
    UpdateAssignmentRequest,
} from "./types";

export const assignmentsApi = {
    // ─── Teacher: template management ────────────────────────────────────────

    list(params?: PaginationParams & { creator_id?: string }) {
        return api
            .get<ListAssignmentsResponse>("/api/v1/assignments", { params })
            .then((res) => res.data);
    },

    get(id: string) {
        return api
            .get<GetAssignmentResponse>(`/api/v1/assignments/${id}`)
            .then((res) => res.data);
    },

    create(data: CreateAssignmentRequest) {
        return api
            .post<{ ID: string }>("/api/v1/assignments", data)
            .then((res) => res.data);
    },

    update(id: string, data: Omit<UpdateAssignmentRequest, "AssignmentID">) {
        return api
            .patch<GetAssignmentResponse>(`/api/v1/assignments/${id}`, {
                AssignmentID: id,
                ...data,
            })
            .then((res) => res.data);
    },

    delete(id: string) {
        return api.delete(`/api/v1/assignments/${id}`);
    },

    // ─── Teacher: submissions & feedback ─────────────────────────────────────

    listSubmissions(
        templateId: string,
        params?: PaginationParams & { status_filter?: number },
    ) {
        return api
            .get<ListAssignmentSubmissionsResponse>(
                `/api/v1/assignments/${templateId}/submissions`,
                { params },
            )
            .then((res) => res.data);
    },

    getSubmission(submissionId: string) {
        return api
            .get<GetStudentSubmissionResponse>(
                `/api/v1/submissions/${submissionId}`,
            )
            .then((res) => res.data);
    },

    provideFeedback(submissionId: string, data: ProvideFeedbackRequest) {
        return api
            .post(`/api/v1/submissions/${submissionId}/feedback`, data)
            .then((res) => res.data);
    },

    // ─── Student: assignment workflow ─────────────────────────────────────────

    listMine(params?: PaginationParams & { status_filter?: number }) {
        return api
            .get<ListMyAssignmentsResponse>("/api/v1/my/assignments", {
                params,
            })
            .then((res) => res.data);
    },

    start(templateId: string) {
        return api
            .post<StartAssignmentResponse>(
                `/api/v1/my/assignments/${templateId}/start`,
            )
            .then((res) => res.data);
    },

    saveDraft(submissionId: string, data: SaveDraftRequest) {
        return api
            .put<SaveDraftResponse>(
                `/api/v1/my/submissions/${submissionId}/draft`,
                data,
            )
            .then((res) => res.data);
    },

    submit(submissionId: string, data: SubmitAssignmentRequest) {
        return api
            .post<SubmitAssignmentResponse>(
                `/api/v1/my/submissions/${submissionId}/submit`,
                data,
            )
            .then((res) => res.data);
    },
};
