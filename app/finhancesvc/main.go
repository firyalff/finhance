package main

import (
	"context"
	"finhancesvc/configs"
	"finhancesvc/drivers"
	"finhancesvc/modules/auth"
	"finhancesvc/modules/cashflow"
	"finhancesvc/modules/statistic"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mkideal/cli"
)

const AppVersion = "v0.0.1"

var DBPool *pgxpool.Pool

func main() {
	os.Exit(cli.Run(new(configs.ServerConfig), func(cliCtx *cli.Context) (err error) {
		ctx := context.Background()

		cfg := cliCtx.Argv().(*configs.ServerConfig)

		router := drivers.InitRouting()

		err = drivers.InitLogger(cfg.SentryDSN, router)
		if err != nil {
			panic(err)
		}

		DBPool, err = drivers.InitDBPool(ctx, cfg.DBURI)
		if err != nil {
			panic(err)
		}
		err = DBPool.Ping(ctx)
		if err != nil {
			panic(err)
		}

		defer DBPool.Close()

		baseRouter := router.Group("")
		v1Router := router.Group("/v1")

		RegisterRoutes(baseRouter)

		auth.InitModule(DBPool, *cfg, v1Router)
		cashflow.InitModule(DBPool, *cfg, v1Router)
		statistic.InitModule(DBPool, *cfg, v1Router)

		return drivers.StartRouteServer(router)
	}))
}
