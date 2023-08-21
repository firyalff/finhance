package statistic

import (
	"context"
	"time"
)

func getUserIncomeExpenseStatistic(ctx context.Context, userID string, periodStartTime, periodEndTime time.Time) (statistic IncomeExpenseStatisticDB, err error) {
	return getIncomeExpenseTotalByPeriod(ctx, userID, periodStartTime, periodEndTime)
}

func getUserCashflowCategoryStatistic(ctx context.Context, userID string, periodStartTime, periodEndTime time.Time) (statistic IncomeExpenseStatisticDB, err error) {
	return getCashflowCategoriesotalByPeriod(ctx, userID, periodStartTime, periodEndTime)
}
