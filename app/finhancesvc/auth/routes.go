package auth

import (
	"finhancesvc/shared"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.POST("/login", loginHandler)
}

type loginValidation struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

func loginHandler(ctx *gin.Context) {
	var payload loginValidation
	ctx.Bind(&payload)

	err := shared.Validator().Struct(payload)
	if err != nil {
		errBody := shared.GenerateErrorResponse("BAD_REQ", shared.ParseValidatorError(err))
		ctx.JSON(400, errBody)
		return
	}

	ctx.JSON(200, gin.H{
		"message": "OK",
	})
}
