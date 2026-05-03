package permissions

const (
	RoleAdmin   = "admin"
	RoleStudent = "student"
	RoleTeacher = "teacher"
)

const (
	ScopeAssignmentsRead   = "assignments:read"
	ScopeAssignmentsWrite  = "assignments:write"
	ScopeAssignmentsDelete = "assignments:delete"
)

const (
	ScopeSubmissionsCreate = "submissions:create"
	ScopeSubmissionsRead   = "submissions:read"
)

const (
	ScopeFeedbackRead  = "feedback:read"
	ScopeFeedbackWrite = "feedback:write"
)

const (
	ScopeWidgetsRead  = "widgets:read"
	ScopeWidgetsWrite = "widgets:write"
)

func GetScopesForRole(role string) []string {
	switch role {
	case RoleAdmin:
		return []string{
			ScopeAssignmentsRead, ScopeAssignmentsWrite, ScopeAssignmentsDelete,
			ScopeSubmissionsCreate, ScopeSubmissionsRead, ScopeFeedbackRead,
			ScopeFeedbackWrite, ScopeWidgetsRead, ScopeWidgetsWrite,
		}
	case RoleTeacher:
		return []string{
			ScopeAssignmentsRead, ScopeAssignmentsWrite, ScopeAssignmentsDelete,
			ScopeSubmissionsRead, ScopeFeedbackRead, ScopeFeedbackWrite,
			ScopeWidgetsRead,
		}
	case RoleStudent:
		return []string{
			ScopeAssignmentsRead, ScopeSubmissionsCreate, ScopeSubmissionsRead,
			ScopeFeedbackRead, ScopeWidgetsRead,
		}
	default:
		return nil
	}
}
