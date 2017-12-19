package postgres

import (
	"log"
	"strings"

	"github.com/jcarley/datica-users/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

type PostgresTokenProvider struct {
}

func NewPostgresTokenProvider() *PostgresTokenProvider {
	return &PostgresTokenProvider{}
}

func (this *PostgresTokenProvider) AddToken(t models.Token) (token models.Token, e error) {

	tokenId := uuid.NewV4().String()

	t.TokenId = tokenId

	statement := `insert into tokens (token_id, email, token)
							  values ($1, $2, $3)`

	tx := startTransaction()

	defer func(tx *sqlx.Tx) {
		if r := recover(); r != nil {
			// did we get a pq error
			if pqerr, ok := r.(*pq.Error); ok {
				// unique_violation
				errCodeName := strings.TrimSpace(pqerr.Code.Name())
				log.Println("pq error: ", errCodeName)
				rollbackTransaction(tx)
				token = models.Token{}
				e = pqerr
			}
		}
	}(tx)

	result := tx.MustExec(statement,
		t.TokenId,
		t.Email,
		t.Token,
	)

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Println("pq error: ", err.Code.Name())
			rollbackTransaction(tx)
			return models.Token{}, err
		}
	}

	if rowsAffected > 0 {
		commitTransaction(tx)
	} else {
		rollbackTransaction(tx)
	}

	return t, nil
}

func (this *PostgresTokenProvider) FindToken(token string) (models.Token, error) {

	db := GetDB()

	statement := `SELECT token_id, email, token
								FROM tokens
								WHERE tokens.token = $1`

	t := models.Token{}
	err := db.Get(&t, statement, token)

	if err != nil {
		return models.Token{}, ErrRecordNotFound
	}

	return t, nil
}
