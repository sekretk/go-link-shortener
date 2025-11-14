package auth_test

import (
	"bytes"
	"encoding/json"
	"go/adv-demo/configs"
	"go/adv-demo/internal/auth"
	"go/adv-demo/internal/user"
	"go/adv-demo/pkg/db"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func bootstrap() (*auth.AuthHandler, sqlmock.Sqlmock, error) {
	database, mock, err := sqlmock.New()

	if err != nil {
		return nil, nil, err
	}

	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: database,
	}))

	if err != nil {
		return nil, nil, err
	}

	userRepo := user.NewUserRepository(&db.Db{
		DB: gormDb,
	})

	handler := auth.AuthHandler{
		Config: &configs.Config{
			Auth: configs.AuthConfig{
				Secret: "sekret",
			},
		},
		AuthService: auth.NewUserService(userRepo),
	}

	return &handler, mock, nil
}

func TestLoginSuccess(t *testing.T) {
	handler, mock, err := bootstrap()

	rows := sqlmock.NewRows([]string{"email", "password"}).AddRow("test@test.test", "$2a$10$y/iUc4X.oxfV4yjLFRNuk.hdliyVlu0ERDgKXvzY7h.37GzbJNv0u")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	if err != nil {
		t.Fatal(err)
		return
	}

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "test@test.test",
		Password: "password",
	})

	reader := bytes.NewReader(data)

	w := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/auth/login", reader)

	handler.Login()(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("got %d, expected %d", w.Code, http.StatusOK)
	}
}

func TestRegisterSucc(t *testing.T) {
	handler, mock, err := bootstrap()
	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"email", "password", "name"}))
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	if err != nil {
		t.Fatal(err)
		return
	}

	data, _ := json.Marshal(&auth.RegisterRequest{
		Email:    "test@test.test",
		Password: "password",
		Name:     "Test",
	})

	reader := bytes.NewReader(data)

	w := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/auth/register", reader)

	handler.Register()(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("got %d, expected %d", w.Code, http.StatusOK)
	}
}
