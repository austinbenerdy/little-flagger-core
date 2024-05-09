package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upCreateFeatureFlagTable, downCreateFeatureFlagTable)
}

func upCreateFeatureFlagTable(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return nil
}

func downCreateFeatureFlagTable(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
