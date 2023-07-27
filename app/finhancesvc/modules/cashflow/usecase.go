package cashflow

import (
	"context"
	"finhancesvc/shared"
	"log"
)

func getUserCashflows(ctx context.Context, userID string, limit, offset int32) (cashflows []cashflowDB, totalCashflow int, err error) {
	tx, err := CashflowModuleInstance.dbPool.Begin(ctx)
	if err != nil {
		log.Print(err)
		return []cashflowDB{}, 0, shared.ErrInternal
	}

	totalCashflow, err = countCashflowByUserID(ctx, tx, userID)
	if err != nil {
		return []cashflowDB{}, 0, shared.ErrInternal
	}

	cashflows, err = getCashflowsByUserID(ctx, tx, userID, limit, offset)
	if err != nil {
		return []cashflowDB{}, 0, shared.ErrInternal
	}

	return
}
