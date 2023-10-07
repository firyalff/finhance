package auth

import (
	"finhancesvc/configs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthModule struct {
	dbPool       *pgxpool.Pool
	serverConfig configs.ServerConfig
}

var AuthModuleInstance AuthModule

func InitModule(dbPool *pgxpool.Pool, serverConfig configs.ServerConfig, router *gin.RouterGroup) {
	AuthModuleInstance = AuthModule{
		dbPool:       dbPool,
		serverConfig: serverConfig,
	}

	registerRoutes(router)
}
