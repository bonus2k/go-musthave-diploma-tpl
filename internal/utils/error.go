package utils

import (
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
