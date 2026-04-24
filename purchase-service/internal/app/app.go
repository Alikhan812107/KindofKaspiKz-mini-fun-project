package app

import (
	"database/sql"
	"log"
	"net"
	"os"

	grpcTransport "purchase-service/internal/transport/grpc"
	httpTransport "purchase-service/internal/transport/http"
	"purchase-service/internal/repository"
	"purchase-service/internal/repository/postgres"
	"purchase-service/internal/usecase"
	orderpb "purchase-service/pkg/pb/order"
	paymentpb "purchase-service/pkg/pb/payment"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Run() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:0000@localhost:5432/purchase_db?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	paymentAddr := os.Getenv("PAYMENT_GRPC_ADDR")
	if paymentAddr == "" {
		paymentAddr = "localhost:9091"
	}

	conn, err := grpc.NewClient(paymentAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to payment service: %v", err)
	}

	paymentClient := paymentpb.NewPaymentServiceClient(conn)

	repo := postgres.NewPurchaseRepo(db)
	uc := usecase.NewPurchaseUsecase(repo, paymentClient)

	go runGRPC(repo)

	handler := httpTransport.NewHandler(uc)
	r := gin.Default()
	handler.RegisterRoutes(r)

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}
	log.Printf("Purchase Service HTTP running on :%s", httpPort)
	r.Run(":" + httpPort)
}

func runGRPC(repo repository.PurchaseRepository) {
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "9090"
	}

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen on grpc port: %v", err)
	}

	s := grpc.NewServer()
	orderpb.RegisterOrderServiceServer(s, grpcTransport.NewOrderServer(repo))

	log.Printf("Purchase Service gRPC running on :%s", grpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("grpc server failed: %v", err)
	}
}
