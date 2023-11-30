package services

import (
	"context"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal"
	mock "github.com/bonus2k/go-musthave-diploma-tpl/internal/mocks"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal/repositories"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"reflect"
	"sort"
	"testing"
	"time"
)

func init() {
	log.Printf("init logger")
	err := internal.InitLogger("info")
	if err != nil {
		log.Printf("err init logger %v", err)
		os.Exit(1)
	}
}

func getStore(t *testing.T) *mock.MockStore {
	ctrl := gomock.NewController(t)
	mockStore := mock.NewMockStore(ctrl)
	return mockStore
}

func TestNewUserService(t *testing.T) {
	mockStore := getStore(t)
	type args struct {
		storage repositories.Store
	}
	tests := []struct {
		name string
		args args
		want *UserService
	}{
		{
			name: "smoke test",
			args: args{storage: mockStore},
			want: &UserService{db: mockStore},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserService(tt.args.storage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_AddOrder(t *testing.T) {
	mockStore := getStore(t)
	mockStore.EXPECT().AddOrder(gomock.Any(), gomock.Any()).Return(&internal.Order{}, nil)
	type args struct {
		ctx     context.Context
		id      string
		orderID string
	}
	tests := []struct {
		name       string
		db         repositories.Store
		args       args
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "add_order_illegal_number_of_order",
			db:   mockStore,
			args: args{
				ctx:     context.Background(),
				id:      uuid.New().String(),
				orderID: "123",
			},
			wantErr:    true,
			wantErrMsg: "illegal order",
		},
		{
			name: "add_order_illegal_userID",
			db:   mockStore,
			args: args{
				ctx:     context.Background(),
				id:      "123",
				orderID: "4539088167512356",
			},
			wantErr:    true,
			wantErrMsg: "invalid UUID",
		},
		{
			name: "add_order",
			db:   mockStore,
			args: args{
				ctx:     context.Background(),
				id:      uuid.New().String(),
				orderID: "4539088167512356",
			},
			wantErr:    false,
			wantErrMsg: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				db: tt.db,
			}
			err := us.AddOrder(tt.args.ctx, tt.args.id, tt.args.orderID)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				assert.Containsf(t, err.Error(), tt.wantErrMsg, "AddOrder() error = %v, wantErr %v", err, tt.wantErrMsg)
			}
		})
	}
}

func TestUserService_AddWithdraw(t *testing.T) {
	mockStore := getStore(t)
	mockStore.EXPECT().SaveWithdrawal(gomock.Any(), gomock.Any()).Return(nil)
	type args struct {
		ctx context.Context
		dto internal.WithdrawDto
		id  string
	}
	tests := []struct {
		name       string
		db         repositories.Store
		args       args
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "add_withdraw_illegal_number_of_order",
			db:   mockStore,
			args: args{
				ctx: context.Background(),
				dto: internal.WithdrawDto{Order: "123"},
				id:  uuid.New().String(),
			},
			wantErr:    true,
			wantErrMsg: "illegal order",
		},
		{
			name: "add_withdraw_illegal_userID",
			db:   mockStore,
			args: args{
				ctx: context.Background(),
				dto: internal.WithdrawDto{Order: "4539088167512356"},
				id:  "123",
			},
			wantErr:    true,
			wantErrMsg: "invalid UUID",
		},
		{
			name: "add_withdraw",
			db:   mockStore,
			args: args{
				ctx: context.Background(),
				dto: internal.WithdrawDto{Order: "4539088167512356"},
				id:  uuid.New().String(),
			},
			wantErr:    false,
			wantErrMsg: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				db: tt.db,
			}
			err := us.AddWithdraw(tt.args.ctx, tt.args.dto, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddWithdraw() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				assert.Containsf(t, err.Error(), tt.wantErrMsg, "AddWithdraw() error = %v, wantErr %v", err, tt.wantErrMsg)
			}
		})
	}
}

