package common

import "testing"

func Test_Jwt(t *testing.T) {
	username := "tiptok"
	password := "password"
	token, e := GenerateToken(username, password)
	if e != nil {
		t.Fatal(e)
	}
	user, e := ParseJWTToken(token)
	if e != nil {
		t.Fatal(e)
	}
	if user.Username != username || user.Password != password {
		t.Fatal("parse jwt token error.")
	}
}
