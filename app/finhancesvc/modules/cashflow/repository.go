package cashflow

import (
	"context"
	"log"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v5"
)

type cashflowDB struct {
	Id               uuid.UUID
	CashflowType     string
	UserID           uuid.UUID
	Amount           int32
	ProofDocumentUrl string
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
