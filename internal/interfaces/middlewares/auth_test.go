package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_readSigned(t *testing.T) {
	sign := []byte{116, 79, 253, 154, 106, 127, 165, 70, 139, 56, 218, 213, 105, 253, 76}

	type args struct {
		cookie    *http.Cookie
		secretKey []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "readSigned wrong name",
			args: args{
				cookie:    &http.Cookie{Name: "uups"},
				secretKey: sign,
			},
			wantErr: true,
		},
		{
			name: "readSigned wrong value",
			args: args{
				cookie: &http.Cookie{
					Name:  "gophermart",
					Value: "123456",
				},
				secretKey: nil,
			},
			wantErr: true,
		},
		{
			name: "readSigned wrong secret key",
			args: args{
				cookie: &http.Cookie{
					Name:  "gophermart",
					Value: "VPnFhpmNlKNCWJqE0g25dR76M2e8mYmKSUM5lKzw8zA0MmYwNTU4Yy0wNGYzLTRlMTEtOWVlMS02ZGU3MTdjYTY5ZTk=",
				},
				secretKey: []byte{1, 2, 3, 4, 5, 6, 7},
			},
			wantErr: true,
		},
		{
			name: "readSigned correct",
			args: args{
				cookie: &http.Cookie{
					Name:  "gophermart",
					Value: "VPnFhpmNlKNCWJqE0g25dR76M2e8mYmKSUM5lKzw8zA0MmYwNTU4Yy0wNGYzLTRlMTEtOWVlMS02ZGU3MTdjYTY5ZTk=",
				},
				secretKey: sign,
			},
			wantErr: false,
			want:    "42f0558c-04f3-4e11-9ee1-6de717ca69e9",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.AddCookie(tt.args.cookie)
			got, err := readSigned(request, tt.args.secretKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("readSigned() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readSigned() got = %v, want %v", got, tt.want)
			}
		})
	}
}
