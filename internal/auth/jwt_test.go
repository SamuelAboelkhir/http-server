package auth_test

import (
	"testing"
	"time"

	"github.com/SamuelAboelkhir/http-server/internal/auth"
	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		userID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
		want        string
		wantErr     bool
	}{
		{
			name:        "valid token",
			userID:      uuid.New(),
			tokenSecret: "secret",
			expiresIn:   time.Minute,
			want:        "token",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := auth.MakeJWT(tt.userID, tt.tokenSecret, tt.expiresIn)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("MakeJWT() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("MakeJWT() succeeded unexpectedly")
			}
		})
	}
}

func TestValidateJwt(t *testing.T) {
	name := uuid.New()
	token, _ := auth.MakeJWT(name, "secret", time.Minute)
	token2, _ := auth.MakeJWT(name, "secret", time.Nanosecond)
	token3, _ := auth.MakeJWT(uuid.New(), "secret", time.Minute)

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		tokenString string
		tokenSecret string
		want        uuid.UUID
		wantErr     bool
	}{
		{
			name:        "valid token",
			tokenString: token,
			tokenSecret: "secret",
			want:        name,
			wantErr:     false,
		},
		{
			name:        "Expired",
			tokenString: token2,
			tokenSecret: "secret",
			want:        name,
			wantErr:     true,
		},
		{
			name:        "invalid token",
			tokenString: token3,
			tokenSecret: "secret2",
			want:        name,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.ValidateJwt(tt.tokenString, tt.tokenSecret)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ValidateJwt() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ValidateJwt() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("ValidateJwt() = %v, want %v", got, tt.want)
			}
		})
	}
}
