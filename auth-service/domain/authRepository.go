package domain

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"local.packages/lib/errs"
	"local.packages/lib/logger"
)

type AuthRepository interface {
	FindBy(string, string) (*Login, *errs.AppError)
	GenerateAndSaveRefreshTokenToStore(AuthToken) (string, *errs.AppError)
	RefreshTokenExists()
}

func (a AuthRepositoryDB) FindBy(username string, password string) (*Login, *errs.AppError) {
	query := `SELECT username, u.customer_id, role, group_concat(a.account_id) as account_numbers FROM users u
                  LEFT JOIN accounts a ON a.customer_id = u.customer_id
                WHERE username = ? and password = ?
                GROUP BY a.customer_id`

	var login Login
	err := a.client.Get(&login, query, username, password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewAuthenticationError("invalid credentials")
		} else {
			logger.Error("Error while verifying login request from database " + err.Error())
			return nil, errs.NewUnexpectedError("unexpected database error")
		}
	}
	return &login, nil
}

func (a AuthRepositoryDB) GenerateAndSaveRefreshTokenToStore(authToken AuthToken) (string, *errs.AppError) {
	refreshToken, appErr := authToken.newRefreshToken()
	if appErr != nil {
		return "", appErr
	}

	sqlInsert := `INSERT INTO refresh_token_store(refresh_token) VALUES(?)`
	_, err := a.client.Exec(sqlInsert, refreshToken)
	if err != nil {
		logger.Error("unexpected database error " + err.Error())
		return "", errs.NewUnexpectedError("Error unexpected database error")
	}

	return refreshToken, nil
}

func (a AuthRepositoryDB) RefreshTokenExists() {
	// TODO implement me
	panic("implement me")
}

type AuthRepositoryDB struct {
	client *sqlx.DB
}

func NewAuthRepositoryDB(client *sqlx.DB) AuthRepository {
	return AuthRepositoryDB{client: client}
}

type AuthRepositoryStub struct{}

func (a AuthRepositoryStub) FindBy(username string, password string) (*Login, *errs.AppError) {
	// TODO implement me
	panic("implement me")
}

func (a AuthRepositoryStub) GenerateAndSaveRefreshTokenToStore(accessToken AuthToken) (string, *errs.AppError) {
	// TODO implement me
	panic("implement me")
}

func (a AuthRepositoryStub) RefreshTokenExists() {
	// TODO implement me
	panic("implement me")
}

func NewAuthRepositoryStub() AuthRepository {
	return AuthRepositoryStub{}
}
