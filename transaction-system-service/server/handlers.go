package server

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Ragontar/transactionSystem/queueManager"
	"github.com/gorilla/mux"
)

const timeout = 5 * time.Minute

type Transaction struct {
	Amount int `json:"amount,omitempty"`
}

func AccountIncBalancePOST(w http.ResponseWriter, r *http.Request) {
	processRequest(queueManager.OperationInc, w, r)
}

func AccountDecBalancePOST(w http.ResponseWriter, r *http.Request) {
	processRequest(queueManager.OperationDec, w, r)
}

func processRequest(operation queueManager.Operation, w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user-id"]
	var err error

	tqm, ok := activeWorkers[userID]
	if !ok {
		tqm, err = queueManager.NewTransactionQueueManager(
			userID,
			historyStorage,
			accountStorage,
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		activeWorkers[userID] = tqm

		tqm.StartQueueProcessor()
	}

	var requestData Transaction
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	json.Unmarshal(b, &requestData)
	t := queueManager.Transaction{
		AccountID: tqm.AccountID,
		Operation: operation,
		Amount: requestData.Amount,
		Response: make(chan string),
		Date: time.Now(),
	}

	if err := tqm.Enqueue(&t); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var transactionStatus string
	// ? Context vs timer
	timer := time.NewTimer(timeout)

	select {
	case resp := <- t.Response:
		transactionStatus = resp
		timer.Stop()
	case <-timer.C:
		w.WriteHeader(http.StatusRequestTimeout)
		w.Write([]byte("Request timeout"))
		return
	}

	if transactionStatus == "rejected" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not enough balance"))
		return
	}
	if transactionStatus == "error" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ok"))
}
