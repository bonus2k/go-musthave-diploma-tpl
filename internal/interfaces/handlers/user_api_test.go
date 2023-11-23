package handlers

import (
	"github.com/bonus2k/go-musthave-diploma-tpl/internal"
	mock "github.com/bonus2k/go-musthave-diploma-tpl/internal/mocks"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal/repositories"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal/services"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func init() {
	err := internal.InitLogger("info")
	if err != nil {
		log.Printf("err Init logger %v", err)
		os.Exit(1)
	}
}

type testData struct {
	mockStore     *mock.MockStore
	handlerUser   *HandlerUser
	userID1       string
	cookie1       http.Cookie
	userID2       string
	cookie2       http.Cookie
	wrongCookie   http.Cookie
	wrongIDCookie http.Cookie
}

func initTestServices(t *testing.T) *testData {
	sign := []byte{116, 79, 253, 154, 106, 127, 165, 70, 139, 56, 218, 213, 105, 253, 76}
	wrongSign := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	ctrl := gomock.NewController(t)
	mockStore := mock.NewMockStore(ctrl)
	service := services.NewUserService(mockStore)
	return &testData{
		mockStore:     mockStore,
		handlerUser:   NewHandlerUser(service, sign),
		userID1:       "98dcfb07-e16f-4e53-9a28-d2a2e4eed026",
		cookie1:       writeSigned("98dcfb07-e16f-4e53-9a28-d2a2e4eed026", sign),
		userID2:       "98dcfb07-e16f-4e53-9a28-d2a2e4eed027",
		cookie2:       writeSigned("98dcfb07-e16f-4e53-9a28-d2a2e4eed027", sign),
		wrongCookie:   writeSigned("123456 ", wrongSign),
		wrongIDCookie: writeSigned("123156", sign),
	}
}

func TestNewHandlerUser(t *testing.T) {
	sign := []byte{116, 79, 253, 154, 106, 127, 165, 70, 139, 56, 218, 213, 105, 253, 76}
	ctrl := gomock.NewController(t)
	mockStore := mock.NewMockStore(ctrl)
	service := services.NewUserService(mockStore)

	type args struct {
		service   *services.UserService
		secretKey []byte
	}
	tests := []struct {
		name string
		args args
		want *HandlerUser
	}{
		{
			name: "smoke test",
			args: args{
				service:   service,
				secretKey: sign,
			},
			want: NewHandlerUser(service, sign),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHandlerUser(tt.args.service, tt.args.secretKey); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("NewHandlerUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandlerUser_AddOrder(t *testing.T) {
	testServices := initTestServices(t)

	testServices.mockStore.EXPECT().
		AddOrder(gomock.Any(), &mock.MatchOrder{Order: &internal.Order{Number: 4539088167512356}}).
		Return(&internal.Order{}, nil).AnyTimes()

	testServices.mockStore.EXPECT().
		AddOrder(gomock.Any(), &mock.MatchOrder{Order: &internal.Order{Number: 3533841638640315}}).
		Return(nil, repositories.ErrOrderIsExistAnotherUser).AnyTimes()

	testServices.mockStore.EXPECT().
		AddOrder(gomock.Any(), &mock.MatchOrder{Order: &internal.Order{Number: 3536137811022331}}).
		Return(nil, repositories.ErrOrderIsExistThisUser).AnyTimes()

	tests := []struct {
		name        string
		body        string
		contentType string
		statusCode  int
		cookie      *http.Cookie
	}{
		{
			name:        "add order 400",
			body:        "4539088167512356",
			contentType: "application/json",
			statusCode:  400,
			cookie:      &testServices.cookie1,
		},
		{
			name:        "add order 401",
			body:        "4539088167512356",
			contentType: "text/plain",
			statusCode:  401,
			cookie:      &testServices.wrongCookie,
		},
		{
			name:        "add order 200",
			body:        "3536137811022331",
			contentType: "text/plain",
			statusCode:  200,
			cookie:      &testServices.cookie1,
		},
		{
			name:        "add order 409",
			body:        "3533841638640315",
			contentType: "text/plain",
			statusCode:  409,
			cookie:      &testServices.cookie1,
		},
		{
			name:        "add order 422",
			body:        "12345",
			contentType: "text/plain",
			statusCode:  422,
			cookie:      &testServices.cookie1,
		},
		{
			name:        "add order 202",
			body:        "4539088167512356",
			contentType: "text/plain",
			statusCode:  202,
			cookie:      &testServices.cookie1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(tt.body))

			request.AddCookie(tt.cookie)
			request.Header.Set("Content-Type", tt.contentType)
			responseRecorder := httptest.NewRecorder()

			testServices.handlerUser.AddOrder(responseRecorder, request)
			result := responseRecorder.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.statusCode, result.StatusCode)
		})
	}
}

func TestHandlerUser_AddWithdraw(t *testing.T) {
	testServices := initTestServices(t)

	testServices.mockStore.EXPECT().
		SaveWithdrawal(gomock.Any(), &mock.MatchWithdraw{Withdraw: &internal.Withdraw{Order: 4539088167512356}}).
		Return(repositories.ErrNotEnoughAmount).AnyTimes()

	testServices.mockStore.EXPECT().
		SaveWithdrawal(gomock.Any(), &mock.MatchWithdraw{Withdraw: &internal.Withdraw{Order: 3533841638640315}}).
		Return(nil).AnyTimes()

	tests := []struct {
		name        string
		body        string
		contentType string
		statusCode  int
		cookie      *http.Cookie
	}{
		{
			name:        "AddWithdraw 401",
			body:        "",
			contentType: "application/json",
			statusCode:  401,
			cookie:      &testServices.wrongCookie,
		},
		{
			name:        "AddWithdraw 500",
			body:        "",
			contentType: "application/json",
			statusCode:  500,
			cookie:      &testServices.cookie1,
		},
		{
			name:        "AddWithdraw 402",
			body:        `{"order": "4539088167512356", "sum": 751}`,
			contentType: "application/json",
			statusCode:  402,
			cookie:      &testServices.cookie1,
		},
		{
			name:        "AddWithdraw 422",
			body:        `{"order": "12345", "sum": 751}`,
			contentType: "application/json",
			statusCode:  422,
			cookie:      &testServices.cookie1,
		},
		{
			name:        "AddWithdraw 200",
			body:        `{"order": "3533841638640315", "sum": 751}`,
			contentType: "application/json",
			statusCode:  200,
			cookie:      &testServices.cookie1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(tt.body))

			request.AddCookie(tt.cookie)
			request.Header.Set("Content-Type", tt.contentType)
			responseRecorder := httptest.NewRecorder()

			testServices.handlerUser.AddWithdraw(responseRecorder, request)
			result := responseRecorder.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.statusCode, result.StatusCode)
		})
	}
}

