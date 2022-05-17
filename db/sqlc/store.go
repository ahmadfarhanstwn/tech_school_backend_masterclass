package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db: db,
	}
} 

func (s *Store) execTransactions(ctx context.Context, fn func(q *Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("err : %v, rollback error : %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

//parameter struct to perform transfer operation
type TransferTransactionsParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
}

//result struct to store data after performing transfer operation
type TransferTransactionsResult struct {
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}

var TxKey = struct{}{}

func (s *Store) TransferTransaction(ctx context.Context, arg TransferTransactionsParams) (TransferTransactionsResult,error) {
	var res TransferTransactionsResult

	err := s.execTransactions(ctx, func(q *Queries) error {
		var err error

		res.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID: arg.ToAccountId,
			Amount : arg.Amount,
		})
		if err != nil {
			return err
		}

		res.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount : -arg.Amount,
		})
		if err != nil {
			return err
		}

		res.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount : arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountId < arg.ToAccountId {
			res.FromAccount, res.ToAccount, err = addMoney(ctx, q, arg.FromAccountId, -arg.Amount, arg.ToAccountId, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			res.ToAccount, res.FromAccount, err = addMoney(ctx, q, arg.ToAccountId, arg.Amount, arg.FromAccountId, -arg.Amount)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return res, err
}

func addMoney (
	ctx context.Context,
	q *Queries,
	accountId1,
	amount1,
	accountId2,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: accountId1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: accountId2,
		Amount: amount2,
	})
	if err != nil {
		return
	}

	return
}