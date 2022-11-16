package queueManager

import (
	"context"
	"fmt"

	"github.com/Ragontar/transactionSystem/env"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DBHistoryStorage struct {
	db *pgxpool.Pool
}

func NewDBHistoryStorage() (*DBHistoryStorage, error) {
	s := &DBHistoryStorage{}
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

func (s *DBHistoryStorage) Save(Transaction) error {
	panic("NOT IMPLEMENTED")
}

func (s *DBHistoryStorage) Load() ([]Transaction, error) {
	panic("NOT IMPLEMENTED")
}
