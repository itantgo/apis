package db

import (
	"database/sql"

	"github.com/itantgo/api/model"
	"github.com/jmoiron/sqlx"
	_"github.com/lib/pq"
)

type Config struct {
	ConnectString string
}

func InitDb(cfg Config) (*pgDb, error) {
	if dbConn, err := sqlx.Connect("postgres", cfg.ConnectString); err != nil {
		return nil, err
	} else {
		p := &pgDb{dbConn: dbConn}
		if err := p.dbConn.Ping(); err != nil {
			return nil, err
		}
		if err := p.createTablesIfNotExist(); err != nil {
			return nil, err
		}
		if err := p.prepareSqlStatements(); err != nil {
			return nil, err
		}
		return p, nil
	}
}

type pgDb struct {
	dbConn *sqlx.DB
	sqlSelectPeople *sqlx.Stmt
}

func (p *pgDb) createTablesIfNotExist() error {
	create_sql := `

       CREATE TABLE IF NOT EXISTS people (
       id SERIAL NOT NULL PRIMARY KEY,
       first TEXT NOT NULL,
       email TEXT NOT NULL);

    `
	if rows, err := p.dbConn.Query(create_sql); err != nil {
		return err
	} else {
		rows.Close()
	}
	return nil
}

func (p *pgDb) prepareSqlStatements() (err error) {

	if p.sqlSelectPeople, err = p.dbConn.Preparex(
		"SELECT id, first, email FROM people",
	); err != nil {
		return err
	}

	return nil
}