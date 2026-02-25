package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	pw1 := "correctPW123"
	pw2 := "wrongPW123"
	hash1, _ := HashPassword(pw1)
	hash2, _ := HashPassword(pw2)

	testCases := []struct {
		name    string
		pw      string
		hash    string
		matchPw bool
		wantErr bool
	}{
		{name: "Correct",
			pw:      pw1,
			hash:    hash1,
			matchPw: true,
			wantErr: false,
		},
		{name: "InCorect",
			pw:      "WrongPWH1",
			hash:    hash1,
			matchPw: false,
			wantErr: false,
		},
		{
			name:    "PW don't match hash",
			pw:      pw1,
			hash:    hash2,
			matchPw: false,
			wantErr: false,
		},
		{
			name:    "Invalid hash",
			pw:      pw1,
			hash:    "Wronghas",
			matchPw: false,
			wantErr: true,
		},
		{
			name:    "Empty pw",
			pw:      "",
			hash:    "invalidhash",
			matchPw: false,
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		matchP, err := CheckPasswordHash(tc.pw, tc.hash)
		if (err != nil) != tc.wantErr {
			t.Errorf("CheckPasswordHash error :%v wantErr %v", err, tc.wantErr)
		}
		if !tc.wantErr && (matchP != tc.matchPw) {
			t.Errorf("CheckPasswordHash expected %v , got %v", tc.matchPw, matchP)
		}
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secret12345"
	expiresIn := 1 * time.Hour
	tokenString, _ := MakeJWT(userID, tokenSecret, expiresIn)

	testCases := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Correct",
			tokenString: tokenString,
			tokenSecret: tokenSecret,
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Incorrect tokenString",
			tokenString: "testing",
			tokenSecret: tokenSecret,
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Incorrect tokensecret",
			tokenString: tokenString,
			tokenSecret: "wrong token secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}
	for _, tc := range testCases {
		userID, err := ValidateJWT(tc.tokenString, tc.tokenSecret)
		if (err != nil) != tc.wantErr {
			t.Errorf("ValidateJWT error :%v wantErr %v", err, tc.wantErr)
		}
		if !tc.wantErr && (userID != tc.wantUserID) {
			t.Errorf("ValidateJWT expected userId: %v , got userId: %v", tc.wantUserID, userID)
		}
	}
}
