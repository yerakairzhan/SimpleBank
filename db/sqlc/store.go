package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

type SQLStore struct {
	db      *sql.DB
	Queries *Queries
}

// AddAccountBalance implements Store.
func (store *SQLStore) AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (Account, error) {
	panic("unimplemented")
}

// CreateAccount implements Store.
func (store *SQLStore) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	panic("unimplemented")
}

// CreateEntry implements Store.
func (store *SQLStore) CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error) {
	panic("unimplemented")
}

// CreateTransfer implements Store.
func (store *SQLStore) CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error) {
	panic("unimplemented")
}

// DeleteAccount implements Store.
func (store *SQLStore) DeleteAccount(ctx context.Context, id int64) error {
	panic("unimplemented")
}

// GetAccount implements Store.
func (store *SQLStore) GetAccount(ctx context.Context, id int64) (Account, error) {
	return store.Queries.GetAccount(ctx, id)
}

// GetAccountForUpdate implements Store.
func (store *SQLStore) GetAccountForUpdate(ctx context.Context, id int64) (Account, error) {
	panic("unimplemented")
}

// GetEntry implements Store.
func (store *SQLStore) GetEntry(ctx context.Context, id int64) (Entry, error) {
	panic("unimplementedEntry")
}

// GetTransfer implements Store.
func (store *SQLStore) GetTransfer(ctx context.Context, id int64) (Transfer, error) {
	panic("unimplemented HERE")
}

// ListAccount implements Store.
func (store *SQLStore) ListAccount(ctx context.Context, arg ListAccountParams) ([]Account, error) {
	panic("unimplemented")
}

// ListEntries implements Store.
func (store *SQLStore) ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error) {
	panic("unimplemented")
}

// ListTransfers implements Store.
func (store *SQLStore) ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfer, error) {
	panic("unimplemented")
}

// UpdateAccount implements Store.
func (store *SQLStore) UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error) {
	panic("unimplemented")
}

// NewStore creates a new Store object
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// GetAccount retrieves an account by ID
/*func (store *Store) GetAccount(ctx context.Context, id int64) (Account, error) {
	return store.Queries.GetAccount(ctx, id)
}

//listAccounts should be there but i deleted it

func (store *Store) ListAccounts(ctx context.Context, arg ListAccountParams) ([]Account, error) {
	return store.Queries.ListAccount(ctx, arg)
}

// GetEntry retrieves an entry by ID
func (store *Store) GetEntry(ctx context.Context, id int64) (Entry, error) {
	return store.Queries.GetEntry(ctx, id)
}

// GetTransfer retrieves a transfer by ID
func (store *Store) GetTransfer(ctx context.Context, id int64) (Transfer, error) {
	return store.Queries.GetTransfer(ctx, id)
}*/

// TransferTxParams contains the input parameters for a transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult contains the result of a transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to another within a transaction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Create the transfer
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// Create the entries
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// Get account balances and update them accordingly
		if arg.FromAccountID < arg.ToAccountID {
			// Lock the "from" account and update balance
			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount,
			})
			if err != nil {
				return err
			}

			// Lock the "to" account and update balance
			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			})
			if err != nil {
				return err
			}
		} else {
			// Lock the "to" account and update balance first
			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			})
			if err != nil {
				return err
			}

			// Lock the "from" account and update balance
			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}
