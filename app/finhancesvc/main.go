package main

import (
	"finhancesvc/drivers"
	"os"

	"github.com/mkideal/cli"
)

const AppVersion = "v0.0.1"

func main() {
	os.Exit(cli.Run(new(ServerConfig), func(ctx *cli.Context) error {
		router := drivers.InitRouting()
		RegisterRoutes(router)

		return drivers.StartRouteServer(router)
	}))
}
