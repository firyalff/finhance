package middlewares

import (
	"finhancesvc/shared"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var UserIDKey = "userID"

func AuthMiddleware(jwtSecret []byte) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		unauthorizedErrBody := shared.GenerateErrorResponse("UNAUTHORIZED", shared.ErrUnauthorized)

		stringToken, err := shared.ExtractAuthorizationToken(ctx)
		if err != nil {
			fmt.Println(err)
			ctx.JSON(http.StatusUnauthorized, unauthorizedErrBody)
			ctx.Abort()
		}

		jwtClaims, err := shared.ExtractJWTClaims(stringToken, jwtSecret)
		if err != nil {
			fmt.Println(err)

			ctx.JSON(http.StatusUnauthorized, unauthorizedErrBody)
			ctx.Abort()
		}

		ctx.Set(UserIDKey, jwtClaims.UserID)
		ctx.Next()
	}

}
