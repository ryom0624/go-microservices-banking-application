module auth

go 1.16

replace local.packages/lib => ./../lib

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/jmoiron/sqlx v1.3.4
	local.packages/lib v0.0.0-00010101000000-000000000000
)
