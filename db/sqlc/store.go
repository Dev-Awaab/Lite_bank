package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db quries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
		Queries: New(db),
	}
}

// exectx executes a function within a database transaction 
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction 
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Amount int64 
}

// TransferTxresult is thebresult of the transfer transaction
type TransferTxResult struct {
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account_id"`
	ToAccount Account `json:"to_account_id"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json: "to_entry"`
}

var txKey = struct{}{}


func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)

		fmt.Println(txName, "Create transfer")

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})

		if err != nil {
			return err
		}

	
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})

		if err != nil {
			return err
		}

	
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})

		if err != nil {
			return err
		}


		// get account -> update it's balance

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(context.Background(), q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(context.Background(), q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
			if err != nil {
				return err
			}

		}
		
		
		
		return nil
	})

	return result, err
}

// func addMoney( ctx context.Context, q *Queries, accountID1 int64, amount1 int64, accountID2 int64, amount2 int64)(
// account1 Account,
// account2 Account,
// err error
// )
// {
//  account1, err := q.AddAccountBalance(ctx, AddAccountBalanceParams{
// 	ID: accountID1,
// 	Amount: amount1
//   })

//   if err  != nil {
// 	return
//   }

//   account2, err := q.AddAccountBalance(ctx, AddAccountBalanceParams{
// 	ID: accountID2,
// 	Amount: amount2
//   })

//   return
// }

func addMoney(
	ctx context.Context, q *Queries, accountID1 int64, amount1 int64, accountID2 int64, amount2 int64,
) (account1 Account,
	account2 Account,
	err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})

	return
}