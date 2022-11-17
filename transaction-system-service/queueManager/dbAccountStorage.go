package queueManager

import (
	"context"
	"fmt"
	"time"

	"github.com/Ragontar/transactionSystem/env"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DBAccountStorage struct {
	db *pgxpool.Pool
	timeout time.Duration
}

func NewDBAccountStorage() (*DBAccountStorage, error) {
	s := &DBAccountStorage{}
	s.timeout = 5 * time.Minute
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
	var balance int
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	row := s.db.QueryRow(ctx, SQL_GET_BALANCE_BY_USER_ID, userID)
	err := row.Scan(&balance) // ? closes automatically

	return balance, err
}

func (s *DBAccountStorage) UpdateBalance(userID string, value int) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	_, err := s.db.Query(ctx, SQL_UPDATE_BALANCE_BY_USER_ID, value, userID)
	return err
}

func (s *DBAccountStorage) GetAccountID(userID string) (string, error) {
	var account_id string
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	row := s.db.QueryRow(ctx, SQL_SELECT_ACCOUNT_ID_BY_USER_ID, userID)
	err := row.Scan(&account_id)

	return account_id, err
}