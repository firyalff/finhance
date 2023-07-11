package auth

import "github.com/jackc/pgx/v5/pgxpool"

type AuthModule struct {
	dbPool *pgxpool.Pool
}

var authModuleInstance AuthModule

func InitModule(dbPool *pgxpool.Pool) {
	authModuleInstance = AuthModule{
		dbPool: dbPool,
	}
}