func TestUserService_CreateNewUser(t *testing.T) {
	mockStore := getStore(t)
	mockStore.EXPECT().AddUser(gomock.Any(), gomock.Any()).Return(nil)
	type args struct {
		ctx  context.Context
		user *internal.UserDto
	}
	tests := []struct {
		name       string
		db         repositories.Store
		args       args
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "add_user_without_password",
			db:   mockStore,
			args: args{
				ctx:  context.Background(),
				user: &internal.UserDto{Login: "user", Pass: " "},
			},
			wantErr:    true,
			wantErrMsg: "illegal user argument",
		},
		{
			name: "add_user_without_username",
			db:   mockStore,
			args: args{
				ctx:  context.Background(),
				user: &internal.UserDto{Login: "   ", Pass: "pass"},
			},
			wantErr:    true,
			wantErrMsg: "illegal user argument",
		},
		{
			name: "add_user",
			db:   mockStore,
			args: args{
				ctx:  context.Background(),
				user: &internal.UserDto{Login: "TestUser", Pass: "Password"},
			},
			wantErr:    false,
			wantErrMsg: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				db: tt.db,
			}
			got, err := us.CreateNewUser(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateNewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.Equalf(t, tt.args.user.Login, got.Login, "CreateNewUser() login = %v, wantLogin %v", tt.args.user.Login, got.Login)
				assert.NotNil(t, got, "CreateNewUser() got must be not nil")
			}

		})
	}
}

