package queueManager

import (
	"context"
	"fmt"
	"time"

	"github.com/Ragontar/transactionSystem/env"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DBHistoryStorage struct {
	db      *pgxpool.Pool
	timeout time.Duration
}

func NewDBHistoryStorage() (*DBHistoryStorage, error) {
	s := &DBHistoryStorage{}
	s.timeout = 5 * time.Minute
	err := s.init()
	return s, err
}

func (s *DBHistoryStorage) init() error {
	if s.db != nil {
		return nil
	}
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s/%s",
		env.DB_USER,
		env.DB_PASSWORD,
		env.DB_ADDR,
		env.DB_DATABASE,
	)

	var err error
	s.db, err = pgxpool.Connect(context.TODO(), dsn)

	return err
}

func (s *DBHistoryStorage) Save(t Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	_, err := s.db.Query(
		ctx,
		SQL_INSERT_TRANSACTION_HISTORY_ENTRY,
		uuid.NewString(),
		t.AccountID,
		t.Operation,
		t.Amount,
		t.Date,
	)

	return err
}

func (s *DBHistoryStorage) Load(userID string) ([]Transaction, error) {
	th := []Transaction{}
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	rows, _ := s.db.Query(ctx, SQL_SELECT_TRANSACTION_HISTORY_ENTRIES_BY_USER_ID, userID)
	defer rows.Close()
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.AccountID, &t.Operation, &t.Amount, &t.Date)
		if err != nil {
			return nil, err
		}
		th = append(th, t)
	}

	return th, nil
}
