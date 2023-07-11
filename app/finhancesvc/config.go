package main

import "github.com/mkideal/cli"

const appVersion = "v0.0.1"

type ServerConfig struct {
	cli.Helper
	DBURI     string `cli:"*dburi" usage:"Application database URI" dft:"$FINHANCESVC_DBURI"`
	JWTSecret string `cli:"*jwtsecret" usage:"Secret JWT token" dft:"$FINHANCESVC_JWT_SECRET"`
}
