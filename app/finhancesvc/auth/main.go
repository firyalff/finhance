package auth

import (
	"finhancesvc/configs"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthModule struct {
	dbPool       *pgxpool.Pool
	serverConfig configs.ServerConfig
}

var authModuleInstance AuthModule

func InitModule(dbPool *pgxpool.Pool, serverConfig configs.ServerConfig) {
	authModuleInstance = AuthModule{
		dbPool:       dbPool,
		serverConfig: serverConfig,
	}
}
