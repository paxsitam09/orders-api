package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

type App struct {
	router http.Handler
	rdb    *redis.Client
}

func New() *App {

	app := &App{
		router: loadRoutes(),
		rdb: redis.NewClient(&redis.Options{
			Addr:     "redis-10670.c340.ap-northeast-2-1.ec2.redns.redis-cloud.com:10670", // Host and Port
			Password: "ocknMiTUWBO52WZRaNBfm8lA2ByPrn5s",                                  // Password
			DB:       0,                                                                   // Default DB
		}),
	}
	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}

	// Test the Redis connection
	// ctx := context.Background()
	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)
	} else {
		fmt.Println("Successfully connected to Redis!")
	}

	defer func() {
		if err := a.rdb.Close(); err != nil {
			fmt.Println("failed to close redis", err)
		}
	}()

	fmt.Println("Starting Server")

	ch := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server : %w", err)
		}
		close(ch)
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return server.Shutdown(timeoutCtx)
	}
}
