package statistic

import (
	"finhancesvc/middlewares"
	"finhancesvc/shared"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (statisticwModule StatisticModule) RegisterRoutes(router *gin.RouterGroup) {
	const STATISTICS_PATH_PREFIX = "/statistics"
	statisticRouteGroup := router.Group(STATISTICS_PATH_PREFIX)

	statisticRouteGroup.Use(middlewares.AuthMiddleware([]byte(statisticwModule.serverConfig.JWTSecret)))

	statisticRouteGroup.GET("/income-expenses", StatisticModuleInstance.incomeExpenseStatisticHandler)
}

type BasicStatisticQuery struct {
	PeriodStartTime string `form:"period_start_time" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	PeriodEndTime   string `form:"period_end_time" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

type IncomeExpenseStatisticResponse struct {
	PeriodStartTime time.Time `json:"period_start_time"`
	PeriodEndTime   time.Time `json:"period_end_time"`
	TotalIncome     int       `json:"total_income"`
	TotalExpense    int       `json:"total_expense"`
	NetCashflow     int       `json:"net_cashflow"`
}

func (statisticwModule StatisticModule) incomeExpenseStatisticHandler(ctx *gin.Context) {
	var payload BasicStatisticQuery
	ctx.BindQuery(&payload)

	err := shared.Validator().Struct(payload)
	if err != nil {
		fmt.Println(err)
		errBody := shared.GenerateErrorResponse("BAD_REQ", shared.ParseValidatorError(err))
		ctx.JSON(400, errBody)
		return
	}

	startTime, _ := shared.StringToRFC3339(payload.PeriodStartTime)
	endTime, _ := shared.StringToRFC3339(payload.PeriodEndTime)

	userID := ctx.GetString(middlewares.UserIDKey)

	statistic, err := getUserIncomeExpenseStatistic(ctx, userID, startTime, endTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, shared.GenerateErrorResponse("INTERNALERR", nil))
		return
	}

	resp := IncomeExpenseStatisticResponse{
		PeriodStartTime: startTime,
		PeriodEndTime:   endTime,
		TotalIncome:     statistic.TotalIncome,
		TotalExpense:    statistic.TotalExpense,
		NetCashflow:     statistic.TotalIncome - statistic.TotalExpense,
	}

	ctx.JSON(200, resp)
}

type CashflowCategoryStatisticComponent struct {
	Percentage float32 `json:"percentage"`
	Amount     int     `json:"amount"`
}

type CashflowCategoryStatisticResponse struct {
	PeriodStartTime            time.Time                                     `json:"period_start_time"`
	PeriodEndTime              time.Time                                     `json:"period_end_time"`
	CashflowCategoryStatistics map[string]CashflowCategoryStatisticComponent `json:"statistics"`
}

func (statisticwModule StatisticModule) cashflowCategoryStatisticHandler(ctx *gin.Context) {
	var payload BasicStatisticQuery
	ctx.BindQuery(&payload)

	err := shared.Validator().Struct(payload)
	if err != nil {
		fmt.Println(err)
		errBody := shared.GenerateErrorResponse("BAD_REQ", shared.ParseValidatorError(err))
		ctx.JSON(400, errBody)
		return
	}

	startTime, _ := shared.StringToRFC3339(payload.PeriodStartTime)
	endTime, _ := shared.StringToRFC3339(payload.PeriodEndTime)

	// userID := ctx.GetString(middlewares.UserIDKey)

	// statistic, err := getUserCashflowCategoryStatistic(ctx, userID, startTime, endTime)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, shared.GenerateErrorResponse("INTERNALERR", nil))
	// 	return
	// }

	resp := CashflowCategoryStatisticResponse{
		PeriodStartTime: startTime,
		PeriodEndTime:   endTime,
	}

	ctx.JSON(200, resp)
}