func TestHandlerUser_GetBalance(t *testing.T) {
	testServices := initTestServices(t)

	withdraws := &[]internal.Withdraw{
		{
			ID:       uuid.MustParse("35e1cbd0-c3ba-44eb-8632-0d91c280dee6"),
			CreateAt: time.Date(2023, 11, 10, 14, 00, 00, 000, time.Local),
			Order:    140672056,
			Sum:      12.64,
			UserID:   uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed026"),
		},
	}

	testServices.mockStore.EXPECT().
		GetWithdrawals(gomock.Any(), uuid.MustParse(testServices.userID1)).
		Return(withdraws, nil).AnyTimes()

	testServices.mockStore.EXPECT().
		GetUser(gomock.Any(), uuid.MustParse(testServices.userID1)).
		Return(&internal.User{Bill: 0.001}, nil).AnyTimes()

	tests := []struct {
		name        string
		body        string
		contentType string
		statusCode  int
		wantBody    string
		cookie      *http.Cookie
	}{
		{
			name:        "GetBalance 401",
			contentType: "application/json",
			statusCode:  401,
			cookie:      &testServices.wrongCookie,
		},
		{
			name:        "GetBalance 200",
			contentType: "application/json",
			statusCode:  200,
			cookie:      &testServices.cookie1,
			wantBody:    `{"current":0.001,"withdrawn":12.64}`,
		},
		{
			name:        "GetBalance 500",
			contentType: "application/json",
			statusCode:  500,
			cookie:      &testServices.wrongIDCookie,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(tt.body))

			request.AddCookie(tt.cookie)
			request.Header.Set("Content-Type", tt.contentType)
			responseRecorder := httptest.NewRecorder()

			testServices.handlerUser.GetBalance(responseRecorder, request)
			result := responseRecorder.Result()

			defer result.Body.Close()
			resBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.statusCode, result.StatusCode)
			if len(resBody) != 0 {
				assert.JSONEq(t, tt.wantBody, string(resBody))
			}

		})
	}
}

