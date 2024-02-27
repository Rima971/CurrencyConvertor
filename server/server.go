package main

import (
	pb "github.com/rima971/currency-convertor/currencyConvertor"
	"github.com/rima971/currency-convertor/services"
	"google.golang.org/grpc"
	"log"
	"net"
)

const port = ":8090"

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("cannot create listener: %v", err)
	}
	serverRegistrar := grpc.NewServer()

	c := services.NewService("constants/conversionFactor.json")

	pb.RegisterCurrencyConvertorServiceServer(serverRegistrar, c)
	err = serverRegistrar.Serve(lis)
	if err != nil {
		log.Fatalf("impossible to serve: %v", err)
	}
}
