package services

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"net/url"
	"log"
	"os"
  "github.com/NikitaLitvishko/Genesis-school-API-service/go/utils"
)

func GetRate(w http.ResponseWriter, r *http.Request) {
	rate, err := GetBitcoinExchangeRate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(rate)
	if err != nil {
		return
	}
}

type RateResponse struct {
	Status struct {
		Timestamp    string `json:"timestamp"`
		ErrorCode    int       `json:"error_code"`
		ErrorMessage string       `json:"error_message"`
		Elapsed      int       `json:"elapsed"`
		CreditCount  int       `json:"credit_count"`
		Notice       string       `json:"notice"`
	} `json:"status"`
	Data struct {
		ID          int       `json:"id"`
		Symbol      string    `json:"symbol"`
		Name        string    `json:"name"`
		Amount      int       `json:"amount"`
		LastUpdated string `json:"last_updated"`
		Quote       struct {
			UAH struct {
				Price       float64   `json:"price"`
				LastUpdated string `json:"last_updated"`
			} `json:"2824"`
		} `json:"quote"`
	} `json:"data"`
}

func GetBitcoinExchangeRate() (float64, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET","https://pro-api.coinmarketcap.com/v2/tools/price-conversion", nil)
  if err != nil {
    log.Print(err)
    os.Exit(1)
  }

  q := url.Values{}
  q.Add("amount", "1")
  q.Add("id", utils.BTCID)
  q.Add("convert_id", utils.UAHID)

  req.Header.Set("Accepts", "application/json")
  req.Header.Add("X-CMC_PRO_API_KEY", utils.APIKey)
  req.URL.RawQuery = q.Encode()


  resp, err := client.Do(req);
  if err != nil {
    fmt.Println("Error sending request to server")
    os.Exit(1)
  }
  respBody, _ := ioutil.ReadAll(resp.Body)

	var rateResponse RateResponse
	err = json.Unmarshal(respBody, &rateResponse)

	return rateResponse.Data.Quote.UAH.Price, nil
}