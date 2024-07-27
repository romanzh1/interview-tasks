package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	order_handler "gitlab.ozon.dev/go/classroom-12/students/week-8-workshop/internal/app/order"
	"gitlab.ozon.dev/go/classroom-12/students/week-8-workshop/internal/infra/repo/order"
	"gitlab.ozon.dev/go/classroom-12/students/week-8-workshop/internal/infra/shard_manager"
)

func main() {
	var (
		ctx       = context.Background()
		cfg       = newConfig()
		databases []*pgxpool.Pool
	)

	for _, dbConf := range cfg.databasePool {
		db, err := pgxpool.New(ctx, dbConf.DSN)
		if err != nil {
			log.Fatal(fmt.Errorf("failed to connect to db: %w", err))
		}
		defer db.Close()

		databases = append(databases, db)
	}

	var (
		sm = shard_manager.New(
			shard_manager.GetMurmur3ShardFn(len(databases)),
			databases,
		)
		orderRepo    = order.NewRepo(sm)
		orderHandler = order_handler.New(orderRepo)
	)

	http.HandleFunc("/create", orderHandler.Create)
	http.HandleFunc("/get", orderHandler.GetByID)
	http.HandleFunc("/list_by_user_id", orderHandler.ListByUserID)
	http.HandleFunc("/list_by_id", orderHandler.ListByID)

	_ = http.ListenAndServe(":8090", nil)
}
