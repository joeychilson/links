package database

import (
	"embed"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate(dbURL string) error {
	dburl, err := url.Parse(dbURL)
	if err != nil {
		return err
	}

	dbmate := dbmate.New(dburl)
	dbmate.FS = migrations
	dbmate.MigrationsDir = []string{"migrations"}
	dbmate.AutoDumpSchema = false

	if _, err = dbmate.Status(false); err != nil {
		return err
	}

	err = dbmate.CreateAndMigrate()
	if err != nil {
		return err
	}
	return nil
}
