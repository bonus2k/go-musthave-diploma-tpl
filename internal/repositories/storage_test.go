package repositories

import (
	"context"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal/repositories/testdata"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

var container *testdata.PostgresContainer

func TestMain(m *testing.M) {
	log.Println("setup test case")
	err := internal.InitLogger("info")
	if err != nil {
		log.Printf("err init logger: %v", err)
		os.Exit(1)
	}

	ctx := context.Background()
	container, err = testdata.NewPostgresContainer(ctx)
	if err != nil {
		log.Printf("err create postgres container: %v", err)
		os.Exit(1)
	}

	err = container.InitDB(ctx)
	if err != nil {
		log.Printf("err init db: %v", err)
		os.Exit(1)
	}

	code := m.Run()
	log.Println("teardown test case")
	err = container.Terminate(ctx)
	if err != nil {
		log.Printf("err terminating container: %v", err)
		os.Exit(1)
	}
	os.Exit(code)
}

func TestNewStore(t *testing.T) {
	db, err := container.InitData()
	if err != nil {
		log.Printf("err TestNewStore %v", err)
		os.Exit(1)
	}

	tests := []struct {
		name string
	}{
		{
			name: "smoke test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewStore(db)
			if store == nil {
				t.Errorf("NewStore() = %v is nil", store)
				return
			}
			if err := store.CheckConnection(); err != nil {
				t.Errorf("NewStore() = %v", store)
			}
		})
	}
}

