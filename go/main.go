package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
  "net/url"
  "os"
	"encoding/json"
	"bufio"
	"strconv"
	"errors"

	"github.com/gorilla/mux"
	"gopkg.in/gomail.v2"
)

const (
	APIKey        = "ae090038-01da-4b72-919a-1522c9f752eb" // coinmarketcap
	BTCID         = "1"
	UAHID         = "2824"
	HOST          = "smtp.gmail.com"
	PORT 					= "465"
	EmailUser     = "boyetsbilly@gmail.com" // Встановіть своє значення
	EmailPassword = "ikhpoypnpdjhjfyq" // Встановіть своє значення
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/rate", getRate).Methods("GET")
	router.HandleFunc("/api/subscribe", subscribeEmail).Methods("POST")
	router.HandleFunc("/api/sendEmails", sendEmails).Methods("POST")

	port := ":3000"
	fmt.Println("Server is listening on port", port)
	log.Fatal(http.ListenAndServe(port, router))
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

func getRate(w http.ResponseWriter, r *http.Request) {
	rate, err := getBitcoinExchangeRate()
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

func getBitcoinExchangeRate() (float64, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET","https://pro-api.coinmarketcap.com/v2/tools/price-conversion", nil)
  if err != nil {
    log.Print(err)
    os.Exit(1)
  }

  q := url.Values{}
  q.Add("amount", "1")
  q.Add("id", BTCID)
  q.Add("convert_id", UAHID)

  req.Header.Set("Accepts", "application/json")
  req.Header.Add("X-CMC_PRO_API_KEY", APIKey)
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

func subscribeEmail(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	if checkEmail(email) {
		w.WriteHeader(http.StatusConflict)
		return
	}

	err := addEmail(email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func checkEmail(email string) bool {
	subscriptions, err := readEmails()
	if err != nil {
		return false
	}

	for _, sub := range subscriptions {
		if strings.EqualFold(sub, email) {
			return true
		}
	}

	return false
}

func addEmail(email string) error {
	file, err := os.OpenFile("subscribed_emails.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, email)
	if err != nil {
		return err
	}

	return nil
}

func readEmails() ([]string, error) {
	file, err := os.Open("subscribed_emails.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	subscriptions := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		subscriptions = append(subscriptions, scanner.Text())
	}

	return subscriptions, nil
}

func sendEmails(w http.ResponseWriter, r *http.Request) {
	rate, err := getBitcoinExchangeRate()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	subscriptions, err := readEmails()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, email := range subscriptions {
		sendEmail(email, rate)
	}

	w.WriteHeader(http.StatusOK)
}

func sendEmail(email string, rate float64) {
	mailer, err := makeMailer()
  if err != nil {
    return
  }

  recipient := email
  if recipient == "" {
    return
  }

  m := gomail.NewMessage()
  m.SetHeader("From", EmailUser)
  m.SetHeader("To", email)
  m.SetHeader("Subject", "Курс BTC/UAH")
  m.SetBody("text/plain", fmt.Sprintf("Актуальний курс BTC до UAH: %f", rate))

  err = mailer.DialAndSend(m)
  if err != nil {
    return
  }
}

func makeMailer() (*gomail.Dialer, error) {
  smtp := HOST
  port, err := strconv.Atoi(PORT)
  if err != nil {
    fmt.Println("Error while parsing smtp port.", err)
    return nil, err
  }
  name := EmailUser
  pwd := EmailPassword
  if port == 0 || smtp == "" && name == "" {
    return nil, errors.New("Invalid mailer parameters")
  }

  return gomail.NewDialer(smtp, port, name, pwd), nil
}