package jwt_test

import (
	"go/adv-demo/pkg/jwt"
	"testing"
)

func TestJWTCreate(t *testing.T) {

	jwtService := jwt.NewJWT("123123123123123123123123")

	token, err := jwtService.Create(jwt.JWTData{
		Email: "test@test.test",
	})

	if err != nil {
		t.Fatal(err)
	}

	isValid, data := jwtService.Parse(token)

	if !isValid {
		t.Fatal("Token int valid")
	}

	if data.Email != "test@test.test" {
		t.Fatal("Emails not matched")
	}

}
