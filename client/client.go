package main

import (
	"context"
	pb "github.com/rima971/currency-convertor/currencyConvertor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("error occured while closing the connection: %s", err)
		}
	}(conn)

	currencyConvertorClient := pb.NewCurrencyConvertorServiceClient(conn)

	ctx := context.Background()

	req := &pb.CurrencyConvertorRequest{
		Money: &pb.Money{
			Value:    20,
			Currency: "INR",
		},
		TargetCurrency: "USD",
	}

	res, err := currencyConvertorClient.Convert(ctx, req)

	if err != nil {
		log.Fatalf("some error occurred at the server while registering user: %s", err)
	}
	log.Printf("response from the body: %s", res)
}
