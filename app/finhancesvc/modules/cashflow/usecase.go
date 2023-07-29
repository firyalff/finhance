package cashflow

import (
	"context"
	"finhancesvc/shared"
	"log"

	"github.com/jackc/pgx"
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

func getUserCashflowByID(ctx context.Context, userID, cashflowID string) (cashflow cashflowDB, err error) {
	cashflow, err = getCashflowByUserIDandID(ctx, userID, cashflowID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return cashflow, shared.ErrNotFound
		}
		return
	}

	return
}

func createNewUserCashflow(ctx context.Context, userID string, cashflowData cashflowCreationPayload) (err error) {
	return createCashflow(ctx, userID, cashflowData)
}
