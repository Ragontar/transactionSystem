package queueManager

import (
	"errors"
	"sync"
)

type Operation string

const (
	OperationInc Operation = "+"
	OperationDec Operation = "-"
)

type TransactionQueueManager struct {
	UserID           string
	Balance          int
	mu               sync.Mutex
	TransactionQueue []Transaction
}

// TODO TBD
func (tqm *TransactionQueueManager) ExecuteNext() bool {
	if len(tqm.TransactionQueue) == 0 {
		return false
	}
	tqm.mu.Lock()
	t := tqm.TransactionQueue[0]
	tqm.TransactionQueue = tqm.TransactionQueue[1:]
	tqm.mu.Unlock()

	println(t)
	//TODO do some execution. Cases: accepted/rejected

	panic("NOT IMPLEMENTED!!!!")
	return true
}

func (tqm *TransactionQueueManager) Enqueue(t Transaction) error {
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

type Transaction struct {
	Operation Operation
	Amount    int
}

type HistoryStorage interface {
	Save() error
	Load() error
}