func TestHandlerUser_GetOrders(t *testing.T) {
	testServices := initTestServices(t)

	orders := &[]internal.Order{
		{
			ID:       uuid.MustParse("334b0360-8222-44fc-bf2e-77ced208f2cd"),
			CreateAt: time.Date(2023, 01, 01, 14, 01, 00, 000, time.Local),
			Number:   4539088167512356,
			Accrual:  100.0,
			Status:   internal.OrderStatusProcessed,
			UserID:   uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed026"),
		},
		{
			ID:       uuid.MustParse("334b0360-8222-44fc-bf2e-77ced208f2ce"),
			CreateAt: time.Date(2023, 01, 01, 14, 02, 00, 000, time.Local),
			Number:   3536137811022331,
			Accrual:  0,
			Status:   internal.OrderStatusNew,
			UserID:   uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed026"),
		},
		{
			ID:       uuid.MustParse("334b0360-8222-44fc-bf2e-77ced208f2cf"),
			CreateAt: time.Date(2023, 01, 01, 14, 03, 00, 000, time.Local),
			Number:   3533841638640315,
			Accrual:  0,
			Status:   internal.OrderStatusInvalid,
			UserID:   uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed026"),
		},
	}

	testServices.mockStore.EXPECT().
		GetOrders(gomock.Any(), uuid.MustParse(testServices.userID2)).
		Return(&[]internal.Order{}, nil).AnyTimes()

	testServices.mockStore.EXPECT().
		GetOrders(gomock.Any(), uuid.MustParse(testServices.userID1)).
		Return(orders, nil).AnyTimes()

	tests := []struct {
		name        string
		body        string
		contentType string
		statusCode  int
		wantBody    string
		cookie      *http.Cookie
	}{
		{
			name:        "GetOrders 401",
			contentType: "application/json",
			statusCode:  401,
			cookie:      &testServices.wrongCookie,
		},
		{
			name:        "GetOrders 500",
			contentType: "application/json",
			statusCode:  500,
			cookie:      &testServices.wrongIDCookie,
		},
		{
			name:        "GetOrders 204",
			contentType: "application/json",
			statusCode:  204,
			cookie:      &testServices.cookie2,
		},
		{
			name:        "GetOrders 200",
			contentType: "application/json",
			statusCode:  200,
			cookie:      &testServices.cookie1,
			wantBody: `[
							{"number":"4539088167512356","status":"PROCESSED","accrual":100,"uploaded_at":"2023-01-01T14:01:00+03:00"},
							{"number":"3536137811022331","status":"NEW","accrual":0,"uploaded_at":"2023-01-01T14:02:00+03:00"},
							{"number":"3533841638640315","status":"INVALID","accrual":0,"uploaded_at":"2023-01-01T14:03:00+03:00"}
						]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(tt.body))

			request.AddCookie(tt.cookie)
			request.Header.Set("Content-Type", tt.contentType)
			responseRecorder := httptest.NewRecorder()

			testServices.handlerUser.GetOrders(responseRecorder, request)
			result := responseRecorder.Result()

			defer result.Body.Close()
			resBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.statusCode, result.StatusCode)
			if len(resBody) != 0 {
				assert.JSONEq(t, tt.wantBody, string(resBody))
			}

		})
	}
}

func TestHandlerUser_GetWithdrawals(t *testing.T) {
	testServices := initTestServices(t)

	withdraws := &[]internal.Withdraw{
		{
			ID:       uuid.MustParse("35e1cbd0-c3ba-44eb-8632-0d91c280dee6"),
			CreateAt: time.Date(2023, 11, 10, 14, 00, 00, 000, time.Local),
			Order:    140672056,
			Sum:      12.64,
			UserID:   uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed026"),
		},
	}

	testServices.mockStore.EXPECT().
		GetWithdrawals(gomock.Any(), uuid.MustParse(testServices.userID1)).
		Return(withdraws, nil).AnyTimes()

	testServices.mockStore.EXPECT().
		GetWithdrawals(gomock.Any(), uuid.MustParse(testServices.userID2)).
		Return(&[]internal.Withdraw{}, nil).AnyTimes()

	tests := []struct {
		name        string
		body        string
		contentType string
		statusCode  int
		wantBody    string
		cookie      *http.Cookie
	}{
		{
			name:        "GetWithdrawals 401",
			contentType: "application/json",
			statusCode:  401,
			cookie:      &testServices.wrongCookie,
		},
		{
			name:        "GetWithdrawals 500",
			contentType: "application/json",
			statusCode:  500,
			cookie:      &testServices.wrongIDCookie,
		},
		{
			name:        "GetWithdrawals 204",
			contentType: "application/json",
			statusCode:  204,
			cookie:      &testServices.cookie2,
		},
		{
			name:        "GetWithdrawals 200",
			contentType: "application/json",
			statusCode:  200,
			cookie:      &testServices.cookie1,
			wantBody:    `[{"order":"140672056","sum":12.64,"processed_at":"2023-11-10T14:00:00+03:00"}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(tt.body))

			request.AddCookie(tt.cookie)
			request.Header.Set("Content-Type", tt.contentType)
			responseRecorder := httptest.NewRecorder()

			testServices.handlerUser.GetWithdrawals(responseRecorder, request)
			result := responseRecorder.Result()

			defer result.Body.Close()
			resBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.statusCode, result.StatusCode)
			if len(resBody) != 0 {
				assert.JSONEq(t, tt.wantBody, string(resBody))
			}

		})
	}
}

