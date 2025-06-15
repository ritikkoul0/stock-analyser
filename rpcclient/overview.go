package overview

import (
	"context"
	"log"
	"stock-analyser/kafka"
	"time"

	pb "github.com/ritikkoul0/stock-rpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Overview() {
	// Dial the gRPC server first
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewStockAnalyserClient(conn)

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		symbol, err := kafka.ReadMessage[string](ctx)
		if err != nil {
			continue
		}
		log.Printf("Received stock symbol from Kafka: %s", symbol)

		req := &pb.Stockrequest{Symbol: symbol}

		resp, err := client.GetStockDetail(ctx, req)
		if err != nil {
			log.Printf("Error calling GetStockDetail for %s: %v", symbol, err)
			continue
		}
		log.Printf("Server responded with: %s", resp.Message)
	}
}
