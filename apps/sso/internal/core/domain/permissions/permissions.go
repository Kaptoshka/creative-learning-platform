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

var RoleScopes = map[string][]string{
	RoleAdmin: {
		ScopeAssignmentsRead,
		ScopeAssignmentsWrite,
		ScopeAssignmentsDelete,
		ScopeSubmissionsCreate,
		ScopeSubmissionsRead,
		ScopeFeedbackRead,
		ScopeFeedbackWrite,
		ScopeWidgetsRead,
		ScopeWidgetsWrite,
	},
	RoleTeacher: {
		ScopeAssignmentsRead,
		ScopeAssignmentsWrite,
		ScopeAssignmentsDelete,
		ScopeSubmissionsRead,
		ScopeFeedbackRead,
		ScopeFeedbackWrite,
		ScopeWidgetsRead,
	},
	RoleStudent: {
		ScopeAssignmentsRead,
		ScopeSubmissionsCreate,
		ScopeSubmissionsRead,
		ScopeFeedbackRead,
		ScopeWidgetsRead,
	},
}
