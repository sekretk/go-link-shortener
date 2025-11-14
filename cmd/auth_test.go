package main

import (
	"bytes"
	"encoding/json"
	"go/adv-demo/internal/auth"
	"go/adv-demo/internal/user"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDb() *gorm.DB {
	err := godotenv.Load(".env")

	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	return db
}

func initData(db *gorm.DB) {
	db.Create(&user.User{
		Email:    "test@test.test",
		Password: "$2a$10$y/iUc4X.oxfV4yjLFRNuk.hdliyVlu0ERDgKXvzY7h.37GzbJNv0u",
		Name:     "Test user",
	})
}

func removeData(db *gorm.DB) {
	db.Unscoped().
		Where("email = ?", "test@test.test").
		Delete(&user.User{})

}

func TestLoginSuccess(t *testing.T) {

	db := initDb()

	initData(db)

	ts := httptest.NewServer(App())

	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "test@test.test",
		Password: "password",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("Expected %d got %d", 200, res.StatusCode)
	}

	removeData(db)
}

func TestLoginFailedForInvalidCreds(t *testing.T) {

	db := initDb()

	initData(db)

	defer removeData(db)

	ts := httptest.NewServer(App())

	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "some@example.com",
		Password: "password",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 401 {
		t.Fatalf("Expected %d got %d", 401, res.StatusCode)
	}

}

func TestLoginReturnToken(t *testing.T) {

	db := initDb()

	initData(db)

	defer removeData(db)

	ts := httptest.NewServer(App())

	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "test@test.test",
		Password: "password",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))

	if err != nil {
		t.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		t.Fatal("Error on get body")
	}

	var response auth.LoginResponse

	err = json.Unmarshal(body, &response)

	if err != nil {
		t.Fatal(err)
	}

	if response.Token == "" {
		t.Fatalf("Expected Token")
	}
}
