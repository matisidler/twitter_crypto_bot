package main

import (
	"cryptobot/price"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Credentials struct {
	AccessToken       string
	ConsumerKey       string
	ConsumerSecret    string
	AccessTokenSecret string
}

type Prices struct {
	CurrentPrice   float64
	Change         float64
	PercentChange  float64
	HighPriceOfDay float64
	LowPriceOfDay  float64
}

func getClient(creds *Credentials) (*twitter.Client, error) {
	config := oauth1.NewConfig("ConsumerKey", "ConsumerSecret")
	token := oauth1.NewToken("AccessToken", "AccessTokenSecret")

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	_, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func ParsingNumbers(data *price.FinnhubResponse, crypto string) (*Prices, string) {
	m := Prices{}
	percent, err := json.Number.Float64(data.PercentChange)
	if err != nil {
		log.Fatalf("error percent: %v", err)
	}
	m.PercentChange = percent

	current, err := json.Number.Float64(data.CurrentPrice)
	if err != nil {
		log.Fatalf("error percent: %v", err)
	}
	m.CurrentPrice = current

	change, err := json.Number.Float64(data.Change)
	if err != nil {
		log.Fatalf("error percent: %v", err)
	}
	m.Change = change

	high, err := json.Number.Float64(data.HighPriceOfDay)
	if err != nil {
		log.Fatalf("error percent: %v", err)
	}
	m.HighPriceOfDay = high

	low, err := json.Number.Float64(data.LowPriceOfDay)
	if err != nil {
		log.Fatalf("error percent: %v", err)
	}
	m.LowPriceOfDay = low

	return &m, crypto

}

func Testing(data *Prices, crypto string) {
	client, err := getClient(&Credentials{})
	if err != nil {
		fmt.Println(err)
	}
	var str string
	var text string
	switch crypto {
	case "BINANCE:BTCUSDT":
		str = "🤑 BITCOIN UPDATE #BTCUSD 🤑"
	case "BINANCE:ETHUSDT":
		str = "🤑 ETHEREUM UPDATE #ETHUSD 🤑"
	case "BINANCE:ADAUSDT":
		str = "🤑 CARDANO UPDATE #ADAUSD 🤑"
	default:
		str = "error"
	}
	if crypto == "BINANCE:ADAUSDT" {
		text = fmt.Sprintf(`🤑 %s 🤑

💵 CURRENT PRICE: %v%.3f 💵
		
️‍🔥 DAY CHANGE: %v%.3f️‍ 🔥
		
📈 PERCENTAGE CHANGE: %.2f%v 📈
	
⬆️ HIGHEST PRICE OF DAY: %v%.3f ⬆️
		
⬇️ LOWEST PRICE OF DAY: %v%.3f ⬇️`, str, "$", data.CurrentPrice, "$", data.Change, data.PercentChange, "%", "$", data.HighPriceOfDay, "$", data.LowPriceOfDay)
	} else {
		text = fmt.Sprintf(`🤑 %s 🤑

💵 CURRENT PRICE: %v%.f 💵

️‍🔥 DAY CHANGE: %v%.f️‍ 🔥

📈 PERCENTAGE CHANGE: %.2f%v 📈

⬆️ HIGHEST PRICE OF DAY: %v%.f ⬆️

⬇️ LOWEST PRICE OF DAY: %v%.f ⬇️`, str, "$", data.CurrentPrice, "$", data.Change, data.PercentChange, "%", "$", data.HighPriceOfDay, "$", data.LowPriceOfDay)
	}
	_, _, err = client.Statuses.Update(text, nil)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

}
func main() {
	for {
		dataBtc, crypto := price.GetPrice("BINANCE:BTCUSDT")
		dataParsedBtc, crypto := ParsingNumbers(dataBtc, crypto)
		Testing(dataParsedBtc, crypto)
		time.Sleep(10 * time.Minute)

		dataBtc, crypto = price.GetPrice("BINANCE:ETHUSDT")
		dataParsedBtc, crypto = ParsingNumbers(dataBtc, crypto)
		Testing(dataParsedBtc, crypto)
		time.Sleep(10 * time.Minute)

		dataBtc, crypto = price.GetPrice("BINANCE:ADAUSDT")
		dataParsedBtc, crypto = ParsingNumbers(dataBtc, crypto)
		Testing(dataParsedBtc, crypto)
		time.Sleep(10 * time.Minute)
	}

}
