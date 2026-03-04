package postgres

import (
	"database/sql"
	"decisionMaker/config"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var GetInstance = sync.OnceValues(func() (*sql.DB, error) {
	return sql.Open("pgx", config.DBUrl)
})
