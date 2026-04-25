package models

type User struct {
	ID         int
	Email      string
	PassHash   []byte
	FirstName  string
	LastName   string
	MiddleName string
}
