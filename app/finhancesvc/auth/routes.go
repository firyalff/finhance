package auth

import (
	"finhancesvc/shared"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

	userRecord, err := getUserByCredentials(ctx, payload)
	if err != nil {
		if err == shared.ErrNotFound || err == shared.ErrUnauthorized {
			ctx.JSON(http.StatusUnauthorized, shared.GenerateErrorResponse("UNAUTHORIZED", nil))
		} else {
			ctx.JSON(http.StatusInternalServerError, shared.GenerateErrorResponse("INTERNALERR", nil))
		}
		return
	}

	tokenString, err := generateAthenticationToken(userRecord.Id.String(), authModule.serverConfig.JWTSecret, authModule.serverConfig.JWTExpireDayCount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, shared.GenerateErrorResponse("INTERNALERR", nil))
		return
	}

	ctx.JSON(200, gin.H{
		"token": tokenString,
	})
}

type registerValidation struct {
	FirstName            string
	LastName             string
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=8,max=64"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=8,max=64"`
}

func (authModule AuthModule) registerHandler(ctx *gin.Context) {
	var payload registerValidation
	ctx.Bind(&payload)

	err := shared.Validator().Struct(payload)
	if err != nil {
		errBody := shared.GenerateErrorResponse("BAD_REQ", shared.ParseValidatorError(err))
		ctx.JSON(400, errBody)
		return
	}

	totalUser, err := countUserByEmail(ctx, payload.Email)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, shared.GenerateErrorResponse("INTERNAL_ERR", nil))
		return
	}

	if totalUser > 0 {
		ctx.JSON(http.StatusConflict, shared.GenerateErrorResponse("RESOURCE_EXIST", nil))
		return
	}

	ctx.JSON(201, gin.H{
		"message": "OK",
	})
}
