package domain

import (
	"banking/errs"
	"banking/logger"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type AccountRepositoryDB struct {
	client *sqlx.DB
}

func (d AccountRepositoryDB) Save(a Account) (*Account, *errs.AppError) {
	sqlInsert := "INSERT INTO accounts(customer_id, opening_date, account_type, amount, status) VALUES(?,?,?,?,?)"

	result, err := d.client.Exec(sqlInsert, a.CustomerID, a.OpeningDate, a.AccountType, a.Amount, a.Status)
	if err != nil {
		logger.Error("error while creating new account " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected error from database")
	}

	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("error while getting last insert id for new account " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected error from database")
	}

	a.AccountID = strconv.FormatInt(id, 10)
	return &a, nil
}

func (d AccountRepositoryDB) SaveTransaction(t Transaction) (*Transaction, *errs.AppError) {
	tx, err := d.client.Begin()
	if err != nil {
		logger.Error("error while start transaction " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected error from database")
	}

	sqlInsert := "INSERT INTO transactions(account_id, amount, transaction_type, transaction_date) VALUES(?,?,?,?)"
	result, err := tx.Exec(sqlInsert, t.AccountID, t.Amount, t.TransactionType, t.TransactionDate)
	if err != nil {
		tx.Rollback()
		logger.Error("error while creating new transaction " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected error from database")
	}

	var updateAmountSql string
	if t.IsWithdrawal() {
		updateAmountSql = "UPDATE accounts SET amount = amount - ? WHERE account_id = ?"
	} else {
		updateAmountSql = "UPDATE accounts SET amount = amount + ? WHERE account_id = ?"
	}

	_, err = tx.Exec(updateAmountSql, t.Amount, t.AccountID)
	if err != nil {
		tx.Rollback()
		logger.Error("error while creating new transaction " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected error from database")
	}

	transactionID, err := result.LastInsertId()
	if err != nil {
		logger.Error("error while getting last insert id for new transaction " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected error from database")
	}
	account, appError := d.Find(t.AccountID)
	if appError != nil {
		return nil, errs.NewNotFoundError("account not found")
	}

	t.ID = strconv.FormatInt(transactionID, 10)
	t.Amount = account.Amount

	tx.Commit()

	return &t, nil
}

func (d AccountRepositoryDB) Find(accountID string) (*Account, *errs.AppError) {
	accountSql := `SELECT account_id, customer_id,  opening_date, account_type, amount, status from accounts WHERE account_id = ?`

	var a Account
	err := d.client.Get(&a, accountSql, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("account not found")
		} else {
			logger.Error("err find account scan mysql: " + err.Error())
			return nil, errs.NewUnexpectedError("unexpected database error")
		}
	}

	return &a, nil
}

func NewAccountRepositoryDB(dbClient *sqlx.DB) AccountRepositoryDB {
	return AccountRepositoryDB{client: dbClient}
}
