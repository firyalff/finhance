package auth

import (
	"finhancesvc/configs"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthModule struct {
	dbPool       *pgxpool.Pool
	serverConfig configs.ServerConfig
}

var AuthModuleInstance AuthModule

func InitModule(dbPool *pgxpool.Pool, serverConfig configs.ServerConfig) {
	AuthModuleInstance = AuthModule{
		dbPool:       dbPool,
		serverConfig: serverConfig,
	}
}
