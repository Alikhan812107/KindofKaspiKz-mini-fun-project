package app

import (
	"database/sql"
	"log"
	"net"
	"os"

	_ "transaction-service/pkg/codec"
	grpcTransport "transaction-service/internal/transport/grpc"
	httpTransport "transaction-service/internal/transport/http"
	"transaction-service/internal/repository/postgres"
	"transaction-service/internal/usecase"
	pb "transaction-service/pkg/pb/payment"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	_ "github.com/lib/pq"
)

func Run() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:0000@localhost:5432/transaction_db?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	repo := postgres.NewPaymentRepo(db)
	uc := usecase.NewPaymentUsecase(repo)

	go runGRPC(uc)

	handler := httpTransport.NewHandler(uc)
	r := gin.Default()
	handler.RegisterRoutes(r)

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8081"
	}
	log.Printf("Payment Service HTTP running on :%s", httpPort)
	r.Run(":" + httpPort)
}

func runGRPC(uc *usecase.PaymentUsecase) {
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "9091"
	}

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen on grpc port: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPaymentServiceServer(s, grpcTransport.NewPaymentServer(uc))

	log.Printf("Payment Service gRPC running on :%s", grpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("grpc server failed: %v", err)
	}
}