func TestStore_AddOrder(t *testing.T) {
	db, err := container.InitData()
	if err != nil {
		log.Printf("err TestStore_AddOrder %v", err)
		os.Exit(1)
	}

	type args struct {
		ctx   context.Context
		order *internal.Order
	}
	tests := []struct {
		name       string
		args       args
		want       *internal.Order
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "add_order_violates_foreign_key_constraint",
			args: args{
				ctx: context.Background(),
				order: &internal.Order{
					ID:       uuid.New(),
					CreateAt: time.Now(),
					Number:   123456789,
					Accrual:  1.1,
					Status:   internal.OrderStatusNew,
					UserID:   uuid.New()},
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "SQLSTATE 23503",
		},
		{
			name: "add_order_is_exist_this_user",
			args: args{
				ctx: context.Background(),
				order: &internal.Order{
					ID:       uuid.New(),
					CreateAt: time.Now(),
					Number:   4539088167512356,
					Accrual:  1.1,
					Status:   internal.OrderStatusNew,
					UserID:   uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed026")},
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "this order is exist the user",
		},
		{
			name: "add_order_is_exist_another_user",
			args: args{
				ctx: context.Background(),
				order: &internal.Order{
					ID:       uuid.New(),
					CreateAt: time.Now(),
					Number:   4539088167512356,
					Accrual:  1.1,
					Status:   internal.OrderStatusNew,
					UserID:   uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed027")},
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "this order is exist another user",
		},
		{
			name: "add_order_is_correct",
			args: args{
				ctx: context.Background(),
				order: &internal.Order{
					ID:       uuid.MustParse("3e23bb5c-5cd6-4ca9-afa5-8d498576a080"),
					CreateAt: time.Now(),
					Number:   6011223604226714,
					Accrual:  1.1,
					Status:   internal.OrderStatusProcessed,
					UserID:   uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed027")},
			},
			want: &internal.Order{
				ID:       uuid.MustParse("3e23bb5c-5cd6-4ca9-afa5-8d498576a080"),
				CreateAt: time.Now(),
				Number:   6011223604226714,
				Accrual:  1.1,
				Status:   internal.OrderStatusProcessed,
				UserID:   uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed027")},
			wantErr:    false,
			wantErrMsg: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &Store{
				db: db,
			}
			got, err := store.AddOrder(tt.args.ctx, tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				assert.Containsf(t, err.Error(), tt.wantErrMsg, "expected error containing %q, got %s", tt.wantErrMsg, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddOrder() got = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestStore_AddUser(t *testing.T) {
	db, err := container.InitData()
	if err != nil {
		log.Printf("err TestStore_AddUser %v", err)
		os.Exit(1)
	}

	type args struct {
		ctx  context.Context
		user *internal.User
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "add_user_is_correct",
			args: args{
				ctx: context.Background(),
				user: &internal.User{
					ID:       uuid.New(),
					CreateAt: time.Now(),
					Login:    "TestUser3",
					Password: "password",
				},
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "add_user_is_exist",
			args: args{
				ctx: context.Background(),
				user: &internal.User{
					ID:       uuid.New(),
					CreateAt: time.Now(),
					Login:    "TestUser1",
					Password: "password",
				},
			},
			wantErr:    true,
			wantErrMsg: "user is exist",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &Store{
				db: db,
			}
			err := store.AddUser(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				assert.Containsf(t, err.Error(), tt.wantErrMsg, "expected error containing %q, got %s", tt.wantErrMsg, err)
			}
		})
	}
}

func TestStore_CheckConnection(t *testing.T) {
	db, err := container.InitData()
	if err != nil {
		log.Printf("err TestStore_CheckConnection %v", err)
		os.Exit(1)
	}
	wrongDB, err := sqlx.Open("pgx", "postgres://user:user@localhost:12345/postgres?sslmode=disable")
	if err != nil {
		log.Printf("err TestStore_CheckConnection %v", err)
		os.Exit(1)
	}

	tests := []struct {
		db      *sqlx.DB
		name    string
		wantErr bool
	}{
		{
			db:      db,
			name:    "check_connection_is_correct",
			wantErr: false,
		},
		{
			db:      wrongDB,
			name:    "check_connection_is_fault",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &Store{
				db: tt.db,
			}
			if err := store.CheckConnection(); (err != nil) != tt.wantErr {
				t.Errorf("CheckConnection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStore_FindUserByLogin(t *testing.T) {
	db, err := container.InitData()
	if err != nil {
		log.Printf("err TestStore_FindUserByLogin %v", err)
		os.Exit(1)
	}

	type args struct {
		ctx   context.Context
		login string
	}
	tests := []struct {
		name       string
		args       args
		want       *internal.User
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "user_is_not_exist",
			args: args{
				ctx:   context.Background(),
				login: "user",
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "user not found",
		},
		{
			name: "user_is_exist",
			args: args{
				ctx:   context.Background(),
				login: "TestUser1",
			},
			want: &internal.User{
				ID:       uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed026"),
				CreateAt: time.Date(2023, 01, 01, 14, 00, 00, 000, time.UTC),
				Login:    "TestUser1",
				Password: "password",
				Bill:     0,
			},
			wantErr:    false,
			wantErrMsg: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &Store{
				db: db,
			}
			got, err := store.FindUserByLogin(tt.args.ctx, tt.args.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindUserByLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				assert.Containsf(t, err.Error(), tt.wantErrMsg, "expected error containing %q, got %s", tt.wantErrMsg, err)
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("FindUserByLogin() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_GetOrders(t *testing.T) {
	db, err := container.InitData()
	if err != nil {
		log.Printf("err TestStore_GetOrders %v", err)
		os.Exit(1)
	}
	type args struct {
		ctx    context.Context
		userID uuid.UUID
	}
	tests := []struct {
		name         string
		args         args
		wantSizeList int
		wantErr      bool
	}{
		{
			name: "get_order_of_user_who_exist",
			args: args{
				ctx:    context.Background(),
				userID: uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed026"),
			},
			wantSizeList: 3,
			wantErr:      false,
		},
		{
			name: "get_order_of_user_who_not_exist",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New(),
			},
			wantSizeList: 0,
			wantErr:      false,
		},
		{
			name: "get_order_is_not_exist",
			args: args{
				ctx:    context.Background(),
				userID: uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed027"),
			},
			wantSizeList: 0,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &Store{
				db: db,
			}
			got, err := store.GetOrders(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(*got) != tt.wantSizeList {
				t.Errorf("GetOrders() got size = %v, want size %v", got, tt.wantSizeList)
			}
		})
	}
}

func TestStore_GetOrdersNotProcessed(t *testing.T) {
	db, err := container.InitData()
	if err != nil {
		log.Printf("TestStore_GetOrdersNotProcessed %v", err)
		os.Exit(1)
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name         string
		args         args
		wantSizeList int
		wantErr      bool
	}{
		{
			name:         "get_orders_not_processed",
			args:         args{ctx: context.Background()},
			wantSizeList: 1,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &Store{
				db: db,
			}
			got, err := store.GetOrdersNotProcessed(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrdersNotProcessed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(*got) != tt.wantSizeList {
				t.Errorf("GetOrdersNotProcessed() got size = %v, want size %v", got, tt.wantSizeList)
			}
		})
	}
}

func TestStore_GetUser(t *testing.T) {
	db, err := container.InitData()
	if err != nil {
		log.Printf("TestStore_GetUser %v", err)
		os.Exit(1)
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		want    *internal.User
		wantErr bool
	}{
		{
			name: "get_user_is_exist",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed027"),
			},
			want: &internal.User{
				ID:       uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed027"),
				CreateAt: time.Date(2023, 01, 01, 14, 00, 00, 000, time.UTC),
				Login:    "TestUser2",
				Password: "password",
				Bill:     100,
			},
			wantErr: false,
		},
		{
			name: "get_user_is_not_exist",
			args: args{
				ctx: context.Background(),
				id:  uuid.New(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &Store{
				db: db,
			}
			got, err := store.GetUser(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_GetWithdrawals(t *testing.T) {
	db, err := container.InitData()
	if err != nil {
		log.Printf("TestStore_GetWithdrawals %v", err)
		os.Exit(1)
	}
	var withdrawals []internal.Withdraw
	type args struct {
		ctx    context.Context
		userID uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		want    *[]internal.Withdraw
		wantErr bool
	}{
		{
			name: "get_withdraws_is_exist",
			args: args{
				ctx:    context.Background(),
				userID: uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed026"),
			},
			want: &[]internal.Withdraw{
				{
					ID:       uuid.MustParse("35e1cbd0-c3ba-44eb-8632-0d91c280dee6"),
					CreateAt: time.Date(2023, 11, 10, 11, 00, 00, 000, time.UTC),
					Order:    140672056,
					Sum:      502.6499938964844,
					UserID:   uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed026"),
				},
			},
			wantErr: false,
		},
		{
			name: "get_withdraws_is_not_exist",
			args: args{
				ctx:    context.Background(),
				userID: uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed027"),
			},
			want:    &withdrawals,
			wantErr: false,
		},
		{
			name: "get_withdraws_who_user_is_not_exist",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New(),
			},
			want:    &withdrawals,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &Store{
				db: db,
			}
			got, err := store.GetWithdrawals(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWithdrawals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("GetWithdrawals() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_SaveWithdrawal(t *testing.T) {
	db, err := container.InitData()
	if err != nil {
		log.Printf("TestStore_SaveWithdrawal %v", err)
		os.Exit(1)
	}
	type args struct {
		ctx        context.Context
		withdrawal *internal.Withdraw
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "save_withdrawal_who_user_enough_amount",
			args: args{
				ctx: context.Background(),
				withdrawal: &internal.Withdraw{
					ID:       uuid.New(),
					CreateAt: time.Now(),
					Order:    123456,
					Sum:      99,
					UserID:   uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed027"),
				},
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "save_withdrawal_who_user_not_enough_amount",
			args: args{
				ctx: context.Background(),
				withdrawal: &internal.Withdraw{
					ID:       uuid.New(),
					CreateAt: time.Now(),
					Order:    123456,
					Sum:      0.00000001,
					UserID:   uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed026"),
				},
			},
			wantErr:    true,
			wantErrMsg: "not enough amount",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &Store{
				db: db,
			}
			err := store.SaveWithdrawal(tt.args.ctx, tt.args.withdrawal)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveWithdrawal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				assert.Containsf(t, err.Error(), tt.wantErrMsg, "expected error containing %q, got %s", tt.wantErrMsg, err)
			}
		})
	}
}

func TestStore_UpdateOrder(t *testing.T) {
	db, err := container.InitData()
	if err != nil {
		log.Printf("TestStore_UpdateOrder %v", err)
		os.Exit(1)
	}
	type args struct {
		ctx   context.Context
		order *internal.Order
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		userBill float32
	}{
		{
			name: "add_order_processing",
			args: args{
				ctx: context.Background(),
				order: &internal.Order{
					ID:       uuid.New(),
					CreateAt: time.Now(),
					Number:   3536137811022331,
					Accrual:  100,
					Status:   internal.OrderStatusProcessing,
					UserID:   uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed026"),
				},
			},
			wantErr:  false,
			userBill: 0,
		},
		{
			name: "add_order_processed",
			args: args{
				ctx: context.Background(),
				order: &internal.Order{
					ID:       uuid.New(),
					CreateAt: time.Now(),
					Number:   3536137811022331,
					Accrual:  99.111,
					Status:   internal.OrderStatusProcessed,
					UserID:   uuid.MustParse("98dcfb07-e16f-4e53-9a28-d2a2e4eed026"),
				},
			},
			wantErr:  false,
			userBill: 99.111,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &Store{
				db: db,
			}
			if err := store.UpdateOrder(tt.args.ctx, tt.args.order); (err != nil) != tt.wantErr {
				t.Errorf("UpdateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
			user, err := store.GetUser(tt.args.ctx, tt.args.order.UserID)
			assert.NoErrorf(t, err, "GetUser() error = %v", err)
			assert.Equalf(t, tt.userBill, user.Bill, "Equal user bill() got = %v, want %v", user.Bill, tt.userBill)
		})
	}
}
