package queueManager

import (
	"errors"
	"time"
)

type Operation string

const (
	OperationInc Operation = "+"
	OperationDec Operation = "-"
)

type HistoryStorage interface {
	Save(Transaction) error
	Load(userID string) ([]Transaction, error)
}

type AccountStorage interface {
	GetBalance(userID string) (int, error)
	UpdateBalance(userID string, value int) error
	GetAccountID(userID string) (string, error)
}

type Transaction struct {
	AccountID string
	Operation Operation
	Amount    int
	Date      time.Time
	Response  chan string
}

type TransactionQueueManager struct {
	UserID             string
	AccountID          string
	Balance            int
	TransactionQueue   chan *Transaction
	TransactionHistory []Transaction

	historyStorage HistoryStorage
	accountStorage AccountStorage
}

func NewTransactionQueueManager(userID string, hs HistoryStorage, as AccountStorage) (*TransactionQueueManager, error) {
	tqm := TransactionQueueManager{UserID: userID, historyStorage: hs, accountStorage: as}

	var err error
	tqm.TransactionHistory, err = tqm.historyStorage.Load(userID)
	if err != nil {
		return nil, err
	}
	tqm.AccountID, err = tqm.accountStorage.GetAccountID(tqm.UserID)
	if err != nil {
		return nil, err
	}
	tqm.Balance, err = tqm.accountStorage.GetBalance(tqm.UserID)
	tqm.TransactionQueue = make(chan *Transaction, 100)

	return &tqm, err
}

func (tqm *TransactionQueueManager) ExecuteNext() bool {
	t := <-tqm.TransactionQueue

	if t.Operation == OperationDec && t.Amount > tqm.Balance {
		tqm.rejectTransaction(t)
		return true
	}
	err := tqm.acceptTransaction(t)
	if err != nil {
		println(err)
	}

	return true
}

func (tqm *TransactionQueueManager) Enqueue(t *Transaction) error {
	if t == nil {
		return errors.New("transaction is nil")
	}
	if t.Amount <= 0 {
		return errors.New("incorrect amount")
	}
	if t.Operation != OperationInc && t.Operation != OperationDec {
		return errors.New("bad operation")
	}

	tqm.TransactionQueue <- t

	return nil
}

func (tqm *TransactionQueueManager) acceptTransaction(t *Transaction) error {
	oldBalance := tqm.Balance
	if t.Operation == OperationDec {
		tqm.Balance = tqm.Balance - t.Amount
	}
	if t.Operation == OperationInc {
		tqm.Balance = tqm.Balance + t.Amount
	}

	err := tqm.accountStorage.UpdateBalance(tqm.UserID, tqm.Balance)
	if err != nil {
		tqm.Balance = oldBalance
		t.Response <- "error"
		return err
	}
	err = tqm.historyStorage.Save(*t)
	if err != nil {
		// ? UpdateBalance error case?
		tqm.Balance = oldBalance
		tqm.accountStorage.UpdateBalance(tqm.UserID, tqm.Balance)
		// ?
		t.Response <- "error"
		return err
	}
	tqm.TransactionHistory = append(tqm.TransactionHistory, *t)

	t.Response <- "accepted"

	return nil
}

func (tqm *TransactionQueueManager) rejectTransaction(t *Transaction) error {
	t.Response <- "rejected"
	//cringe
	return nil
}

func (tqm *TransactionQueueManager) StartQueueProcessor() {
	go func() {
		for {
			tqm.ExecuteNext()
		}
	}()
}
