package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal/loggers"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal/models/dto"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal/models/entities"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal/repositories"
	"github.com/bonus2k/go-musthave-diploma-tpl/cmd/gophermart/internal/utils"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

type UserService struct {
	db *repositories.Store
}

func NewUserService(storage *repositories.Store) *UserService {
	return &UserService{db: storage}
}

func (us *UserService) CreateNewUser(ctx context.Context, user *dto.User) (*entities.User, error) {
	if isIllegalUserArgument(user) {
		return nil, ErrIllegalUserArgument
	}
	loggers.Log.Debug("create user")
	password, err := utils.SignPassword(user.Pass)
	if err != nil {
		return nil, fmt.Errorf("can't create password, %w", err)
	}
	entity := &entities.User{ID: uuid.New(), CreateAt: time.Now(), Login: user.Login, Password: password}
	if err := us.db.AddUser(ctx, entity); err != nil {
		return nil, err
	}
	return entity, nil
}

func (us *UserService) AddOrder(ctx context.Context, id string, orderID string) error {
	luna, ok := isLuna(orderID)
	if !ok {
		return ErrIllegalOrder
	}
	userID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	order := &entities.Order{ID: uuid.New(), CreateAt: time.Now(), Number: int64(luna), Status: entities.NEW, UserID: userID}
	_, err = us.db.AddOrder(ctx, order)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) LoginUser(ctx context.Context, user dto.User) (*uuid.UUID, error) {
	login, err := us.db.FindUserByLogin(ctx, user.Login)
	if err != nil {
		return nil, err
	}
	ok, err := utils.CheckPassword(user.Pass, login.Password)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrWrongAuth
	}
	return &login.ID, nil
}

func (us *UserService) GetOrders(ctx context.Context, id string) (*[]entities.Order, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	orders, err := us.db.GetOrders(ctx, userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func isLuna(order string) (int, bool) {
	number, err := strconv.Atoi(order)
	if err != nil {
		return number, false
	}
	return number, (number%10+checksum(number/10))%10 == 0
}

func checksum(number int) int {
	var luhn int
	for i := 0; number > 0; i++ {
		cur := number % 10
		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}

func isIllegalUserArgument(user *dto.User) bool {
	trimLogin := strings.TrimSpace(user.Login)
	trimPassword := strings.TrimSpace(user.Pass)
	if len(trimLogin) == 0 || len(trimPassword) == 0 {
		loggers.Logf.Errorf("parameters is wrong; login=%s, password=%s",
			user.Login, user.Pass)
		return true
	}
	return false
}

var ErrIllegalUserArgument = errors.New("illegal user argument")
var ErrIllegalOrder = errors.New("illegal order")
var ErrWrongAuth = errors.New("wrong authorization")
