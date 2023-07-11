package auth

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type userDB struct {
	Id        uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Password  string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func getUserByEmail(ctx context.Context, dbPool *pgxpool.Pool, email string) (user userDB, err error) {
	query := "SELECT * FROM users WHERE email=$1"
	row := dbPool.QueryRow(ctx, query, email)
	if err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Status, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("can't find any item")
		}
		return userDB{}, err
	}
	return
}

func validateUserPassword(inputPassword, userPassword string) (isValid bool, err error) {
	return
}
