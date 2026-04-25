package driving

import (
	"context"
)

type AuthService interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
		firstName string,
		lastName string,
		middleName string,
	) (userID int64, err error)
}
