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

var Admins = newAdminsTable("public", "admins", "")

type adminsTable struct {
	postgres.Table

	// Columns
	ID        postgres.ColumnInteger
	DateUntil postgres.ColumnTimestampz

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type AdminsTable struct {
	adminsTable

	EXCLUDED adminsTable
}

// AS creates new AdminsTable with assigned alias
func (a AdminsTable) AS(alias string) *AdminsTable {
	return newAdminsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new AdminsTable with assigned schema name
func (a AdminsTable) FromSchema(schemaName string) *AdminsTable {
	return newAdminsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new AdminsTable with assigned table prefix
func (a AdminsTable) WithPrefix(prefix string) *AdminsTable {
	return newAdminsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new AdminsTable with assigned table suffix
func (a AdminsTable) WithSuffix(suffix string) *AdminsTable {
	return newAdminsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newAdminsTable(schemaName, tableName, alias string) *AdminsTable {
	return &AdminsTable{
		adminsTable: newAdminsTableImpl(schemaName, tableName, alias),
		EXCLUDED:    newAdminsTableImpl("", "excluded", ""),
	}
}

func newAdminsTableImpl(schemaName, tableName, alias string) adminsTable {
	var (
		IDColumn        = postgres.IntegerColumn("id")
		DateUntilColumn = postgres.TimestampzColumn("date_until")
		allColumns      = postgres.ColumnList{IDColumn, DateUntilColumn}
		mutableColumns  = postgres.ColumnList{IDColumn, DateUntilColumn}
	)

	return adminsTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:        IDColumn,
		DateUntil: DateUntilColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
