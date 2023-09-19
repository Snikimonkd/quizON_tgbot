//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var RegistrationsDraft = newRegistrationsDraftTable("public", "registrations_draft", "")

type registrationsDraftTable struct {
	postgres.Table

	// Columns
	UserID    postgres.ColumnInteger
	GameID    postgres.ColumnInteger
	TeamID    postgres.ColumnInteger
	TeamName  postgres.ColumnString
	CreatedAt postgres.ColumnTimestampz
	UpdatedAt postgres.ColumnTimestampz

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type RegistrationsDraftTable struct {
	registrationsDraftTable

	EXCLUDED registrationsDraftTable
}

// AS creates new RegistrationsDraftTable with assigned alias
func (a RegistrationsDraftTable) AS(alias string) *RegistrationsDraftTable {
	return newRegistrationsDraftTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new RegistrationsDraftTable with assigned schema name
func (a RegistrationsDraftTable) FromSchema(schemaName string) *RegistrationsDraftTable {
	return newRegistrationsDraftTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new RegistrationsDraftTable with assigned table prefix
func (a RegistrationsDraftTable) WithPrefix(prefix string) *RegistrationsDraftTable {
	return newRegistrationsDraftTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new RegistrationsDraftTable with assigned table suffix
func (a RegistrationsDraftTable) WithSuffix(suffix string) *RegistrationsDraftTable {
	return newRegistrationsDraftTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newRegistrationsDraftTable(schemaName, tableName, alias string) *RegistrationsDraftTable {
	return &RegistrationsDraftTable{
		registrationsDraftTable: newRegistrationsDraftTableImpl(schemaName, tableName, alias),
		EXCLUDED:                newRegistrationsDraftTableImpl("", "excluded", ""),
	}
}

func newRegistrationsDraftTableImpl(schemaName, tableName, alias string) registrationsDraftTable {
	var (
		UserIDColumn    = postgres.IntegerColumn("user_id")
		GameIDColumn    = postgres.IntegerColumn("game_id")
		TeamIDColumn    = postgres.IntegerColumn("team_id")
		TeamNameColumn  = postgres.StringColumn("team_name")
		CreatedAtColumn = postgres.TimestampzColumn("created_at")
		UpdatedAtColumn = postgres.TimestampzColumn("updated_at")
		allColumns      = postgres.ColumnList{UserIDColumn, GameIDColumn, TeamIDColumn, TeamNameColumn, CreatedAtColumn, UpdatedAtColumn}
		mutableColumns  = postgres.ColumnList{UserIDColumn, GameIDColumn, TeamIDColumn, TeamNameColumn, CreatedAtColumn, UpdatedAtColumn}
	)

	return registrationsDraftTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		UserID:    UserIDColumn,
		GameID:    GameIDColumn,
		TeamID:    TeamIDColumn,
		TeamName:  TeamNameColumn,
		CreatedAt: CreatedAtColumn,
		UpdatedAt: UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}