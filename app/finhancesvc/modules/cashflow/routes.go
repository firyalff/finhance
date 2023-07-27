package cashflow

import (
	"finhancesvc/middlewares"
	"finhancesvc/shared"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (cashflowModule CashflowModule) RegisterRoutes(router *gin.RouterGroup) {
	const PATH_PREFIX = "/cashflows"
	routeGroup := router.Group(PATH_PREFIX)

	routeGroup.Use(middlewares.AuthMiddleware([]byte(cashflowModule.serverConfig.JWTSecret)))

	routeGroup.GET("", CashflowModuleInstance.listCashflowHandler)
	routeGroup.POST("", CashflowModuleInstance.createCashflowHandler)

}

type ListCashflowResponse struct {
	Id           string    `json:"id"`
	Amount       int       `json:"amount"`
	CashflowType string    `json:"type"`
	CreatedAt    time.Time `json:"created_at"`
}

func (cashflowModule CashflowModule) listCashflowHandler(ctx *gin.Context) {
	var payload shared.PagedRequest
	ctx.Bind(&payload)

	err := shared.Validator().Struct(payload)
	if err != nil {
		errBody := shared.GenerateErrorResponse("BAD_REQ", shared.ParseValidatorError(err))
		ctx.JSON(400, errBody)
		return
	}

	userID := ctx.GetString(middlewares.UserIDKey)

	cashflows, total, err := getUserCashflows(ctx, userID, payload.Limit, payload.Offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, shared.GenerateErrorResponse("INTERNALERR", nil))
		return
	}

	cashflowResp := make([]ListCashflowResponse, len(cashflows))

	for index, cashflow := range cashflows {
		cashflowResp[index] = ListCashflowResponse{
			Id:           cashflow.Id.String(),
			Amount:       int(cashflow.Amount),
			CashflowType: cashflow.CashflowType,
			CreatedAt:    cashflow.CreatedAt.Time,
		}
	}

	response := shared.GeneratePagedResponse(total, int(payload.Limit), int(payload.Offset), cashflowResp)

	ctx.JSON(200, response)
}

type cashflowCreationPayload struct {
}

func (cashflowModule CashflowModule) createCashflowHandler(ctx *gin.Context) {
	var payload cashflowCreationPayload
	ctx.Bind(&payload)

	ctx.JSON(201, gin.H{
		"message": "OK",
	})
}

func (cashflowModule CashflowModule) detailCashflowHandler(ctx *gin.Context) {
	var payload shared.DetailByIDRequest
	ctx.Bind(&payload)

	ctx.JSON(200, gin.H{
		"message": "OK",
	})
}