func TestUserService_GetBalance(t *testing.T) {
	mockStore := getStore(t)
	withdraws := &[]internal.Withdraw{
		{
			ID:       uuid.MustParse("35e1cbd0-c3ba-44eb-8632-0d91c280dee6"),
			CreateAt: time.Date(2023, 11, 10, 14, 00, 00, 000, time.Local),
			Order:    140672056,
			Sum:      12.64,
			UserID:   uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed026"),
		},
		{
			ID:       uuid.MustParse("35e1cbd0-c3ba-44eb-8632-0d91c280dee7"),
			CreateAt: time.Date(2023, 12, 10, 14, 00, 00, 000, time.Local),
			Order:    140672057,
			Sum:      27.385,
			UserID:   uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed026"),
		},
		{
			ID:       uuid.MustParse("35e1cbd0-c3ba-44eb-8632-0d91c280dee8"),
			CreateAt: time.Date(2023, 11, 11, 14, 00, 00, 000, time.Local),
			Order:    140672058,
			Sum:      0.11111111,
			UserID:   uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed026"),
		},
	}
	mockStore.EXPECT().GetWithdrawals(gomock.Any(), uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed026")).Return(withdraws, nil)
	mockStore.EXPECT().GetUser(gomock.Any(), uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed026")).Return(&internal.User{Bill: 0.001}, nil)
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		db      repositories.Store
		args    args
		want    *internal.Balance
		wantErr bool
	}{
		{
			name: "get_balance",
			db:   mockStore,
			args: args{
				ctx: context.Background(),
				id:  "98dcfb07-e16f-4e53-9a28-d2a2e4eed026",
			},
			want:    &internal.Balance{Current: 0.001, Withdrawn: 40.136112},
			wantErr: false,
		},
		{
			name: "get_balance_for_wrong_userID",
			db:   mockStore,
			args: args{
				ctx: context.Background(),
				id:  "123",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				db: tt.db,
			}
			got, err := us.GetBalance(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBalance() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_GetOrders(t *testing.T) {
	mockStore := getStore(t)
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
	ordersDTO := &[]internal.OrderDto{
		{
			Number:  "4539088167512356",
			Status:  "PROCESSED",
			Accrual: 100,
			Upload:  time.Date(2023, 01, 01, 14, 01, 00, 000, time.Local),
		},
		{
			Number:  "3536137811022331",
			Status:  "NEW",
			Accrual: 0,
			Upload:  time.Date(2023, 01, 01, 14, 02, 00, 000, time.Local),
		},
		{
			Number:  "3533841638640315",
			Status:  "INVALID",
			Accrual: 0,
			Upload:  time.Date(2023, 01, 01, 14, 03, 00, 000, time.Local),
		},
	}
	mockStore.EXPECT().GetOrders(gomock.Any(), gomock.Any()).Return(orders, nil)
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		db      repositories.Store
		args    args
		want    *[]internal.OrderDto
		wantErr bool
	}{
		{
			name: "get_orders_for_wrong_userID",
			db:   mockStore,
			args: args{
				ctx: context.Background(),
				id:  "123",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "get_orders",
			db:   mockStore,
			args: args{
				ctx: context.Background(),
				id:  "98dcfb07-e16f-4e53-9a28-d2a2e4eed026",
			},
			want:    ordersDTO,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				db: tt.db,
			}
			got, err := us.GetOrders(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				sort.Slice((*got), func(i, j int) bool {
					return (*got)[i].Number > (*got)[j].Number
				})
				sort.Slice((*tt.want), func(i, j int) bool {
					return (*tt.want)[i].Number > (*tt.want)[j].Number
				})
			}

			if !cmp.Equal(got, tt.want) {
				t.Errorf("GetOrders() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_GetOrdersNotProcessed(t *testing.T) {
	mockStore1 := getStore(t)
	mockStore2 := getStore(t)
	orders := &[]internal.Order{
		{
			ID:       uuid.MustParse("334b0360-8222-44fc-bf2e-77ced208f2ce"),
			CreateAt: time.Date(2023, 01, 01, 14, 02, 00, 000, time.Local),
			Number:   3536137811022331,
			Accrual:  0,
			Status:   internal.OrderStatusNew,
			UserID:   uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed026"),
		},
	}
	mockStore1.EXPECT().GetOrdersNotProcessed(gomock.Any()).Return(orders, nil)
	mockStore2.EXPECT().GetOrdersNotProcessed(gomock.Any()).Return(&[]internal.Order{}, nil)
	tests := []struct {
		name       string
		db         repositories.Store
		want       []string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:       "get_orders_not_processed",
			db:         mockStore1,
			want:       []string{"3536137811022331"},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:       "get_orders_not_processed_empty_list",
			db:         mockStore2,
			want:       nil,
			wantErr:    true,
			wantErrMsg: "empty list",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				db: tt.db,
			}
			got, err := us.GetOrdersNotProcessed()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrdersNotProcessed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equalf(t, got, tt.want, "GetOrdersNotProcessed() got = %v, want %v", got, tt.want)
			if tt.wantErr {
				assert.Containsf(t, err.Error(), tt.wantErrMsg, "GetOrdersNotProcessed() error = %v, wantErr %v", err, tt.wantErrMsg)
			}
		})
	}
}

func TestUserService_GetWithdrawals(t *testing.T) {
	mockStore := getStore(t)
	withdraws := &[]internal.Withdraw{
		{
			ID:       uuid.MustParse("35e1cbd0-c3ba-44eb-8632-0d91c280dee6"),
			CreateAt: time.Date(2023, 11, 10, 14, 00, 00, 000, time.Local),
			Order:    140672056,
			Sum:      12.64,
			UserID:   uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed026"),
		},
		{
			ID:       uuid.MustParse("35e1cbd0-c3ba-44eb-8632-0d91c280dee7"),
			CreateAt: time.Date(2023, 12, 10, 14, 00, 00, 000, time.Local),
			Order:    140672057,
			Sum:      27.385,
			UserID:   uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed026"),
		},
		{
			ID:       uuid.MustParse("35e1cbd0-c3ba-44eb-8632-0d91c280dee8"),
			CreateAt: time.Date(2023, 11, 11, 14, 00, 00, 000, time.Local),
			Order:    140672058,
			Sum:      0.11111111,
			UserID:   uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed026"),
		},
	}
	mockStore.EXPECT().GetWithdrawals(gomock.Any(), gomock.Any()).Return(withdraws, nil)
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		db      repositories.Store
		args    args
		want    *[]internal.WithdrawDto
		wantErr bool
	}{
		{
			name: "get_withdrawals",
			db:   mockStore,
			args: args{
				ctx: context.Background(),
				id:  "98dcfb07-e16f-4e53-9a28-d2a2e4eed026",
			},
			want: &[]internal.WithdrawDto{
				{
					Order:    "140672056",
					Sum:      12.64,
					CreateAt: time.Date(2023, 11, 10, 14, 00, 00, 000, time.Local),
				},
				{
					Order:    "140672057",
					Sum:      27.385,
					CreateAt: time.Date(2023, 12, 10, 14, 00, 00, 000, time.Local),
				},
				{
					Order:    "140672058",
					Sum:      0.11111111,
					CreateAt: time.Date(2023, 11, 11, 14, 00, 00, 000, time.Local),
				},
			},
			wantErr: false,
		},
		{
			name: "get_withdrawals_wrong_userID",
			db:   mockStore,
			args: args{
				ctx: context.Background(),
				id:  "123",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				db: tt.db,
			}
			got, err := us.GetWithdrawals(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWithdrawals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				sort.Slice(*got, func(i, j int) bool {
					return (*got)[i].Order > (*got)[j].Order
				})
				sort.Slice(*tt.want, func(i, j int) bool {
					return (*tt.want)[i].Order > (*tt.want)[j].Order
				})
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("GetWithdrawals() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_LoginUser(t *testing.T) {
	mockStore := getStore(t)
	userID := uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed026")
	user := &internal.User{
		ID:       userID,
		CreateAt: time.Date(2023, 01, 01, 14, 00, 00, 000, time.Local),
		Login:    "TestUser1",
		Password: "1bf92ee0af9687162f7f9c861a1d2cbfdaf2e3ab5ec70335e0d68f5455b54d6dfd631dd94175d250",
		Bill:     0,
	}
	mockStore.EXPECT().FindUserByLogin(gomock.Any(), gomock.Any()).Return(user, nil).AnyTimes()
	type args struct {
		ctx  context.Context
		user internal.UserDto
	}
	tests := []struct {
		name       string
		db         repositories.Store
		args       args
		want       *uuid.UUID
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "login_user_wrong_password",
			db:   mockStore,
			args: args{
				ctx: context.Background(),
				user: internal.UserDto{
					Login: "TestUser1",
					Pass:  "password1",
				},
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "wrong authorization",
		},
		{
			name: "login_user",
			db:   mockStore,
			args: args{
				ctx: context.Background(),
				user: internal.UserDto{
					Login: "TestUser1",
					Pass:  "Password",
				},
			},
			want:       &userID,
			wantErr:    false,
			wantErrMsg: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				db: tt.db,
			}
			got, err := us.LoginUser(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoginUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_UpdateOrder(t *testing.T) {
	mockStore := getStore(t)
	mockStore.EXPECT().UpdateOrder(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	tests := []struct {
		name       string
		db         repositories.Store
		accrual    *internal.AccrualDto
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "update_order",
			db:   mockStore,
			accrual: &internal.AccrualDto{
				Order:   "4539088167512356",
				Status:  "NEW",
				Accrual: 0.01,
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "update_order_wrong_order_number",
			db:   mockStore,
			accrual: &internal.AccrualDto{
				Order:   " ",
				Status:  "NEW",
				Accrual: 0.01,
			},
			wantErr:    true,
			wantErrMsg: "parse accrual number",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				db: tt.db,
			}
			err := us.UpdateOrder(tt.accrual)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				assert.Containsf(t, err.Error(), tt.wantErrMsg, "UpdateOrder() error = %v, wantErr %v", err, tt.wantErrMsg)
			}
		})
	}
}

func Test_checksum(t *testing.T) {

	tests := []struct {
		name     string
		checksum map[int]int
	}{
		{
			name:     "test_checksum_batch",
			checksum: map[int]int{1: 2, 10: 1, 2000000000: 2, 4539088167512356: 8, 3536137811022331: 1, 3533841638640315: 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.checksum {
				if got := checksum(k); got != v {
					t.Errorf("checksum(%v) = %v, want %v", k, got, v)
				}
			}

		})
	}
}

func Test_isIllegalUserArgument(t *testing.T) {
	tests := []struct {
		name string
		user *internal.UserDto
		want bool
	}{
		{
			name: "wrong_password",
			user: &internal.UserDto{
				Login: "user",
				Pass:  "",
			},
			want: true,
		},
		{
			name: "wrong_login",
			user: &internal.UserDto{
				Login: " ",
				Pass:  "password",
			},
			want: true,
		},
		{
			name: "wrong_login_and_password",
			user: &internal.UserDto{
				Login: "    ",
				Pass:  "  ",
			},
			want: true,
		},
		{
			name: "correct user",
			user: &internal.UserDto{
				Login: "login",
				Pass:  "password",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isIllegalUserArgument(tt.user); got != tt.want {
				t.Errorf("isIllegalUserArgument() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isLuna(t *testing.T) {
	type order struct {
		str string
		num int
	}
	tests := []struct {
		name   string
		orders map[order]bool
	}{
		{
			name: "test_isLuna_batch",
			orders: map[order]bool{
				order{"1", 1}:                               false,
				order{"1234567890", 1234567890}:             false,
				order{"4539088167512356", 4539088167512356}: true,
				order{"3536137811022331", 3536137811022331}: true,
				order{"112345678912345", 112345678912345}:   false,
				order{"140672056", 140672056}:               true,
				order{"10000000000000", 10000000000000}:     false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for order, b := range tt.orders {
				got, ok := isLuna(order.str)
				if got != order.num || b != ok {
					t.Errorf("isLuna(%v) = %v, %v, want %v, %v", order.str, order.num, b, got, ok)
				}
			}

		})
	}
}
