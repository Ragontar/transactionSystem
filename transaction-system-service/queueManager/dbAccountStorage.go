package queueManager

import (
	"context"
	"fmt"

	"github.com/Ragontar/transactionSystem/env"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DBAccountStorage struct {
	db *pgxpool.Pool
}

func NewDBAccountStorage() (*DBAccountStorage, error) {
	s := &DBAccountStorage{}
	err := s.init()
	return s, err
}

func (s *DBAccountStorage) init() error {
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

func (s *DBAccountStorage) GetBalance(userID string) (int, error) {
	panic("NOT IMPLEMENTED")
}

func (s *DBAccountStorage) UpdateBalance(userID string, value int) error {
	panic("NOT IMPLEMENTED")
}