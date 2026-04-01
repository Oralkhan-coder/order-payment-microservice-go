package main

import (
	"context"
	"log"
)

func main() {
	ctx := context.Background()
	cfg := pkg.Config{
		Database: "done_db",
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "postgres",
	}

	if err := pkg.RunMigrations(cfg); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	db, err := pkg.NewDB(ctx, cfg)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer db.Pool.Close()

	orderRepo := repository.NewOrderRepository(db.Pool)
	orderService := service.NewOrderService(orderRepo)
	server := internalHttp.NewSimpleServer(orderService)

	log.Println("starting the server on :8080")

	server.Run(ctx)
}
