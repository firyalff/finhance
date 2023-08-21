package statistic

import (
	"finhancesvc/configs"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StatisticModule struct {
	dbPool       *pgxpool.Pool
	serverConfig configs.ServerConfig
}

var StatisticModuleInstance StatisticModule

func InitModule(dbPool *pgxpool.Pool, serverConfig configs.ServerConfig) {
	StatisticModuleInstance = StatisticModule{
		dbPool:       dbPool,
		serverConfig: serverConfig,
	}
}
