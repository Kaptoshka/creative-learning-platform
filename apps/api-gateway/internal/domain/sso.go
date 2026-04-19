package domain

type RegisterRequest struct {
	Email      string
	Password   string
	FirstName  string
	LastName   string
	MiddleName string
}

type RegisterResponse struct {
	UserID int64
}

type LoginRequest struct {
	Email    string
	Password string
	AppID    int32
}

type LoginResponse struct {
	Token string
}

type LogoutRequest struct {
	Token string
}

type LogoutResponse struct {
	Success bool
}
