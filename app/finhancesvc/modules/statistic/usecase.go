package statistic

import (
	"context"
	"time"
)

func getUserIncomeExpenseStatistic(ctx context.Context, userID string, periodStartTime, periodEndTime time.Time) (statistic IncomeExpenseStatisticDB, err error) {
	return getIncomeExpenseTotalByPeriod(ctx, userID, periodStartTime, periodEndTime)
}

func getUserCashflowCategoryStatistics(ctx context.Context, userID, categoryType string, periodStartTime, periodEndTime time.Time) (statistic []CashflowCategoryStatisticDB, err error) {
	return getCashflowCategoriesotalByPeriod(ctx, userID, categoryType, periodStartTime, periodEndTime)
}
