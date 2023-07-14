package auth

import (
	"context"
	"log"
	"time"

	"github.com/gofrs/uuid"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v5"
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
		if err != pgx.ErrNoRows {
			log.Print(err)
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

func countUserByEmail(ctx context.Context, email string) (totalRow int, err error) {
	query := "SELECT count(id) FROM users WHERE email=$1"

	row := authModuleInstance.dbPool.QueryRow(ctx, query, email)
	err = row.Scan(&totalRow)

	return
}

func createUser(ctx context.Context, registrationRecord registerValidation) (err error) {
	userID, err := uuid.NewV7()
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registrationRecord.Password), 8)
	if err != nil {
		return err
	}

	query := `INSERT INTO users(id, first_name, last_name, email, password) VALUES ($1, $2, $3, $4, $5)`
	_, err = authModuleInstance.dbPool.Exec(ctx, query, userID.String(), registrationRecord.FirstName, registrationRecord.LastName, registrationRecord.Email, hashedPassword)

	return
}
