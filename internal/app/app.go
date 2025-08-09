package app

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Aleksei-D/go-loyalty-system/internal/config"
	"github.com/Aleksei-D/go-loyalty-system/internal/domain/infrastructure"
	"github.com/Aleksei-D/go-loyalty-system/internal/logger"
	"github.com/Aleksei-D/go-loyalty-system/internal/orders_agent"
	"github.com/Aleksei-D/go-loyalty-system/internal/router"
	"github.com/Aleksei-D/go-loyalty-system/internal/service"
	"github.com/Aleksei-D/go-loyalty-system/pkg/datasource"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	db  *sql.DB
	cfg *config.Config
}

func NewApp(cfg *config.Config) (*App, error) {
	db, err := datasource.NewDatabase(*cfg.DatabaseUri)
	if err != nil {
		return nil, err
	}
	return &App{
		db:  db,
		cfg: cfg,
	}, nil
}

func (app *App) Run() error {
	serviceApp := service.NewService(
		infrastructure.NewPostgresBalanceRepository(app.db),
		infrastructure.NewPostgresOrderRepository(app.db),
		infrastructure.NewPostgresUserRepository(app.db),
		infrastructure.NewPostgresWithdrawalRepository(app.db),
	)

	r := router.NewRouter(serviceApp, *app.cfg.SecretKey)

	server := &http.Server{Addr: *app.cfg.ServerAddr, Handler: r}
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, time.Duration(*app.cfg.Wait)*time.Second)
		defer shutdownCancel()

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(context.DeadlineExceeded, shutdownCtx.Err()) {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	err := server.ListenAndServe()
	if err != nil && !errors.Is(http.ErrServerClosed, err) {
		logger.Log.Fatal("cannot start server", zap.Error(err))
	}

	orderAgent := orders_agent.NewOrdersAgent(service.NewOrderService(infrastructure.NewPostgresOrderRepository(app.db)), app.cfg)
	go orderAgent.Run(serverCtx)

	<-serverCtx.Done()

	return nil
}
