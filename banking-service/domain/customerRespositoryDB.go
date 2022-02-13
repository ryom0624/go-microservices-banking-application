package domain

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"local.packages/lib/errs"
	"local.packages/lib/logger"
)

type CustomerRepositoryDB struct {
	client *sqlx.DB
}

func (d CustomerRepositoryDB) FindAll(status string) ([]Customer, *errs.AppError) {
	customers := make([]Customer, 0)

	var findAllSql string
	var err error
	if status == "" {
		findAllSql = `SELECT customer_id, name, city, zipcode, date_of_birth, status from customers`
		err = d.client.Select(&customers, findAllSql)
	} else {
		findAllSql = `SELECT customer_id, name, city, zipcode, date_of_birth, status from customers WHERE status = ?`
		err = d.client.Select(&customers, findAllSql, status)
	}

	if err != nil {
		// if err == sql.ErrNoRows {
		// 	return nil, errs.NewNotFoundError("customer does not exist")
		// }
		logger.Error("err find all mysql: " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected error")
	}
	return customers, nil
}

func (d CustomerRepositoryDB) ByID(id string) (*Customer, *errs.AppError) {
	customerSql := `SELECT customer_id, name, city, zipcode, date_of_birth, status from customers WHERE customer_id = ?`

	var c Customer
	err := d.client.Get(&c, customerSql, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("customer not found")
		} else {
			logger.Error("err find by id scan mysql: " + err.Error())
			return nil, errs.NewUnexpectedError("unexpected database error")
		}
	}

	return &c, nil
}

func NewCustomerRepositoryDB(dbClient *sqlx.DB) CustomerRepositoryDB {
	return CustomerRepositoryDB{client: dbClient}
}