func TestHandlerUser_Login(t *testing.T) {
	testServices := initTestServices(t)

	user := &internal.User{
		ID:       uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed026"),
		CreateAt: time.Date(2023, 01, 01, 14, 00, 00, 000, time.Local),
		Login:    "TestUser1",
		Password: "1bf92ee0af9687162f7f9c861a1d2cbfdaf2e3ab5ec70335e0d68f5455b54d6dfd631dd94175d250",
		Bill:     0,
	}

	testServices.mockStore.EXPECT().
		FindUserByLogin(gomock.Any(), "TestUser1").
		Return(user, nil).AnyTimes()

	tests := []struct {
		name        string
		body        string
		contentType string
		statusCode  int
		wantCookie  *http.Cookie
	}{
		{
			name:        "User Login 400",
			contentType: "text",
			statusCode:  400,
		},
		{
			name:        "User Login 500",
			contentType: "application/json",
			statusCode:  500,
		},
		{
			name:        "User Login 200",
			contentType: "application/json",
			statusCode:  200,
			body: `{
						"login": "TestUser1",
						"password": "Password"
					}`,
			wantCookie: &testServices.cookie1,
		},
		{
			name:        "User Login 401",
			contentType: "application/json",
			statusCode:  401,
			body: `{
						"login": "TestUser1",
						"password": "password"
					}`,
			wantCookie: &testServices.cookie1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(tt.body))

			request.Header.Set("Content-Type", tt.contentType)
			responseRecorder := httptest.NewRecorder()

			testServices.handlerUser.Login(responseRecorder, request)
			result := responseRecorder.Result()
			defer result.Body.Close()
			cookies := result.Cookies()

			assert.Equal(t, tt.statusCode, result.StatusCode)
			if len(cookies) != 0 {
				assert.Equal(t, tt.wantCookie.Value, cookies[0].Value)
			}

		})
	}
}

func TestHandlerUser_RegisterUser(t *testing.T) {
	testServices := initTestServices(t)

	testServices.mockStore.EXPECT().
		AddUser(gomock.Any(), &mock.MatchUser{User: &internal.User{Login: "TestUser1"}}).
		Return(nil).AnyTimes()

	testServices.mockStore.EXPECT().
		AddUser(gomock.Any(), &mock.MatchUser{User: &internal.User{Login: "TestUser2"}}).
		Return(repositories.ErrUserIsExist).AnyTimes()

	tests := []struct {
		name        string
		body        string
		contentType string
		statusCode  int
		wantCookie  *http.Cookie
	}{
		{
			name:        "RegisterUser 400",
			contentType: "text",
			statusCode:  400,
		},
		{
			name:        "RegisterUser 500",
			contentType: "application/json",
			statusCode:  500,
		},
		{
			name:        "RegisterUser 200",
			contentType: "application/json",
			statusCode:  200,
			body: `{
						"login": "TestUser1",
						"password": "password"
					}`,
			wantCookie: &testServices.cookie1,
		},
		{
			name:        "RegisterUser 409",
			contentType: "application/json",
			statusCode:  409,
			body: `{
						"login": "TestUser2",
						"password": "password"
					}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(tt.body))

			request.Header.Set("Content-Type", tt.contentType)
			responseRecorder := httptest.NewRecorder()

			testServices.handlerUser.RegisterUser(responseRecorder, request)
			result := responseRecorder.Result()
			defer result.Body.Close()
			cookies := result.Cookies()

			assert.Equal(t, tt.statusCode, result.StatusCode)
			if len(cookies) != 0 {
				assert.Equal(t, tt.wantCookie.Name, cookies[0].Name)
				assert.NotEqual(t, cookies[0].Value, "")
			}

		})
	}
}

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

func Test_writeSigned(t *testing.T) {
	sign := []byte{116, 79, 253, 154, 106, 127, 165, 70, 139, 56, 218, 213, 105, 253, 76}
	type args struct {
		value  string
		secret []byte
	}
	tests := []struct {
		name string
		args args
		want http.Cookie
	}{
		{
			name: "get cookie is correct",
			args: args{
				value:  "42f0558c-04f3-4e11-9ee1-6de717ca69e9",
				secret: sign,
			},
			want: http.Cookie{
				Name:     "gophermart",
				Value:    "VPnFhpmNlKNCWJqE0g25dR76M2e8mYmKSUM5lKzw8zA0MmYwNTU4Yy0wNGYzLTRlMTEtOWVlMS02ZGU3MTdjYTY5ZTk=",
				Path:     "/",
				MaxAge:   3600,
				HttpOnly: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := writeSigned(tt.args.value, tt.args.secret); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("writeSigned() = %v, want %v", got, tt.want)
			}
		})
	}
}
