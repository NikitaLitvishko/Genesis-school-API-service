package services

import (
	"fmt"
	"net/http"
	"strings"
  "os"
	"bufio"
)


func SubscribeEmail(w http.ResponseWriter, r *http.Request) {
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
	subscriptions, err := ReadEmails()
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
	file, err := os.OpenFile("./subscribed_emails.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
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

func ReadEmails() ([]string, error) {
	file, err := os.Open("./subscribed_emails.txt")
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