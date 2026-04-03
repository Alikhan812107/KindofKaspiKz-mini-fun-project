package app

import (
	"database/sql"
	"log"

	httpTransport "transaction-service/internal/transport/http"

	"transaction-service/internal/repository/postgres"
	"transaction-service/internal/usecase"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func Run() {
	// DB connection
	db, err := sql.Open("postgres", "postgres://postgres:password@localhost:5432/transaction_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// layers
	repo := postgres.NewTransactionRepo(db)
	uc := usecase.NewTransactionUsecase(repo)
	handler := httpTransport.NewHandler(uc)

	// server
	r := gin.Default()
	handler.RegisterRoutes(r)

	log.Println("Transaction Service running on :8081")
	r.Run(":8081")
}
