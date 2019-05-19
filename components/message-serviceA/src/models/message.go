package models

import (
	"errors"
	"strings"
)

type MessagePayload struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

func(message MessagePayload) Validate() error {

	if message.Amount > 100000000 && message.Amount < -100000000 {
		return errors.New("Invalid amount")
	}

	if strings.ToUpper(message.Currency) != "EUR" {
		return errors.New("Unsuported currency")
	}
	return nil
}

func(message *MessagePayload) ConvertAmount(){
	minimalCurrency := message.Amount * 100
	message.Amount = float64(minimalCurrency)
}