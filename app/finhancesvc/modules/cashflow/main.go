package cashflow

import (
	"finhancesvc/configs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CashflowModule struct {
	dbPool       *pgxpool.Pool
	serverConfig configs.ServerConfig
}

var CashflowModuleInstance CashflowModule

func InitModule(dbPool *pgxpool.Pool, serverConfig configs.ServerConfig, router *gin.RouterGroup) {
	CashflowModuleInstance = CashflowModule{
		dbPool:       dbPool,
		serverConfig: serverConfig,
	}

	registerRoutes(router)
}
