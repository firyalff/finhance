package cashflow

import (
	"context"
	"log"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v5"
)

type cashflowDB struct {
	Id               uuid.UUID
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
	query := "SELECT id, cashflow_type, amount, created_at FROM cashflows WHERE user_id=$1 ORDER BY id DESC LIMIT $2 OFFSET $3"

	rows, err := tx.Query(ctx, query, userID, limit, offset)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var row cashflowDB
		err := rows.Scan(&row.Id, &row.CashflowType, &row.Amount, &row.CreatedAt)
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

func getCashflowByUserIDandID(ctx context.Context, userID, cashflowID string) (cashflow cashflowDB, err error) {
	query := "SELECT id, cashflow_type, amount, proof_document_url, created_at, updated_at FROM cashflows WHERE user_id=$1 AND id=$2"
	row := CashflowModuleInstance.dbPool.QueryRow(ctx, query, userID, cashflowID)

	if err := row.Scan(&cashflow.Id, &cashflow.CashflowType, &cashflow.Amount, &cashflow.ProofDocumentUrl, &cashflow.CreatedAt, &cashflow.UpdatedAt); err != nil {
		if err != pgx.ErrNoRows {
			log.Print(err)
		}
		return cashflowDB{}, err
	}
	return
}

func createCashflow(ctx context.Context, userID string, cashflowData cashflowCreationPayload) (err error) {
	cashflowIDUUID, err := uuid.NewV7()
	if err != nil {
		log.Print(err)
		return
	}

	query := `INSERT INTO cashflows(id, user_id, amount, name, notes, cashflow_type, proof_document_url) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = CashflowModuleInstance.dbPool.Exec(ctx, query, cashflowIDUUID.String(), userID, cashflowData.Amount, cashflowData.Name, cashflowData.Notes, cashflowData.CashflowType, cashflowData.ProofDocumentUrl)
	if err != nil {
		log.Print(err)
	}
	return
}

func getUserCashflowByIDInTx(ctx context.Context, tx pgx.Tx, userID, cashflowID string) (cashflow cashflowDB, err error) {
	query := "SELECT id, cashflow_type, amount, proof_document_url, created_at, updated_at FROM cashflows WHERE user_id=$1 AND id=$2"
	row := tx.QueryRow(ctx, query, userID, cashflowID)

	if err := row.Scan(&cashflow.Id, &cashflow.CashflowType, &cashflow.Amount, &cashflow.ProofDocumentUrl, &cashflow.CreatedAt, &cashflow.UpdatedAt); err != nil {
		if err != pgx.ErrNoRows {
			log.Print(err)
		}
		return cashflowDB{}, err
	}
	return
}

func updateCashflowByID(ctx context.Context, tx pgx.Tx, cashflowID string, updateData cashflowUpdatePayload) (err error) {
	query := `UPDATE cashflows SET amount= $1, name= $2, notes= $3, cashflow_type= $4, proof_document_url= $5, updated_at= $6 WHERE id= $7`

	_, err = tx.Exec(ctx, query, updateData.Amount, updateData.Name, updateData.Notes, updateData.CashflowType, updateData.ProofDocumentUrl, time.Now(), cashflowID)
	if err != nil {
		log.Print(err)
	}
	return

}
