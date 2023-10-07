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
	CategoryID   string
	CategoryName string
	CategoryType string
	TotalAmount  int
}

func getCashflowCategoriesotalByPeriod(ctx context.Context, userID, cashflowType string, startTime, endTime time.Time) (statistics []CashflowCategoryStatisticDB, err error) {
	categoryTypeQuery := ""

	queryArgs := []any{userID, startTime, endTime}

	if cashflowType != "" {
		categoryTypeQuery = " AND cashflow_category_type = $4 "
		queryArgs = append(queryArgs, cashflowType)
	}

	query := `SELECT category_id , cashflow_categories."name" category_name, cashflow_category_type category_type, COALESCE(sum(amount), 0) total_amount
	FROM cashflows 
	JOIN cashflow_categories ON cashflow_categories.id = category_id 
	WHERE cashflow_categories.user_id = $1 AND (cashflows.created_at BETWEEN $2 AND $3) ` + categoryTypeQuery + `
	GROUP BY category_id, cashflow_categories."name", cashflow_category_type;`

	rows, err := StatisticModuleInstance.dbPool.Query(ctx, query, queryArgs...)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var row CashflowCategoryStatisticDB
		err := rows.Scan(&row.CategoryID, &row.CategoryName, &row.CategoryType, &row.TotalAmount)
		if err != nil {
			log.Print(err)
		}
		statistics = append(statistics, row)
	}

	if err := rows.Err(); err != nil {
		log.Print(err)
	}

	return
}
