package app

import (
	"database/sql"
	"log"

	httpTransport "purchase-service/internal/transport/http"
	"purchase-service/internal/repository/postgres"
	"purchase-service/internal/usecase"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func Run() {
	db, err := sql.Open("postgres", "postgres://postgres:0000@localhost:5432/purchase_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	repo := postgres.NewPurchaseRepo(db)
	uc := usecase.NewPurchaseUsecase(repo)
	handler := httpTransport.NewHandler(uc)

	r := gin.Default()
	handler.RegisterRoutes(r)

	log.Println("Purchase Service running on :8080")
	r.Run(":8080")
}
