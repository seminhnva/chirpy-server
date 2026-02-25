package auth

import "testing"

func TestCheckPasswordHas(t *testing.T) {
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
			wantErr: false,
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
