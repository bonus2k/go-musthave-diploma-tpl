package utils

import (
	"errors"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal"
	"testing"
)

func TestRetryAfterError(t *testing.T) {

	internal.InitLogger("debug")

	tests := []struct {
		name    string
		f       func() error
		wantErr bool
	}{
		{
			name: "func with err",
			f: func() error {
				return errors.New("test err")
			},
			wantErr: true,
		},
		{
			name: "func without err",
			f: func() error {
				return nil
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RetryAfterError(tt.f); (err != nil) != tt.wantErr {
				t.Errorf("RetryAfterError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
