package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"math/rand"
	"net/http"
	"time"

	"github.com/danilomartinelli/motor-credito/internal/response"
)

type CreditLimitResult struct {
	PaymentValue float64
	Installments int
	DebtValue    float64
}

func (app *application) status(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Status": "OK",
	}

	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) checkCredit(w http.ResponseWriter, r *http.Request) {
	debtID := chi.URLParam(r, "debtId")

	if debtID == "" {
		app.badRequest(w, r, fmt.Errorf("debtId is required"))
		return
	}

	creditResult, err := app.verifyCreditLimit(debtID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := map[string]interface{}{
		"DebtID":       debtID,
		"PaymentValue": fmt.Sprintf("%.2f", creditResult.PaymentValue),
		"Installments": creditResult.Installments,
		"DebtValue":    fmt.Sprintf("%.2f", creditResult.DebtValue),
	}

	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) verifyCreditLimit(debtID string) (CreditLimitResult, error) {
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)

	debtValue := float64(len(debtID)) * 200.0

	discountRate := 0.05 + rnd.Float64()*(0.20-0.05)
	newValueToPay := debtValue * (1 - discountRate)

	if newValueToPay >= debtValue {
		newValueToPay = debtValue * 0.95
	}

	installments := int(newValueToPay / 1000)
	if installments < 1 {
		installments = 1
	}
	if installments > 12 {
		installments = 12
	}

	return CreditLimitResult{
		PaymentValue: newValueToPay,
		Installments: installments,
		DebtValue:    debtValue,
	}, nil
}
