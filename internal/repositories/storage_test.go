package repositories

import (
	"context"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal/repositories/testdata"
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

	type args struct {
		db *sqlx.DB
	}
	tests := []struct {
		name string
		args args
		want *Store
	}{
		{
			name: "smoke test",
			args: args{db: db},
			want: NewStore(db),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewStore(tt.args.db)
			if !reflect.DeepEqual(store, tt.want) {
				t.Errorf("NewStore() = %v, want %v", store, tt.want)
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
			if (err != nil) == tt.wantErr {
				assert.Containsf(t, err.Error(), tt.wantErrMsg, "expected error containing %q, got %s", tt.wantErrMsg, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddOrder() got = %v, want %v", got, tt.want)
			}

		})
	}
}

//
//func TestStore_AddUser(t *testing.T) {
//	type fields struct {
//		db *sqlx.DB
//	}
//	type args struct {
//		ctx  context.Context
//		user *internal.User
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			store := &Store{
//				db: tt.fields.db,
//			}
//			if err := store.AddUser(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
//				t.Errorf("AddUser() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestStore_CheckConnection(t *testing.T) {
//	type fields struct {
//		db *sqlx.DB
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			store := &Store{
//				db: tt.fields.db,
//			}
//			if err := store.CheckConnection(); (err != nil) != tt.wantErr {
//				t.Errorf("CheckConnection() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestStore_FindUserByLogin(t *testing.T) {
//	type fields struct {
//		db *sqlx.DB
//	}
//	type args struct {
//		ctx   context.Context
//		login string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *internal.User
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			store := &Store{
//				db: tt.fields.db,
//			}
//			got, err := store.FindUserByLogin(tt.args.ctx, tt.args.login)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("FindUserByLogin() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("FindUserByLogin() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestStore_GetOrder(t *testing.T) {
//	type fields struct {
//		db *sqlx.DB
//	}
//	type args struct {
//		ctx    context.Context
//		login  string
//		number int64
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *internal.Order
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			store := &Store{
//				db: tt.fields.db,
//			}
//			got, err := store.GetOrder(tt.args.ctx, tt.args.login, tt.args.number)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetOrder() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetOrder() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestStore_GetOrders(t *testing.T) {
//	type fields struct {
//		db *sqlx.DB
//	}
//	type args struct {
//		ctx    context.Context
//		userID uuid.UUID
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *[]internal.Order
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			store := &Store{
//				db: tt.fields.db,
//			}
//			got, err := store.GetOrders(tt.args.ctx, tt.args.userID)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetOrders() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetOrders() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestStore_GetOrdersNotProcessed(t *testing.T) {
//	type fields struct {
//		db *sqlx.DB
//	}
//	type args struct {
//		ctx context.Context
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *[]internal.Order
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			store := &Store{
//				db: tt.fields.db,
//			}
//			got, err := store.GetOrdersNotProcessed(tt.args.ctx)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetOrdersNotProcessed() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetOrdersNotProcessed() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestStore_GetUser(t *testing.T) {
//	type fields struct {
//		db *sqlx.DB
//	}
//	type args struct {
//		ctx context.Context
//		id  uuid.UUID
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *internal.User
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			store := &Store{
//				db: tt.fields.db,
//			}
//			got, err := store.GetUser(tt.args.ctx, tt.args.id)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestStore_GetWithdrawals(t *testing.T) {
//	type fields struct {
//		db *sqlx.DB
//	}
//	type args struct {
//		ctx    context.Context
//		userID uuid.UUID
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *[]internal.Withdraw
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			store := &Store{
//				db: tt.fields.db,
//			}
//			got, err := store.GetWithdrawals(tt.args.ctx, tt.args.userID)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetWithdrawals() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetWithdrawals() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestStore_SaveWithdrawal(t *testing.T) {
//	type fields struct {
//		db *sqlx.DB
//	}
//	type args struct {
//		ctx        context.Context
//		withdrawal *internal.Withdraw
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			store := &Store{
//				db: tt.fields.db,
//			}
//			if err := store.SaveWithdrawal(tt.args.ctx, tt.args.withdrawal); (err != nil) != tt.wantErr {
//				t.Errorf("SaveWithdrawal() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestStore_UpdateOrder(t *testing.T) {
//	type fields struct {
//		db *sqlx.DB
//	}
//	type args struct {
//		ctx   context.Context
//		order *internal.Order
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			store := &Store{
//				db: tt.fields.db,
//			}
//			if err := store.UpdateOrder(tt.args.ctx, tt.args.order); (err != nil) != tt.wantErr {
//				t.Errorf("UpdateOrder() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
