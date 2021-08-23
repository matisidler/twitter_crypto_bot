package price

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go"
)

type FinnhubResponse struct {
	CurrentPrice   json.Number `json:"c"`
	Change         json.Number `json:"d"`
	PercentChange  json.Number `json:"dp"`
	HighPriceOfDay json.Number `json:"h"`
	LowPriceOfDay  json.Number `json:"l"`
}

func GetPrice(crypto string) (*FinnhubResponse, string) {

	cfg := finnhub.NewConfiguration()
	cfg.AddDefaultHeader("X-Finnhub-Token", "TOKEN")
	finnhubClient := finnhub.NewAPIClient(cfg).DefaultApi

	ahora := time.Now()
	btcvolumeInfo, _, err := finnhubClient.CryptoCandles(context.Background(), crypto, "D", ahora.Add(24*time.Hour*-1).Unix(), time.Now().Unix())

	if err != nil {
		fmt.Println(err)
	}
	btcvolume := btcvolumeInfo.V

	ticker := fmt.Sprintf("http://finnhub.io/api/v1/quote?symbol=%s&token=TOKEN", crypto)
	res, err := http.Get(ticker)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	data := FinnhubResponse{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		log.Fatalf(err.Error())
		return nil, ""
	}

	fmt.Printf("----------BINANCE:BTCUSDT-----------\n, current price: %s\n, change: %s\n, percent change: %s\n, highest price of day: %s\n, lowest price of day: %s\n, volume: %v\n \n---------------------------\n", data.CurrentPrice, data.Change, data.PercentChange, data.HighPriceOfDay, data.LowPriceOfDay, btcvolume)
	return &data, crypto
}
