package testutil

import (
	"database/sql"
	"testing"

	"github.com/jinzhu/gorm"
)

// DbEntity godoc
type DbEntity struct {
	TableName string
	Keyname   string
	Key       interface{}
}

// SetupDBSandbox godoc
func SetupDBSandbox(t *testing.T, db *gorm.DB) {
	const hookName = "dbCleanup"

	entries := []DbEntity{}

	db.LogMode(false)
	db.Callback().Create().After("gorm:create").
		Register(hookName, func(scope *gorm.Scope) {
			entries = append(entries, DbEntity{
				TableName: scope.TableName(),
				Keyname:   scope.PrimaryKey(),
				Key:       scope.PrimaryKeyValue(),
			})
		})

	t.Cleanup(func() {
		defer db.Close()
		defer db.Callback().Create().Remove(hookName)

		_, inTransaction := db.CommonDB().(*sql.Tx)
		tx := db
		if !inTransaction {
			tx = db.Begin()
		}

		for _, entry := range entries {
			tx.Table(entry.TableName).Where(entry.Keyname+" = ?", entry.Key).Delete("")
		}

		if !inTransaction {
			tx.Commit()
		}
	})
}
