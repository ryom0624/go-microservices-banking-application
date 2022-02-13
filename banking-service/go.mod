module banking

go 1.16

replace local.packages/errs => ./../lib/errs

replace local.packages/logger => ./../lib/logger

require (
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gorilla/mux v1.8.0
	github.com/jmoiron/sqlx v1.3.4
	local.packages/errs v0.0.0-00010101000000-000000000000
	local.packages/logger v0.0.0-00010101000000-000000000000
)
