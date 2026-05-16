package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"autopark/internal/application"
	"autopark/internal/infrastructure/config"
	httpapi "autopark/internal/infrastructure/http"
	"autopark/internal/infrastructure/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	for i := 0; i < 30; i++ {
		if err = pool.Ping(ctx); err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		log.Fatal(err)
	}

	store := postgres.NewStore(pool)
	crud := application.NewCRUDService(store)
	reports := application.NewReportService(store)
	auth := application.NewAuthService(store, cfg.AuthSecret)
	router := httpapi.NewRouter(crud, reports, auth)

	log.Printf("backend listening on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
