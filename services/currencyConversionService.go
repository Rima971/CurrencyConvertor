package services

import (
	"context"
	"encoding/json"
	pb "github.com/rima971/currency-convertor/currencyConvertor"
	"io/ioutil"
	"log"
	"math"
	"os"
)

type CurrencyConversionService struct {
	pb.UnimplementedCurrencyConvertorServiceServer
	Path string
}

func NewService(path string) *CurrencyConversionService {
	return &CurrencyConversionService{Path: path}
}

func (c *CurrencyConversionService) Convert(ctx context.Context, req *pb.CurrencyConvertorRequest) (*pb.Money, error) {
	jsonFile, err := os.Open(c.Path)
	if err != nil {
		return nil, err
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			log.Fatalf("error occurred while closing json file: %s", err)
		}
	}(jsonFile)

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var factors = map[string]float64{}
	err = json.Unmarshal([]byte(byteValue), &factors)
	if err != nil {
		return nil, err
	}

	calculatedMoneyValue := req.Money.Value * factors[req.TargetCurrency] / factors[req.Money.Currency]
	calculatedMoneyValue = math.Round(calculatedMoneyValue*100) / 100

	return &pb.Money{
		Value:    calculatedMoneyValue,
		Currency: req.TargetCurrency,
	}, nil
}
