// Package service contains the services that will be used in the application.
package service

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/NutriPocket/UserService/model"
	"github.com/NutriPocket/UserService/repository"
)

type UserService struct {
	repository repository.IUserRepository
}

func NewUserService(userRepository repository.IUserRepository) (*UserService, error) {
	var err error

	if userRepository == nil {
		userRepository, err = repository.NewUserRepository(nil)
		if err != nil {
			log.Errorf("Failed to create user repository: %v", err)
			return nil, err
		}
	}

	return &UserService{repository: userRepository}, nil
}

func (service *UserService) EncodePassword(password string) string {
	hashPasswordBytes := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hashPasswordBytes[:])

}

func (service *UserService) CreateUser(userData *model.BaseUser) (model.User, error) {
	userData.Password = service.EncodePassword(userData.Password)

	return service.repository.CreateUser(userData)
}

func (service *UserService) Login(userData *model.LoginUser) (model.User, error) {
	savedUser, err := service.repository.GetUserWithPassword(userData.EmailOrUsername)

	if err != nil {
		return model.User{}, err
	}

	if savedUser == (model.SavedUser{}) {
		return model.User{}, &model.AuthenticationError{
			Title:  "Credentials don't match",
			Detail: "User identification or password are wrong, please try again",
		}
	}

	userData.Password = service.EncodePassword(userData.Password)

	if userData.Password != savedUser.Password {
		return model.User{}, &model.AuthenticationError{
			Title:  "Credentials don't match",
			Detail: "User identification or password are wrong, please try again",
		}
	}

	return model.User{Id: savedUser.Id, Username: savedUser.Username, Email: savedUser.Email}, nil
}

func (service *UserService) GetAllUsers(params model.GetUsersParams) ([]model.User, error) {
	return service.repository.GetAllUsers(params)
}

func (service *UserService) GetUser(username string) (model.User, error) {
	user, err := service.repository.GetUser(username)

	if err != nil {
		return user, err
	}

	if user == (model.User{}) {
		return user, &model.NotFoundError{Title: "User not found", Detail: "The user with the username " + username + " was not found"}
	}

	return user, nil
}

func (service *UserService) UpdateUser(username string, userData *model.EditableUser) (model.User, error) {
	user, err := service.repository.GetUser(username)

	if err != nil {
		return user, err
	}

	if user == (model.User{}) {
		return user, &model.NotFoundError{Title: "User not found", Detail: "The user with the username " + username + " was not found"}
	}

	if userData.Picture != "" {
		user.EditableUser.Picture = userData.Picture
	}

	return service.repository.UpdateUser(user.Id, userData)
}