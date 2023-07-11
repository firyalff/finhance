package main

import (
	"context"
	"finhancesvc/auth"
	"finhancesvc/drivers"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mkideal/cli"
)

const AppVersion = "v0.0.1"

var DBPool *pgxpool.Pool

func main() {
	os.Exit(cli.Run(new(ServerConfig), func(cliCtx *cli.Context) (err error) {
		ctx := context.Background()

		router := drivers.InitRouting()
		RegisterRoutes(router)
		auth.RegisterRoutes(router)

		DBPool, err = drivers.InitDBPool(ctx, "")
		if err != nil {
			panic(err)
		}

		return drivers.StartRouteServer(router)
	}))
}
