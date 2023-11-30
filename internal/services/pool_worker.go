package services

import (
	"errors"
	"fmt"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal"
	"github.com/bonus2k/go-musthave-diploma-tpl/internal/interfaces/clients"
	"go.uber.org/zap"
	"time"
)

const timeoutErrTooManyRequests = 2 * time.Minute

type PoolWorker struct {
	client      *clients.ClientAccrual
	serviceUser *UserService
	orderIn     chan string
	Err         chan error
}

func NewPoolWorker(client *clients.ClientAccrual, serviceUser *UserService) *PoolWorker {
	ordersIn := make(chan string, 10)
	err := make(chan error)
	return &PoolWorker{client: client, serviceUser: serviceUser, orderIn: ordersIn, Err: err}
}

func (p *PoolWorker) StarIntegration(countWorker int, requestTime *time.Ticker) {

	pauses := make([]chan struct{}, 0)
	for i := 0; i < countWorker; i++ {
		name := i
		pause := make(chan struct{})
		p.worker(name, pause)
		pauses = append(pauses, pause)
	}

	go func() {
		for range requestTime.C {
			numbers, err := p.serviceUser.GetOrdersNotProcessed()
			if err != nil {
				internal.Log.Error("error start workers of integration")
				break
			}
			for _, n := range numbers {
				number := n
				p.orderIn <- number
			}
		}
	}()

	for err := range p.Err {
		internal.Log.Error("error integration", zap.Error(err))
		if errors.Is(err, clients.ErrTooManyRequests) {
			go func() {
				for _, pause := range pauses {
					ch := pause
					ch <- struct{}{}
				}
			}()
		}
	}
}

func (p *PoolWorker) worker(nameWorker int, pause chan struct{}) {
	go func() {
		defer close(pause)
		for {
			select {
			case order := <-p.orderIn:
				internal.Logf.Debugf("worker %d, order %s send request to accrual services", nameWorker, order)
				accrual, err := p.client.CheckAccrual(order)
				if err != nil {
					p.Err <- fmt.Errorf("error worker %d %w", nameWorker, err)
					break
				}
				internal.Logf.Debugf("worker %d, save %v in order", nameWorker, accrual)
				err = p.serviceUser.UpdateOrder(accrual)
				if err != nil {
					p.Err <- fmt.Errorf("error worker %d %w", nameWorker, err)
				}
			case <-pause:
				internal.Logf.Debugf("worker %d do pause", nameWorker)
				time.Sleep(timeoutErrTooManyRequests)
			}
		}
	}()
}
