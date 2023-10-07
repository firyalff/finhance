package statistic

import (
	"finhancesvc/configs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StatisticModule struct {
	dbPool       *pgxpool.Pool
	serverConfig configs.ServerConfig
}

var StatisticModuleInstance StatisticModule

func InitModule(dbPool *pgxpool.Pool, serverConfig configs.ServerConfig, router *gin.RouterGroup) {
	StatisticModuleInstance = StatisticModule{
		dbPool:       dbPool,
		serverConfig: serverConfig,
	}
	registerRoutes(router)

}
