package agent

import (
	"context"
	"github.com/Aleksei-D/go-loyalty-system/internal/config"
	"github.com/Aleksei-D/go-loyalty-system/internal/logger"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
	"github.com/Aleksei-D/go-loyalty-system/internal/service"
	"go.uber.org/zap"
	"time"
)

const orderChunkSizeDefault = 5

type OrdersAgent struct {
	orderService   *service.OrderService
	httpClient     *StatusUpdaterClient
	config         *config.Config
	orderChunkSize uint
}

func NewOrdersAgent(orderService *service.OrderService, config *config.Config) *OrdersAgent {
	return &OrdersAgent{
		orderService:   orderService,
		config:         config,
		httpClient:     NewClientAgent(*config.AccrualSystemAddress),
		orderChunkSize: orderChunkSizeDefault,
	}
}

func (o *OrdersAgent) Run(ctx context.Context) {
	doneCh := make(chan struct{})
	defer close(doneCh)
	pollTicker := time.NewTicker(time.Duration(*o.config.PollInterval) * time.Second)
	defer pollTicker.Stop()

	errorCh := make(chan error)
	defer close(errorCh)

	orderNumberCh := o.orderGenerator(ctx, doneCh, errorCh, pollTicker)
	updateDataOrderCh := o.OrdersStatusGenerator(ctx, doneCh, orderNumberCh, errorCh)
	go o.updateStatusOrder(ctx, doneCh, updateDataOrderCh, errorCh)

	for err := range errorCh {
		logger.Log.Warn(err.Error(), zap.Error(err))
	}
}

func (o *OrdersAgent) updateStatusOrder(ctx context.Context, doneCh chan struct{}, updateDataOrderCh <-chan *models.Order, errorCh chan<- error) {
	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Worker: Context done, exiting.")
			return
		case <-doneCh:
			return
		case order := <-updateDataOrderCh:
			err := o.orderService.UpdateStatus(ctx, order)
			if err != nil {
				errorCh <- err
			}
		}
	}
}

func (o *OrdersAgent) orderGenerator(ctx context.Context, doneCh chan struct{}, errorCh chan<- error, pollTicker *time.Ticker) <-chan *models.Order {
	orderCh := make(chan *models.Order, o.orderChunkSize)
	go func() {
		defer close(orderCh)
	newOrderLoop:
		for {
			select {
			case <-ctx.Done():
				logger.Log.Info("Worker: Context done, exiting.")
				return
			case <-doneCh:
				return
			case <-pollTicker.C:
				orders, err := o.orderService.GetNotAcceptedOrderNumbers(ctx, o.orderChunkSize, *o.config.UpdateTimeout)
				if err != nil {
					errorCh <- err
					continue newOrderLoop
				}

				for _, order := range orders {
					orderCh <- order
				}
			}
		}
	}()
	return orderCh
}

func (o *OrdersAgent) OrdersStatusGenerator(ctx context.Context, doneCh chan struct{}, orderCh <-chan *models.Order, errorCh chan<- error) <-chan *models.Order {
	orderNewStatusCh := make(chan *models.Order, o.orderChunkSize)
	go func() {
		defer close(orderNewStatusCh)
		for {
			select {
			case <-ctx.Done():
				logger.Log.Info("Worker: Context done, exiting.")
				return
			case <-doneCh:
				return
			default:
				for w := 1; w <= int(*o.config.RateLimit); w++ {
					go o.getOrdersStatus(ctx, doneCh, orderCh, orderNewStatusCh, errorCh)
				}
			}
		}
	}()
	return orderNewStatusCh
}

func (o *OrdersAgent) getOrdersStatus(ctx context.Context, doneCh chan struct{}, orderCh <-chan *models.Order, orderNewStatusCh chan<- *models.Order, errorCh chan<- error) {
	select {
	case <-ctx.Done():
		logger.Log.Info("Worker: Context done, exiting.")
		return
	case <-doneCh:
		return
	case newOrder := <-orderCh:
		order, err := o.httpClient.getOrderStatus(newOrder.Number)
		if err != nil {
			errorCh <- err
			return
		}
		orderNewStatusCh <- order.ToOrder()
	}
}
