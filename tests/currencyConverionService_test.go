package main

import (
	"context"
	pb "github.com/rima971/currency-convertor/currencyConvertor"
	"github.com/rima971/currency-convertor/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

func server(ctx context.Context) (pb.CurrencyConvertorServiceClient, func()) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer()

	c := services.NewService("/Users/madhurima.poddar_ftc/GolandProjects/currencyConvertor/constants/conversionFactor_test.json")
	pb.RegisterCurrencyConvertorServiceServer(baseServer, c)
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("error serving server: %s", err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %s", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
	}

	client := pb.NewCurrencyConvertorServiceClient(conn)

	return client, closer
}

func TestCurrencyConvertor_convert(t *testing.T) {
	ctx := context.Background()

	client, closer := server(ctx)
	defer closer()

	type expectation struct {
		out *pb.Money
		err error
	}
	tests := map[string]struct {
		in       *pb.CurrencyConvertorRequest
		expected expectation
	}{
		"success: 0 USD sent - expect 0 INR": {
			in: &pb.CurrencyConvertorRequest{
				TargetCurrency: "INR",
				Money: &pb.Money{
					Value:    0,
					Currency: "USD",
				},
			},
			expected: expectation{
				out: &pb.Money{
					Value:    0,
					Currency: "INR",
				},
				err: nil,
			},
		},
		"success: 20 USD sent - expect 1666.67 INR": {
			in: &pb.CurrencyConvertorRequest{
				TargetCurrency: "INR",
				Money: &pb.Money{
					Value:    20,
					Currency: "USD",
				},
			},
			expected: expectation{
				out: &pb.Money{
					Value:    1666.67,
					Currency: "INR",
				},
				err: nil,
			},
		},
		"success: 20 INR sent - expect 1666.67 USD": {
			in: &pb.CurrencyConvertorRequest{
				TargetCurrency: "USD",
				Money: &pb.Money{
					Value:    20,
					Currency: "INR",
				},
			},
			expected: expectation{
				out: &pb.Money{
					Value:    0.24,
					Currency: "USD",
				},
				err: nil,
			},
		},
		"success: 12.34 USD sent - expect 1028.33 INR": {
			in: &pb.CurrencyConvertorRequest{
				TargetCurrency: "INR",
				Money: &pb.Money{
					Value:    12.34,
					Currency: "USD",
				},
			},
			expected: expectation{
				out: &pb.Money{
					Value:    1028.33,
					Currency: "INR",
				},
				err: nil,
			},
		},
		"success: 12.34 INR sent - expect 0.15 USD": {
			in: &pb.CurrencyConvertorRequest{
				TargetCurrency: "USD",
				Money: &pb.Money{
					Value:    12.34,
					Currency: "INR",
				},
			},
			expected: expectation{
				out: &pb.Money{
					Value:    0.15,
					Currency: "USD",
				},
				err: nil,
			},
		},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.Convert(ctx, tt.in)
			if err != nil {
				if tt.expected.err == nil || tt.expected.err.Error() != err.Error() {
					t.Errorf("Err -> \nWant: %q\nGot: %q\n", tt.expected.err, err)
				}
			} else {
				if tt.expected.out.Value != out.Value ||
					tt.expected.out.Currency != out.Currency {
					t.Errorf("Out -> \nWant: %q\nGot : %q", tt.expected.out, out)
				}
			}

		})
	}
}
