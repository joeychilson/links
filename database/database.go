package database

import (
	"context"
	"embed"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
)

//go:embed migrations/*.sql
var migrations embed.FS

type DBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

type Queries struct {
	db DBTX
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Migrate(dbURL string) error {
	dburl, err := url.Parse(dbURL)
	if err != nil {
		return err
	}

	dbmate := dbmate.New(dburl)
	dbmate.FS = migrations
	dbmate.MigrationsDir = []string{"migrations"}
	dbmate.AutoDumpSchema = false

	if _, err = dbmate.Status(true); err != nil {
		return err
	}

	err = dbmate.CreateAndMigrate()
	if err != nil {
		return err
	}
	return nil
}

func (q *Queries) WithTx(tx pgx.Tx) *Queries {
	return &Queries{
		db: tx,
	}
}
