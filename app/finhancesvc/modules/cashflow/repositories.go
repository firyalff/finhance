package cashflow

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v5"
)

type cashflowDB struct {
	Id               uuid.UUID
	CategoryID       *string
	CategoryName     *string
	CashflowType     string
	UserID           uuid.UUID
	Amount           int32
	ProofDocumentUrl *string
	CreatedAt        pgtype.Timestamptz
	UpdatedAt        pgtype.Timestamptz
}

func countCashflowByUserID(ctx context.Context, tx pgx.Tx, userID string) (total int, err error) {
	query := "SELECT count(id) FROM cashflows WHERE user_id=$1"

	row := tx.QueryRow(ctx, query, userID)
	err = row.Scan(&total)
	if err != nil {
		log.Print(err)
	}

	return
}

func getCashflowsByUserID(ctx context.Context, tx pgx.Tx, userID string, limit, offset int32) (cashflows []cashflowDB, err error) {
	query := `SELECT cashflows.id, category_id, cashflow_categories.name category_name, cashflow_categories.cashflow_category_type cashflow_type, amount, cashflows.created_at FROM cashflows 
	JOIN cashflow_categories ON category_id = cashflow_categories.id 
	WHERE cashflows.user_id=$1 ORDER BY id DESC LIMIT $2 OFFSET $3`

	rows, err := tx.Query(ctx, query, userID, limit, offset)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var row cashflowDB
		err := rows.Scan(&row.Id, &row.CategoryID, &row.CategoryName, &row.CashflowType, &row.Amount, &row.CreatedAt)
		if err != nil {
			log.Print(err)
		}
		cashflows = append(cashflows, row)
	}

	if err := rows.Err(); err != nil {
		log.Print(err)
	}

	return
}

func getUserCashflowByID(ctx context.Context, userID, cashflowID string) (cashflow cashflowDB, err error) {
	query := `SELECT cashflows.id, category_id, cashflow_categories.name category_name, cashflow_categories.cashflow_category_type cashflow_type, amount, proof_document_url, cashflows.created_at, cashflows.updated_at 
	FROM cashflows 
	JOIN cashflow_categories ON category_id = cashflow_categories.id 
	WHERE cashflows.user_id=$1 AND cashflows.id=$2`
	row := CashflowModuleInstance.dbPool.QueryRow(ctx, query, userID, cashflowID)

	if err := row.Scan(&cashflow.Id, &cashflow.CategoryID, &cashflow.CategoryName, &cashflow.CashflowType, &cashflow.Amount, &cashflow.ProofDocumentUrl, &cashflow.CreatedAt, &cashflow.UpdatedAt); err != nil {
		if err != pgx.ErrNoRows {
			log.Print(err)
		}
		return cashflowDB{}, err
	}
	return
}

func createCashflow(ctx context.Context, tx pgx.Tx, userID string, cashflowData cashflowCreationPayload) (err error) {
	cashflowIDUUID, err := uuid.NewV7()
	if err != nil {
		log.Print(err)
		return
	}

	query := `INSERT INTO cashflows(id, user_id, amount, name, notes, category_id, proof_document_url) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = CashflowModuleInstance.dbPool.Exec(ctx, query, cashflowIDUUID.String(), userID, cashflowData.Amount, cashflowData.Name, cashflowData.Notes, cashflowData.CashflowCategoryID, cashflowData.ProofDocumentUrl)
	if err != nil {
		log.Print(err)
	}
	return
}

func getUserCashflowByIDInTx(ctx context.Context, tx pgx.Tx, userID, cashflowID string) (cashflow cashflowDB, err error) {
	query := `SELECT id, category_id, cashflow_categories.name category_name, amount, proof_document_url, created_at, updated_at 
	FROM cashflows 
	JOIN cashflow_categories ON category_id = cashflow_categories.id
	WHERE cashflows.user_id=$1 AND id=$2`
	row := tx.QueryRow(ctx, query, userID, cashflowID)

	if err := row.Scan(&cashflow.Id, &cashflow.CategoryID, &cashflow.CategoryName, &cashflow.Amount, &cashflow.ProofDocumentUrl, &cashflow.CreatedAt, &cashflow.UpdatedAt); err != nil {
		if err != pgx.ErrNoRows {
			log.Print(err)
		}
		return cashflowDB{}, err
	}
	return
}

func updateCashflowByID(ctx context.Context, tx pgx.Tx, cashflowID string, updateData cashflowUpdatePayload) (err error) {
	query := `UPDATE cashflows SET amount= $1, name= $2, notes= $3, category_id= $4, proof_document_url= $5, updated_at= $6 WHERE id= $7`

	_, err = tx.Exec(ctx, query, updateData.Amount, updateData.Name, updateData.Notes, updateData.CashflowCategoryID, updateData.ProofDocumentUrl, time.Now(), cashflowID)
	if err != nil {
		log.Print(err)
	}
	return

}

type cashflowCategoryDB struct {
	Id                   uuid.UUID
	UserID               uuid.UUID
	CashflowCategoryType string
	Name                 string
	CreatedAt            pgtype.Timestamptz
	UpdatedAt            pgtype.Timestamptz
	DeletedAt            pgtype.Timestamptz
}

func createCashflowCategory(ctx context.Context, userID string, categoryData cashflowCategoryCreatePayload) (err error) {
	categoryIDUUID, err := uuid.NewV7()
	if err != nil {
		log.Print(err)
		return
	}

	query := `INSERT INTO cashflow_categories(id, user_id, cashflow_category_type, name) VALUES ($1, $2, $3, $4)`
	_, err = CashflowModuleInstance.dbPool.Exec(ctx, query, categoryIDUUID.String(), userID, categoryData.CategoryType, strings.ToLower(categoryData.Name))
	if err != nil {
		log.Print(err)
	}
	return
}

func countCashflowCategoriesByUserID(ctx context.Context, tx pgx.Tx, userID string) (total int, err error) {
	query := "SELECT count(id) FROM cashflow_categories WHERE user_id=$1 AND deleted_at IS NULL"

	row := tx.QueryRow(ctx, query, userID)
	err = row.Scan(&total)
	if err != nil {
		log.Print(err)
	}

	return
}

func countUserCashflowCategoriesByID(ctx context.Context, tx pgx.Tx, userID, categoryID string) (total int, err error) {
	query := "SELECT count(id) FROM cashflow_categories WHERE id=$1 AND user_id=$2 AND deleted_at IS NULL"

	row := tx.QueryRow(ctx, query, categoryID, userID)
	err = row.Scan(&total)
	if err != nil {
		log.Print(err)
	}

	return
}

func getCashflowCategoriesByUserID(ctx context.Context, tx pgx.Tx, userID string, limit, offset int32) (cashflowCategories []cashflowCategoryDB, err error) {
	query := "SELECT id, cashflow_category_type, name, created_at FROM cashflow_categories WHERE user_id=$1 AND deleted_at IS NULL ORDER BY id DESC LIMIT $2 OFFSET $3"

	rows, err := tx.Query(ctx, query, userID, limit, offset)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var row cashflowCategoryDB
		err := rows.Scan(&row.Id, &row.CashflowCategoryType, &row.Name, &row.CreatedAt)
		if err != nil {
			log.Print(err)
		}
		cashflowCategories = append(cashflowCategories, row)
	}

	if err := rows.Err(); err != nil {
		log.Print(err)
	}

	return
}
