package shared

import (
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGRepository struct {
	Transaction pgx.Tx
	Pool        *pgxpool.Pool
}
