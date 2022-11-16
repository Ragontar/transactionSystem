package server

import "github.com/Ragontar/transactionSystem/queueManager"

var activeWorkers = make(map[string]*queueManager.TransactionQueueManager)
var accountStorage *queueManager.DBAccountStorage
var historyStorage *queueManager.DBHistoryStorage

// ! Two separate connection poolers to the same DB. Looks cringe, but intended.
func init() {
	var err error
	accountStorage, err = queueManager.NewDBAccountStorage()
	if err != nil {
		panic(err)
	}
	historyStorage, err = queueManager.NewDBHistoryStorage()
	if err != nil {
		panic(err)
	}
}