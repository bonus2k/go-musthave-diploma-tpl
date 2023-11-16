package utils

import "testing"

func TestCheckPassword(t *testing.T) {
	tests := []struct {
		name     string
		sign     string
		password string
		want     bool
	}{
		{
			name:     "password is valid",
			sign:     "password",
			password: "password",
			want:     true,
		},
		{
			name:     "password is invalid",
			sign:     "password",
			password: "password1",
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sign, _ := SignPassword(tt.password)
			got, _ := CheckPassword(tt.sign, sign)
			if got != tt.want {
				t.Errorf("CheckPassword() got = %v, want %v", got, tt.want)
			}
		})
	}
}
