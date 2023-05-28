package services

import (
	"fmt"
	"net/http"
	"strconv"
	"errors"
  "github.com/NikitaLitvishko/Genesis-school-API-service/go/utils"
	"gopkg.in/gomail.v2"
)


func SendEmails(w http.ResponseWriter, r *http.Request) {
	rate, err := GetBitcoinExchangeRate()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	subscriptions, err := ReadEmails()
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
  m.SetHeader("From", utils.EmailUser)
  m.SetHeader("To", email)
  m.SetHeader("Subject", "Курс BTC/UAH")
  m.SetBody("text/plain", fmt.Sprintf("Актуальний курс BTC до UAH: %f", rate))

  err = mailer.DialAndSend(m)
  if err != nil {
    return
  }
}

func makeMailer() (*gomail.Dialer, error) {
  smtp := utils.HOST
  port, err := strconv.Atoi(utils.PORT)
  if err != nil {
    fmt.Println("Error while parsing smtp port.", err)
    return nil, err
  }
  name := utils.EmailUser
  pwd := utils.EmailPassword
  if port == 0 || smtp == "" && name == "" {
    return nil, errors.New("Invalid mailer parameters")
  }

  return gomail.NewDialer(smtp, port, name, pwd), nil
}