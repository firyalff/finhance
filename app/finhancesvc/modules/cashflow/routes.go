package cashflow

import (
	"finhancesvc/middlewares"
	"finhancesvc/shared"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (cashflowModule CashflowModule) RegisterRoutes(router *gin.RouterGroup) {
	const CASHFLOW_PATH_PREFIX = "/cashflows"
	cashflowRouteGroup := router.Group(CASHFLOW_PATH_PREFIX)

	cashflowRouteGroup.Use(middlewares.AuthMiddleware([]byte(cashflowModule.serverConfig.JWTSecret)))

	cashflowRouteGroup.GET("", CashflowModuleInstance.listCashflowHandler)
	cashflowRouteGroup.GET("/:id", CashflowModuleInstance.detailCashflowHandler)
	cashflowRouteGroup.POST("", CashflowModuleInstance.createCashflowHandler)
	cashflowRouteGroup.PUT("/:id", CashflowModuleInstance.updateCashflowHandler)

	const CATEGORY_PATH_PREFIX = "/cashflow-categories"
	categoryRouteGroup := router.Group(CATEGORY_PATH_PREFIX)

	categoryRouteGroup.Use(middlewares.AuthMiddleware([]byte(cashflowModule.serverConfig.JWTSecret)))
	categoryRouteGroup.POST("", CashflowModuleInstance.createCashflowCategoryHandler)
	categoryRouteGroup.GET("", CashflowModuleInstance.listCashflowCategoryHandler)
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
	Amount             int     `json:"amount" validate:"required,gt=0"`
	Name               string  `json:"name" validate:"required,min=3"`
	Notes              string  `json:"notes"`
	CashflowCategoryID string  `json:"cashflow_category_id" validate:"required"`
	ProofDocumentUrl   *string `json:"proof_document_url"`
}

func (cashflowModule CashflowModule) createCashflowHandler(ctx *gin.Context) {
	var payload cashflowCreationPayload
	ctx.Bind(&payload)

	err := shared.Validator().Struct(payload)
	if err != nil {
		errBody := shared.GenerateErrorResponse("BAD_REQ", shared.ParseValidatorError(err))
		ctx.JSON(400, errBody)
		return
	}

	userID := ctx.GetString(middlewares.UserIDKey)

	err = createNewUserCashflow(ctx, userID, payload)
	if err != nil {
		errBody := shared.GenerateErrorResponse("INTERNALERR", err)
		ctx.JSON(500, errBody)
		return
	}

	ctx.JSON(201, gin.H{
		"message": "OK",
	})
}

type DetailedCashflowResponse struct {
	Id               string     `json:"id"`
	Amount           int        `json:"amount"`
	CashflowType     string     `json:"type"`
	ProofDocumentUrl *string    `json:"proof_document_url"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
}

func (cashflowModule CashflowModule) detailCashflowHandler(ctx *gin.Context) {
	var payload shared.DetailByIDRequest
	ctx.BindUri(&payload)

	userID := ctx.GetString(middlewares.UserIDKey)

	cashflow, err := getUserCashflowByID(ctx, userID, payload.ID)
	if err != nil {
		if err == shared.ErrNotFound {
			ctx.JSON(http.StatusNotFound, shared.GenerateErrorResponse("NOT_FOUND", nil))
		} else {
			ctx.JSON(http.StatusInternalServerError, shared.GenerateErrorResponse("INTERNALERR", nil))
		}
		return
	}
	resp := DetailedCashflowResponse{
		Id:               cashflow.Id.String(),
		Amount:           int(cashflow.Amount),
		CashflowType:     cashflow.CashflowType,
		ProofDocumentUrl: cashflow.ProofDocumentUrl,
		CreatedAt:        cashflow.CreatedAt.Time,
		UpdatedAt:        nil,
	}

	if &cashflow.UpdatedAt != nil {
		resp.UpdatedAt = &cashflow.UpdatedAt.Time
	}

	ctx.JSON(200, resp)
}

type cashflowUpdatePayload struct {
	Amount           int     `json:"amount" validate:"required,gt=0"`
	Name             string  `json:"name" validate:"required,min=3"`
	Notes            string  `json:"notes"`
	CashflowType     string  `json:"type" validate:"required,oneof=income expense"`
	ProofDocumentUrl *string `json:"proof_document_url"`
}

func (cashflowModule CashflowModule) updateCashflowHandler(ctx *gin.Context) {
	cashflowID := ctx.Param("id")
	var payload cashflowUpdatePayload
	ctx.Bind(&payload)

	err := shared.Validator().Struct(payload)
	if err != nil {
		errBody := shared.GenerateErrorResponse("BAD_REQ", shared.ParseValidatorError(err))
		ctx.JSON(400, errBody)
		return
	}

	userID := ctx.GetString(middlewares.UserIDKey)

	err = updateUserCashflow(ctx, userID, cashflowID, payload)
	if err != nil {
		if err == shared.ErrNotFound {
			ctx.JSON(http.StatusNotFound, shared.GenerateErrorResponse("NOT_FOUND", nil))
		} else {
			ctx.JSON(http.StatusInternalServerError, shared.GenerateErrorResponse("INTERNALERR", nil))
		}
		return
	}

	ctx.JSON(200, gin.H{
		"message": "OK",
	})
}

type cashflowCategoryCreatePayload struct {
	Name         string `json:"name"`
	CategoryType string `json:"type"`
}

func (cashflowModule CashflowModule) createCashflowCategoryHandler(ctx *gin.Context) {
	var payload cashflowCategoryCreatePayload
	ctx.Bind(&payload)

	err := shared.Validator().Struct(payload)
	if err != nil {
		errBody := shared.GenerateErrorResponse("BAD_REQ", shared.ParseValidatorError(err))
		ctx.JSON(400, errBody)
		return
	}

	userID := ctx.GetString(middlewares.UserIDKey)

	err = createUserCashflowCategory(ctx, userID, payload)
	if err != nil {
		if err == shared.ErrNotFound {
			ctx.JSON(http.StatusNotFound, shared.GenerateErrorResponse("NOT_FOUND", nil))
		} else if err == shared.ErrExist {
			ctx.JSON(http.StatusConflict, shared.GenerateErrorResponse("RESOURCE_EXIST", nil))
		} else {
			ctx.JSON(http.StatusInternalServerError, shared.GenerateErrorResponse("INTERNALERR", nil))
		}
		return
	}

	ctx.JSON(201, gin.H{
		"message": "OK",
	})
}

type ListCashflowCategoryResponse struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

func (cashflowModule CashflowModule) listCashflowCategoryHandler(ctx *gin.Context) {
	var payload shared.PagedRequest
	ctx.Bind(&payload)

	err := shared.Validator().Struct(payload)
	if err != nil {
		errBody := shared.GenerateErrorResponse("BAD_REQ", shared.ParseValidatorError(err))
		ctx.JSON(400, errBody)
		return
	}

	userID := ctx.GetString(middlewares.UserIDKey)

	cashflowCategories, total, err := getUserCashflowCategories(ctx, userID, payload.Limit, payload.Offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, shared.GenerateErrorResponse("INTERNALERR", nil))
		return
	}

	cashflowCategoriesResp := make([]ListCashflowCategoryResponse, len(cashflowCategories))

	for index, cashflowCategory := range cashflowCategories {
		cashflowCategoriesResp[index] = ListCashflowCategoryResponse{
			Id:        cashflowCategory.Id.String(),
			Name:      cashflowCategory.Name,
			Type:      cashflowCategory.CashflowCategoryType,
			CreatedAt: cashflowCategory.CreatedAt.Time,
		}
	}

	response := shared.GeneratePagedResponse(total, int(payload.Limit), int(payload.Offset), cashflowCategoriesResp)

	ctx.JSON(200, response)
}
