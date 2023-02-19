package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/http"

	route "lamoda-test/api/routes"
	config "lamoda-test/internal/config"
	"lamoda-test/pkg/logging"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"
)

type App struct {
	cfg        *config.Config
	router     *gin.Engine
	httpServer *http.Server
	pgClient   *sql.DB
}

func NewApp(ctx context.Context, config *config.Config) (*App, error) {
	logging.GetLogger(ctx).Info("router initializing")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPass, config.DBName,
	)
	println(dsn)

	pgClient, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = pgClient.Ping(); err != nil {
		return nil, err
	}

	router := route.NewRouter(pgClient)

	logging.GetLogger(ctx).Info("swagger docs initializing")

	return &App{
		cfg:      config,
		router:   router,
		pgClient: pgClient,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	logging.GetLogger(ctx).Info("application initialized and started")
	defer func() {
		if err := a.pgClient.Close(); err != nil {
			logging.GetLogger(ctx).Error(err)
		}
	}()

	grp, ctx := errgroup.WithContext(ctx)

	grp.Go(func() error {
		return a.startHTTP(ctx)
	})

	return grp.Wait()
}

func (a *App) startHTTP(ctx context.Context) error {
	logging.GetLogger(ctx).WithFields(map[string]interface{}{
		"IP":   a.cfg.IP,
		"PORT": a.cfg.Port,
	})

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", a.cfg.IP, a.cfg.Port))
	if err != nil {
		logging.GetLogger(ctx).WithError(err).Fatal("failed to create listener")
	}

	handler := a.router

	a.httpServer = &http.Server{
		Handler: handler,
	}

	logging.GetLogger(ctx).Info("application completely initialized and started")

	if err = a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logging.GetLogger(ctx).Warning("server shutdown")
		default:
			logging.GetLogger(ctx).Fatal(err)
		}
	}

	err = a.httpServer.Shutdown(context.Background())
	if err != nil {
		logging.GetLogger(ctx).Error(err)
	}

	return err
}
