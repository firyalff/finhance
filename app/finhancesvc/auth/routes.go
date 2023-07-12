package auth

import (
	"finhancesvc/shared"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func RegisterRoutes(router *gin.Engine) {
	router.POST("/login", authModuleInstance.loginHandler)
}

type loginValidation struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

func (authModule AuthModule) loginHandler(ctx *gin.Context) {
	var payload loginValidation
	ctx.Bind(&payload)

	err := shared.Validator().Struct(payload)
	if err != nil {
		errBody := shared.GenerateErrorResponse("BAD_REQ", shared.ParseValidatorError(err))
		ctx.JSON(400, errBody)
		return
	}

	userRecord, err := getUserByEmail(ctx, payload.Email)
	if err != nil {
		errBody := shared.GenerateErrorResponse("UNAUTHORIZED", nil)
		httpCode := http.StatusUnauthorized
		if err != pgx.ErrNoRows {
			errBody = shared.GenerateErrorResponse("INTERNAL_ERR", nil)
			httpCode = http.StatusInternalServerError
			log.Print(err)
		}

		ctx.JSON(httpCode, errBody)
		return
	}

	err = validateUserPassword(payload.Password, userRecord.Password)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusUnauthorized, shared.GenerateErrorResponse("UNAUTHORIZED", nil))
		return
	}

	tokenString, err := generateJWT(userRecord.Id.String(), authModule.serverConfig.JWTSecret, authModule.serverConfig.JWTExpireDayCount)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusUnauthorized, shared.GenerateErrorResponse("UNAUTHORIZED", nil))
		return
	}

	ctx.JSON(200, gin.H{
		"message": "OK",
		"token":   tokenString,
	})
}
