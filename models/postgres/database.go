package postgres

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/hashicorp/hcl"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	db *sqlx.DB
)

type DatabaseInfo struct {
	Username     string
	Password     string
	Host         string
	DatabaseName string
}

func ReadConfig(r io.Reader) (DatabaseInfo, error) {
	buffer, err := ioutil.ReadAll(r)
	if err != nil {
		return DatabaseInfo{}, err
	}

	var d DatabaseInfo
	err = hcl.Unmarshal(buffer, &d)
	if err != nil {
		return DatabaseInfo{}, err
	}

	return d, nil
}

func Connect(d DatabaseInfo) {

	if d.Username == "" {
		log.Fatalln("Username is required to connect to database")
	}

	if d.Password == "" {
		log.Fatalln("Password is required to connect to database")
	}

	if d.Host == "" {
		log.Fatalln("Host is required to connect to database")
	}

	if d.DatabaseName == "" {
		log.Fatalln("DatabaseName is required to connect to database")
	}

	connectionString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", d.Username, d.Password, d.Host, d.DatabaseName)
	connection, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	log.Printf("Successfully connected to database: %s", d.DatabaseName)
	db = connection
}

func Disconnect() error {
	log.Printf("Disconnected from database")
	return db.Close()
}

func GetDB() *sqlx.DB {
	return db
}

func startTransaction() *sqlx.Tx {
	db := GetDB()
	return db.MustBegin()
}

func commitTransaction(tx *sqlx.Tx) error {
	return tx.Commit()
}

func rollbackTransaction(tx *sqlx.Tx) error {
	return tx.Rollback()
}

func NewDbTime() time.Time {
	return time.Now().UTC()
}

func NewFormattedDbTime() string {
	t := NewDbTime()
	return FormatDbTime(t)
}

func FromFormattedDbTime(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}

func FormatDbTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func TimeStamps() (createdAt time.Time, updatedAt time.Time) {
	newTime := NewDbTime()
	createdAt = newTime
	updatedAt = newTime
	return
}
