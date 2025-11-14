package auth_test

import (
	"go/adv-demo/internal/auth"
	"go/adv-demo/internal/user"
	"testing"
)

type MockUserRepository struct{}

func (repo *MockUserRepository) Create(dto *user.User) (*user.User, error) {
	return &user.User{
		Email: dto.Email,
		Name:  dto.Name,
	}, nil
}

func (repo *MockUserRepository) FindByEmail(email string) (*user.User, error) {
	return &user.User{
		Email: email,
		Name:  "test user",
	}, nil
}

func TestRegisterSuccess(t *testing.T) {
	AuthService := auth.NewUserService(&MockUserRepository{})
	email, err := AuthService.Register("1@1.com", "", "")

	if err != nil {
		t.Fatal(err)
	}

	if email != "1@1.com" {
		t.Fatal("Email fatal")
	}
}
