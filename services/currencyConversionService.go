package services

import (
	"context"
	"encoding/json"
	"errors"
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
	if req.Money.Value < 0 {
		return nil, errors.New("bad request - money value cannot be negative")
	}
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
	var factors = map[string]*float64{}
	err = json.Unmarshal([]byte(byteValue), &factors)
	if err != nil {
		return nil, err
	}
	targetCurrencyConversionFactor := factors[req.TargetCurrency]
	givenCurrencyConversionFactor := factors[req.Money.Currency]
	if targetCurrencyConversionFactor == nil || givenCurrencyConversionFactor == nil {
		return nil, errors.New("bad request - unsupported currency provided")
	}
	calculatedMoneyValue := req.Money.Value * *targetCurrencyConversionFactor / *givenCurrencyConversionFactor
	calculatedMoneyValue = math.Round(calculatedMoneyValue*100) / 100

	return &pb.Money{
		Value:    calculatedMoneyValue,
		Currency: req.TargetCurrency,
	}, nil
}
