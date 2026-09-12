// --- Enums ---

export enum SubmissionStatus {
    Unspecified = 0,
    NotStarted = 1,
    InProgress = 2,
    Submitted = 3,
    Graded = 4,
    Returned = 5,
}

// --- SSO ---

export interface RegisterRequest {
    email: string;
    password: string;
    first_name: string;
    last_name: string;
    middle_name?: string;
}

export interface RegisterResponse {
    user_id: string;
}

export interface LoginRequest {
    email: string;
    password: string;
    app_id: string;
}

export interface LoginResponse {
    access_token: string;
    refresh_token: string;
}

export interface RefreshRequest {
    refresh_token: string;
}

export interface RefreshResponse {
    access_token: string;
    refresh_token: string;
}

export interface LogoutRequest {
    refresh_token: string;
}

// --- Assignments: Entities ---

export interface AssignmentTemplateLight {
    id: string;
    title: string;
    WidgetType: string;
    due_date?: string;
}

export interface AssignmentTemplate {
    id: string;
    CreatorID: string;
    title: string;
    description: string;
    widget_id: string;
    widget_config: Record<string, unknown>;
    due_date?: string;
    created_at: string;
    updated_at: string;
}

export interface AssignmentTarget {
    group_id?: string;
    student_id?: string;
}

export interface SubmissionVersion {
    id: string;
    version_number: number;
    payload: Record<string, unknown>;
    time_spent: string;
    is_autosave: boolean;
    created_at: string;
    updated_at: string;
}

export interface Submission {
    id: string;
    template_id: string;
    student_id: string;
    status: SubmissionStatus;
    started_at?: string;
    submitted_at?: string;
    latest_version?: SubmissionVersion;
}

export interface Feedback {
    id: string;
    version_id: string;
    grader_id: string;
    text_content: string;
    payload: Record<string, unknown>;
    is_published: boolean;
    created_at: string;
    updated_at: string;
}

export interface MyAssignmentItem {
    template: AssignmentTemplateLight;
    status: SubmissionStatus;
    has_feedback: boolean;
}

// --- Assignments: Requests ---

export interface CreateAssignmentRequest {
    title: string;
    description: string;
    widget_id: string;
    widget_config: Record<string, unknown>;
    due_date?: string;
    targets: AssignmentTarget[];
}

export interface UpdateAssignmentRequest {
    assignment_id: string;
    template: Partial<AssignmentTemplate>;
    update_mask: string[];
    targets: AssignmentTarget[];
}

export interface ProvideFeedbackRequest {
    version_id: string;
    text_content: string;
    payload?: Record<string, unknown>;
    is_published: boolean;
}

export interface SaveDraftRequest {
    payload: Record<string, unknown>;
    time_spent: string; // duration string e.g. "300s"
}

export interface SubmitAssignmentRequest {
    payload: Record<string, unknown>;
    time_spent: string;
}

// --- Assignments: Responses ---

export interface ListAssignmentsResponse {
    items: AssignmentTemplateLight[];
    next_page_token: string;
}

export interface GetAssignmentResponse {
    template: AssignmentTemplate;
    targets: AssignmentTarget[];
}

export interface ListAssignmentSubmissionsResponse {
    items: Submission[];
    next_page_token: string;
}

export interface GetStudentSubmissionResponse {
    template: AssignmentTemplate;
    submission: Submission;
    history: SubmissionVersion[];
    feedback: Feedback[];
}

export interface ListMyAssignmentsResponse {
    items: MyAssignmentItem[];
    next_page_token: string;
}

export interface StartAssignmentResponse {
    submission_id: string;
    started_at: string;
}

export interface SaveDraftResponse {
    version_id: string;
    saved_at: string;
}

export interface SubmitAssignmentResponse {
    version_id: string;
    status: SubmissionStatus;
}

// --- Pagination ---

export interface PaginationParams {
    page_size?: number;
    page_token?: string;
}
