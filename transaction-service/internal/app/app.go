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
	db, err := sql.Open("postgres", "postgres://postgres:0000@localhost:5432/transaction_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	repo := postgres.NewPaymentRepo(db)
	uc := usecase.NewPaymentUsecase(repo)
	handler := httpTransport.NewHandler(uc)

	r := gin.Default()
	handler.RegisterRoutes(r)

	log.Println("Payment Service running on :8081")
	r.Run(":8081")
}
