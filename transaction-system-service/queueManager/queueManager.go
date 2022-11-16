package queueManager

import (
	"errors"
	"sync"
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
}

type Transaction struct {
	Operation Operation
	Amount    int
	Date      time.Time
	Response  chan string
}

type TransactionQueueManager struct {
	UserID             string
	Balance            int
	mu                 sync.Mutex
	TransactionQueue   []*Transaction
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
	tqm.Balance, err = tqm.accountStorage.GetBalance(tqm.UserID)

	return &tqm, err
}

func (tqm *TransactionQueueManager) ExecuteNext() bool {
	if len(tqm.TransactionQueue) == 0 {
		return false
	}
	tqm.mu.Lock()
	t := tqm.TransactionQueue[0]
	tqm.TransactionQueue = tqm.TransactionQueue[1:]
	tqm.mu.Unlock()

	if t.Operation == OperationDec && t.Amount > tqm.Balance {
		tqm.rejectTransaction(t)
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

	tqm.mu.Lock()
	tqm.TransactionQueue = append(tqm.TransactionQueue, t)
	tqm.mu.Unlock()

	return nil
}

func (tqm *TransactionQueueManager) acceptTransaction(t *Transaction) error {
	if t.Operation == OperationDec {
		tqm.Balance = tqm.Balance - t.Amount
	}
	if t.Operation == OperationInc{
		tqm.Balance = tqm.Balance + t.Amount
	}

	err := tqm.accountStorage.UpdateBalance(tqm.UserID, tqm.Balance)
	if err != nil {
		t.Response <- "error"
		return err
	}
	err = tqm.historyStorage.Save(*t)
	if err != nil {
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