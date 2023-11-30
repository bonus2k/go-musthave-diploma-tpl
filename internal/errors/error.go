package errors

import (
	"errors"
	"fmt"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal"
	"time"
)

func RetryAfterError(f func() error) (err error) {
	attempts := 3
	sleep := 1 * time.Second
	for i := 1; ; i++ {
		if err = f(); err == nil {
			return
		}

		if i > attempts {
			break
		}

		time.Sleep(sleep)
		internal.Logf.Debugf("Attempt %d, retrying after error: %v", i, err)
		sleep = sleep + 2*time.Second

	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}

// storage errors
var ErrUserIsExist = errors.New("user is exist")
var ErrUserNotFound = errors.New("user not found")
var ErrOrderIsExistThisUser = errors.New("this order is exist the user")
var ErrOrderIsExistAnotherUser = errors.New("this order is exist another user")
var ErrNotEnoughAmount = errors.New("not enough amount")

// service errors
var ErrIllegalUserArgument = errors.New("illegal user argument")
var ErrIllegalOrder = errors.New("illegal order")
var ErrWrongAuth = errors.New("wrong authorization")

// auth error
var ErrInvalidValue = errors.New("invalid cookie value")
