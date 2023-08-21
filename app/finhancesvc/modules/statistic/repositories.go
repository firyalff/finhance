package statistic

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx"
)

type IncomeExpenseStatisticDB struct {
	TotalIncome  int
	TotalExpense int
}

func getIncomeExpenseTotalByPeriod(ctx context.Context, userID string, startTime, endTime time.Time) (statistic IncomeExpenseStatisticDB, err error) {
	totalIncomeQuery := `SELECT COALESCE(SUM(amount), 0) 
	FROM cashflows 
	JOIN cashflow_categories ON category_id = cashflow_categories.id
	WHERE cashflow_category_type = 'income' AND cashflow_categories.user_id = $1 AND (cashflows.created_at BETWEEN $2 AND $3)`

	totalExpenseQuery := `SELECT COALESCE(SUM(amount), 0) 
	FROM cashflows 
	JOIN cashflow_categories ON category_id = cashflow_categories.id
	WHERE cashflow_category_type = 'expense' AND cashflow_categories.user_id = $1 AND (cashflows.created_at BETWEEN $2 AND $3)`

	query := "SELECT (" + totalIncomeQuery + ") total_income, (" + totalExpenseQuery + ") total_expense"
	row := StatisticModuleInstance.dbPool.QueryRow(ctx, query, userID, startTime, endTime)

	if err = row.Scan(&statistic.TotalIncome, &statistic.TotalExpense); err != nil {
		if err != pgx.ErrNoRows {
			log.Print(err)
		}
	}
	return
}

type CashflowCategoryStatisticDB struct {
	CashflowCategoryID   string
	CashflowCategoryName string
	TotalAmount          int
}

func getCashflowCategoriesotalByPeriod(ctx context.Context, userID string, startTime, endTime time.Time) (statistic IncomeExpenseStatisticDB, err error) {
	totalIncomeQuery := "SELECT COALESCE(SUM(amount), 0) FROM cashflows WHERE cashflow_type = 'income' AND user_id = $1 AND (created_at BETWEEN $2 AND $3)"
	totalExpenseQuery := "SELECT COALESCE(SUM(amount), 0) FROM cashflows WHERE cashflow_type = 'expense' AND user_id = $1 AND (created_at BETWEEN $2 AND $3)"
	query := "SELECT (" + totalIncomeQuery + ") total_income, (" + totalExpenseQuery + ") total_expense"
	row := StatisticModuleInstance.dbPool.QueryRow(ctx, query, userID, startTime, endTime)

	if err = row.Scan(&statistic.TotalIncome, &statistic.TotalExpense); err != nil {
		if err != pgx.ErrNoRows {
			log.Print(err)
		}
	}
	return
}
