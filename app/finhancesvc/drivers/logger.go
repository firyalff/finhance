package drivers

import (
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func InitLogger(sentryDSN string, serverRouter *gin.Engine) (err error) {
	err = sentry.Init(sentry.ClientOptions{
		Dsn:           sentryDSN,
		EnableTracing: false,
	})
	if err != nil {
		return err
	}

	serverRouter.Use(sentrygin.New(sentrygin.Options{}))

	return
}
