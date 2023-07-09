package main

import "github.com/mkideal/cli"

const appVersion = "v0.0.1"

type ServerConfig struct {
	cli.Helper
	DBURI     string `cli:"*dburi" usage:"Application database URI" dft:"$FINHANCE_DBURI"`
	JWTSecret string `cli:"*jwtsecret" usage:"Secret JWT token" dft:"$FINHANCE_JWT_SECRET"`
}
