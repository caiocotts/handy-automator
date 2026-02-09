package postgres

import (
	"database/sql"
	"os"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var GetInstance = sync.OnceValues(func() (*sql.DB, error) {
	return sql.Open("pgx", os.Getenv("DB_URL"))
})
