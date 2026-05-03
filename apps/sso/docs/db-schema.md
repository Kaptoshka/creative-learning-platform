# Database Schema Documentation

## Overview

This schema describes a service for managing users, groups, their relationships, permissions and jwt tokens.

Main entities:

- **Users** — stores user account information
- **Groups** — logical collections of users
- **Roles** — defines roles that can be assigned to users (RBAC)
- **Permissions** — represents fine-grained access rights to system resources
- **Apps** — registered services or client applications within the system
- **Signing keys** — cryptographic keys used for signing and verifying JWTs (with support for key rotation)
- **Refresh tokens** — stores issued refresh tokens used to obtain new access tokens

---

## ER Diagram

```mermaid
erDiagram
    SIGNING_KEYS {
        UUID id PK
        VARCHAR algorithm
        TEXT public_key
        TEXT private_key
        BOOLEAN is_active
        TIMESTAMP created_at
        TIMESTAMP rotated_at
        %% INDEX(is_active) WHERE is_active = true
    }

    USERS {
        UUID id PK
        VARCHAR email
        VARCHAR first_name
        VARCHAR last_name
        VARCHAR middle_name
        VARCHAR pass_hash
        TIMESTAMP created_at
        TIMESTAMP updated_at
        %% UNIQUE(email)
    }

    GROUPS {
        UUID id PK
        VARCHAR name
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }

    USER_GROUPS {
        UUID user_id FK
        UUID group_id FK
        %% PK(user_id, group_id)
        %% INDEX(user_id)
    }

    ROLES {
        UUID id PK
        VARCHAR role
        %% UNIQUE(role)
    }

    PERMISSIONS {
        UUID id PK
        VARCHAR slug
        VARCHAR description
        VARCHAR resource_group
        TIMESTAMP created_at
        %% UNIQUE(slug)
    }

    USER_ROLES {
        UUID user_id FK
        UUID role_id FK
        %% PK(user_id, role_id)
    }

    ROLE_PERMISSIONS {
        UUID role_id FK
        UUID permission_id FK
        %% PK(role_id, permission_id)
    }

    APPS {
        UUID id PK
        VARCHAR name
        VARCHAR secret
        VARCHAR description
        BOOLEAN is_active
        TIMESTAMP created_at
        TIMESTAMP updated_at
        %% UNIQUE(name)
        %% UNIQUE(secret)
    }

    REFRESH_TOKENS {
        UUID id PK
        UUID user_id FK
        VARCHAR token_hash
        TIMESTAMP expires_at
        TIMESTAMP revoked_at
        TIMESTAMP created_at
        %% UNIQUE(token_hash)
        %% INDEX(user_id)
        %% INDEX(expires_at) WHERE revoked_at IS NULL
    }

    %% Relationships
    USERS ||--o{ USER_GROUPS : "user_id"
    GROUPS ||--o{ USER_GROUPS : "group_id"

    USERS ||--o{ USER_ROLES : "user_id"
    ROLES ||--o{ USER_ROLES : "role_id"

    ROLES ||--o{ ROLE_PERMISSIONS : "role_id"
    PERMISSIONS ||--o{ ROLE_PERMISSIONS : "permission_id"

    USERS ||--o{ REFRESH_TOKENS : "user_id"
```
