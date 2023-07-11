package drivers

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDBPool(ctx context.Context, dbURI string) (pool *pgxpool.Pool, err error) {
	return pgxpool.New(ctx, dbURI)
}
