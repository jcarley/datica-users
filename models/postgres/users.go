package postgres

import (
	"errors"
	"log"
	"strings"

	"github.com/jcarley/datica-users/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

var (
	ErrRecordNotFound = errors.New("Record Not Found")
)

type PostgresUserProvider struct {
}

func NewPostgresUserProvider() *PostgresUserProvider {
	return &PostgresUserProvider{}
}

func (this *PostgresUserProvider) AddUser(u models.User) (user models.User, e error) {

	userId := uuid.NewV4().String()
	createdAt, updatedAt := TimeStamps()

	u.UserId = userId
	u.CreatedAt = createdAt
	u.UpdatedAt = updatedAt

	statement := `insert into users (user_id, name, email, salt, password, created_at, updated_at)
							  values ($1, $2, $3, $4, $5, $6, $7)`

	tx := startTransaction()

	defer func(tx *sqlx.Tx) {
		if r := recover(); r != nil {
			// did we get a pq error
			if pqerr, ok := r.(*pq.Error); ok {
				// unique_violation
				errCodeName := strings.TrimSpace(pqerr.Code.Name())
				log.Println("pq error: ", errCodeName)
				rollbackTransaction(tx)
				user = models.User{}
				e = pqerr
			}
		}
	}(tx)

	result := tx.MustExec(statement,
		u.UserId,
		u.Name,
		u.Email,
		u.Salt,
		u.Password,
		u.CreatedAt,
		u.UpdatedAt,
	)

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Println("pq error: ", err.Code.Name())
			rollbackTransaction(tx)
			return models.User{}, err
		}
	}

	if rowsAffected > 0 {
		commitTransaction(tx)
	} else {
		rollbackTransaction(tx)
	}

	return u, nil
}

func (this *PostgresUserProvider) DeleteUser(userId string) (e error) {

	statement := `delete from users where users.user_id = $1`

	tx := startTransaction()

	defer func(tx *sqlx.Tx) {
		if r := recover(); r != nil {
			// did we get a pq error
			if pqerr, ok := r.(*pq.Error); ok {
				errCodeName := strings.TrimSpace(pqerr.Code.Name())
				log.Println("pq error: ", errCodeName)
				rollbackTransaction(tx)
				e = pqerr
			}
		}
	}(tx)

	result := tx.MustExec(statement, userId)

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Println("pq error: ", err.Code.Name())
			rollbackTransaction(tx)
			return err
		}
	}

	if rowsAffected > 0 {
		commitTransaction(tx)
	} else {
		rollbackTransaction(tx)
	}

	return nil
}

func (this *PostgresUserProvider) FindByUsername(username string) (models.User, error) {

	db := GetDB()

	statement := `SELECT user_id, name, email, salt, password, created_at, updated_at
								FROM users
								WHERE users.email = $1`

	user := models.User{}
	err := db.Get(&user, statement, username)

	if err != nil {
		return models.User{}, ErrRecordNotFound
	}

	return user, nil
}

func (this *PostgresUserProvider) FindUser(userId string) (models.User, error) {

	db := GetDB()

	statement := `SELECT user_id, name, email, salt, password, created_at, updated_at
								FROM users
								WHERE users.user_id = $1`

	user := models.User{}
	err := db.Get(&user, statement, userId)

	if err != nil {
		return models.User{}, ErrRecordNotFound
	}

	return user, nil
}

func (this *PostgresUserProvider) UpdateUser(userId string, u models.User) (user models.User, e error) {

	_, updatedAt := TimeStamps()

	u.UpdatedAt = updatedAt

	statement := `update users
								set name = $2,
								    updated_at = $3
	              where users.user_id = $1`

	tx := startTransaction()

	defer func(tx *sqlx.Tx) {
		if r := recover(); r != nil {
			// did we get a pq error
			if pqerr, ok := r.(*pq.Error); ok {
				// unique_violation
				errCodeName := strings.TrimSpace(pqerr.Code.Name())
				log.Println("pq error: ", errCodeName)
				rollbackTransaction(tx)
				user = models.User{}
				e = pqerr
			}
		}
	}(tx)

	result := tx.MustExec(statement,
		u.UserId,
		u.Name,
		u.UpdatedAt,
	)

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Println("pq error: ", err.Code.Name())
			rollbackTransaction(tx)
			return models.User{}, err
		}
	}

	if rowsAffected > 0 {
		commitTransaction(tx)
	} else {
		rollbackTransaction(tx)
	}

	return u, nil
}
