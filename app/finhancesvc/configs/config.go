package configs

import "github.com/mkideal/cli"

const appVersion = "v0.0.1"

type ServerConfig struct {
	cli.Helper
	DBURI              string `cli:"*dburi" usage:"Application database URI" dft:"$FINHANCESVC_DBURI"`
	JWTSecret          string `cli:"*jwtsecret" usage:"Secret JWT token" dft:"$FINHANCESVC_JWT_SECRET"`
	JWTExpireDayCount  int    `cli:"*jwtexpdaycount" usage:"JWT expiration in days" dft:"$FINHANCESVC_JWT_EXP_DAY"`
	SentryDSN          string `cli:"*sentrydsn" usage:"Sentry DSN" dft:"$FINHANCESVC_SENTRY_DSN"`
	SMTPHostURL        string `cli:"*smtphosturl" usage:"" dft:"$FINHANCESVC_SMTP_HOST_URL"`
	SMTPHostPORT       string `cli:"*smtphostport" usage:"" dft:"$FINHANCESVC_SMTP_HOST_PORT"`
	SMTPUsername       string `cli:"*smtpusername" usage:"" dft:"$FINHANCESVC_SMTP_USERNAME"`
	SMTPPassword       string `cli:"*smtppassword" usage:"" dft:"$FINHANCESVC_SMTP_PASSWORD"`
	EmailDefaultSender string `cli:"*emaildefaultsender" usage:"" dft:"$FINHANCESVC_EMAIL_SENDER"`
}
