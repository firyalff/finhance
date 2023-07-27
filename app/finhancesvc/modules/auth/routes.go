package auth

import (
	"finhancesvc/shared"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (authModule AuthModule) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/login", AuthModuleInstance.loginHandler)
	router.POST("/register", AuthModuleInstance.registerHandler)
	router.GET("/account-activation", AuthModuleInstance.accountActivationHandler)
}

type loginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

func (authModule AuthModule) loginHandler(ctx *gin.Context) {
	var payload loginPayload
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

type registerPayload struct {
	FirstName            string `json:"first_name" validate:"required"`
	LastName             string `json:"last_name" validate:"required"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=8,max=64"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=8,max=64,eqfield=Password"`
}

func (authModule AuthModule) registerHandler(ctx *gin.Context) {
	var payload registerPayload
	ctx.Bind(&payload)

	err := shared.Validator().Struct(payload)
	if err != nil {
		errBody := shared.GenerateErrorResponse("BAD_REQ", shared.ParseValidatorError(err))
		ctx.JSON(400, errBody)
		return
	}

	err = validateUniqueEmailRegistration(ctx, payload.Email)
	if err != nil {
		errBody := shared.GenerateErrorResponse("INTERNALERR", err)
		httpCode := http.StatusInternalServerError
		if err == shared.ErrExist {
			errBody = shared.GenerateErrorResponse("RESOURCE_EXIST", err)
			httpCode = http.StatusConflict
		}

		ctx.JSON(httpCode, errBody)
		return
	}

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}

	err = registerNewUser(ctx, payload, scheme+"://"+ctx.Request.Host)
	if err != nil {
		errBody := shared.GenerateErrorResponse("INTERNALERR", err)
		ctx.JSON(http.StatusInternalServerError, errBody)
		return
	}

	ctx.JSON(201, gin.H{
		"message": "OK",
	})
	return
}

type accountActivationPayload struct {
	ActivationToken string `form:"activation_token" validate:"required"`
}

func (authModule AuthModule) accountActivationHandler(ctx *gin.Context) {
	var payload accountActivationPayload

	ctx.BindQuery(&payload)

	err := shared.Validator().Struct(payload)
	if err != nil {
		errBody := shared.GenerateErrorResponse("BAD_REQ", shared.ParseValidatorError(err))
		ctx.JSON(400, errBody)
		return
	}

	err = activateAccount(ctx, payload.ActivationToken)
	if err != nil {
		errBody := shared.GenerateErrorResponse("INTERNALERR", err)
		httpStatus := http.StatusInternalServerError
		if err == shared.ErrNotFound {
			errBody = shared.GenerateErrorResponse("NOT_FOUND", err)
			httpStatus = http.StatusNotFound
		}
		ctx.JSON(httpStatus, errBody)

		return
	}

	ctx.JSON(200, gin.H{
		"message": "OK",
	})
	return
}
