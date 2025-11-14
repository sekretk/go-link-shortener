package auth

import (
	"errors"
	"fmt"
	"go/adv-demo/internal/user"
	"go/adv-demo/pkg/di"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository di.IUserRepository
}

func NewUserService(userRepository di.IUserRepository) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
	}
}

func (service *AuthService) Register(email, password, name string) (string, error) {
	existedUser, _ := service.UserRepository.FindByEmail(email)

	if existedUser != nil {
		return "", errors.New(ErrUserExists)
	}

	hanshedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //TODO: add salt

	if err != nil {
		return "", err
	}

	user := &user.User{
		Email:    email,
		Password: string(hanshedPwd),
		Name:     name,
	}

	_, err = service.UserRepository.Create(user)

	if err != nil {
		fmt.Printf("Error")
		return "", err
	}

	return user.Email, nil
}

func (service *AuthService) Login(email, password string) (string, error) {
	fmt.Println(email)
	existedUser, err := service.UserRepository.FindByEmail(email)

	if err != nil {
		return "", errors.New(ErrWrongCredentials)
	}

	err = bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(password))

	if err != nil {
		return "", errors.New(ErrWrongCredentials)
	}

	return existedUser.Email, nil
}
