package auth

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type userDB struct {
	Id        uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Password  string
	Status    string
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

func getUserByEmail(ctx context.Context, email string) (user userDB, err error) {
	query := "SELECT * FROM users WHERE email=$1"
	row := authModuleInstance.dbPool.QueryRow(ctx, query, email)
	if err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Status, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			log.Debug().Err(err).Msg("can't find any item")
		}
		return userDB{}, err
	}
	return
}

func validateUserPassword(inputPassword, userPassword string) (err error) {
	return bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(inputPassword))
}

type UserJWTClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func generateJWT(userID string, secret string, dayToExpire int) (string, error) {
	claims := UserJWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{time.Now().AddDate(0, 0, dayToExpire)},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
