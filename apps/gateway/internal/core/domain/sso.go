package domain

type RegisterRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
}

type RegisterResponse struct {
	UserID string `json:"user_id"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	AppID    string `json:"app_id"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type LogoutResponse struct{}

type LogoutAllRequest struct {
	UserID string `json:"user_id"`
}

type LogoutAllResponse struct{}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type DeactivateAppRequest struct {
	AppID string `json:"app_id"`
}

type DeactivateAppResponse struct{}

type RegisterAppRequest struct {
	Name        string `json:"name"`
	Secret      string `json:"secret"`
	Description string `json:"description"`
}

type RegisterAppResponse struct {
	AppID string `json:"app_id"`
}
