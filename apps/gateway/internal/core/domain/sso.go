package domain

type RegisterRequest struct {
	Email      string
	Password   string
	FirstName  string
	LastName   string
	MiddleName string
}

type RegisterResponse struct {
	UserID string
}

type LoginRequest struct {
	Email    string
	Password string
	AppID    string
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
}

type LogoutRequest struct {
	RefreshToken string
}

type LogoutResponse struct{}

type LogoutAllRequest struct {
	UserID string
}

type LogoutAllResponse struct{}

type RefreshRequest struct {
	RefreshToken string
}

type RefreshResponse struct {
	AccessToken  string
	RefreshToken string
}

type DeactivateAppRequest struct {
	AppID string
}

type DeactivateAppResponse struct{}

type RegisterAppRequest struct {
	Name        string
	Secret      string
	Description string
}

type RegisterAppResponse struct {
	AppID string
}
