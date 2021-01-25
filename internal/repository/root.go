package repository

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	DB *sqlx.DB
}

func NewRepository() Repository {
	DB, err := sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		logrus.Fatal(err)
	}
	return Repository{DB: DB}
}

func (r Repository) InsertSessionCookie(value string) (err error) {
	statement := `INSERT INTO session (value) values($1)`
	_, err = r.DB.Exec(statement, value)
	if err != nil {
		logrus.Error(err)
	}
	return
}
