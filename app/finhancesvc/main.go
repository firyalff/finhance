package main

import (
	"context"
	"finhancesvc/auth"
	"finhancesvc/configs"
	"finhancesvc/drivers"
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

		DBPool, err = drivers.InitDBPool(ctx, cfg.DBURI)
		if err != nil {
			panic(err)
		}
		err = DBPool.Ping(ctx)
		if err != nil {
			panic(err)
		}

		auth.InitModule(DBPool, *cfg)

		router := drivers.InitRouting()

		err = drivers.InitLogger(cfg.SentryDSN, router)
		if err != nil {
			panic(err)
		}

		RegisterRoutes(router)
		auth.RegisterRoutes(router)

		return drivers.StartRouteServer(router)
	}))
}
