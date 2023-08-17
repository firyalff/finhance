package cashflow

import (
	"context"
	"finhancesvc/shared"

	"github.com/jackc/pgx"
)

func getUserCashflows(ctx context.Context, userID string, limit, offset int32) (cashflows []cashflowDB, totalCashflow int, err error) {
	tx, err := CashflowModuleInstance.dbPool.Begin(ctx)
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
		tx.Commit(ctx)
	}()

	if err != nil {
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
	tx, err := CashflowModuleInstance.dbPool.Begin(ctx)
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
		tx.Commit(ctx)
	}()
	if err != nil {
		return shared.ErrInternal
	}

	total, err := countUserCashflowCategoriesByID(ctx, tx, userID, cashflowData.CashflowCategoryID)
	if err != nil {
		return shared.ErrInternal
	}
	if total < 1 {
		return shared.ErrNotFound
	}

	return createCashflow(ctx, tx, userID, cashflowData)
}

func updateUserCashflow(ctx context.Context, userID, cashflowID string, payload cashflowUpdatePayload) (err error) {
	tx, err := CashflowModuleInstance.dbPool.Begin(ctx)

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
		tx.Commit(ctx)
	}()

	_, err = getUserCashflowByIDInTx(ctx, tx, userID, cashflowID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return shared.ErrNotFound
		}
		return
	}

	err = updateCashflowByID(ctx, tx, cashflowID, payload)

	return
}

func createUserCashflowCategory(ctx context.Context, userID string, payload cashflowCategoryCreatePayload) (err error) {
	err = createCashflowCategory(ctx, userID, payload)
	if shared.PGDuplicationError(err) {
		err = shared.ErrExist
	}
	return
}

func getUserCashflowCategories(ctx context.Context, userID string, limit, offset int32) (cashflows []cashflowCategoryDB, totalCashflowCategories int, err error) {
	tx, err := CashflowModuleInstance.dbPool.Begin(ctx)
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
		tx.Commit(ctx)
	}()

	if err != nil {

		return []cashflowCategoryDB{}, 0, shared.ErrInternal
	}

	totalCashflowCategories, err = countCashflowCategoriesByUserID(ctx, tx, userID)
	if err != nil {
		return []cashflowCategoryDB{}, 0, shared.ErrInternal
	}

	cashflows, err = getCashflowCategoriesByUserID(ctx, tx, userID, limit, offset)
	if err != nil {
		return []cashflowCategoryDB{}, 0, shared.ErrInternal
	}

	return
}
