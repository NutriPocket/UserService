// Package model contains the structs types that will be used in the application.
package model

// BaseUser is a struct that contains the base user data received from the client
type BaseUser struct {
	Username string
	Email    string
	Password string
}

type EditableUser struct {
	Picture string `json:"picture"`
}

// User is a struct that contains the user data that will be sent to the client
type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	EditableUser
}

// LoginUser is a struct that contains the user data received from the client when logging in
type LoginUser struct {
	EmailOrUsername string
	Password        string
}

// SavedUser is a struct that combines BaseUser fields with an additional Id field
type SavedUser struct {
	BaseUser
	Id string
}
