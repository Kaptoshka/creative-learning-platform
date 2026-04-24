# 📚 Database Schema Documentation

## Overview

This scheme describes a service for managing assignments, their templates, submissions, solution versions, and feedback.

Main entities:

- **Widgets** — types of assignments (e.g., forms, tests, etc.)
- **Assignment Templates** — assignment templates
- **Assignment Targets** — whom the assignment is assigned to
- **Submissions** — attempts to complete the assignment
- **Submission Versions** — versions of solutions (including autosaves)
- **Feedbacks** — feedback from reviewers

---

## 📊 ER Diagram

```mermaid
erDiagram
    WIDGETS {
        UUID id PK
        VARCHAR type
        INTEGER version
        JSONB config_schema
        JSONB submission_schema
        %% UNIQUE(type, version)
    }

    ASSIGNMENT_TEMPLATES {
        UUID id PK
        UUID creator_id
        VARCHAR title
        TEXT description
        UUID widget_id FK
        JSONB widget_config
        TIMESTAMP due_date
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }

    ASSIGNMENT_TARGETS {
        UUID id PK
        UUID template_id FK
        UUID group_id
        UUID student_id
        TIMESTAMP created_at
        TIMESTAMP updated_at
        %% CHECK(group_id IS NOT NULL OR student_id IS NOT NULL)
    }

    SUBMISSIONS {
        UUID id PK
        UUID template_id FK
        UUID student_id
        VARCHAR status
        TIMESTAMP started_at
        TIMESTAMP submitted_at
        %% UNIQUE(student_id, template_id)
        %% INDEX(template_id, status)
        %% INDEX(student_id)
    }

    SUBMISSION_VERSIONS {
        UUID id PK
        UUID submission_id FK
        INTEGER version_number
        JSONB payload
        INTEGER time_spent_seconds
        BOOLEAN is_autosave
        TIMESTAMP created_at
        TIMESTAMP updated_at
        %% INDEX(submission_id, version_number DESC)
    }

    FEEDBACKS {
        UUID id PK
        UUID version_id FK
        UUID grader_id
        TEXT text_content
        JSONB payload
        BOOLEAN is_published
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }

    %% Relationships
    WIDGETS ||--o{ ASSIGNMENT_TEMPLATES : "widget_id"
    ASSIGNMENT_TEMPLATES ||--o{ ASSIGNMENT_TARGETS : "template_id"
    ASSIGNMENT_TEMPLATES ||--o{ SUBMISSIONS : "template_id"
    SUBMISSIONS ||--o{ SUBMISSION_VERSIONS : "submission_id"
    SUBMISSION_VERSIONS ||--o{ FEEDBACKS : "version_id"
```
